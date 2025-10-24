package usecase

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"golang.org/x/crypto/bcrypt"
)

func (u *appUsecaseImpl) generateJWTToken(userID string, username string) (string, error) {
	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 240).Unix(), // Token expires in 240 hours
	})

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(u.deps.Config.JWTSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *appUsecaseImpl) Login(req api.V1LoginRequest) (*api.V1LoginResponse, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	user, err := repo.GetUserByLogin(req.Username)
	if err != nil {
		return nil, models.ErrWrongCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, models.ErrWrongCredentials
	}

	token, err := u.generateJWTToken(user.ID.String(), req.Username)
	if err != nil {
		return nil, err
	}

	return &api.V1LoginResponse{Token: token}, nil
}
