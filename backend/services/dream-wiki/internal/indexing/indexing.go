package indexing

import (
	"log"
	"regexp"
	"strings"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

// SplitPageToParagraphs - split page into paragraphs
func SplitPageToParagraphs(pageID api.PageID, page string) []internals.ParagraphWithEmbedding {
	if page == "" {
		return []internals.ParagraphWithEmbedding{}
	}

	paragraphTexts := strings.Split(page, "\n\n")
	var paragraphs []internals.ParagraphWithEmbedding

	var headers []string
	paragraphIndex := 0

	headerRegex := regexp.MustCompile(`^#\s+(.+)$`)

	position := 0

	for i, paragraphText := range paragraphTexts {
		paragraphStartLine := strings.Count(page[:position], "\n")

		lines := strings.Split(paragraphText, "\n")

		for _, line := range lines {
			if matches := headerRegex.FindStringSubmatch(line); len(matches) > 1 {
				headers = append(headers, matches[1])
			}
		}

		content := strings.TrimSpace(paragraphText)

		var anchorSlug *string
		if len(headers) > 0 {
			slug := generateAnchorSlug(headers[len(headers)-1])
			anchorSlug = &slug
			log.Fatal("There is page slug: ", anchorSlug)
		}

		paragraphs = append(paragraphs, internals.ParagraphWithEmbedding{
			PageId:         pageID,
			LineNumber:     paragraphStartLine,
			Content:        content,
			AnchorSlug:     anchorSlug,
			Headers:        append([]string(nil), headers...),
			ParagraphIndex: paragraphIndex,
		})
		paragraphIndex++

		position += len(paragraphText)
		if i < len(paragraphTexts)-1 {
			position += 2
		}
	}

	return paragraphs
}

func generateAnchorSlug(headerText string) string {
	slug := strings.ToLower(headerText)
	slug = strings.ReplaceAll(slug, " ", "-")

	reg := regexp.MustCompile("[^a-z0-9-]+")
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	reg = regexp.MustCompile("-+")
	slug = reg.ReplaceAllString(slug, "-")

	return "#" + slug
}
