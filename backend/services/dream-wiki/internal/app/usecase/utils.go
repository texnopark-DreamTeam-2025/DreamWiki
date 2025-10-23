package usecase

import "strings"

func extractSlugFromURL(pageURL string) string {
	const prefix = "https://wiki.yandex.ru/"

	return strings.TrimSuffix(strings.TrimPrefix(pageURL, prefix), "/")
}
