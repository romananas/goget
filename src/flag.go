package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

type options struct {
	RateLimit  int
	Urls       []string
	Mirror     bool
	Directory  string
	OutputFile string
	Background bool
	InputFile  string
}

func ParseCLA() options {
	var opts options
	flag.CommandLine.SortFlags = true

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "CyberPangolin's Goget V0.0.1\n\n")
		fmt.Fprintf(os.Stderr, "Usage: GoGet <OPTION> [URLS]\n\n")
		flag.PrintDefaults()
	}
	flag.IntVar(&opts.RateLimit, "rate-limit", 0, "limit download rate by x files/sec")
	flag.BoolVarP(&opts.Mirror, "mirror", "m", false, "mirror the whole website starting from the base url (can't parse multiple urls)")
	flag.BoolVarP(&opts.Background, "background", "B", false, "launche the download in background")
	flag.StringVarP(&opts.Directory, "path", "P", "", "download all file on the choosen directory")
	flag.StringVarP(&opts.OutputFile, "output-file", "O", "", "download all result in a choosen file (all values will be written in the same files)")
	flag.StringVarP(&opts.InputFile, "input-file", "i", "", "Read URLs from a local or external file.  If - is specified as file, URLs are read from the standard input.  Use ./- to read from a file literally named -.")
	flag.Parse()
	opts.Urls = flag.Args()
	return opts
}
