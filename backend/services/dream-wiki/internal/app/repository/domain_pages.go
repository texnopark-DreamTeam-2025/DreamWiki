package repository

import (
	"github.com/google/uuid"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func (r *appRepositoryImpl) RetrievePageByID(pageID uuid.UUID) (*api.Page, error) {
	yql := `
	SELECT
		page_id,
		title,
		content
	FROM Page WHERE page_id=$pageID;
	`

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$pageID", types.UuidValue(pageID)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var retrievedPageID uuid.UUID
	var title string
	var content string
	if err = result.FetchExactlyOne(&retrievedPageID, &title, &content); err != nil {
		return nil, err
	}

	return &api.Page{
		PageId:  retrievedPageID,
		Content: content,
		Title:   title,
	}, nil
}

func (r *appRepositoryImpl) DeleteAllPages() error {
	yql := `
		DELETE FROM Paragraph;
		DELETE FROM Page;
	`

	result, err := r.ydbClient.InTX().Execute(yql)
	if err != nil {
		return err
	}
	defer result.Close()

	r.log.Info("All pages and paragraphs deleted successfully")
	return nil
}

func (r *appRepositoryImpl) GetPageBySlug(yWikiSlug string) (*api.Page, error) {
	yql := `
	SELECT
		page_id,
		title,
		content
	FROM Page WHERE ywiki_slug=$yWikiSlug;
	`

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$yWikiSlug", types.TextValue(yWikiSlug)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var pageID api.PageID
	var title string
	var content string
	if err = result.FetchExactlyOne(&pageID, &title, &content); err != nil {
		return nil, err
	}

	return &api.Page{
		PageId:  pageID,
		Title:   title,
		Content: content,
	}, nil
}

func (r *appRepositoryImpl) CreatePage(yWikiSlug string, title string, content string) (*api.PageID, error) {
	panic("unimplemented")
}

func (r *appRepositoryImpl) UpsertPage(page api.Page, yWikiSlug string) (pageID *uuid.UUID, err error) {
	yql1 := `
	UPDATE Page
	SET
		title=$title,
		content=$content
	WHERE ywiki_slug=$yWikiSlug
	RETURNING page_id;`

	yql2 := `
	INSERT INTO Page(page_id, title, ywiki_slug, content)
	VALUES (
		RandomUUID(4),
		$title,
		$yWikiSlug,
		$content
	)
	RETURNING page_id;`

	pageID = new(uuid.UUID)

	parameters := []table.ParameterOption{
		table.ValueParam("$title", types.TextValue(page.Title)),
		table.ValueParam("$content", types.TextValue(page.Content)),
		table.ValueParam("$yWikiSlug", types.TextValue(yWikiSlug)),
	}

	result1, err := r.ydbClient.InTX().Execute(yql1, parameters...)
	if err != nil {
		return nil, err
	}
	defer result1.Close()

	if result1.RowCount() == 1 {
		err := result1.FetchExactlyOne(&pageID)
		if err != nil {
			return nil, err
		}
		return pageID, nil
	}

	result2, err := r.ydbClient.InTX().Execute(yql2, parameters...)
	if err != nil {
		return nil, err
	}
	defer result2.Close()

	err = result2.FetchExactlyOne(&pageID)
	if err != nil {
		return nil, err
	}

	r.log.Debug("Inserted page with id ", pageID)

	return pageID, nil
}

func (r *appRepositoryImpl) DeletePageBySlug(yWikiSlug string) error {
	yql := `
	DELETE FROM Page WHERE ywiki_slug=$yWikiSlug;`

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$yWikiSlug", types.TextValue(yWikiSlug)))
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}

func (r *appRepositoryImpl) GetAllPageDigests() ([]api.PageDigest, error) {
	yql := `SELECT page_id, title FROM Page`

	result, err := r.ydbClient.InTX().Execute(yql)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	pages := make([]api.PageDigest, 0, result.RowCount())
	for result.NextRow() {
		var page api.PageDigest
		err := result.FetchRow(&page.PageId, &page.Title)
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func (r *appRepositoryImpl) AppendPageRevision(pageID api.PageID, newContent string) (*internals.RevisionID, error) {
	panic("unimplemented")
}

func (r *appRepositoryImpl) GetPageByID(pageID api.PageID) (*api.Page, *internals.PageAdditionalInfo, error) {
	panic("unimplemented")
}

func (r *appRepositoryImpl) SetPageActualRevision(pageID api.PageID, revisionID internals.RevisionID) error {
	panic("unimplemented")
}

func (r *appRepositoryImpl) SetPageTitle(pageID api.PageID, newTitle string) error {
	panic("unimplemented")
}
