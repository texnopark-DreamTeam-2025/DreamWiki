package usecase

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (u *appUsecaseImpl) GetPagesTree(activePagesIDs []api.PageID) ([]api.TreeItem, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	items, err := repo.GetAllPageDigests()
	if err != nil {
		return nil, err
	}

	result := make([]api.TreeItem, 0, len(items))
	for _, item := range items {
		result = append(result, api.TreeItem{
			PageDigest: item,
			Children:   nil,
			Expanded:   false,
		})
	}

	return result, nil
}

func (u *appUsecaseImpl) GetDiagnosticInfo(req api.V1DiagnosticInfoGetRequest) (*api.V1DiagnosticInfoGetResponse, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	page, _, err := repo.GetPageByID(req.PageId)
	if err != nil {
		return nil, err
	}

	return &api.V1DiagnosticInfoGetResponse{
		Page: *page,
	}, nil
}
