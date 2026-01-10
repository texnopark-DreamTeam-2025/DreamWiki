package indexing

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"golang.org/x/text/unicode/norm"
)

var (
	headerRegex = regexp.MustCompile(`^(#{1,})\s+(.+)$`)
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
				slug := headerToSlug(headerText)
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
			isHeader := slices.Contains(headerLines, lineIndex)
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
					slug := headerToSlug(lastHeader)
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

func headerToSlug(title string) string {
	re := regexp.MustCompile(`\(([^)]+)\)`)
	title = re.ReplaceAllStringFunc(title, func(match string) string {
		inner := strings.TrimSpace(match[1 : len(match)-1])
		allowed := regexp.MustCompile(`[^\p{L}\p{N}+\-]+`)
		cleaned := allowed.ReplaceAllString(inner, " ")
		return " " + cleaned + " "
	})

	title = norm.NFD.String(title)
	title = strings.ToLower(title)

	translitMap := map[rune]string{
		'а': "a",
		'б': "b",
		'в': "v",
		'г': "g",
		'д': "d",
		'е': "e",
		'ё': "yo",
		'ж': "zh",
		'з': "z",
		'и': "i",
		'й': "j",
		'к': "k",
		'л': "l",
		'м': "m",
		'н': "n",
		'о': "o",
		'п': "p",
		'р': "r",
		'с': "s",
		'т': "t",
		'у': "u",
		'ф': "f",
		'х': "h",
		'ц': "c",
		'ч': "ch",
		'ш': "sh",
		'щ': "sh",
		'ъ': "",
		'ы': "y",
		'ь': "",
		'э': "e",
		'ю': "yu",
		'я': "ya",
	}

	var sb strings.Builder
	for _, r := range title {
		if repl, ok := translitMap[r]; ok {
			sb.WriteString(repl)
		} else if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '+' || r == '-' {
			sb.WriteRune(r)
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) {
			sb.WriteRune(r)
		} else {
			sb.WriteRune(' ')
		}
	}
	title = sb.String()

	space := regexp.MustCompile(`\s+`)
	title = space.ReplaceAllString(title, " ")

	title = strings.ReplaceAll(title, " ", "-")

	title = strings.Trim(title, "-")

	dash := regexp.MustCompile(`-+`)
	title = dash.ReplaceAllString(title, "-")

	fmt.Println(title)
	return title
}
