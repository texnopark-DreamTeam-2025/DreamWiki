package usecase

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/indexing"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (u *appUsecaseImpl) IndexatePage(pageID api.PageID) (*api.V1IndexatePageResponse, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	return u.indexatePageInTransaction(repo, pageID)
}

func (u *appUsecaseImpl) indexatePageInTransaction(repo repository.AppRepository, pageID api.PageID) (*api.V1IndexatePageResponse, error) {
	err := repo.RemovePageIndexation(pageID)
	if err != nil {
		return nil, err
	}

	page, _, err := repo.GetPageByID(pageID)
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
			PageID:     pageID,
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
		PageId: pageID,
	}, nil
}
