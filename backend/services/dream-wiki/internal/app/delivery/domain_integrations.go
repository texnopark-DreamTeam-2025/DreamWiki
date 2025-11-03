package delivery

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (d *AppDelivery) GithubAccountPR(ctx context.Context, request api.GithubAccountPRRequestObject) (api.GithubAccountPRResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)

	_, err := usecase.GithubAccountPRAsync(request.Body.PrUrl)
	if err != nil {
		d.log.Error(err.Error())
		return api.GithubAccountPR500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}
	return api.GithubAccountPR200JSONResponse{}, nil
}

func (d *AppDelivery) IntegrationLogsGet(ctx context.Context, request api.IntegrationLogsGetRequestObject) (api.IntegrationLogsGetResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)

	logs, nextInfo, err := usecase.GetIntegrationLogs(request.Body.IntegrationId, request.Body.Cursor)
	if err != nil {
		d.log.Error(err.Error())
		return api.IntegrationLogsGet500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}
	return api.IntegrationLogsGet200JSONResponse{LogFields: logs, NextInfo: *nextInfo}, nil
}

func (d *AppDelivery) YwikiAddPage(ctx context.Context, request api.YwikiAddPageRequestObject) (api.YwikiAddPageResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	err := usecase.FetchPageFromYWiki(request.Body.PageUrl)
	if err != nil {
		d.log.Error(err.Error())
		return api.YwikiAddPage500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}
	return api.YwikiAddPage200JSONResponse{}, nil
}

func (d *AppDelivery) YwikiFetchAll(ctx context.Context, request api.YwikiFetchAllRequestObject) (api.YwikiFetchAllResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)

	_, err := usecase.YwikiFetchAllAsync()
	if err != nil {
		d.log.Error(err.Error())
		return api.YwikiFetchAll500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}
	return api.YwikiFetchAll200Response{}, nil
}
