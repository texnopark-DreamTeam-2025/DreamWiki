package models

import "github.com/google/uuid"

type ParagraphWithEmbedding struct {
	PageID     uuid.UUID
	LineNumber int
	Content    string
	Embedding  string
}
