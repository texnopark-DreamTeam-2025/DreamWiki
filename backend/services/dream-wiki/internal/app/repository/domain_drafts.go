package repository

import (
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func (r *appRepositoryImpl) GetDraftByID(draftID api.DraftID) (*api.Draft, error) {
	yql := `
	SELECT
		d.draft_id,
		d.draft_title,
		d.content,
		d.status,
		d.created_at,
		d.updated_at,
		d.page_revision_id,
		r.revision_id,
		p.page_id,
		p.title
	FROM Draft d
	JOIN PageRevision r ON r.revision_id=d.page_revision_id
	JOIN Page p ON r.page_id = p.current_revision_id
	WHERE d.draft_id = $draftID;
	`

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$draftID", types.UuidValue(draftID)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var originalContent string
	var retrievedDraftID api.DraftID
	var draftTitle string
	var content string
	var status string
	var createdAt time.Time
	var updatedAt time.Time
	var pageRevisionID int64
	var baseRevisionID int64
	var pageID api.PageID
	var pageTitle string

	err = result.FetchExactlyOne(
		&retrievedDraftID,
		&draftTitle,
		&content,
		&status,
		&createdAt,
		&updatedAt,
		&pageRevisionID,
		&baseRevisionID,
		&pageID,
		&pageTitle,
	)
	if err != nil {
		return nil, err
	}

	if status == string(api.Active) && pageRevisionID != baseRevisionID {
		status = string(api.NeedsRebase)
	}

	return &api.Draft{
		Content:   content,
		CreatedAt: createdAt,
		DraftDigest: api.DraftDigest{
			DraftId:    retrievedDraftID,
			DraftTitle: draftTitle,
			PageDigest: api.PageDigest{
				PageId: pageID,
				Title:  pageTitle,
			},
			Status: api.DraftStatus(status),
		},
		OriginalPageContent: &originalContent,
		UpdatedAt:           updatedAt,
	}, nil
}

func (r *appRepositoryImpl) ListDrafts(cursor *api.Cursor, limit int64) ([]api.DraftDigest, *api.Cursor, error) {
	yql := `
	SELECT
		d.draft_id,
		d.draft_title,
		d.status,
		p.page_id,
		p.title
	FROM Draft d
	JOIN Page p ON d.page_revision_id = p.current_revision_id
	ORDER BY d.updated_at DESC
	LIMIT $limit;
	`

	parameters := []table.ParameterOption{
		table.ValueParam("$limit", types.Int64Value(limit)),
	}

	result, err := r.ydbClient.InTX().Execute(yql, parameters...)
	if err != nil {
		return nil, nil, err
	}
	defer result.Close()

	drafts := make([]api.DraftDigest, 0, result.RowCount())
	for result.NextRow() {
		var draftID api.DraftID
		var draftTitle string
		var status string
		var pageID api.PageID
		var pageTitle string

		err := result.FetchRow(&draftID, &draftTitle, &status, &pageID, &pageTitle)
		if err != nil {
			return nil, nil, err
		}

		drafts = append(drafts, api.DraftDigest{
			DraftId:    draftID,
			DraftTitle: draftTitle,
			PageDigest: api.PageDigest{
				PageId: pageID,
				Title:  pageTitle,
			},
			Status: api.DraftStatus(status),
		})
	}

	// For simplicity, we're not implementing cursor pagination here
	// In a real implementation, you would need to handle cursor-based pagination
	return drafts, nil, nil
}

func (r *appRepositoryImpl) RemoveDraft(draftID api.DraftID) error {
	yql := `
	DELETE FROM Draft WHERE draft_id = $draftID;
	`

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$draftID", types.UuidValue(draftID)))
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}

func (r *appRepositoryImpl) SetDraftBaseRevision(draftID api.DraftID, newRevisionID internals.RevisionID) error {
	yql := `
	UPDATE Draft
	SET page_revision_id = $newRevisionID, updated_at = CurrentUtcDatetime()
	WHERE draft_id = $draftID;
	`

	parameters := []table.ParameterOption{
		table.ValueParam("$draftID", types.UuidValue(draftID)),
		table.ValueParam("$newRevisionID", types.Int64Value(newRevisionID)),
	}

	result, err := r.ydbClient.InTX().Execute(yql, parameters...)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}

func (r *appRepositoryImpl) SetDraftContent(draftID api.DraftID, newContent string) error {
	yql := `
	UPDATE Draft
	SET content = $newContent, updated_at = CurrentUtcDatetime()
	WHERE draft_id = $draftID;
	`

	parameters := []table.ParameterOption{
		table.ValueParam("$draftID", types.UuidValue(draftID)),
		table.ValueParam("$newContent", types.TextValue(newContent)),
	}

	result, err := r.ydbClient.InTX().Execute(yql, parameters...)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}

func (r *appRepositoryImpl) SetDraftStatus(draftID api.DraftID, newStatus api.DraftStatus) error {
	yql := `
	UPDATE Draft
	SET status = $newStatus, updated_at = CurrentUtcDatetime()
	WHERE draft_id = $draftID;
	`

	parameters := []table.ParameterOption{
		table.ValueParam("$draftID", types.UuidValue(draftID)),
		table.ValueParam("$newStatus", types.TextValue(string(newStatus))),
	}

	result, err := r.ydbClient.InTX().Execute(yql, parameters...)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}

func (r *appRepositoryImpl) SetDraftTitle(draftID api.DraftID, newTitle string) error {
	yql := `
	UPDATE Draft
	SET draft_title = $newTitle, updated_at = CurrentUtcDatetime()
	WHERE draft_id = $draftID;
	`

	parameters := []table.ParameterOption{
		table.ValueParam("$draftID", types.UuidValue(draftID)),
		table.ValueParam("$newTitle", types.TextValue(newTitle)),
	}

	result, err := r.ydbClient.InTX().Execute(yql, parameters...)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}

func (r *appRepositoryImpl) CreateDraft(pageID api.PageID, draftTitle string, draftContent string) (*api.DraftID, error) {
	// First, get the current revision ID for the page
	pageYql := `
	SELECT current_revision_id
	FROM Page
	WHERE page_id = $pageID;
	`

	pageResult, err := r.ydbClient.InTX().Execute(pageYql, table.ValueParam("$pageID", types.UuidValue(pageID)))
	if err != nil {
		return nil, err
	}
	defer pageResult.Close()

	var currentRevisionID int64
	if err = pageResult.FetchExactlyOne(&currentRevisionID); err != nil {
		return nil, err
	}

	// Create the draft
	yql := `
	INSERT INTO Draft (draft_id, page_revision_id, status, draft_title, content, created_at, updated_at)
	VALUES (RandomUuid(4), $pageRevisionID, 'active', $draftTitle, $content, CurrentUtcDatetime(), CurrentUtcDatetime())
	RETURNING draft_id;
	`

	parameters := []table.ParameterOption{
		table.ValueParam("$pageRevisionID", types.Int64Value(currentRevisionID)),
		table.ValueParam("$draftTitle", types.TextValue(draftTitle)),
		table.ValueParam("$content", types.TextValue(draftContent)),
	}

	result, err := r.ydbClient.InTX().Execute(yql, parameters...)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var draftID api.DraftID
	err = result.FetchExactlyOne(&draftID)
	if err != nil {
		return nil, err
	}

	return &draftID, nil
}
