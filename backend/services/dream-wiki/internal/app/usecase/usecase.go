package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/indexing"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"golang.org/x/crypto/bcrypt"
)

type AppUsecaseImpl struct {
	ctx  context.Context
	deps *deps.Deps
	log  logger.Logger
}

func NewAppUsecaseImpl(ctx context.Context, deps *deps.Deps) *AppUsecaseImpl {
	return &AppUsecaseImpl{ctx: ctx, deps: deps, log: deps.Logger}
}

func (u *AppUsecaseImpl) generateJWTToken(userID string, username string) (string, error) {
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

func (u *AppUsecaseImpl) Login(req api.V1LoginRequest) (*api.V1LoginResponse, error) {
	// Get user from repository by login
	repo := repository.StartTransaction(u.ctx, u.deps)
	defer repo.Rollback()

	user, err := repo.GetUserByLogin(req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate a JWT token
	token, err := u.generateJWTToken(user.ID.String(), req.Username)
	if err != nil {
		return nil, err
	}

	return &api.V1LoginResponse{Token: token}, nil
}

func (u *AppUsecaseImpl) Search(req api.V1SearchRequest) (*api.V1SearchResponse, error) {
	embedding, err := u.deps.InferenceClient.GenerateEmbedding(u.ctx, req.Query)
	if err != nil {
		return nil, err
	}
	repo := repository.StartTransaction(u.ctx, u.deps)
	defer repo.Rollback()

	results, err := repo.SearchByEmbedding(req.Query, embedding)
	if err != nil {
		return nil, err
	}

	apiResults := make([]api.SearchResultItem, len(results))
	for i, result := range results {
		apiResults[i] = api.SearchResultItem{
			PageId:      result.PageId,
			Title:       result.Title,
			Description: result.Description,
		}
	}

	return &api.V1SearchResponse{
		ResultItems: apiResults,
	}, nil
}

func (u *AppUsecaseImpl) GetDiagnosticInfo(req api.V1DiagnosticInfoGetRequest) (*api.V1DiagnosticInfoGetResponse, error) {
	repo := repository.StartTransaction(u.ctx, u.deps)
	defer repo.Rollback()

	page, err := repo.RetrievePageByID(req.PageId)
	if err != nil {
		return nil, err
	}

	return &api.V1DiagnosticInfoGetResponse{
		Page: *page,
	}, nil
}

func (u *AppUsecaseImpl) IndexatePage(req api.V1IndexatePageRequest) (*api.V1IndexatePageResponse, error) {
	repo := repository.StartTransaction(u.ctx, u.deps)
	defer repo.Rollback()

	err := repo.RemovePageIndexation(req.PageId)
	if err != nil {
		return nil, err
	}

	page, err := repo.RetrievePageByID(req.PageId)
	if err != nil {
		return nil, err
	}

	paragraphs := indexing.SplitPageToParagraphs(page.Content)

	for i, paragraph := range paragraphs {
		embedding, err := u.deps.InferenceClient.GenerateEmbedding(u.ctx, paragraph)
		if err != nil {
			return nil, err
		}
		paragraphWithEmbedding := models.ParagraphWithEmbedding{
			PageID:     req.PageId,
			LineNumber: i,
			Content:    paragraph,
			Embedding:  embedding,
		}

		err = repo.AddIndexedParagraph(paragraphWithEmbedding)
		if err != nil {
			return nil, err
		}
	}

	repo.Commit()
	return &api.V1IndexatePageResponse{
		PageId: req.PageId,
	}, nil
}

// func (u *AppUsecaseImpl) FetchFromExternalSource() (*api.V1FetchFromExternalSourceResponse, error) {
// 	repo := repository.StartTransaction(u.ctx, u.deps)
// 	defer repo.Rollback()

// 	// удаляем все pages и paragraphs
// 	err := repo.DeleteAllPages()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return nil, nil
// }
