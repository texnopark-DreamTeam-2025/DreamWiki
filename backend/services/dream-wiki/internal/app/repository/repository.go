package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/ydb_wrapper"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

type (
	AppRepository interface {
		Commit()
		Rollback()

		SearchByEmbedding(query string, queryEmbedding models.Embedding) ([]api.SearchResultItem, error)
		RetrievePageByID(pageID uuid.UUID) (*api.Page, error)
		RemovePageIndexation(pageID uuid.UUID) error
		AddIndexedParagraph(paragraph models.ParagraphWithEmbedding) error
		DeleteAllPages() error
		GetAllPageDigests() ([]api.PageDigest, error)
		GetUserByLogin(login string) (*models.User, error)
		WriteIntegrationLogField(integrationID api.IntegrationID, logText string) error
		GetIntegrationLogFields(integrationID string, cursor *string, limit int) (fields []api.IntegrationLogField, newCursor string, err error)
		GetPageBySlug(yWikiSlug string) (*api.Page, error)
		UpsertPage(page api.Page, ywikiSlug string) (pageID *uuid.UUID, err error)
		DeletePageBySlug(yWikiSlug string) error
	}

	appRepositoryImpl struct {
		ctx       context.Context
		ydbClient ydb_wrapper.YDBWrapper
		log       logger.Logger
	}
)

func NewAppRepository(ctx context.Context, deps *deps.Deps) AppRepository {
	return &appRepositoryImpl{
		ctx:       ctx,
		log:       deps.Logger,
		ydbClient: ydb_wrapper.NewYDBWrapper(ctx, deps, true),
	}
}

func (r *appRepositoryImpl) Commit() {
	r.ydbClient.Commit()
}

func (r *appRepositoryImpl) Rollback() {
	r.ydbClient.Rollback()
}

func (r *appRepositoryImpl) SearchByEmbedding(query string, queryEmbedding models.Embedding) ([]api.SearchResultItem, error) {
	yql := `
		$K = 20;
		$targetEmbedding = Knn::ToBinaryStringFloat($queryEmbedding);

		SELECT
			par.page_id,
			page.title,
			par.content,
			Unwrap(Knn::CosineDistance(Unwrap(par.embedding), $targetEmbedding)) As CosineDistance
		FROM Paragraph par
		JOIN Page page USING(page_id)
		ORDER BY Knn::CosineDistance(embedding, $targetEmbedding)
		LIMIT $K;
	`

	yqlEmbedding := embeddingToYDBList(queryEmbedding)

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$queryEmbedding", yqlEmbedding))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	searchResult := make([]api.SearchResultItem, 0)
	for result.NextRow() {
		var retrievedPageID uuid.UUID
		var title string
		var pageContent string
		var distance float32
		err = result.FetchRow(&retrievedPageID, &title, &pageContent, &distance)
		if err != nil {
			return nil, err
		}
		r.log.Debug("Distance is ", distance)
		searchResult = append(searchResult, api.SearchResultItem{
			PageId:      retrievedPageID,
			Title:       title,
			Description: pageContent,
		})
	}

	return searchResult, nil
}

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

func (r *appRepositoryImpl) RemovePageIndexation(pageID uuid.UUID) error {
	yql := `
		DELETE FROM Paragraph WHERE page_id=$pageID;
	`

	result, err := r.ydbClient.InTX().Execute(yql,
		table.ValueParam("$pageID", types.UuidValue(pageID)),
	)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}

func (r *appRepositoryImpl) AddIndexedParagraph(paragraph models.ParagraphWithEmbedding) error {
	yql := `
		INSERT INTO Paragraph (page_id, line_number, content, embedding, anchor_link_slug)
		VALUES (
			$pageID,
			$lineNumber,
			$content,
			Untag(Knn::ToBinaryStringFloat($embedding), "FloatVector"),
			$anchorLineSlug
		);
	`

	result, err := r.ydbClient.InTX().Execute(yql,
		table.ValueParam("$pageID", types.UuidValue(paragraph.PageID)),
		table.ValueParam("$lineNumber", types.Int32Value(int32(paragraph.LineNumber))),
		table.ValueParam("$content", types.TextValue(paragraph.Content)),
		table.ValueParam("$embedding", embeddingToYDBList(paragraph.Embedding)),
		table.ValueParam("$anchorLineSlug", types.TextValue("")), // TODO
	)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
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

func (r *appRepositoryImpl) GetUserByLogin(login string) (*models.User, error) {
	yql := `
	SELECT
		user_id,
		login,
		password_hash_bcrypt
	FROM User WHERE login=$login;
	`

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$login", types.TextValue(login)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var userID uuid.UUID
	var userLogin string
	var passwordHash string
	if err = result.FetchExactlyOne(result, &userID, &userLogin, &passwordHash); err != nil {
		return nil, err
	}

	return &models.User{
		ID:           userID,
		Login:        userLogin,
		PasswordHash: passwordHash,
	}, nil
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

	var pageID uuid.UUID
	var title string
	var content string
	if err = result.FetchExactlyOne(result, &pageID, &title, &content); err != nil {
		return nil, err
	}

	return &api.Page{
		PageId:  pageID,
		Title:   title,
		Content: content,
	}, nil
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
	defer result.Close()

	return err
}

func (r *appRepositoryImpl) WriteIntegrationLogField(integrationID api.IntegrationID, logText string) error {
	yql := `INSERT INTO IntegrationLogField (integration_id, log_text, created_at)
	VALUES ($integrationID, $logText, CurrentUtcDatetime())`

	result, err := r.ydbClient.OutsideTX().Execute(yql,
		table.ValueParam("$integrationID", types.TextValue(string(integrationID))),
		table.ValueParam("$logText", types.TextValue(logText)),
	)
	defer result.Close()

	return err
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

func (r *appRepositoryImpl) GetIntegrationLogFields(integrationID string, cursor *string, limit int) (fields []api.IntegrationLogField, newCursor string, err error) {
	yql := `
		SELECT field_id, log_text, created_at
		FROM IntegrationLogField
		WHERE integration_id=$integrationID
			AND (
				created_at > $timeFrom
				OR (created_at = $timeFrom AND field_id > $idFrom)
			)
		ORDER BY created_at DESC
		LIMIT $limit
	`

	timeFrom, idFrom := decodeCursor(cursor)

	result, err := r.ydbClient.InTX().Execute(yql,
		table.ValueParam("$integrationID", types.TextValue(integrationID)),
		table.ValueParam("$limit", types.Int32Value(int32(limit))),
		table.ValueParam("$timeFrom", types.TimestampValueFromTime(timeFrom)),
		table.ValueParam("$idFrom", types.Int64Value(idFrom)),
	)
	if err != nil {
		return nil, "", err
	}
	defer result.Close()

	fields = make([]api.IntegrationLogField, 0, limit)
	if result.RowCount() == 0 {
		if cursor == nil {
			return fields, "", nil
		}
		return fields, *cursor, nil
	}

	newIDFrom := int64(0)
	for result.NextRow() {
		var content string
		var createdAt time.Time
		err := result.FetchRow(&newIDFrom, &content, &createdAt)
		if err != nil {
			return nil, "", err
		}
		fields = append(fields, api.IntegrationLogField{
			Content:   content,
			CreatedAt: createdAt,
		})
	}

	newTimeFrom := fields[len(fields)-1].CreatedAt
	return fields, encodeCursor(newTimeFrom, newIDFrom), nil
}
