package repository

import (
	"sort"
	"testing"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

func TestMergeOverlappingParagraphs(t *testing.T) {
	tests := []struct {
		name     string
		input    []internals.ParagraphWithContext
		expected []internals.ParagraphWithContext
	}{
		{
			name:     "Empty input",
			input:    []internals.ParagraphWithContext{},
			expected: []internals.ParagraphWithContext{},
		},
		{
			name: "Single paragraph",
			input: []internals.ParagraphWithContext{
				{
					PageId:          api.PageID{},
					ParagraphIndex:  0,
					StartLineNumber: 1,
					EndLineNumber:   1,
					Content:         "Single paragraph",
				},
			},
			expected: []internals.ParagraphWithContext{
				{
					PageId:          api.PageID{},
					ParagraphIndex:  0,
					StartLineNumber: 1,
					EndLineNumber:   1,
					Content:         "Single paragraph",
				},
			},
		},
		{
			name: "Non-overlapping paragraphs",
			input: []internals.ParagraphWithContext{
				{
					PageId:          api.PageID{},
					ParagraphIndex:  0,
					StartLineNumber: 1,
					EndLineNumber:   1,
					Content:         "First paragraph",
				},
				{
					PageId:          api.PageID{},
					ParagraphIndex:  1,
					StartLineNumber: 3,
					EndLineNumber:   3,
					Content:         "Second paragraph",
				},
			},
			expected: []internals.ParagraphWithContext{
				{
					PageId:          api.PageID{},
					ParagraphIndex:  0,
					StartLineNumber: 1,
					EndLineNumber:   1,
					Content:         "First paragraph",
				},
				{
					PageId:          api.PageID{},
					ParagraphIndex:  1,
					StartLineNumber: 3,
					EndLineNumber:   3,
					Content:         "Second paragraph",
				},
			},
		},
		{
			name: "Adjacent paragraphs",
			input: []internals.ParagraphWithContext{
				{
					PageId:          api.PageID{},
					ParagraphIndex:  0,
					StartLineNumber: 1,
					EndLineNumber:   1,
					Content:         "First paragraph",
				},
				{
					PageId:          api.PageID{},
					ParagraphIndex:  1,
					StartLineNumber: 2,
					EndLineNumber:   2,
					Content:         "Second paragraph",
				},
			},
			expected: []internals.ParagraphWithContext{
				{
					PageId:          api.PageID{},
					ParagraphIndex:  0,
					StartLineNumber: 1,
					EndLineNumber:   2,
					Content:         "First paragraph\n\nSecond paragraph",
				},
			},
		},
		{
			name: "Overlapping paragraphs",
			input: []internals.ParagraphWithContext{
				{
					PageId:          api.PageID{},
					ParagraphIndex:  0,
					StartLineNumber: 1,
					EndLineNumber:   3,
					Content:         "First paragraph",
				},
				{
					PageId:          api.PageID{},
					ParagraphIndex:  1,
					StartLineNumber: 2,
					EndLineNumber:   4,
					Content:         "Second paragraph",
				},
				{
					PageId:          api.PageID{},
					ParagraphIndex:  1,
					StartLineNumber: 2,
					EndLineNumber:   4,
					Content:         "Second paragraph",
				},
				{
					PageId:          api.PageID{},
					ParagraphIndex:  2,
					StartLineNumber: 5,
					EndLineNumber:   5,
					Content:         "Third paragraph",
				},
			},
			expected: []internals.ParagraphWithContext{
				{
					PageId:          api.PageID{},
					ParagraphIndex:  0,
					StartLineNumber: 1,
					EndLineNumber:   4,
					Content:         "First paragraph\n\nSecond paragraph\n\nThird paragraph",
				},
			},
		},
		{
			name: "Multiple overlapping paragraphs",
			input: []internals.ParagraphWithContext{
				{
					PageId:          api.PageID{},
					ParagraphIndex:  0,
					StartLineNumber: 1,
					EndLineNumber:   2,
					Content:         "First paragraph",
				},
				{
					PageId:          api.PageID{},
					ParagraphIndex:  2,
					StartLineNumber: 4,
					EndLineNumber:   5,
					Content:         "Third paragraph",
				},
				{
					PageId:          api.PageID{},
					ParagraphIndex:  3,
					StartLineNumber: 5,
					EndLineNumber:   6,
					Content:         "Fourth paragraph",
				},
			},
			expected: []internals.ParagraphWithContext{
				{
					PageId:          api.PageID{},
					ParagraphIndex:  0,
					StartLineNumber: 1,
					EndLineNumber:   3,
					Content:         "First paragraph",
				},
				{
					PageId:          api.PageID{},
					ParagraphIndex:  2,
					StartLineNumber: 4,
					EndLineNumber:   6,
					Content:         "Third paragraph\n\nFourth paragraph",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputCopy := make([]internals.ParagraphWithContext, len(tt.input))
			copy(inputCopy, tt.input)

			sort.Slice(inputCopy, func(i, j int) bool {
				return inputCopy[i].StartLineNumber < inputCopy[j].StartLineNumber
			})

			result := mergeOverlappingParagraphs(inputCopy)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d paragraphs, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i].StartLineNumber != expected.StartLineNumber {
					t.Errorf("Paragraph %d: expected StartLineNumber %d, got %d", i, expected.StartLineNumber, result[i].StartLineNumber)
				}
				if result[i].EndLineNumber != expected.EndLineNumber {
					t.Errorf("Paragraph %d: expected EndLineNumber %d, got %d", i, expected.EndLineNumber, result[i].EndLineNumber)
				}
				if result[i].Content != expected.Content {
					t.Errorf("Paragraph %d: expected Content %q, got %q", i, expected.Content, result[i].Content)
				}
			}
		})
	}
}
