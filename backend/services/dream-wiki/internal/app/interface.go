package app

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

type AppRepository interface {
	Search(ctx context.Context, query string) ([]models.SearchResult, error)
	GetDiagnosticInfo(ctx context.Context, pageID string) (*models.DiagnosticInfo, error)
	IndexatePage(ctx context.Context, pageID string) error
}

type AppUsecase interface {
	Search(ctx context.Context, req api.V1SearchRequest) (*api.V1SearchResponse, error)
	GetDiagnosticInfo(ctx context.Context, req api.V1DiagnosticInfoGetRequest) (*api.V1DiagnosticInfoGetResponse, error)
	IndexatePage(ctx context.Context, req api.V1IndexatePageRequest) (*api.V1IndexatePageResponse, error)
}
