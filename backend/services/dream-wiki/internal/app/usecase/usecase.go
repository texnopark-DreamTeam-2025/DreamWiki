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

type (
	AppUsecase interface {
		Search(req api.V1SearchRequest) (*api.V1SearchResponse, error)
		GetDiagnosticInfo(req api.V1DiagnosticInfoGetRequest) (*api.V1DiagnosticInfoGetResponse, error)
		IndexatePage(req api.V1IndexatePageRequest) (*api.V1IndexatePageResponse, error)
		Login(req api.V1LoginRequest) (*api.V1LoginResponse, error)
		FetchPageFromYWiki(pageURL string) error
		AccountGitHubPullRequest(pullRequestURL string) error
		GetPagesTree(activePagesIDs []api.PageID) ([]api.TreeItem, error)
		GetIntegrationLogs(integrationID api.IntegrationID, cursor *string) (fields []api.IntegrationLogField, newCursor string, err error)
		ApplyDraft(draftID api.DraftID) error
		CancelTask(taskID api.TaskID) error
		CreateDraft(pageURL string) (api.DraftDigest, error)
		DeleteDraft(draftID api.DraftID) error
		GetDraft(draftID api.DraftID) (api.Draft, error)
		GetTaskDetails(taskID api.TaskID) (api.Task, error)
		ListDrafts(cursor *string) ([]api.DraftDigest, error)
		ListTasks(cursor *string) (tasks []api.TaskDigest, newCursor string, err error)
		RetryTask(taskID api.TaskID) error
		UpdateDraft(draftID api.DraftID, newContent *string, newTitle *string) error
	}

	appUsecaseImpl struct {
		ctx  context.Context
		deps *deps.Deps
		log  logger.Logger
	}
)

func NewAppUsecaseImpl(ctx context.Context, deps *deps.Deps) AppUsecase {
	return &appUsecaseImpl{ctx: ctx, deps: deps, log: deps.Logger}
}

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

