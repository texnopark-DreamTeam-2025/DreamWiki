package delivery

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (d *AppDelivery) Login(ctx context.Context, request api.LoginRequestObject) (api.LoginResponseObject, error) {
	usecase := usecase.NewAppUsecaseImpl(ctx, d.deps)
	resp, err := usecase.Login(*request.Body)
	if err != nil {
		d.log.Error(err.Error())
		return api.Login500JSONResponse{Message: "Internal server error"}, nil
	}

	return api.Login200JSONResponse(*resp), nil
}
