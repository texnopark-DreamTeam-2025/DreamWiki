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

func (d *AppDelivery) Login(ctx context.Context, request api.LoginRequestObject) (api.LoginResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.Login(*request.Body)
	if err != nil {
		d.log.Error(err.Error())
		return api.Login500JSONResponse{Message: "Internal server error"}, nil
	}

	return api.Login200JSONResponse(*resp), nil
}

func (d *AppDelivery) Logout(ctx context.Context, request api.LogoutRequestObject) (api.LogoutResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	err := usecase.Logout()
	if err != nil {
		d.log.Error(err.Error())
		return api.Logout500JSONResponse{Message: "Internal server error"}, nil
	}

	return api.Logout200Response{}, nil
}

// YwikiFetchAll implements api.StrictServerInterface.
func (d *AppDelivery) YwikiFetchAll(ctx context.Context, request api.YwikiFetchAllRequestObject) (api.YwikiFetchAllResponseObject, error) {
	panic("unimplemented")
}

func NewAppDelivery(deps *deps.Deps) *AppDelivery {
	return &AppDelivery{deps: deps, log: deps.Logger}
}

func (d *AppDelivery) Search(ctx context.Context, request api.SearchRequestObject) (api.SearchResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.Search(*request.Body)
	if err != nil {
		d.log.Error(err.Error())
		return api.Search500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: "internal error"}}, nil
	}

	return api.Search200JSONResponse(*resp), nil
}

func (d *AppDelivery) GetDiagnosticInfo(ctx context.Context, request api.GetDiagnosticInfoRequestObject) (api.GetDiagnosticInfoResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.GetDiagnosticInfo(*request.Body)
	if err != nil {
		d.log.Error(err.Error())
		return api.GetDiagnosticInfo500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: "internal error"}}, nil
	}

	return api.GetDiagnosticInfo200JSONResponse(*resp), nil
}

func (d *AppDelivery) IndexatePage(ctx context.Context, request api.IndexatePageRequestObject) (api.IndexatePageResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.IndexatePage(*request.Body)
	if err != nil {
		d.log.Error(err.Error())
		return api.IndexatePage500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: "internal error"}}, nil
	}

	return api.IndexatePage200JSONResponse(*resp), nil
}

func (d *AppDelivery) GithubAccountPR(ctx context.Context, request api.GithubAccountPRRequestObject) (api.GithubAccountPRResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) IntegrationLogsGet(ctx context.Context, request api.IntegrationLogsGetRequestObject) (api.IntegrationLogsGetResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) YwikiAddPage(ctx context.Context, request api.YwikiAddPageRequestObject) (api.YwikiAddPageResponseObject, error) {
	return api.YwikiAddPage500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: "not implemented"}}, nil
}
