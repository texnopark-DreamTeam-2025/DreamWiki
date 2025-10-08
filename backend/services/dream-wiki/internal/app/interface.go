package app

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

var (
	_ AppRepository = (*repository.AppRepositoryImpl)(nil)
	_ AppUsecase    = (*usecase.AppUsecaseImpl)(nil)
)

type AppRepository interface {
	Search(query string) ([]models.SearchResult, error)
	RetrievePageByID(pageID string) (*api.Page, error)
	RemovePageIndexation(pageID string) error
	AddIndexedParagraph(paragraph models.ParagraphWithEmbedding) error
}

type AppUsecase interface {
	Search(req api.V1SearchRequest) (*api.V1SearchResponse, error)
	GetDiagnosticInfo(req api.V1DiagnosticInfoGetRequest) (*api.V1DiagnosticInfoGetResponse, error)
	IndexatePage(req api.V1IndexatePageRequest) (*api.V1IndexatePageResponse, error)
}
