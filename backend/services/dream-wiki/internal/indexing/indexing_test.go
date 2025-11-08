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

		// Check second paragraph
		require.Equal(t, 6, paragraphs[2].LineNumber)
		require.Equal(t, "This paragraph should be contentful.", paragraphs[2].Content)

		// Check headers for both paragraphs
		require.Equal(t, 1, len(paragraphs[0].Headers))
		require.Equal(t, "Flowers", paragraphs[0].Headers[0])
		require.Equal(t, 1, len(paragraphs[1].Headers))
		require.Equal(t, "Flowers", paragraphs[1].Headers[0])

		// Check anchor slugs
		require.NotNil(t, paragraphs[0].AnchorSlug)
		require.Equal(t, "#flowers", *paragraphs[0].AnchorSlug)
		require.NotNil(t, paragraphs[1].AnchorSlug)
		require.Equal(t, "#flowers", *paragraphs[1].AnchorSlug)
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
