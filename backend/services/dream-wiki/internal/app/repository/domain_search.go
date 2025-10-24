package repository

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

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
		var retrievedPageID api.PageID
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
