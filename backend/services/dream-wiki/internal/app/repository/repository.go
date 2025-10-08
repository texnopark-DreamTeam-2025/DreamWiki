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
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

type AppRepositoryImpl struct {
	ctx     context.Context
	tx      table.TransactionActor
	log     *logger.Logger
	success chan bool
	closed  int32
}

func StartTransaction(ctx context.Context, deps *deps.Deps) *AppRepositoryImpl {
	success := make(chan bool)
	txRetriever := make(chan table.TransactionActor)

	action := func(ctx context.Context, tx table.TransactionActor) error {
		txRetriever <- tx
		shouldCommit := <-success
		close(success)
		if shouldCommit {
			deps.Logger.Info("transaction committed")
			return nil
		}
		deps.Logger.Info("transaction rolled back")
		return fmt.Errorf("transaction rolled back")
	}

	go deps.DB.DoTx(context.Background(), action)
	tx := <-txRetriever
	close(txRetriever)
	return &AppRepositoryImpl{
		ctx: ctx, tx: tx, success: success, log: deps.Logger}
}

func (r *AppRepositoryImpl) Commit() {
	if n := atomic.AddInt32(&(r.closed), 1); n == 1 {
		r.success <- true
	}
}

func (r *AppRepositoryImpl) Rollback() {
	if n := atomic.AddInt32(&(r.closed), 1); n == 1 {
		r.success <- false
	}
}

func (r *AppRepositoryImpl) Search(query string) ([]models.SearchResult, error) {
	results := []models.SearchResult{
		{
			Title:       "Результат поиска 1",
			Description: "Описание результата поиска по запросу: " + query,
			PageID:      "page-1",
		},
		{
			Title:       "Результат поиска 2",
			Description: "Еще один результат для: " + query,
			PageID:      "page-2",
		},
	}
	return results, nil
}

func (r *AppRepositoryImpl) RetrievePageByID(pageID string) (*api.Page, error) {
	yql := `
		SELECT CAST(page_id AS String), content FROM Page WHERE page_id=$pageID;
	`

	pageUUID, err := uuid.Parse(pageID)
	if err != nil {
		return nil, err
	}

	result, err := r.tx.Execute(r.ctx, yql, table.NewQueryParameters(
		table.ValueParam("$pageID", types.UuidValue(pageUUID)),
	))
	if err != nil {
		return nil, err
	}
	err = result.Err()
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var s1, s2 string
	result.NextResultSet(r.ctx)
	rowCount := result.CurrentResultSet().RowCount()
	if rowCount != 1 {
		r.log.Errorf("Invalid row count: %d, expected 1", rowCount)
		return nil, fmt.Errorf("invalid row count: %d, expected 1", rowCount)
	}
	for result.NextRow() {
		err = result.Scan(&s1, &s2)
		if err != nil {
			return nil, err
		}
	}

	return &api.Page{
		PageId:  s1,
		Content: s2,
		Title:   "Заголовок страницы " + pageID,
	}, nil
}

func (r *AppRepositoryImpl) RemovePageIndexation(pageID string) error {
	yql := `
		DELETE FROM Paragraph WHERE page_id=$pageID;
	`

	pageUUID, err := uuid.Parse(pageID)
	if err != nil {
		return err
	}

	result, err := r.tx.Execute(r.ctx, yql, table.NewQueryParameters(
		table.ValueParam("$pageID", types.UuidValue(pageUUID)),
	))
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

func (r *AppRepositoryImpl) AddIndexedParagraph(paragraph models.ParagraphWithEmbedding) error {
	yql := `
		INSERT INTO Paragraph (paragraph_id, page_id, line_number, content, embedding)
		VALUES ($paragraphID, $pageID, $lineNumber, $content, $embedding);
	`

	paragraphUUID, err := uuid.Parse(paragraph.ParagraphID)
	if err != nil {
		return err
	}

	pageUUID, err := uuid.Parse(paragraph.PageID)
	if err != nil {
		return err
	}

	result, err := r.tx.Execute(r.ctx, yql, table.NewQueryParameters(
		table.ValueParam("$paragraphID", types.UuidValue(paragraphUUID)),
		table.ValueParam("$pageID", types.UuidValue(pageUUID)),
		table.ValueParam("$lineNumber", types.Int32Value(int32(paragraph.LineNumber))),
		table.ValueParam("$content", types.TextValue(paragraph.Content)),
		table.ValueParam("$embedding", types.TextValue(paragraph.Embedding)),
	))
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
