package indexing

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func TestSplitPageToParagraphs(t *testing.T) {
	t.Parallel()

	pageID := api.PageID(uuid.New())

	t.Run("basic functionality with headers", func(t *testing.T) {
		t.Parallel()

		content := `# Introduction

This is the first paragraph of the introduction. It contains more than three words to be considered contentful.

This is the second paragraph of the introduction. It also contains more than three words.

# Main Content

This is the first paragraph of the main content. It has enough words to be contentful.

This is the second paragraph of the main content. It also has enough words.

# Conclusion

This is the conclusion paragraph. It should be contentful as well.`

		paragraphs := SplitPageToParagraphs(pageID, content)

		// Should have 8 paragraphs (including header-only paragraphs)
		require.Equal(t, 8, len(paragraphs))

		// Check first paragraph
		require.Equal(t, 0, paragraphs[0].LineNumber)
		require.Equal(t, 0, paragraphs[0].ParagraphIndex)
		require.Equal(t, "# Introduction", paragraphs[0].Content)
		require.True(t, paragraphs[0].IsHeader)

		// Check that headers are tracked
		require.Equal(t, 1, len(paragraphs[0].Headers))
		require.Equal(t, "Introduction", paragraphs[0].Headers[0])

		// Check anchor slug
		require.NotNil(t, paragraphs[0].AnchorSlug)
		require.Equal(t, "#introduction", *paragraphs[0].AnchorSlug)

		// Check that later paragraphs have accumulated headers
		require.Equal(t, 3, len(paragraphs[7].Headers))
		require.Equal(t, "Introduction", paragraphs[7].Headers[0])
		require.Equal(t, "Main Content", paragraphs[7].Headers[1])
		require.Equal(t, "Conclusion", paragraphs[7].Headers[2])
	})

	t.Run("list items in paragraph", func(t *testing.T) {
		t.Parallel()

		content := `# Flowers

Ромашки:
* Белые
* Черные

This paragraph should be contentful.`

		paragraphs := SplitPageToParagraphs(pageID, content)

		// Should have 3 paragraphs (no longer filtering empty paragraphs)
		require.Equal(t, 3, len(paragraphs))

		// Check first paragraph (the list)
		require.Equal(t, 0, paragraphs[0].LineNumber)
		require.Equal(t, "# Flowers", paragraphs[0].Content)
		require.True(t, paragraphs[0].IsHeader)

		// Check second paragraph
		require.Equal(t, 6, paragraphs[2].LineNumber)
		require.Equal(t, "This paragraph should be contentful.", paragraphs[2].Content)
		require.False(t, paragraphs[2].IsHeader)

		// Check headers for both paragraphs
		require.Equal(t, 1, len(paragraphs[0].Headers))
		require.Equal(t, "Flowers", paragraphs[0].Headers[0])
		require.Equal(t, 1, len(paragraphs[1].Headers))
		require.Equal(t, "Flowers", paragraphs[1].Headers[0])

		// Check anchor slugs
		require.NotNil(t, paragraphs[1].AnchorSlug)
		require.Equal(t, "#flowers", *paragraphs[1].AnchorSlug)
		require.NotNil(t, paragraphs[2].AnchorSlug)
		require.Equal(t, "#flowers", *paragraphs[2].AnchorSlug)
	})

	t.Run("non contentful paragraphs", func(t *testing.T) {
		t.Parallel()

		content := `# Header

Too few words.

This is a contentful paragraph with more than three words.

Also too few.`

		paragraphs := SplitPageToParagraphs(pageID, content)

		// Should have 4 paragraphs (no longer filtering empty paragraphs)
		require.Equal(t, 4, len(paragraphs))

		// Check content
		require.Equal(t, "# Header", paragraphs[0].Content)
		require.True(t, paragraphs[0].IsHeader)

		// Check that headers are still tracked
		require.Equal(t, 1, len(paragraphs[0].Headers))
		require.Equal(t, "Header", paragraphs[0].Headers[0])
	})

	t.Run("empty content", func(t *testing.T) {
		t.Parallel()

		content := ""

		paragraphs := SplitPageToParagraphs(pageID, content)

		// Should have 0 paragraphs
		require.Equal(t, 0, len(paragraphs))
	})

	t.Run("only headers", func(t *testing.T) {
		t.Parallel()

		content := `# Header1

# Header2

# Header3`

		paragraphs := SplitPageToParagraphs(pageID, content)

		// Should have 3 paragraphs (no longer filtering empty paragraphs)
		require.Equal(t, 3, len(paragraphs))

		// Check that all paragraphs are headers
		for _, p := range paragraphs {
			require.True(t, p.IsHeader)
		}

	})

	t.Run("special characters in headers", func(t *testing.T) {
		t.Parallel()
	})

	t.Run("special characters in headers", func(t *testing.T) {
		t.Parallel()

		content := `# Header with специальные characters и symbols!

This is a contentful paragraph.`

		paragraphs := SplitPageToParagraphs(pageID, content)

		// Should have 2 paragraphs (no longer filtering empty paragraphs)
		require.Equal(t, 2, len(paragraphs))

		// Check anchor slug generation
		require.NotNil(t, paragraphs[0].AnchorSlug)
		require.Equal(t, "#header-with-characters-symbols", *paragraphs[0].AnchorSlug)
	})
	t.Run("headers with multiple levels", func(t *testing.T) {
		content := `# Main Header

Some content here.

## Sub Header

More content.

### Deep Header

Deep content.

Back to main content.`

		paragraphs := SplitPageToParagraphs(pageID, content)
		require.Equal(t, 7, len(paragraphs))

		// Check that header paragraphs have IsHeader = true
		require.True(t, paragraphs[0].IsHeader) // # Main Header
		require.True(t, paragraphs[2].IsHeader) // ## Sub Header
		require.True(t, paragraphs[4].IsHeader) // ### Deep Header

		// Check that content paragraphs have IsHeader = false
		require.False(t, paragraphs[1].IsHeader) // Some content here.
		require.False(t, paragraphs[3].IsHeader) // More content.
		require.False(t, paragraphs[5].IsHeader) // Deep content.
		require.False(t, paragraphs[6].IsHeader) // Back to main content.

		require.Equal(t, 1, len(paragraphs[0].Headers))
		require.Equal(t, "Main Header", paragraphs[0].Headers[0])

		require.Equal(t, 2, len(paragraphs[2].Headers))
		require.Equal(t, "Main Header", paragraphs[2].Headers[0])
		require.Equal(t, "Sub Header", paragraphs[2].Headers[1])

		require.Equal(t, 3, len(paragraphs[4].Headers))
		require.Equal(t, "Main Header", paragraphs[4].Headers[0])
		require.Equal(t, "Sub Header", paragraphs[4].Headers[1])
		require.Equal(t, "Deep Header", paragraphs[4].Headers[2])
	})

	t.Run("headers without content paragraphs", func(t *testing.T) {
		content := `# Header 1

# Header 2

# Header 3`

		paragraphs := SplitPageToParagraphs(pageID, content)
		require.Equal(t, 3, len(paragraphs))

		for i, p := range paragraphs {
			require.Equal(t, i*2, p.LineNumber)
			require.Equal(t, i, p.ParagraphIndex)
			// All paragraphs should be headers
			require.True(t, p.IsHeader)
		}
	})

	t.Run("mixed content with lists and headers", func(t *testing.T) {
		content := `# Getting Started

This is introduction.

## Prerequisites

- Item 1
- Item 2

## Installation

To install, run:

    code example

Final content.`

		paragraphs := SplitPageToParagraphs(pageID, content)
		require.Equal(t, 8, len(paragraphs))

		require.Equal(t, 0, paragraphs[0].LineNumber) // # Getting Started
		require.True(t, paragraphs[0].IsHeader)
		require.Equal(t, 2, paragraphs[1].LineNumber) // This is introduction.
		require.False(t, paragraphs[1].IsHeader)
		require.Equal(t, 4, paragraphs[2].LineNumber) // ## Prerequisites
		require.True(t, paragraphs[2].IsHeader)
		require.Equal(t, 6, paragraphs[3].LineNumber) // - Item 1
		require.False(t, paragraphs[3].IsHeader)
		require.Equal(t, 9, paragraphs[4].LineNumber) // ## Installation
		require.True(t, paragraphs[4].IsHeader)
		require.Equal(t, 11, paragraphs[5].LineNumber) // Final content.
		require.False(t, paragraphs[5].IsHeader)
	})

	t.Run("header with special symbols", func(t *testing.T) {
		content := `# Header with: symbols, and - punctuation!

Content here.`

		paragraphs := SplitPageToParagraphs(pageID, content)
		require.Equal(t, 2, len(paragraphs))

		// Check that the header paragraph has IsHeader = true
		require.True(t, paragraphs[0].IsHeader)
		// Check that the content paragraph has IsHeader = false
		require.False(t, paragraphs[1].IsHeader)

		require.Equal(t, "#header-with-symbols-and-punctuation", *paragraphs[0].AnchorSlug)
	})
}

func TestGenerateAnchorSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		headerText string
		expected   string
	}{
		{
			name:       "simple header",
			headerText: "Introduction",
			expected:   "#introduction",
		},
		{
			name:       "header with spaces",
			headerText: "Getting Started",
			expected:   "#getting-started",
		},
		{
			name:       "header with special characters",
			headerText: "FAQ & Help",
			expected:   "#faq-help",
		},
		{
			name:       "header with numbers",
			headerText: "Chapter 1: Basics",
			expected:   "#chapter-1-basics",
		},
		{
			name:       "header with russian text",
			headerText: "Приветствие и консультация",
			expected:   "#",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := generateAnchorSlug(tt.headerText)
			require.Equal(t, tt.expected, result)
		})
	}
}
