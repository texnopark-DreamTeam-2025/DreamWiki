// Package indexing provides functionality for splitting pages into paragraphs
// and determining if paragraphs are contentful.
package indexing

import (
	"strings"
)

// SplitPageToParagraphs splits a page into paragraphs based on double newlines.
func SplitPageToParagraphs(page string) []string {
	// Split the page by double newlines to get paragraphs
	paragraphs := strings.Split(page, "\n\n")

	// Filter out non-contentful paragraphs
	contentfulParagraphs := make([]string, 0)
	for _, paragraph := range paragraphs {
		if isParagraphContentful(paragraph) {
			contentfulParagraphs = append(contentfulParagraphs, strings.TrimSpace(paragraph))
		}
	}

	return contentfulParagraphs
}

// isParagraphContentful checks if a paragraph contains more than 3 words.
func isParagraphContentful(paragraph string) bool {
	// Trim whitespace and check if empty
	trimmed := strings.TrimSpace(paragraph)
	if trimmed == "" {
		return false
	}

	// Split by whitespace to count words
	words := strings.Fields(trimmed)
	return len(words) > 3
}
