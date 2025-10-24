package delivery

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (d *AppDelivery) GithubAccountPR(ctx context.Context, request api.GithubAccountPRRequestObject) (api.GithubAccountPRResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) IntegrationLogsGet(ctx context.Context, request api.IntegrationLogsGetRequestObject) (api.IntegrationLogsGetResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)

	logs, newCursor, err := usecase.GetIntegrationLogs(request.Body.IntegrationId, request.Body.Cursor)
	if err != nil {
		d.log.Error(err.Error())
		return api.IntegrationLogsGet500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}
	return api.IntegrationLogsGet200JSONResponse{LogFields: logs, Cursor: newCursor}, nil
}

func (d *AppDelivery) YwikiAddPage(ctx context.Context, request api.YwikiAddPageRequestObject) (api.YwikiAddPageResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	err := usecase.FetchPageFromYWiki(request.Body.PageUrl)
	if err != nil {
		d.log.Error(err.Error())
		return api.YwikiAddPage500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: err.Error()}}, nil
	}
	return api.YwikiAddPage200JSONResponse{}, nil
}

func (d *AppDelivery) YwikiFetchAll(ctx context.Context, request api.YwikiFetchAllRequestObject) (api.YwikiFetchAllResponseObject, error) {
	panic("unimplemented")
}
