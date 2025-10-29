package delivery

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (d *AppDelivery) CreateDraft(ctx context.Context, request api.CreateDraftRequestObject) (api.CreateDraftResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	result, err := usecase.CreateDraft(request.Body.PageUrl)
	if err != nil {
		d.log.Error(err.Error())
		return api.CreateDraft500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.CreateDraft200JSONResponse{DraftId: result.DraftId}, nil
}

func (d *AppDelivery) DeleteDraft(ctx context.Context, request api.DeleteDraftRequestObject) (api.DeleteDraftResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	err := usecase.DeleteDraft(request.Body.DraftId)
	if err != nil {
		d.log.Error(err.Error())
		return api.DeleteDraft500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.DeleteDraft200JSONResponse{}, nil
}

func (d *AppDelivery) GetDraft(ctx context.Context, request api.GetDraftRequestObject) (api.GetDraftResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	result, err := usecase.GetDraft(request.Body.DraftId)
	if err != nil {
		d.log.Error(err.Error())
		return api.GetDraft500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.GetDraft200JSONResponse{Draft: *result}, nil
}

func (d *AppDelivery) ApplyDraft(ctx context.Context, request api.ApplyDraftRequestObject) (api.ApplyDraftResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	err := usecase.ApplyDraft(request.Body.DraftId)
	if err != nil {
		d.log.Error(err.Error())
		return api.ApplyDraft500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.ApplyDraft200JSONResponse{}, nil
}

func (d *AppDelivery) ListDrafts(ctx context.Context, request api.ListDraftsRequestObject) (api.ListDraftsResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	result, newCursor, err := usecase.ListDrafts(request.Body.Cursor)
	if err != nil {
		d.log.Error(err.Error())
		return api.ListDrafts500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.ListDrafts200JSONResponse{
		Drafts: result,
		Cursor: *newCursor,
	}, nil
}

func (d *AppDelivery) UpdateDraft(ctx context.Context, request api.UpdateDraftRequestObject) (api.UpdateDraftResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	err := usecase.UpdateDraft(request.Body.DraftId, request.Body.NewContent, request.Body.NewTitle)
	if err != nil {
		d.log.Error(err.Error())
		return api.UpdateDraft500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.UpdateDraft200JSONResponse{}, nil
}
