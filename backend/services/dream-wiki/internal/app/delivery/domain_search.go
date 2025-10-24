package delivery

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (d *AppDelivery) Search(ctx context.Context, request api.SearchRequestObject) (api.SearchResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.Search(*request.Body)
	if err != nil {
		d.log.Error(err.Error())
		return api.Search500JSONResponse{ErrorResponseJSONResponse: api.ErrorResponseJSONResponse{Message: internalErrorMessage}}, nil
	}

	return api.Search200JSONResponse(*resp), nil
}
