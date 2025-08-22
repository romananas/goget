package download

import (
	"fmt"
	"os"
	"path/filepath"
)

func AllDone(d []Download) bool {
	for _, x := range d {
		if !x.IsDone() {
			return false
		}
	}
	return true
}

// Génère un nom unique à la wget
func UniqueFilename(path string) string {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	candidate := filepath.Join(dir, base)
	i := 1
	for {
		if _, err := os.Stat(candidate); os.IsNotExist(err) {
			return candidate
		}
		candidate = filepath.Join(dir, fmt.Sprintf("%s%s.%d", name, ext, i))
		i++
	}
}

// Applique sur un tableau de chemins
func DeduplicateFilenames(paths []string) []string {
	result := make([]string, len(paths))
	for i, p := range paths {
		result[i] = UniqueFilename(p)
	}
	return result
}
