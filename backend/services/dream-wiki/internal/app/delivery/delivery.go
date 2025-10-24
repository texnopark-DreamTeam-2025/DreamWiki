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

var (
	_ api.StrictServerInterface = &AppDelivery{}
)

const (
	internalErrorMessage string = "internal error"
)

func NewAppDelivery(deps *deps.Deps) *AppDelivery {
	return &AppDelivery{deps: deps, log: deps.Logger}
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

func (d *AppDelivery) YwikiFetchAll(ctx context.Context, request api.YwikiFetchAllRequestObject) (api.YwikiFetchAllResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) Search(ctx context.Context, request api.SearchRequestObject) (api.SearchResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.Search(*request.Body)
	if err != nil {
		d.log.Error(err.Error())
		return api.Search500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.Search200JSONResponse(*resp), nil
}

func (d *AppDelivery) GetDiagnosticInfo(ctx context.Context, request api.GetDiagnosticInfoRequestObject) (api.GetDiagnosticInfoResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.GetDiagnosticInfo(*request.Body)
	if err != nil {
		d.log.Error(err.Error())
		return api.GetDiagnosticInfo500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.GetDiagnosticInfo200JSONResponse(*resp), nil
}

func (d *AppDelivery) IndexatePage(ctx context.Context, request api.IndexatePageRequestObject) (api.IndexatePageResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.IndexatePage(request.Body.PageId)
	if err != nil {
		d.log.Error(err.Error())
		return api.IndexatePage500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.IndexatePage200JSONResponse(*resp), nil
}

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

func (d *AppDelivery) PagesTreeGet(ctx context.Context, request api.PagesTreeGetRequestObject) (api.PagesTreeGetResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	result, err := usecase.GetPagesTree(request.Body.ActivePageIds)
	if err != nil {
		d.log.Error(err.Error())
		return api.PagesTreeGet500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}
	return api.PagesTreeGet200JSONResponse{Tree: result}, nil
}

func (d *AppDelivery) ApplyDraft(ctx context.Context, request api.ApplyDraftRequestObject) (api.ApplyDraftResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) CancelTask(ctx context.Context, request api.CancelTaskRequestObject) (api.CancelTaskResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) CreateDraft(ctx context.Context, request api.CreateDraftRequestObject) (api.CreateDraftResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) DeleteDraft(ctx context.Context, request api.DeleteDraftRequestObject) (api.DeleteDraftResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) GetDraft(ctx context.Context, request api.GetDraftRequestObject) (api.GetDraftResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) GetTaskDetails(ctx context.Context, request api.GetTaskDetailsRequestObject) (api.GetTaskDetailsResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) ListDrafts(ctx context.Context, request api.ListDraftsRequestObject) (api.ListDraftsResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) ListTasks(ctx context.Context, request api.ListTasksRequestObject) (api.ListTasksResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) RecreateTask(ctx context.Context, request api.RecreateTaskRequestObject) (api.RecreateTaskResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) UpdateDraft(ctx context.Context, request api.UpdateDraftRequestObject) (api.UpdateDraftResponseObject, error) {
	panic("unimplemented")
}

func (d *AppDelivery) GetTaskInternalState(ctx context.Context, request api.GetTaskInternalStateRequestObject) (api.GetTaskInternalStateResponseObject, error) {
	panic("unimplemented")
}
