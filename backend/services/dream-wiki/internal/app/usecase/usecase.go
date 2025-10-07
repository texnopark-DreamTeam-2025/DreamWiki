package usecase

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

type AppUsecaseImpl struct {
	ctx  context.Context
	deps *deps.Deps
}

func NewAppUsecaseImpl(ctx context.Context, deps *deps.Deps) *AppUsecaseImpl {
	return &AppUsecaseImpl{ctx: ctx, deps: deps}
}

func (u *AppUsecaseImpl) Search(req api.V1SearchRequest) (*api.V1SearchResponse, error) {
	repo := repository.StartTransaction(u.ctx, u.deps)
	defer repo.Rollback()

	results, err := repo.Search(req.Query)
	if err != nil {
		return nil, err
	}

	apiResults := make([]api.SearchResultItem, len(results))
	for i, result := range results {
		apiResults[i] = api.SearchResultItem{
			Title:       result.Title,
			Description: result.Description,
			PageId:      result.PageID,
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

	// err := repo.IndexatePage(u.ctx, req.PageId)
	// if err != nil {
	// 	repo.Rollback()
	// 	return nil, err
	// }

	return &api.V1IndexatePageResponse{
		PageId: req.PageId,
	}, nil
}
