package usecase

import (
	"context"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

type AppUsecase struct {
	repo *repository.AppRepository
}

func NewAppUsecase(repo *repository.AppRepository) *AppUsecase {
	return &AppUsecase{repo: repo}
}

func (u *AppUsecase) Search(ctx context.Context, req api.V1SearchRequest) (*api.V1SearchResponse, error) {
	results, err := u.repo.Search(ctx, req.Query)
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

func (u *AppUsecase) GetDiagnosticInfo(ctx context.Context, req api.V1DiagnosticInfoGetRequest) (*api.V1DiagnosticInfoGetResponse, error) {
	info, err := u.repo.GetDiagnosticInfo(ctx, req.PageId)
	if err != nil {
		return nil, err
	}

	return &api.V1DiagnosticInfoGetResponse{
		PageId:    info.PageID,
		Content:   info.Content,
		Title:     info.Title,
		CreatedAt: time.Now(),
	}, nil
}

func (u *AppUsecase) IndexatePage(ctx context.Context, req api.V1IndexatePageRequest) (*api.V1IndexatePageResponse, error) {
	err := u.repo.IndexatePage(ctx, req.PageId)
	if err != nil {
		return nil, err
	}

	return &api.V1IndexatePageResponse{
		PageId: req.PageId,
	}, nil
}
