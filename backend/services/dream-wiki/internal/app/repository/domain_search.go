package repository

import (
	"strings"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func (r *appRepositoryImpl) SearchByEmbedding(query string, queryEmbedding internals.Embedding) ([]internals.SearchResultItem, error) {
	yql := `
		$targetEmbedding = Knn::ToBinaryStringFloat($queryEmbedding);

		SELECT
			par.page_id,
			page.ywiki_slug,
			par.paragraph_index,
			page.title,
			par.content,
			par.anchor_link_slug,
			par.headers,
			Unwrap(Knn::CosineDistance(Unwrap(par.embedding), $targetEmbedding)) As CosineDistance
		FROM Paragraph par
		JOIN Page page USING(page_id)
		ORDER BY Knn::CosineDistance(embedding, $targetEmbedding)
		LIMIT $limit;
	`

	yqlEmbedding := embeddingToYDBList(queryEmbedding)

	result, err := r.ydbClient.InTX().Execute(yql, table.ValueParam("$queryEmbedding", yqlEmbedding),
		table.ValueParam("$limit", types.Uint64Value(uint64(limit))),
	)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	searchResult := make([]internals.SearchResultItem, 0)
	for result.NextRow() {
		var retrievedPageID api.PageID
		var pageSlug string
		var paragraphIndex int64
		var title string
		var pageContent string
		var anchorLinkSlug *string
		var headers string
		var distance float32
		err = result.FetchRow(
			&retrievedPageID,
			&pageSlug,
			&paragraphIndex,
			&title,
			&pageContent,
			&anchorLinkSlug,
			&headers,
			&distance,
		)
		if err != nil {
			return nil, err
		}
		r.log.Debug("Distance is ", distance)
		searchResult = append(searchResult, internals.SearchResultItem{
			PageId:           retrievedPageID,
			PageSlug:         pageSlug,
			PageTitle:        title,
			ParagraphIndex:   int(paragraphIndex),
			ParagraphContent: pageContent,
			AnchorSlug:       anchorLinkSlug,
			Headers:          strings.Split(headers, "\n"),
		})
	}

	return searchResult, nil
}

func (r *appRepositoryImpl) SearchByEmbeddingWithContext(query string, queryEmbedding internals.Embedding, contextSize int) ([]internals.ParagraphWithContext, error) {
	initialResults, err := r.SearchByEmbedding(query, queryEmbedding)
	if err != nil {
		return nil, err
	}

	var allParagraphs []internals.ParagraphWithContext
	allRetrievedParagraphs := make([]internals.ParagraphWithContext, 0)

	for _, result := range initialResults {
		paragraphIndex := result.ParagraphIndex

		startIndex := int32(paragraphIndex - contextSize)
		endIndex := int32(paragraphIndex + contextSize)

		yql := `
			SELECT
				page_id,
				line_number,
				content,
				paragraph_index
			FROM Paragraph
			WHERE page_id = $page_id AND paragraph_index >= $start_index AND paragraph_index <= $end_index
			ORDER BY paragraph_index
		`

		result, err := r.ydbClient.InTX().Execute(yql,
			table.ValueParam("$page_id", types.UuidValue(result.PageId)),
			table.ValueParam("$start_index", types.Int32Value(startIndex)),
			table.ValueParam("$end_index", types.Int32Value(endIndex)))
		if err != nil {
			return nil, err
		}

		for result.NextRow() {
			var pageID api.PageID
			var lineNumber int64
			var content string
			var paragraphIndex int64

			err = result.FetchRow(&pageID, &lineNumber, &content, &paragraphIndex)
			if err != nil {
				result.Close()
				return nil, err
			}

			paragraphWithContext := internals.ParagraphWithContext{
				PageId:          pageID,
				ParagraphIndex:  int(paragraphIndex),
				StartLineNumber: int(lineNumber),
				EndLineNumber:   int(lineNumber),
				Content:         content,
			}

			allRetrievedParagraphs = append(allRetrievedParagraphs, paragraphWithContext)
		}
		result.Close()
	}

	pageParagraphs := groupParagraphsByPages(allRetrievedParagraphs)

	for _, paragraphs := range pageParagraphs {
		merged := mergeOverlappingParagraphs(paragraphs)
		allParagraphs = append(allParagraphs, merged...)
	}

	return allParagraphs, nil
}
