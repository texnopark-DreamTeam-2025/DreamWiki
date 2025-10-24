package usecase

import "strings"

func extractYWikiSlugFromURL(pageURL string) string {
	const prefix = "https://wiki.yandex.ru/"

	return strings.TrimSuffix(strings.TrimPrefix(pageURL, prefix), "/")
}
