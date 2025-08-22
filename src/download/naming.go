package download

import (
	"net/url"
	"path"
)

func IsFile(u url.URL) bool {
	return path.Ext(u.Path) != ""
}

// func ExtractDirs(u url.URL) (dirs []string, file string) {
// }
