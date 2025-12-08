package repository

import (
	"github.com/google/uuid"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func (r *appRepositoryImpl) GetUserByLogin(username string) (*models.User, error) {
	yql := `
	SELECT
		user_id,
		username,
		password_hash_bcrypt
	FROM User WHERE username=$username;
	`

	result, err := r.tx.InTX().Execute(yql, table.ValueParam("$username", types.TextValue(username)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var userID uuid.UUID
	var usernameFromDB string
	var passwordHash string
	if err = result.FetchExactlyOne(&userID, &usernameFromDB, &passwordHash); err != nil {
		return nil, err
	}

	return &models.User{
		ID:           userID,
		Login:        usernameFromDB,
		PasswordHash: passwordHash,
	}, nil
}
