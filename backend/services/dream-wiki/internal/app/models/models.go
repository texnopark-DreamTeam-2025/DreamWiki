package models

import (
	"github.com/google/uuid"
)

type Embedding []float32

type ParagraphWithEmbedding struct {
	PageID     uuid.UUID
	LineNumber int
	Content    string
	Embedding  Embedding
}

type User struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
}
