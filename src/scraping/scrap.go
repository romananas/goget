package scraping

import (
	"main/download"
	"net/url"
	"path"
)

const (
	PENDING = iota
	DOWNLOADING
	DONE
	PARSED
	ERROR
)

// SCRAP PART

// scrap represents a single scraping task, containing the target URL, the file name to save the content,
// the directory where the file will be stored, and the current state of the scraping process.
type scrap struct {
	url   url.URL
	file  string
	dir   string
	state int
}

func (self *scrap) download() (download.Download, error) {
	downloadStatus, err := download.Get(self.url, path.Join(".", self.dir, self.file), 0)
	if err != nil {
		return downloadStatus, err
	}
	return downloadStatus, nil
}

// parse reads and parses the file specified by the scrap struct's directory and file fields.
// It returns a slice of strings containing the parsed data, or an error if the operation fails.
// If the scrap is in the ERROR state, it returns nil without error.
// On successful parsing, it updates the scrap's state to PARSED.
func (self *scrap) parse() ([]string, error) {
	if self.state == ERROR {
		return nil, nil
	}
	data, err := GetPaths(path.Join(".", self.dir, self.file))
	if err != nil {
		return nil, err
	}
	self.state = PARSED
	return data, nil
}
