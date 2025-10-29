package usecase

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
)

func (u *appUsecaseImpl) CreateDraft(pageURL string) (*api.DraftDigest, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	// Extract the YWiki slug from the page URL
	slug := extractYWikiSlugFromURL(pageURL)

	// Get the page by its slug
	page, err := repo.GetPageBySlug(slug)
	if err != nil {
		return nil, err
	}

	// Create a draft with the page ID, title, and content
	draftID, err := repo.CreateDraft(page.PageId, page.Title, page.Content)
	if err != nil {
		return nil, err
	}

	// Get the created draft to return its digest
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

	// Remove the draft by its ID
	err := repo.RemoveDraft(draftID)
	if err != nil {
		return err
	}

	return repo.Commit()
}

func (u *appUsecaseImpl) GetDraft(draftID api.DraftID) (*api.Draft, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	// Get the draft by its ID
	draft, err := repo.GetDraftByID(draftID)
	if err != nil {
		return nil, err
	}

	return draft, nil
}

func (u *appUsecaseImpl) ApplyDraft(draftID api.DraftID) error {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	// Get the draft by its ID
	draft, err := repo.GetDraftByID(draftID)
	if err != nil {
		return err
	}

	// Append a new revision to the page with the draft content
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

func (u *appUsecaseImpl) ListDrafts(cursor *string) ([]api.DraftDigest, *api.Cursor, error) {
	repo := repository.NewAppRepository(u.ctx, u.deps)
	defer repo.Rollback()

	// Convert string cursor to *api.Cursor
	var apiCursor *api.Cursor
	if cursor != nil {
		apiCursor = (*api.Cursor)(cursor)
	}

	// List drafts with a reasonable limit (e.g., 50)
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
