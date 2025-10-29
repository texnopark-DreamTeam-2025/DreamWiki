package delivery

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (d *AppDelivery) CancelTask(ctx context.Context, request api.CancelTaskRequestObject) (api.CancelTaskResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	err := usecase.CancelTask(request.Body.TaskId)
	if err != nil {
		d.log.Error(err.Error())
		return api.CancelTask500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.CancelTask200JSONResponse{}, nil
}

func (d *AppDelivery) GetTaskDetails(ctx context.Context, request api.GetTaskDetailsRequestObject) (api.GetTaskDetailsResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	result, err := usecase.GetTaskDetails(request.Body.TaskId)
	if err != nil {
		d.log.Error(err.Error())
		return api.GetTaskDetails500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.GetTaskDetails200JSONResponse{Task: result}, nil
}

func (d *AppDelivery) ListTasks(ctx context.Context, request api.ListTasksRequestObject) (api.ListTasksResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	result, newCursor, err := usecase.ListTasks(request.Body.Cursor)
	if err != nil {
		d.log.Error(err.Error())
		return api.ListTasks500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.ListTasks200JSONResponse{
		Tasks:  result,
		Cursor: newCursor,
	}, nil
}

func (d *AppDelivery) RecreateTask(ctx context.Context, request api.RecreateTaskRequestObject) (api.RecreateTaskResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	result, err := usecase.RecreateTask(request.Body.TaskId)
	if err != nil {
		d.log.Error(err.Error())
		return api.RecreateTask500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.RecreateTask200JSONResponse{NewTaskId: *result}, nil
}

func (d *AppDelivery) GetTaskInternalState(ctx context.Context, request api.GetTaskInternalStateRequestObject) (api.GetTaskInternalStateResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	result, err := usecase.GetTaskInternalState(request.Body.TaskId)
	if err != nil {
		d.log.Error(err.Error())
		return api.GetTaskInternalState500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.GetTaskInternalState200JSONResponse{
		Actions:   result.Actions,
		TaskId:    result.TaskId,
		TaskState: result.TaskState,
	}, nil
}
