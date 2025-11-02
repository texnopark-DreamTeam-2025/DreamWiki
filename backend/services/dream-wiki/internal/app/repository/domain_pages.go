package repository

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

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
		p.page_id,
		p.title,
		p.ywiki_slug,
		p.current_revision_id,
		r.content
	FROM Page p
	JOIN PageRevision r ON p.current_revision_id=r.revision_id
	WHERE p.ywiki_slug=$yWikiSlug;
	`

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$yWikiSlug", types.TextValue(yWikiSlug)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var pageID api.PageID
	var title string
	var ywikiSlug string
	var currentRevisionID int64
	var content string
	if err = result.FetchExactlyOne(&pageID, &title, &ywikiSlug, &currentRevisionID, &content); err != nil {
		return nil, err
	}

	return &api.Page{
		PageId:    pageID,
		Title:     title,
		Content:   content,
		YwikiSlug: ywikiSlug,
	}, nil
}

func (r *appRepositoryImpl) CreatePage(yWikiSlug string, title string, content string) (*api.PageID, error) {
	yql := `
	INSERT INTO Page(page_id, title, ywiki_slug)
	VALUES (
		RandomUuid(4),
		$title,
		$yWikiSlug
	)
	RETURNING page_id;`

	parameters := []table.ParameterOption{
		table.ValueParam("$title", types.TextValue(title)),
		table.ValueParam("$yWikiSlug", types.TextValue(yWikiSlug)),
	}

	result, err := r.ydbClient.InTX().Execute(yql, parameters...)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var pageID api.PageID
	err = result.FetchExactlyOne(&pageID)
	if err != nil {
		return nil, err
	}

	r.log.Debug("Inserted page with id ", pageID)

	_, err = r.AppendPageRevision(pageID, content)
	if err != nil {
		return nil, err
	}

	return &pageID, nil
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
	yql1 := `
	SELECT current_revision_id
	FROM Page
	WHERE page_id = $pageID;`

	result1, err := r.ydbClient.InTX().Execute(yql1, table.ValueParam("$pageID", types.UuidValue(pageID)))
	if err != nil {
		return nil, err
	}
	defer result1.Close()

	var currentRevisionID *int64
	err = result1.FetchExactlyOne(&currentRevisionID)
	if err != nil {
		return nil, err
	}

	var ydbRevisionID types.Value
	if currentRevisionID != nil {
		ydbRevisionID = types.OptionalValue(types.Int64Value(*currentRevisionID))
	} else {
		ydbRevisionID = types.NullValue(types.TypeInt64)
	}

	yql2 := `
	INSERT INTO PageRevision(page_id, previous_revision_id, content)
	VALUES (
		$pageID,
		$previousRevisionID,
		$content
	)
	RETURNING revision_id;`

	parameters := []table.ParameterOption{
		table.ValueParam("$pageID", types.UuidValue(pageID)),
		table.ValueParam("$previousRevisionID", ydbRevisionID),
		table.ValueParam("$content", types.TextValue(newContent)),
	}

	result2, err := r.ydbClient.InTX().Execute(yql2, parameters...)
	if err != nil {
		return nil, err
	}
	defer result2.Close()

	var revisionID internals.RevisionID
	err = result2.FetchExactlyOne(&revisionID)
	if err != nil {
		return nil, err
	}

	yql3 := `
	UPDATE Page
	SET current_revision_id = $revisionID
	WHERE page_id = $pageID;`

	parameters3 := []table.ParameterOption{
		table.ValueParam("$pageID", types.UuidValue(pageID)),
		table.ValueParam("$revisionID", types.Int64Value(revisionID)),
	}

	result3, err := r.ydbClient.InTX().Execute(yql3, parameters3...)
	if err != nil {
		return nil, err
	}
	defer result3.Close()

	r.log.Debug("Appended page ", pageID, " revision with id ", revisionID)

	return &revisionID, nil
}

func (r *appRepositoryImpl) GetPageByID(pageID api.PageID) (*api.Page, *internals.PageAdditionalInfo, error) {
	yql := `
	SELECT
		p.page_id,
		p.title,
		p.ywiki_slug,
		p.current_revision_id,
		r.content
	FROM Page p
	JOIN PageRevision r ON p.current_revision_id=r.revision_id
	WHERE p.page_id=$pageID;
	`

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$pageID", types.UuidValue(pageID)))
	if err != nil {
		return nil, nil, err
	}
	defer result.Close()

	var retrievedPageID api.PageID
	var title string
	var ywikiSlug string
	var currentRevisionID int64
	var content string
	if err = result.FetchExactlyOne(&retrievedPageID, &title, &ywikiSlug, &currentRevisionID, &content); err != nil {
		return nil, nil, err
	}

	page := &api.Page{
		PageId:    retrievedPageID,
		Title:     title,
		Content:   content,
		YwikiSlug: ywikiSlug,
	}

	pageAdditionalInfo := &internals.PageAdditionalInfo{
		CurrentRevisionId: &currentRevisionID,
	}

	return page, pageAdditionalInfo, nil
}

func (r *appRepositoryImpl) SetPageTitle(pageID api.PageID, newTitle string) error {
	yql := `
	UPDATE Page
	SET title = $newTitle
	WHERE page_id = $pageID;
	`

	parameters := []table.ParameterOption{
		table.ValueParam("$pageID", types.UuidValue(pageID)),
		table.ValueParam("$newTitle", types.TextValue(newTitle)),
	}

	result, err := r.ydbClient.InTX().Execute(yql, parameters...)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}
