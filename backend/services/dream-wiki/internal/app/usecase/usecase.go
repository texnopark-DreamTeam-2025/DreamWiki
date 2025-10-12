package usecase

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/indexing"
	inference_client "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/inference"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

type AppUsecaseImpl struct {
	ctx  context.Context
	deps *deps.Deps
	log  logger.Logger
}

func NewAppUsecaseImpl(ctx context.Context, deps *deps.Deps) *AppUsecaseImpl {
	return &AppUsecaseImpl{ctx: ctx, deps: deps, log: deps.Logger}
}

func (u *AppUsecaseImpl) Search(req api.V1SearchRequest) (*api.V1SearchResponse, error) {
	embedding, err := u.deps.InferenceClient.GenerateEmbeddingWithResponse(u.ctx,
		inference_client.GenerateEmbeddingJSONRequestBody{Text: req.Query})
	if err != nil {
		return nil, err
	}
	repo := repository.StartTransaction(u.ctx, u.deps)
	defer repo.Rollback()

	results, err := repo.SearchByEmbedding(req.Query, embedding.JSON200.Embedding)
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
		embedding, err := u.deps.InferenceClient.GenerateEmbeddingWithResponse(u.ctx,
			inference_client.GenerateEmbeddingJSONRequestBody{Text: paragraph},
		)
		if err != nil {
			return nil, err
		}
		paragraphWithEmbedding := models.ParagraphWithEmbedding{
			PageID:     req.PageId,
			LineNumber: i,
			Content:    paragraph,
			Embedding:  embedding.JSON200.Embedding,
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

func (u *AppUsecaseImpl) FetchFromExternalSource() (*api.V1FetchFromExternalSourceResponse, error) {
	repo := repository.StartTransaction(u.ctx, u.deps)
	defer repo.Rollback()

	// удаляем все pages и paragraphs
	err := repo.DeleteAllPages()
	if err != nil {
		return nil, err
	}

	return nil, nil
}
