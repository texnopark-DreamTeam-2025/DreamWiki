package repository

import (
	"sort"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

func groupParagraphsByPages(paragraphs []internals.ParagraphWithContext) map[string][]internals.ParagraphWithContext {
	pageParagraphs := make(map[string][]internals.ParagraphWithContext)

	for _, paragraph := range paragraphs {
		pageKey := string(paragraph.PageId[:])
		pageParagraphs[pageKey] = append(pageParagraphs[pageKey], paragraph)
	}

	return pageParagraphs
}

func mergeOverlappingParagraphs(paragraphs []internals.ParagraphWithContext) []internals.ParagraphWithContext {
	if len(paragraphs) == 0 {
		return paragraphs
	}

	sort.Slice(paragraphs, func(i, j int) bool {
		return paragraphs[i].ParagraphIndex < paragraphs[j].ParagraphIndex
	})

	var merged []internals.ParagraphWithContext
	current := paragraphs[0]

	for i := 1; i < len(paragraphs); i++ {
		next := paragraphs[i]

		if next.ParagraphIndex == current.ParagraphIndex {
			continue
		}

		if next.ParagraphIndex == current.ParagraphIndex+1 {
			current.Content += "\n\n" + next.Content
			current.ParagraphIndex = next.ParagraphIndex
			current.EndLineNumber = next.EndLineNumber
		} else {
			merged = append(merged, current)
			current = next
		}
	}

	merged = append(merged, current)

	return merged
}
