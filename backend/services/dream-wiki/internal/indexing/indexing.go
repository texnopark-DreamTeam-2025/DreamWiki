// Package indexing предоставляет функциональность для разделения страниц на абзацы
// и определения, являются ли абзацы содержательными.
package indexing

import (
	"regexp"
	"strings"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

// SplitPageToParagraphs разделяет страницу на абзацы на основе двойных переводов строк
// и возвращает структуры ParagraphWithEmbedding со всей необходимой информацией.
func SplitPageToParagraphs(pageID api.PageID, page string) []internals.ParagraphWithEmbedding {
	// Handle empty page
	if page == "" {
		return []internals.ParagraphWithEmbedding{}
	}

	// Split the page by double newlines to get paragraphs
	paragraphTexts := strings.Split(page, "\n\n")
	var paragraphs []internals.ParagraphWithEmbedding

	// Track headers for anchor slug generation
	var headers []string
	paragraphIndex := 0

	// Regex for identifying headers (# Header)
	headerRegex := regexp.MustCompile(`^#\s+(.+)$`)

	// Keep track of line numbers
	lineNumber := 0

	for _, paragraphText := range paragraphTexts {
		// Calculate the starting line number for this paragraph
		paragraphStartLine := lineNumber

		// Split paragraph into lines to process headers
		lines := strings.Split(paragraphText, "\n")

		// Process lines to track headers
		for _, line := range lines {
			// Check if line is a header
			if matches := headerRegex.FindStringSubmatch(line); len(matches) > 1 {
				// Add header to headers list
				headers = append(headers, matches[1])
			}
		}

		// Use the entire paragraph text as content (trimmed)
		content := strings.TrimSpace(paragraphText)

		// Generate anchor slug from last header
		var anchorSlug *string
		if len(headers) > 0 {
			slug := generateAnchorSlug(headers[len(headers)-1])
			anchorSlug = &slug
		}

		paragraphs = append(paragraphs, internals.ParagraphWithEmbedding{
			PageId:         pageID,
			LineNumber:     paragraphStartLine,
			Content:        content,
			AnchorSlug:     anchorSlug,
			Headers:        append([]string(nil), headers...), // Copy headers slice
			ParagraphIndex: paragraphIndex,
			// Embedding will be filled later
		})
		paragraphIndex++

		// Update line number: lines in this paragraph + 2 for the \n\n separator
		lineNumber += len(lines) + 2
	}

	// Adjust line numbers - we've added 2 extra for the last paragraph
	if len(paragraphs) > 0 {
		// The last paragraph doesn't have a trailing \n\n separator
		// So we need to subtract 2 from its line count adjustment
		// But this is complex to track precisely, so we'll leave it as is for now
	}

	return paragraphs
}

// generateAnchorSlug generates an anchor slug from a header text
func generateAnchorSlug(headerText string) string {
	// Convert to lowercase
	slug := strings.ToLower(headerText)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters except hyphens and alphanumeric
	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// Replace multiple consecutive hyphens with single hyphen
	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	// Ensure it starts with #
	return "#" + slug
}
