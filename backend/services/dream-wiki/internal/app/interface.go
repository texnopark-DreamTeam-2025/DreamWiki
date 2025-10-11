package app

import (
	"github.com/google/uuid"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/local_model"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

var (
	_ AppRepository = (*repository.AppRepositoryImpl)(nil)
	_ AppUsecase    = (*usecase.AppUsecaseImpl)(nil)
)

type AppRepository interface {
	SearchByEmbedding(query string, queryEmbedding local_model.Embedding) ([]models.ParagraphWithEmbedding, error)
	RetrievePageByID(pageID uuid.UUID) (*api.Page, error)
	RemovePageIndexation(pageID uuid.UUID) error
	AddIndexedParagraph(paragraph models.ParagraphWithEmbedding) error
	DeletePage() error
}

type AppUsecase interface {
	Search(req api.V1SearchRequest) (*api.V1SearchResponse, error)
	GetDiagnosticInfo(req api.V1DiagnosticInfoGetRequest) (*api.V1DiagnosticInfoGetResponse, error)
	IndexatePage(req api.V1IndexatePageRequest) (*api.V1IndexatePageResponse, error)
}
