package repository

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
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
			name: "Non-adjacent paragraphs",
			input: []internals.ParagraphWithContext{
				{
					PageId:         api.PageID{},
					ParagraphIndex: 0,
					Content:        "First paragraph",
				},
				{
					PageId:         api.PageID{},
					ParagraphIndex: 2,
					Content:        "Second paragraph",
				},
			},
			expected: []internals.ParagraphWithContext{
				{
					PageId:         api.PageID{},
					ParagraphIndex: 0,
					Content:        "First paragraph",
				},
				{
					PageId:         api.PageID{},
					ParagraphIndex: 2,
					Content:        "Second paragraph",
				},
			},
		},
		{
			name: "Adjacent paragraphs",
			input: []internals.ParagraphWithContext{
				{
					PageId:         api.PageID{},
					ParagraphIndex: 0,
					Content:        "First paragraph",
				},
				{
					PageId:         api.PageID{},
					ParagraphIndex: 1,
					Content:        "Second paragraph",
				},
			},
			expected: []internals.ParagraphWithContext{
				{
					PageId:         api.PageID{},
					ParagraphIndex: 0,
					Content:        "First paragraph\n\nSecond paragraph",
				},
			},
		},
		{
			name: "Overlapping paragraphs",
			input: []internals.ParagraphWithContext{
				{
					PageId:         api.PageID{},
					ParagraphIndex: 1,
					Content:        "Second paragraph",
				},
				{
					PageId:         api.PageID{},
					ParagraphIndex: 2,
					Content:        "Third paragraph",
				},
				{
					PageId:         api.PageID{},
					ParagraphIndex: 0,
					Content:        "First paragraph",
				},
				{
					PageId:         api.PageID{},
					ParagraphIndex: 1,
					Content:        "Second paragraph",
				},
			},
			expected: []internals.ParagraphWithContext{
				{
					PageId:         api.PageID{},
					ParagraphIndex: 0,
					Content:        "First paragraph\n\nSecond paragraph\n\nThird paragraph",
				},
			},
		},
		{
			name: "With line numbers",
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
					EndLineNumber:   2,
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

			assert.Equal(t, len(tt.expected), len(result), "Expected %d paragraphs, got %d", len(tt.expected), len(result))

			for i, expected := range tt.expected {
				assert.Equal(t, expected.StartLineNumber, result[i].StartLineNumber, "Paragraph %d: expected StartLineNumber %d, got %d", i, expected.StartLineNumber, result[i].StartLineNumber)
				assert.Equal(t, expected.EndLineNumber, result[i].EndLineNumber, "Paragraph %d: expected EndLineNumber %d, got %d", i, expected.EndLineNumber, result[i].EndLineNumber)
				assert.Equal(t, expected.Content, result[i].Content, "Paragraph %d: expected Content %q, got %q", i, expected.Content, result[i].Content)
			}
		})
	}
}
