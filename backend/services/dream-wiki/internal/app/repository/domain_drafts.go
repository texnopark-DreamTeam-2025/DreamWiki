package repository

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

func (r *appRepositoryImpl) GetDraftByID(draftID api.DraftID) (*api.Draft, error) {
	panic("unimplemented")
}

func (r *appRepositoryImpl) ListDrafts(cursor *api.Cursor, limit int64) ([]api.DraftDigest, *api.Cursor, error) {
	panic("unimplemented")
}

func (r *appRepositoryImpl) RemoveDraft(draftID api.DraftID) error {
	panic("unimplemented")
}

func (r *appRepositoryImpl) SetDraftBaseRevision(draftID api.DraftID, newRevisionID internals.RevisionID) error {
	panic("unimplemented")
}

func (r *appRepositoryImpl) SetDraftContent(draftID api.DraftID, newContent string) error {
	panic("unimplemented")
}

func (r *appRepositoryImpl) SetDraftStatus(draftID api.DraftID, newStatus api.DraftStatus) error {
	panic("unimplemented")
}

func (r *appRepositoryImpl) SetDraftTitle(draftID api.DraftID, newTitle string) error {
	panic("unimplemented")
}

func (r *appRepositoryImpl) CreateDraft(pageID api.PageID, draftTitle string, draftContent string) (*api.DraftID, error) {
	panic("unimplemented")
}
