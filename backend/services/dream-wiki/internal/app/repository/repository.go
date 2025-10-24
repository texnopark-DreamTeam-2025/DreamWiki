package repository

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/ydb_wrapper"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

type (
	AppRepository interface {
		Commit()
		Rollback()

		// domain_integration_logs.go
		WriteIntegrationLogField(integrationID api.IntegrationID, logText string) error
		GetIntegrationLogFields(integrationID api.IntegrationID, cursor *string, limit int64) (fields []api.IntegrationLogField, newCursor string, err error)

		// domain_page_indexation.go
		RemovePageIndexation(pageID api.PageID) error
		AddIndexedParagraph(paragraph models.ParagraphWithEmbedding) error

		// domain_pages.go
		GetPageBySlug(yWikiSlug string) (*api.Page, error)
		UpsertPage(page api.Page, ywikiSlug string) (pageID *api.PageID, err error)
		DeletePageBySlug(yWikiSlug string) error
		GetAllPageDigests() ([]api.PageDigest, error)
		RetrievePageByID(pageID api.PageID) (*api.Page, error)
		DeleteAllPages() error

		// domain_search.go
		SearchByEmbedding(query string, queryEmbedding models.Embedding) ([]api.SearchResultItem, error)

		// domain_users.go
		GetUserByLogin(username string) (*models.User, error)
	}

	appRepositoryImpl struct {
		ctx       context.Context
		ydbClient ydb_wrapper.YDBWrapper
		log       logger.Logger
	}
)

var (
	_ AppRepository = &appRepositoryImpl{}
)

func NewAppRepository(ctx context.Context, deps *deps.Deps) AppRepository {
	return &appRepositoryImpl{
		ctx:       ctx,
		log:       deps.Logger,
		ydbClient: ydb_wrapper.NewYDBWrapper(ctx, deps, true),
	}
}

func (r *appRepositoryImpl) Commit() {
	r.ydbClient.Commit()
}

func (r *appRepositoryImpl) Rollback() {
	r.ydbClient.Rollback()
}
