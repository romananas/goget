package scraping

import (
	"main/progress"
	"net/url"
	u "net/url"
	"strings"
)

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

func Scrap(url u.URL) error {

	var scrapper scraps = Init()
	var manager progress.Manager[uint] = progress.New[uint](50, "=>-")
	scrapper.Add(url)

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

		urls, err := scrapper.parse()
		if err != nil {
			return err
		}

		if len(urls) == 0 {
			continue
		}

		for _, new := range urls {
			if len(new) == 0 {
				continue
			}

			new, err := url.Parse(new)
			if err != nil {
				return err
			}
			if new.Host == url.Host {
				scrapper.Add(*new)
			}
		}

	}

	return nil
}

func ScrapMulti(urls []u.URL) error {
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
