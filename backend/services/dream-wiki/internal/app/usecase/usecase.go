package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/indexing"
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

	u.log.Info("usecase is ready")

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

	// Remove old indexing
	err := repo.RemovePageIndexation(req.PageId)
	if err != nil {
		return nil, err
	}

	// Retrieve page content
	page, err := repo.RetrievePageByID(req.PageId)
	if err != nil {
		return nil, err
	}

	// Split page into paragraphs
	paragraphs := indexing.SplitPageToParagraphs(page.Content)

	// Load paragraphs into database
	for i, paragraph := range paragraphs {
		paragraphWithEmbedding := models.ParagraphWithEmbedding{
			ParagraphID: uuid.New().String(),
			PageID:      req.PageId,
			LineNumber:  i,
			Content:     paragraph,
			Embedding:   "", // Placeholder for embedding
		}

		err := repo.AddIndexedParagraph(paragraphWithEmbedding)
		if err != nil {
			return nil, err
		}
	}

	repo.Commit()
	return &api.V1IndexatePageResponse{
		PageId: req.PageId,
	}, nil
}
