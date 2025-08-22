package download

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	p "path"
	"strconv"
	"time"
)

var STATUS_ERROR error = fmt.Errorf("check status")

type Download struct {
	Url        url.URL
	Length     uint
	Downloaded chan uint
	done       chan struct{}
	Status     string
	StatusCode int
}

func (p Download) IsDone() bool {
	select {
	case <-p.done:
		return true
	default:
		return false
	}
}

func Get(u url.URL, path string, limit int) (Download, error) {
	path = UniqueFilename(path)
	fmt.Printf("sending request, awaiting response...")
	resp, err := http.Get(u.String())
	if err != nil {
		return Download{}, fmt.Errorf("%s : %s", err, &u)
	}
	fmt.Printf("status %s\n", resp.Status)
	if resp.StatusCode != http.StatusOK {
		ds := Download{Length: 0, Downloaded: nil, Status: resp.Status, StatusCode: resp.StatusCode, Url: u}
		return ds, STATUS_ERROR
	}

	sizeStr := resp.Header.Get("Content-Length")
	sizeInt, err := strconv.ParseUint(sizeStr, 10, 64)
	var length uint
	if err != nil || sizeStr == "" {
		length = 1
	} else {
		length = uint(sizeInt)
	}

	progressChan := make(chan uint)
	status := Download{
		Url:        u,
		Length:     length,
		Downloaded: progressChan,
		done:       make(chan struct{}),
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
	}
	os.MkdirAll(p.Dir(path), os.ModePerm)

	go download(&status, resp, path, uint(limit))

	return status, nil
}

func download(d *Download, r *http.Response, path string, limit uint) {
	if limit < 1 {
		limit = 32 * 1024
	}
	bucket := limit
	delay := time.Now()
	defer close(d.Downloaded)
	defer r.Body.Close()
	defer close(d.done)
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	buf := make([]byte, bucket) // 32 KB
	var total uint
	for {
		n, err := r.Body.Read(buf)
		if n > 0 {
			_, writeErr := file.Write(buf[:n])
			if writeErr != nil {
				fmt.Println(writeErr)
				return
			}
			total += uint(n)

			if d.Length > 1 {
				d.Downloaded <- total
			}
		}
		bucket -= uint(n)

		if err != nil {
			break
		}
		if time.Now().After(delay.Add(time.Second)) {
			bucket = limit
		}
	}

	if d.Length == 1 {
		d.Downloaded <- 1
	}
}
