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

// SplitPageToParagraphs - split page into paragraphs
func SplitPageToParagraphs(pageID api.PageID, page string) []internals.ParagraphWithEmbedding {
	if page == "" {
		return []internals.ParagraphWithEmbedding{}
	}

	paragraphTexts := strings.Split(page, "\n\n")
	var paragraphs []internals.ParagraphWithEmbedding

	var headers []string
	paragraphIndex := 0

	// Подсчитываем общее количество строк до текущего параграфа
	currentPos := 0

	for i, paragraphText := range paragraphTexts {
		paragraphStartLine := strings.Count(page[:currentPos], "\n")

		// Проверяем, есть ли заголовки в этом параграфе
		lines := strings.Split(paragraphText, "\n")
		var headerLines []int
		var nonHeaderLines []string

		// Сначала находим все заголовки
		for lineIndex, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if matches := headerRegex.FindStringSubmatch(trimmedLine); len(matches) > 1 {
				// Это заголовок
				headerLines = append(headerLines, lineIndex)
			}
		}

		// Обрабатываем заголовки
		for _, lineIndex := range headerLines {
			line := lines[lineIndex]
			trimmedLine := strings.TrimSpace(line)
			if matches := headerRegex.FindStringSubmatch(trimmedLine); len(matches) > 1 {
				// Это заголовок - добавляем его как отдельный параграф
				headerText := strings.TrimSpace(matches[2])
				headers = append(headers, headerText)

				// Находим номер строки этого заголовка
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

		// Собираем все не-заголовочные строки
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

		// Создаем параграф из оставшихся строк (не заголовков)
		if len(nonHeaderLines) > 0 {
			content := strings.TrimSpace(strings.Join(nonHeaderLines, "\n"))
			if content != "" {
				var anchorSlug *string
				// Для контентных параграфов используем slug последнего заголовка
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

		// Обновляем позицию для следующего параграфа
		currentPos += len(paragraphText)
		if i < len(paragraphTexts)-1 {
			currentPos += 2 // +2 для \n\n
		}
	}

	return paragraphs
}

// findLineIndex находит индекс строки в слайсе
func findLineIndex(lines []string, target string) int {
	for i, line := range lines {
		if strings.TrimSpace(line) == strings.TrimSpace(target) {
			return i
		}
	}
	return 0
}

// isHeaderLine проверяет, является ли строка заголовком
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

// extractHeaderText извлекает текст заголовка из строки
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

// generateAnchorSlug генерирует slug для якорной ссылки
func generateAnchorSlug(headerText string) string {
	// Приводим к нижнему регистру
	slug := strings.ToLower(headerText)
	// Заменяем пробелы на дефисы
	slug = strings.ReplaceAll(slug, " ", "-")
	// Удаляем все символы, кроме букв, цифр и дефисов
	slug = anchorSlugReg.ReplaceAllString(slug, "-")
	// Удаляем начальные и конечные дефисы
	slug = strings.Trim(slug, "-")
	// Заменяем множественные дефисы на одинарные
	slug = multipleDashReg.ReplaceAllString(slug, "-")

	// Если остались какие-то символы, добавляем #
	if slug != "" {
		return "#" + slug
	}
	return "#"
}

// isListItemLine проверяет, является ли строка элементом списка
func isListItemLine(text string) bool {
	trimmed := strings.TrimSpace(text)
	// Проверяем маркированные списки
	if strings.HasPrefix(trimmed, "* ") || strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "+ ") {
		return true
	}
	// Проверяем нумерованные списки
	if len(trimmed) > 2 {
		// Ищем шаблон "1. ", "2. ", и т.д.
		parts := strings.SplitN(trimmed, ". ", 2)
		if len(parts) == 2 {
			// Проверяем, что первая часть состоит только из цифр
			_, err := strconv.Atoi(parts[0])
			return err == nil
		}
	}
	return false
}
