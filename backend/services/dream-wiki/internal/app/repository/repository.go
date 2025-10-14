package repository

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result/indexed"
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
		GetUserByLogin(login string) (*models.User, error)
		WriteIntegrationLogField(integrationID api.IntegrationID, logText string) error
		GetPageBySlug(yWikiSlug string) (*api.Page, error)
		UpsertPage(page api.Page, ywikiSlug string) (pageID *uuid.UUID, err error)
		DeletePageBySlug(yWikiSlug string) error
	}

	appRepositoryImpl struct {
		ctx     context.Context
		tx      table.TransactionActor
		log     logger.Logger
		success chan bool
		closed  int32
	}
)

func StartTransaction(ctx context.Context, deps *deps.Deps) AppRepository {
	deps.Logger.Debug("starting YDB transaction...")
	success := make(chan bool)
	txRetriever := make(chan table.TransactionActor)

	action := func(ctx context.Context, tx table.TransactionActor) error {
		txRetriever <- tx
		deps.Logger.Infof("YDB transaction %s started", tx.ID())

		shouldCommit := <-success
		close(success)
		if shouldCommit {
			deps.Logger.Infof("YDB transaction %s committed", tx.ID())
			return nil
		}
		deps.Logger.Infof("YDB transaction %s rolled back", tx.ID())
		return fmt.Errorf("transaction rolled back")
	}

	go deps.DB.DoTx(context.Background(), action)
	tx := <-txRetriever
	close(txRetriever)
	return &appRepositoryImpl{
		ctx: ctx, tx: tx, success: success, log: deps.Logger}
}

func (r *appRepositoryImpl) Commit() {
	if n := atomic.AddInt32(&(r.closed), 1); n == 1 {
		r.success <- true
	}
}

func (r *appRepositoryImpl) Rollback() {
	if n := atomic.AddInt32(&(r.closed), 1); n == 1 {
		r.success <- false
	}
}

func (r *appRepositoryImpl) execute(yql string, opts ...table.ParameterOption) (result result.Result, err error) {
	r.log.Debug("executing yql: ", yql, opts)
	result, err = r.tx.Execute(r.ctx, yql, table.NewQueryParameters(opts...))
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	err = result.Err()
	if err != nil {
		r.log.Error(err)
		return nil, err
	}
	return
}

func (r *appRepositoryImpl) nextResultSet(result result.Result) error {
	ok := result.NextResultSet(r.ctx)
	if ok {
		r.log.Debug("Result set has ", result.CurrentResultSet().RowCount(), " rows")
	} else {
		r.log.Debug("Result set requested, but not exists")
		return fmt.Errorf("result set requested, but not exists")
	}
	return nil
}

func (r *appRepositoryImpl) scanSingleRowResultSet(result result.Result, scanTo ...indexed.RequiredOrOptional) error {
	err := r.nextResultSet(result)
	if err != nil {
		return err
	}
	if result.CurrentResultSet().RowCount() != 1 {
		return fmt.Errorf("expected 1 row in result set, got %d", result.CurrentResultSet().RowCount())
	}
	if result.NextRow() {
		err := result.Scan(scanTo...)
		if err != nil {
			return err
		}
	}
	return nil
}

func embeddingToYDBList(embedding models.Embedding) types.Value {
	embeddingValues := make([]types.Value, len(embedding))
	for i := range embedding {
		embeddingValues[i] = types.FloatValue(embedding[i])
	}
	return types.ListValue(embeddingValues...)
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

	embeddingValues := make([]types.Value, len(queryEmbedding))
	for i := range queryEmbedding {
		embeddingValues[i] = types.FloatValue(queryEmbedding[i])
	}
	yqlEmbedding := embeddingToYDBList(queryEmbedding)
	result, err := r.execute(yql,
		table.ValueParam("$queryEmbedding", yqlEmbedding))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if err = r.nextResultSet(result); err != nil {
		return nil, err
	}

	searchResult := make([]api.SearchResultItem, 0)
	for result.NextRow() {
		var retrievedPageID uuid.UUID
		var title string
		var pageContent string
		var distance float32
		err = result.Scan(&retrievedPageID, &title, &pageContent, &distance)
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

	result, err := r.execute(yql, table.ValueParam("$pageID", types.UuidValue(pageID)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var retrievedPageID uuid.UUID
	var title string
	var content string
	if err = r.scanSingleRowResultSet(result, &retrievedPageID, &title, &content); err != nil {
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

	result, err := r.execute(yql,
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

	result, err := r.execute(yql,
		table.ValueParam("$pageID", types.UuidValue(paragraph.PageID)),
		table.ValueParam("$lineNumber", types.Int32Value(int32(paragraph.LineNumber))),
		table.ValueParam("$content", types.TextValue(paragraph.Content)),
		table.ValueParam("$embedding", embeddingToYDBList(paragraph.Embedding)),
		table.ValueParam("$anchorLineSlug", types.TextValue("")), // TODO
	)
	if err != nil {
		return err
	}
	err = result.Err()
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

	result, err := r.execute(yql)
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

	result, err := r.execute(yql, table.ValueParam("$login", types.TextValue(login)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var userID uuid.UUID
	var userLogin string
	var passwordHash string
	if err = r.scanSingleRowResultSet(result, &userID, &userLogin, &passwordHash); err != nil {
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

	result, err := r.execute(yql, table.ValueParam("$yWikiSlug", types.TextValue(yWikiSlug)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var pageID uuid.UUID
	var title string
	var content string
	if err = r.scanSingleRowResultSet(result, &pageID, &title, &content); err != nil {
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

	result1, err := r.execute(yql1, parameters...)
	if err != nil {
		return nil, err
	}
	defer result1.Close()
	if err = r.nextResultSet(result1); err != nil {
		return nil, err
	}
	if result1.CurrentResultSet().RowCount() == 1 {
		result1.NextRow()
		err := result1.Scan(pageID)
		if err != nil {
			return nil, err
		}
		return pageID, nil
	}

	result2, err := r.execute(yql2, parameters...)
	if err := r.scanSingleRowResultSet(result2, pageID); err != nil {
		return nil, err
	}
	r.log.Debug("Inserted page with id ", pageID)

	return pageID, nil
}

func (r *appRepositoryImpl) DeletePageBySlug(yWikiSlug string) error {
	yql := `
	DELETE FROM Page WHERE ywiki_slug=$yWikiSlug;
	`

	_, err := r.execute(yql, table.ValueParam("$yWikiSlug", types.TextValue(yWikiSlug)))
	return err
}

func (r *appRepositoryImpl) WriteIntegrationLogField(integrationID api.IntegrationID, logText string) error {
	yql := `INSERT INTO IntegrationLogField (integration_id, log_text, created_at)
	VALUES ($integrationID, $logText, CurrentUtcDate())`

	_, err := r.execute(yql,
		table.ValueParam("$integrationID", types.TextValue(string(integrationID))),
		table.ValueParam("$logText", types.TextValue(logText)),
	)
	return err
}
