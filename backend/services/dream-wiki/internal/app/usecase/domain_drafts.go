package usecase

import "github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"

func (u *appUsecaseImpl) CreateDraft(pageURL string) (api.DraftDigest, error) {
	panic("unimplemented")
}

func (u *appUsecaseImpl) DeleteDraft(draftID api.DraftID) error {
	panic("unimplemented")
}

func (u *appUsecaseImpl) GetDraft(draftID api.DraftID) (api.Draft, error) {
	panic("unimplemented")
}

func (u *appUsecaseImpl) ApplyDraft(draftID api.DraftID) error {
	panic("unimplemented")
}

func (u *appUsecaseImpl) ListDrafts(cursor *string) ([]api.DraftDigest, error) {
	panic("unimplemented")
}

func (u *appUsecaseImpl) UpdateDraft(draftID api.DraftID, newContent *string, newTitle *string) error {
	panic("unimplemented")
}
