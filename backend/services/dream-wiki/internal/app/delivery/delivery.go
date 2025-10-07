package delivery

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

type AppDelivery struct {
	deps *deps.Deps
}

func NewAppDelivery(deps *deps.Deps) *AppDelivery {
	return &AppDelivery{deps: deps}
}

func (d *AppDelivery) Search(ctx context.Context, request api.SearchRequestObject) (api.SearchResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.Search(*request.Body)
	if err != nil {
		return api.Search200JSONResponse{}, nil
	}

	return api.Search200JSONResponse(*resp), nil
}

func (d *AppDelivery) GetDiagnosticInfo(ctx context.Context, request api.GetDiagnosticInfoRequestObject) (api.GetDiagnosticInfoResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.GetDiagnosticInfo(*request.Body)
	if err != nil {
		return api.GetDiagnosticInfo200JSONResponse{}, nil
	}

	return api.GetDiagnosticInfo200JSONResponse(*resp), nil
}

func (d *AppDelivery) Indexate(ctx context.Context, request api.IndexateRequestObject) (api.IndexateResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.IndexatePage(*request.Body)
	if err != nil {
		return api.Indexate200JSONResponse{}, nil
	}

	return api.Indexate200JSONResponse(*resp), nil
}
