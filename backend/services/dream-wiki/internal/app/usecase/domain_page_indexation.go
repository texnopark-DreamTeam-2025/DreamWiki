package usecase

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (u *appUsecaseImpl) IndexatePage(pageID api.PageID) (*api.V1IndexatePageResponse, error) {
	_, err := u.CreatePageReindexationTask([]api.PageID{pageID})
	if err != nil {
		return nil, err
	}

	return &api.V1IndexatePageResponse{
		PageId: pageID,
	}, nil
}
