package usecase

import (
	"fmt"
	"strings"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

func (u *appUsecaseImpl) Search(req api.V1SearchRequest) (*api.V1SearchResponse, error) {
	embedding, err := u.deps.InferenceClient.GenerateEmbedding(u.ctx, req.Query)
	if err != nil {
		return nil, err
	}
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	embeddingResults, err := repo.SearchByEmbedding(req.Query, internals.Embedding(embedding), 5)
	if err != nil {
		return nil, err
	}

	terms := strings.Fields(req.Query)
	termResults, err := repo.SearchByTerms(terms, 5)
	if err != nil {
		return nil, err
	}

	u.log.Info("Search: ", len(termResults), " term results, ", len(embeddingResults), " embedding results")

	results := append(embeddingResults, termResults...)

	seen := make(map[string]bool)
	uniqueResults := make([]internals.SearchResultItem, 0)
	for _, result := range results {
		key := fmt.Sprintf("%s-%d", result.PageId, result.ParagraphIndex)
		if !seen[key] {
			seen[key] = true
			uniqueResults = append(uniqueResults, result)
		}
	}

	apiResults := make([]api.SearchResultItem, len(uniqueResults))
	for i, result := range uniqueResults {
		apiResults[i] = api.SearchResultItem{
			PageId:      result.PageId,
			Title:       result.PageTitle,
			Description: result.ParagraphContent,
		}
	}

	return &api.V1SearchResponse{
		ResultItems: apiResults,
	}, nil
}
