package scraping

import (
	"main/download"
	"net/url"
	"path"
)

const (
	WAITING = iota
	DOWNLOADING
	DONE
	PARSED
	ERROR
)

// SCRAP PART

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
