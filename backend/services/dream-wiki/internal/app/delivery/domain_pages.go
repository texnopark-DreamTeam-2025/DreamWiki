package delivery

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (d *AppDelivery) GetDiagnosticInfo(ctx context.Context, request api.GetDiagnosticInfoRequestObject) (api.GetDiagnosticInfoResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.GetDiagnosticInfo(*request.Body)
	if err != nil {
		d.log.Error(err.Error())
		return api.GetDiagnosticInfo500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.GetDiagnosticInfo200JSONResponse(*resp), nil
}

func (d *AppDelivery) PagesTreeGet(ctx context.Context, request api.PagesTreeGetRequestObject) (api.PagesTreeGetResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	result, err := usecase.GetPagesTree(request.Body.ActivePageIds)
	if err != nil {
		d.log.Error(err.Error())
		return api.PagesTreeGet500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}
	return api.PagesTreeGet200JSONResponse{Tree: result}, nil
}
