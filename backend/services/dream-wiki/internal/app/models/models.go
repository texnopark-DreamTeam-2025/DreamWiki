package models

import (
	"fmt"

	"github.com/google/uuid"
)

type (
	Embedding []float32

	ParagraphWithEmbedding struct {
		PageID     uuid.UUID
		LineNumber int
		Content    string
		Embedding  Embedding
	}

	User struct {
		ID           uuid.UUID
		Login        string
		PasswordHash string
	}
)

var (
	ErrWrongCredentials error = fmt.Errorf("wrong credentials")
	ErrNoAccess         error = fmt.Errorf("no access")
	ErrNotFound         error = fmt.Errorf("not found")
)
