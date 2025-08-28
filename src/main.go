package main

import (
	"fmt"
	"main/download"
	"main/progress"
	"main/scraping"
	"net/url"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

func startup() error {
	fmt.Printf("start at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

func end() {
	fmt.Printf("finished at %s\n", time.Now().Format("2006-01-02 15:04:05"))
}

func ParseInputFile(opts *options) error {
	file, err := os.OpenFile(opts.InputFile, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	info, err := os.Stat(opts.InputFile)
	if err != nil {
		return err
	}
	var buff = make([]byte, info.Size())
	_, err = file.Read(buff)
	if err != nil {
		return err
	}
	parts := strings.Fields(string(buff))
	opts.Urls = append(opts.Urls, parts...)
	return nil
}
func ParseUrls(opts options) ([]url.URL, error) {
	var urls []url.URL
	for _, surl := range opts.Urls {
		u, err := url.Parse(surl)
		if err != nil {
			return nil, err
		}
		urls = append(urls, *u)
	}
	return urls, nil
}

func SimpleDownload(urls []url.URL, rate int, fname string) error {
	var manager progress.Manager[uint] = progress.New[uint](50, "=>-")
	var current []download.Download
	for _, u := range urls {
		var p string
		if fname != "" {
			p = fname
		} else {
			p = path.Base(u.Path)
		}
		status, err := download.Get(u, p, rate)
		if err != nil && err != download.STATUS_ERROR {
			return err
		}
		current = append(current, status)
	}
	for _, c := range current {
		if c.StatusCode == 200 {
			manager.Add(c.Downloaded, c.Length, c.Url.Path)
		}
	}
	for !download.AllDone(current) {
	}
	return nil
}

func main() {
	opts := ParseCLA()
	if len(opts.Urls) == 0 && opts.InputFile == "" {
		fmt.Println("\033[31mnothing to do\033[0m")
		return
	}
	if opts.Background && !IsBackground() {
		if pid, err := IntoBackground(); err == nil {
			fmt.Fprintf(os.Stdout, "Continuing in background, pid %d\n", pid)
			os.Exit(0)
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if IsBackground() {
		var err error
		if runtime.GOOS == "windows" {
			err = progress.SetOutput("NUL")
		} else {
			err = progress.SetOutput("/dev/null")
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	if len(opts.InputFile) > 0 {
		if err := ParseInputFile(&opts); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}
	if len(opts.Directory) > 0 {
		if err := os.Chdir(opts.Directory); err != nil {
			fmt.Fprintf(os.Stderr, "failed to change directory: %v\n", err)
			os.Exit(1)
		}
	}
	urls, err := ParseUrls(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	startup()
	switch {
	case opts.Mirror == true:
		if err := scraping.Scrap(urls); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		if err = SimpleDownload(urls, opts.RateLimit, opts.OutputFile); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	end()
}