func (u *appUsecaseImpl) Search(req api.V1SearchRequest) (*api.V1SearchResponse, error) {
	embedding, err := u.deps.InferenceClient.GenerateEmbedding(u.ctx, req.Query)
	if err != nil {
		return nil, err
	}
	repo := repository.NewAppRepository(u.ctx, u.deps)
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

func (u *appUsecaseImpl) GetDiagnosticInfo(req api.V1DiagnosticInfoGetRequest) (*api.V1DiagnosticInfoGetResponse, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	page, err := repo.RetrievePageByID(req.PageId)
	if err != nil {
		return nil, err
	}

	return &api.V1DiagnosticInfoGetResponse{
		Page: *page,
	}, nil
}

func (u *appUsecaseImpl) IndexatePage(req api.V1IndexatePageRequest) (*api.V1IndexatePageResponse, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	return u.indexatePageInTransaction(repo, req)
}

func (u *appUsecaseImpl) indexatePageInTransaction(repo repository.AppRepository, req api.V1IndexatePageRequest) (*api.V1IndexatePageResponse, error) {
	err := repo.RemovePageIndexation(req.PageId)
	if err != nil {
		return nil, err
	}

	page, err := repo.RetrievePageByID(req.PageId)
	if err != nil {
		return nil, err
	}

	paragraphs := indexing.SplitPageToParagraphs(page.Content)

	// Generate embeddings for all paragraphs in batch
	embeddings, err := u.deps.InferenceClient.GenerateEmbeddings(u.ctx, paragraphs)
	if err != nil {
		return nil, err
	}

	for i, paragraph := range paragraphs {
		paragraphWithEmbedding := models.ParagraphWithEmbedding{
			PageID:     req.PageId,
			LineNumber: int64(i),
			Content:    paragraph,
			Embedding:  embeddings[i],
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

func (u *appUsecaseImpl) FetchPageFromYWiki(pageURL string) error {
	// 1. Parse URL, extract slug like:
	// https://wiki.yandex.ru/homepage/required-notification/ -> required-notification
	slug := extractSlugFromURL(pageURL)

	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	// Write integration log
	err := repo.WriteIntegrationLogField("ywiki", fmt.Sprintf("Fetching page with slug: %s", slug))
	if err != nil {
		u.log.Errorf("Failed to write integration log: %v", err)
	}

	// 2. Go to YWiki client from deps, fetch page
	pageResponse, err := u.deps.YWikiClient.GetPage(u.ctx, slug)
	if err != nil {
		// TODO: If fetch return 404, delete page
		// Will be done later
		// err := repo.DeletePageBySlug(slug)
		// if err != nil {
		// 	u.log.Errorf("Failed to delete page with slug %s: %v", slug, err)
		// 	return err
		// }
		// repo.Commit()
		u.log.Errorf("Failed to fetch page from YWiki: %w", err)
		return err
	}

	// 3. Upsert page to repository
	page := api.Page{
		Content:   *pageResponse.Content,
		Title:     pageResponse.Title,
		YwikiSlug: &pageResponse.Slug,
	}
	pageID, err := repo.UpsertPage(page, slug)
	if err != nil {
		u.log.Errorf("Failed to upsert page: %w", err)
		return err
	}

	// 4. Indexate page using usecase function
	indexateReq := api.V1IndexatePageRequest{
		PageId: *pageID,
	}
	_, err = u.indexatePageInTransaction(repo, indexateReq)
	if err != nil {
		u.log.Errorf("Failed to indexate page: %w", err)
		return err
	}

	// 5. Commit transaction
	repo.Commit()
	return nil
}

func (u *appUsecaseImpl) AccountGitHubPullRequest(pullRequestURL string) error {
	panic("unimplemented")
}

func (u *appUsecaseImpl) GetPagesTree(activePagesIDs []api.PageID) ([]api.TreeItem, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	items, err := repo.GetAllPageDigests()
	if err != nil {
		return nil, err
	}

	result := make([]api.TreeItem, 0, len(items))
	for _, item := range items {
		result = append(result, api.TreeItem{
			PageDigest: item,
			Children:   nil,
			Expanded:   false,
		})
	}

	return result, nil
}

func (u *appUsecaseImpl) GetIntegrationLogs(integrationID api.IntegrationID, cursor *string) (fields []api.IntegrationLogField, newCursor string, err error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	fields, newCursor, err = repo.GetIntegrationLogFields(string(integrationID), cursor, 50)
	return
}

func (u *appUsecaseImpl) ApplyDraft(draftID api.DraftID) error {
	panic("unimplemented")
}

func (u *appUsecaseImpl) CancelTask(taskID api.TaskID) error {
	panic("unimplemented")
}

func (u *appUsecaseImpl) CreateDraft(pageURL string) (api.DraftDigest, error) {
	panic("unimplemented")
}

func (u *appUsecaseImpl) DeleteDraft(draftID api.DraftID) error {
	panic("unimplemented")
}

func (u *appUsecaseImpl) GetDraft(draftID api.DraftID) (api.Draft, error) {
	panic("unimplemented")
}

func (u *appUsecaseImpl) GetTaskDetails(taskID api.TaskID) (api.Task, error) {
	panic("unimplemented")
}

func (u *appUsecaseImpl) ListDrafts(cursor *string) ([]api.DraftDigest, error) {
	panic("unimplemented")
}

func (u *appUsecaseImpl) ListTasks(cursor *string) (tasks []api.TaskDigest, newCursor string, err error) {
	panic("unimplemented")
}

func (u *appUsecaseImpl) RetryTask(taskID api.TaskID) error {
	panic("unimplemented")
}

func (u *appUsecaseImpl) UpdateDraft(draftID api.DraftID, newContent *string, newTitle *string) error {
	panic("unimplemented")
}
