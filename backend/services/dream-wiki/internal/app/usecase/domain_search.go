package usecase

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

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
