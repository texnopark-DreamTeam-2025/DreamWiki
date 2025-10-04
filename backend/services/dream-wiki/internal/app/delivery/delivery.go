package delivery

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

type AppDelivery struct {
	usecase *usecase.AppUsecase
}

func NewAppDelivery(usecase *usecase.AppUsecase) *AppDelivery {
	return &AppDelivery{usecase: usecase}
}

func (d *AppDelivery) Search(ctx context.Context, request api.SearchRequestObject) (api.SearchResponseObject, error) {
	resp, err := d.usecase.Search(ctx, *request.Body)
	if err != nil {
		return api.Search200JSONResponse{}, nil
	}

	return api.Search200JSONResponse(*resp), nil
}

func (d *AppDelivery) GetDiagnosticInfo(ctx context.Context, request api.GetDiagnosticInfoRequestObject) (api.GetDiagnosticInfoResponseObject, error) {
	resp, err := d.usecase.GetDiagnosticInfo(ctx, *request.Body)
	if err != nil {
		return api.GetDiagnosticInfo200JSONResponse{}, nil
	}

	return api.GetDiagnosticInfo200JSONResponse(*resp), nil
}

func (d *AppDelivery) Indexate(ctx context.Context, request api.IndexateRequestObject) (api.IndexateResponseObject, error) {
	resp, err := d.usecase.IndexatePage(ctx, *request.Body)
	if err != nil {
		return api.Indexate200JSONResponse{}, nil
	}

	return api.Indexate200JSONResponse(*resp), nil
}
