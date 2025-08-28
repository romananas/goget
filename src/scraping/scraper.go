package scraping

import (
	"main/progress"
	"net/url"
	u "net/url"
	"strings"
)

// IntoAbsolute resolves a relative URL against a base (current) URL and returns the absolute URL as a string pointer.
// If the relative URL is already absolute, it returns the original relative URL.
// Returns an error if either URL cannot be parsed.
//
// Parameters:
//   - current: The base URL as a string.
//   - relatif: The relative or absolute URL as a string.
//
// Returns:
//   - *string: Pointer to the resulting absolute URL string.
//   - error:   Error if URL parsing fails.
func IntoAbsolute(current string, relatif string) (*string, error) {
	relatif_url, err := u.Parse(relatif)
	if err != nil {
		return nil, err
	}
	if relatif_url.IsAbs() {
		return &relatif, nil
	}
	current_url, err := u.Parse(current)
	if err != nil {
		return nil, err
	}
	absolute := current_url.ResolveReference(relatif_url)
	absoluteStr := absolute.String()
	return &absoluteStr, nil
}

// ValidatePath checks if the provided path is either relative or belongs to the same domain as the current URL.
// If so, it converts the path to an absolute URL using IntoAbsolute. Otherwise, it returns the original path.
// Returns a pointer to the validated path string and an error if URL parsing fails.
func ValidatePath(current string, path string) (*string, error) {
	currentUrl, err := u.Parse(current)
	if err != nil {
		return nil, err
	}

	pathUrl, err := u.Parse(path)
	if err != nil {
		return nil, err
	}

	// Si le path est relatif ou du même domaine, on le transforme
	if !pathUrl.IsAbs() || pathUrl.Host == currentUrl.Host {
		return IntoAbsolute(current, path)
	}
	return &path, nil
}

// Scrap crawls and downloads content from the provided list of URLs.
// It manages the download progress, handles URL normalization, and ensures that only URLs
// belonging to the original hosts are followed. The function continues scraping until all
// reachable content is processed or an error occurs. Returns an error if any download or parsing
// operation fails.
func Scrap(urls []u.URL) error {
	var scrapper scraps = Init()
	var manager progress.Manager[uint] = progress.New[uint](50, "=>-")

	// Ajout de toutes les URLs initiales
	for _, url := range urls {
		scrapper.Add(url)
	}

	for !scrapper.IsFullDone() {
		statuses, err := scrapper.Download()
		if err != nil {
			return err
		}

		for _, status := range statuses {
			p := status.Url.Path
			if strings.HasSuffix(p, "/") {
				p += "index.html"
			}
			p = strings.TrimPrefix(p, "/")
			if status.StatusCode == 200 {
				manager.Add(status.Downloaded, status.Length, p)
			} else {
				scrapper.Delete(status.Url)
			}
		}

		newUrls, err := scrapper.parse()
		if err != nil {
			return err
		}

		if len(newUrls) == 0 {
			continue
		}

		for _, newStr := range newUrls {
			if len(newStr) == 0 {
				continue
			}

			newUrl, err := url.Parse(newStr)
			if err != nil {
				return err
			}

			// On ajoute seulement si le host correspond à l’un des hosts initiaux
			for _, original := range urls {
				if newUrl.Host == original.Host {
					scrapper.Add(*newUrl)
					break
				}
			}
		}
	}

	return nil
}
