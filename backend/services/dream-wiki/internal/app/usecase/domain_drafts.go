package usecase

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (u *appUsecaseImpl) CreateDraft(originalPageID api.PageID) (*api.DraftDigest, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	u.log.Info("PAGE ID: #############", originalPageID)

	page, _, err := repo.GetPageByID(originalPageID)
	if err != nil {
		return nil, err
	}

	draftID, err := repo.CreateDraft(page.PageId, page.Title, page.Content)
	if err != nil {
		return nil, err
	}

	draft, err := repo.GetDraftByID(*draftID)
	if err != nil {
		return nil, err
	}

	err = repo.Commit()
	if err != nil {
		return nil, err
	}

	return &draft.DraftDigest, nil
}

func (u *appUsecaseImpl) DeleteDraft(draftID api.DraftID) error {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	err := repo.RemoveDraft(draftID)
	if err != nil {
		return err
	}

	return repo.Commit()
}

func (u *appUsecaseImpl) GetDraft(draftID api.DraftID) (*api.Draft, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	draft, err := repo.GetDraftByID(draftID)
	if err != nil {
		return nil, err
	}

	return draft, nil
}

func (u *appUsecaseImpl) ApplyDraft(draftID api.DraftID) error {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	draft, err := repo.GetDraftByID(draftID)
	if err != nil {
		return err
	}

	_, err = repo.AppendPageRevision(draft.DraftDigest.PageDigest.PageId, draft.Content)
	if err != nil {
		return err
	}

	err = repo.SetDraftStatus(draftID, api.Merged)
	if err != nil {
		return err
	}

	return repo.Commit()
}

func (u *appUsecaseImpl) ListDrafts(cursor *string) ([]api.DraftDigest, *api.NextInfo, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	var apiCursor *api.Cursor
	if cursor != nil {
		apiCursor = (*api.Cursor)(cursor)
	}

	drafts, newCursor, err := repo.ListDrafts(apiCursor, 50)
	if err != nil {
		return nil, nil, err
	}

	err = repo.Commit()
	if err != nil {
		return nil, nil, err
	}
	return drafts, newCursor, nil
}

func (u *appUsecaseImpl) UpdateDraft(draftID api.DraftID, newContent *string, newTitle *string) error {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	if newContent != nil {
		err := repo.SetDraftContent(draftID, *newContent)
		if err != nil {
			return err
		}
	}

	if newTitle != nil {
		err := repo.SetDraftTitle(draftID, *newTitle)
		if err != nil {
			return err
		}
	}

	return repo.Commit()
}
