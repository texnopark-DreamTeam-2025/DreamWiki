package delivery

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

type AppDelivery struct {
	deps *deps.Deps
	log  logger.Logger
}

func NewAppDelivery(deps *deps.Deps) *AppDelivery {
	return &AppDelivery{deps: deps, log: deps.Logger}
}

func (d *AppDelivery) Search(ctx context.Context, request api.SearchRequestObject) (api.SearchResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.Search(*request.Body)
	if err != nil {
		d.log.Error(err.Error())
		return api.Search500JSONResponse{InternalErrorResponseJSONResponse: api.InternalErrorResponseJSONResponse{Message: "internal error"}}, nil
	}

	return api.Search200JSONResponse(*resp), nil
}

func (d *AppDelivery) GetDiagnosticInfo(ctx context.Context, request api.GetDiagnosticInfoRequestObject) (api.GetDiagnosticInfoResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.GetDiagnosticInfo(*request.Body)
	if err != nil {
		d.log.Error(err.Error())
		return api.GetDiagnosticInfo500JSONResponse{InternalErrorResponseJSONResponse: api.InternalErrorResponseJSONResponse{Message: "internal error"}}, nil
	}

	return api.GetDiagnosticInfo200JSONResponse(*resp), nil
}

func (d *AppDelivery) IndexatePage(ctx context.Context, request api.IndexatePageRequestObject) (api.IndexatePageResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.IndexatePage(*request.Body)
	if err != nil {
		d.log.Error(err.Error())
		return api.IndexatePage500JSONResponse{InternalErrorResponseJSONResponse: api.InternalErrorResponseJSONResponse{Message: "internal error"}}, nil
	}

	return api.IndexatePage200JSONResponse(*resp), nil
}

func (d *AppDelivery) FetchFromExternalSource(ctx context.Context, request api.FetchFromExternalSourceRequestObject) (api.FetchFromExternalSourceResponseObject, error) {
	// usecase:= usecase.NewAppUsecaseImpl(ctx, d.deps)
	var resp map[string]interface{}
	return api.FetchFromExternalSource200JSONResponse(resp), nil
}
