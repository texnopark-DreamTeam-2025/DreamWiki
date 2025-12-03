package indexing

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

var (
	headerRegex     = regexp.MustCompile(`^(#{1,})\s+(.+)$`)
	anchorSlugReg   = regexp.MustCompile("[^a-z0-9-]+")
	multipleDashReg = regexp.MustCompile("-+")
)

func SplitPageToParagraphs(pageID api.PageID, page string) []internals.ParagraphWithEmbedding {
	if page == "" {
		return []internals.ParagraphWithEmbedding{}
	}

	paragraphTexts := strings.Split(page, "\n\n")
	var paragraphs []internals.ParagraphWithEmbedding

	var headers []string
	paragraphIndex := 0

	currentPos := 0

	for i, paragraphText := range paragraphTexts {
		paragraphStartLine := strings.Count(page[:currentPos], "\n")

		lines := strings.Split(paragraphText, "\n")
		var headerLines []int
		var nonHeaderLines []string

		for lineIndex, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if matches := headerRegex.FindStringSubmatch(trimmedLine); len(matches) > 1 {
				headerLines = append(headerLines, lineIndex)
			}
		}

		for _, lineIndex := range headerLines {
			line := lines[lineIndex]
			trimmedLine := strings.TrimSpace(line)
			if matches := headerRegex.FindStringSubmatch(trimmedLine); len(matches) > 1 {
				headerText := strings.TrimSpace(matches[2])
				headers = append(headers, headerText)

				lineNumber := paragraphStartLine + lineIndex

				var anchorSlug *string
				slug := generateAnchorSlug(headerText)
				anchorSlug = &slug

				paragraphs = append(paragraphs, internals.ParagraphWithEmbedding{
					PageId:         pageID,
					LineNumber:     lineNumber,
					Content:        trimmedLine,
					AnchorSlug:     anchorSlug,
					Headers:        append([]string(nil), headers...),
					ParagraphIndex: paragraphIndex,
					IsHeader:       true,
				})
				paragraphIndex++
			}
		}

		for lineIndex, line := range lines {
			isHeader := false
			for _, headerLineIndex := range headerLines {
				if lineIndex == headerLineIndex {
					isHeader = true
					break
				}
			}
			if !isHeader {
				nonHeaderLines = append(nonHeaderLines, line)
			}
		}

		if len(nonHeaderLines) > 0 {
			content := strings.TrimSpace(strings.Join(nonHeaderLines, "\n"))
			if content != "" {
				var anchorSlug *string
				if len(headers) > 0 {
					lastHeader := headers[len(headers)-1]
					slug := generateAnchorSlug(lastHeader)
					anchorSlug = &slug
				}

				paragraphs = append(paragraphs, internals.ParagraphWithEmbedding{
					PageId:         pageID,
					LineNumber:     paragraphStartLine,
					Content:        content,
					AnchorSlug:     anchorSlug,
					Headers:        append([]string(nil), headers...),
					ParagraphIndex: paragraphIndex,
					IsHeader:       false,
				})
				paragraphIndex++
			}
		}

		currentPos += len(paragraphText)
		if i < len(paragraphTexts)-1 {
			currentPos += 2
		}
	}

	return paragraphs
}

func findLineIndex(lines []string, target string) int {
	for i, line := range lines {
		if strings.TrimSpace(line) == strings.TrimSpace(target) {
			return i
		}
	}
	return 0
}

func isHeaderLine(text string) bool {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if headerRegex.MatchString(line) {
			return true
		}
	}
	return false
}

func extractHeaderText(text string) string {
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if matches := headerRegex.FindStringSubmatch(line); len(matches) > 1 {
			return strings.TrimSpace(matches[2])
		}
	}
	return ""
}

func generateAnchorSlug(headerText string) string {
	slug := strings.ToLower(headerText)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = anchorSlugReg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	slug = multipleDashReg.ReplaceAllString(slug, "-")

	if slug != "" {
		return "#" + slug
	}
	return "#"
}

func isListItemLine(text string) bool {
	trimmed := strings.TrimSpace(text)
	if strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "+ ") {
		return true
	}
	if len(trimmed) > 2 {
		parts := strings.SplitN(trimmed, ". ", 2)
		if len(parts) == 2 {
			_, err := strconv.Atoi(parts[0])
			return err == nil
		}
	}
	return false
}
