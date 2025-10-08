package models

type SearchResult struct {
	Title       string
	Description string
	PageID      string
}

type ParagraphWithEmbedding struct {
	ParagraphID string
	PageID      string
	LineNumber  int
	Content     string
	Embedding   string
}
