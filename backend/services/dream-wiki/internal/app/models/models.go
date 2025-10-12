package models

import (
	"github.com/google/uuid"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/local_model"
)

type ParagraphWithEmbedding struct {
	PageID     uuid.UUID
	LineNumber int
	Content    string
	Embedding  local_model.Embedding
}

type User struct {
	ID           uuid.UUID
	Login        string
	PasswordHash string
}
