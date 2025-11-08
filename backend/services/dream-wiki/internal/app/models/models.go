package models

import (
	"fmt"

	"github.com/google/uuid"
)

type (
	Embedding []float32

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
	ErrNoRows           error = fmt.Errorf("%w: no rows", ErrNotFound)
)
