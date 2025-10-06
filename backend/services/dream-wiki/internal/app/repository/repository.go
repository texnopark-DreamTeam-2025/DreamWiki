package repository

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

type AppRepository struct {
	ctx     context.Context
	tx      table.TransactionActor
	success chan bool
	closed  int32
}

func StartTransaction(ctx context.Context, deps *deps.Deps) *AppRepository {
	success := make(chan bool)
	txRetriever := make(chan table.TransactionActor)

	action := func(ctx context.Context, tx table.TransactionActor) error {
		txRetriever <- tx
		shouldCommit := <-success
		close(success)
		if shouldCommit {
			return nil
		}
		return fmt.Errorf("transaction cancelled by user")
	}

	go deps.DB.DoTx(context.Background(), action)
	tx := <-txRetriever
	return &AppRepository{
		ctx: ctx, tx: tx, success: success}
}

func (r *AppRepository) Commit() {
	if n := atomic.AddInt32(&(r.closed), 1); n == 1 {
		r.success <- true
	}
}

func (r *AppRepository) Rollback() {
	if n := atomic.AddInt32(&(r.closed), 1); n == 1 {
		r.success <- false
	}
}

func (r *AppRepository) Search(ctx context.Context, query string) ([]models.SearchResult, error) {
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

func (r *AppRepository) GetDiagnosticInfo(ctx context.Context, pageID string) (*models.DiagnosticInfo, error) {
	fmt.Println("GetDiagnosticInfo")
	_, err := r.tx.Execute(r.ctx, "SELECT 1;", nil)
	if err != nil {
		return nil, fmt.Errorf("AppRepository.GetDiagnosticInfo: %w", err)
	}
	var v int
	return &models.DiagnosticInfo{
		PageID:    pageID,
		Content:   "Содержимое страницы " + fmt.Sprintf("%d", v),
		Title:     "Заголовок страницы " + pageID,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}, nil
}

func (r *AppRepository) IndexatePage(ctx context.Context, pageID string) error {
	return nil
}
