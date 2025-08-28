package scraping

import (
	"main/download"
	"net/url"
	p "path"
	"time"
)

type scraps struct {
	scrap []*scrap
}

func Init() scraps {
	return scraps{scrap: make([]*scrap, 0)}
}

func (self *scraps) Get(url url.URL) *scrap {
	for _, scrap := range self.scrap {
		if url == scrap.url {
			return scrap
		}
	}
	return nil
}

func (self *scraps) Delete(url url.URL) {
	for i, scrap := range self.scrap {
		if url == scrap.url {
			self.scrap = append(self.scrap[:i], self.scrap[i+1:]...)
		}
	}
}

func (self *scraps) Add(url url.URL) {
	var file string = p.Base(url.Path)
	var dirs string = p.Dir(url.Path)
	if p.Ext(file) == "" {
		// ! dirs = p.Join(dirs, file)
		file = "index.html"
	}
	if self.Get(url) == nil {
		self.scrap = append(self.scrap, &scrap{url: url, state: PENDING, file: file, dir: dirs})
	}
}

func (self *scraps) IsDone() bool {
	for _, x := range self.scrap {
		if x.state != DONE && x.state != PARSED {
			return false
		}
	}
	return true
}

func (self *scraps) IsFullDone() bool {
	for _, x := range self.scrap {
		if x.state != PARSED && x.state != ERROR {
			return false
		}
	}
	return true
}

func (self *scraps) GetUnparsed() []*scrap {
	var unparsed []*scrap = make([]*scrap, 0)
	for _, x := range self.scrap {
		if x.state == DONE {
			unparsed = append(unparsed, x)
		}
	}
	return unparsed
}

func (self *scraps) Count(state int) int {
	count := 0
	for _, x := range self.scrap {
		if x.state == state {
			count++
		}
	}
	return count
}

func (self *scraps) Download() ([]download.Download, error) {
	var Downloades []download.Download = make([]download.Download, 0)
	statuesMap := make(map[*scrap]download.Download)

	for _, x := range self.scrap {
		if x.state == PENDING {
			Download, err := x.download()
			if err != nil {
				if err == download.STATUS_ERROR {
					// ignorer les erreurs HTTP (ex: 404), mais marquer comme fini
					x.state = ERROR
					continue
				}
				return nil, err
			}

			x.state = DOWNLOADING
			statuesMap[x] = Download

			if Download.StatusCode != 200 {
				x.state = DONE
				continue
			}

			Downloades = append(Downloades, Download)
		}
	}

	go func() {
		for len(statuesMap) > 0 {
			for x, status := range statuesMap {
				if status.IsDone() {
					x.state = DONE
					delete(statuesMap, x)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return Downloades, nil
}

func (self *scraps) parse() ([]string, error) {
	var parsed []string = make([]string, 0)
	for _, unparsed := range self.GetUnparsed() {
		part, err := unparsed.parse()
		if err != nil {
			return nil, err
		}
		validated_part := []string{}
		for _, unvalidated := range part {
			validated, err := ValidatePath(unparsed.url.String(), unvalidated)
			if err != nil {
				return nil, err
			}
			validated_part = append(validated_part, *validated)
		}
		parsed = append(parsed, validated_part...)
	}
	return parsed, nil
}
