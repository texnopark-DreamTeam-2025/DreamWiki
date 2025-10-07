package repository

import (
	"cmp"
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

type AppRepositoryImpl struct {
	ctx     context.Context
	tx      table.TransactionActor
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
		deps.Logger.Info("transaction completed")
		if shouldCommit {
			return nil
		}
		return fmt.Errorf("transaction cancelled by user")
	}

	go deps.DB.DoTx(context.Background(), action)
	tx := <-txRetriever
	close(txRetriever)
	return &AppRepositoryImpl{
		ctx: ctx, tx: tx, success: success}
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

func (r *AppRepositoryImpl) RetrievePageByID(pageID string) (*models.Page, error) {
	fmt.Println("GetDiagnosticInfo")

	result, err := r.tx.Execute(r.ctx, "SELECT * FROM Page;", nil)
	err1 := result.Err()
	if err2 := cmp.Or(err, err1); err2 != nil {
		return nil, fmt.Errorf("AppRepository.GetDiagnosticInfo: %w", err2)
	}
	defer result.Close()

	var s1, s2 string
	for result.NextResultSet(r.ctx) {
		for result.NextRow() {
			result.Scan(&s1, &s2)
		}
	}

	return &models.Page{
		PageID:    s2,
		Content:   s1,
		Title:     "Заголовок страницы " + pageID,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}, nil
}
