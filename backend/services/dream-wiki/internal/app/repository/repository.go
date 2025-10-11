package repository

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/local_model"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

type AppRepositoryImpl struct {
	ctx     context.Context
	tx      table.TransactionActor
	log     logger.Logger
	success chan bool
	closed  int32
}

func StartTransaction(ctx context.Context, deps *deps.Deps) *AppRepositoryImpl {
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

func (r *AppRepositoryImpl) execute(yql string, opts ...table.ParameterOption) (result result.Result, err error) {
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

func (r *AppRepositoryImpl) SearchByEmbedding(query string, queryEmbedding local_model.Embedding) ([]models.ParagraphWithEmbedding, error) {
	yql := `
		$K = 20;
		$TargetEmbedding = Knn::ToBinaryStringFloat($queryEmbedding);

		SELECT
			page_id,
			line_number,
			Knn::CosineDistance(embedding, $TargetEmbedding) As CosineDistance
		FROM Paragraph
		ORDER BY Knn::CosineDistance(embedding, $TargetEmbedding)
		LIMIT $K;
	`

	embeddingValues := make([]types.Value, len(queryEmbedding))
	for i := range queryEmbedding {
		embeddingValues[i] = types.FloatValue(queryEmbedding[i])
	}
	yqlEmbedding := types.ListValue(embeddingValues...)

	result, err := r.tx.Execute(r.ctx, yql, table.NewQueryParameters(
		table.ValueParam("$queryEmbedding", yqlEmbedding),
	))
	// todo obrabotka oshibok

	defer result.Close()

	var retrievedPageID uuid.UUID
	var pageContent string
	if ok := result.NextResultSet(r.ctx); !ok {
		return nil, fmt.Errorf("no result set")
	}
	rowCount := result.CurrentResultSet().RowCount()
	if rowCount != 1 {
		r.log.Errorf("Invalid row count: %d, expected 1", rowCount)
		return nil, fmt.Errorf("invalid row count: %d, expected 1", rowCount)
	}
	for result.NextRow() {
		err = result.Scan(&retrievedPageID, &pageContent)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil // TODO
}

func (r *AppRepositoryImpl) RetrievePageByID(pageID uuid.UUID) (*api.Page, error) {
	yql := `
		SELECT CAST(page_id AS String), content FROM Page WHERE page_id=$pageID;
	`

	result, err := r.execute(yql, table.ValueParam("$pageID", types.UuidValue(pageID)))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if ok := result.NextResultSet(r.ctx); !ok {
		return nil, fmt.Errorf("no result set")
	}
	rowCount := result.CurrentResultSet().RowCount()
	if rowCount != 1 {
		r.log.Errorf("Invalid row count: %d, expected 1", rowCount)
		return nil, fmt.Errorf("invalid row count: %d, expected 1", rowCount)
	}

	var retrievedPageID uuid.UUID
	var pageContent string
	for result.NextRow() {
		err = result.Scan(&retrievedPageID, &pageContent)
		if err != nil {
			return nil, err
		}
	}

	return &api.Page{
		PageId:  retrievedPageID,
		Content: pageContent,
		Title:   "Заголовок страницы",
	}, nil
}

func (r *AppRepositoryImpl) RemovePageIndexation(pageID uuid.UUID) error {
	yql := `
		DELETE FROM Paragraph WHERE page_id=$pageID;
	`

	result, err := r.tx.Execute(r.ctx, yql, table.NewQueryParameters(
		table.ValueParam("$pageID", types.UuidValue(pageID)),
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
		INSERT INTO Paragraph (page_id, line_number, content, embedding)
		VALUES ($pageID, $lineNumber, $content, $embedding);
	`

	result, err := r.tx.Execute(r.ctx, yql, table.NewQueryParameters(
		table.ValueParam("$pageID", types.UuidValue(paragraph.PageID)),
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

func (r *AppRepositoryImpl) DeletePage() error {
	yql := `
		DELETE FROM Paragraph;
		DELETE FROM Page;
	`

	result, err := r.tx.Execute(r.ctx, yql, table.NewQueryParameters())
	if err != nil {
		return err
	}
	err = result.Err()
	if err != nil {
		return err
	}
	defer result.Close()

	r.log.Info("All pages and paragraphs deleted successfully")
	return nil
}
