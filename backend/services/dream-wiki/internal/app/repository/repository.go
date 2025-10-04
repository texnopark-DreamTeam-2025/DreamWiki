package repository

import (
	"context"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/ydbiface"
)

type AppRepository struct {
	db ydbiface.YDBIface
}

func NewAppRepository() *AppRepository {
	return &AppRepository{db: nil}
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
	return &models.DiagnosticInfo{
		PageID:    pageID,
		Content:   "Содержимое страницы " + pageID,
		Title:     "Заголовок страницы " + pageID,
		CreatedAt: time.Now().Add(-24 * time.Hour),
	}, nil
}

func (r *AppRepository) IndexatePage(ctx context.Context, pageID string) error {
	return nil
}
