package delivery

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (d *AppDelivery) IndexatePage(ctx context.Context, request api.IndexatePageRequestObject) (api.IndexatePageResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.IndexatePage(request.Body.PageId)
	if err != nil {
		d.log.Error(err.Error())
		return api.IndexatePage500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.IndexatePage200JSONResponse(*resp), nil
}
