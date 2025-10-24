package repository

import (
	"github.com/google/uuid"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func (r *appRepositoryImpl) GetUserByLogin(login string) (*models.User, error) {
	yql := `
	SELECT
		user_id,
		login,
		password_hash_bcrypt
	FROM User WHERE login=$login;
	`

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$login", types.TextValue(login)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var userID uuid.UUID
	var userLogin string
	var passwordHash string
	if err = result.FetchExactlyOne(&userID, &userLogin, &passwordHash); err != nil {
		return nil, err
	}

	return &models.User{
		ID:           userID,
		Login:        userLogin,
		PasswordHash: passwordHash,
	}, nil
}
