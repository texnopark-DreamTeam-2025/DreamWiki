package repository

import (
	"sort"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func (r *appRepositoryImpl) SearchByEmbedding(query string, queryEmbedding models.Embedding) ([]internals.SearchResultItem, error) {
	yql := `
		$K = 20;
		$targetEmbedding = Knn::ToBinaryStringFloat($queryEmbedding);

		SELECT
			par.page_id,
			page.slug,
			par.paragraph_index,
			page.title,
			par.content,
			par.anchor_link_slug,
			par.headers,
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

	searchResult := make([]internals.SearchResultItem, 0)
	for result.NextRow() {
		var retrievedPageID api.PageID
		var pageSlug string
		var paragraphIndex int
		var title string
		var pageContent string
		var anchorLinkSlug *string
		var headers []string
		var distance float32
		err = result.FetchRow(&retrievedPageID, &pageSlug, &paragraphIndex, &title, &pageContent, &anchorLinkSlug, &headers, &distance)
		if err != nil {
			return nil, err
		}
		r.log.Debug("Distance is ", distance)
		searchResult = append(searchResult, internals.SearchResultItem{
			PageId:           retrievedPageID,
			PageSlug:         pageSlug,
			PageTitle:        title,
			ParagraphIndex:   paragraphIndex,
			ParagraphContent: pageContent,
			AnchorSlug:       anchorLinkSlug,
			Headers:          headers,
		})
	}

	return searchResult, nil
}

func (r *appRepositoryImpl) SearchByEmbeddingWithContext(query string, queryEmbedding models.Embedding, contextSize int) ([]internals.ParagraphWithContext, error) {
	initialResults, err := r.SearchByEmbedding(query, queryEmbedding)
	if err != nil {
		return nil, err
	}

	var allParagraphs []internals.ParagraphWithContext
	pageParagraphs := make(map[string][]internals.ParagraphWithContext)

	for _, result := range initialResults {
		paragraphIndex := result.ParagraphIndex

		startIndex := int32(paragraphIndex - contextSize)
		if startIndex < 0 {
			startIndex = 0
		}
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

			pageKey := string(pageID[:])
			pageParagraphs[pageKey] = append(pageParagraphs[pageKey], paragraphWithContext)
		}
		result.Close()
	}

	for _, paragraphs := range pageParagraphs {
		merged := r.mergeOverlappingParagraphs(paragraphs)
		allParagraphs = append(allParagraphs, merged...)
	}

	return allParagraphs, nil
}

func (r *appRepositoryImpl) mergeOverlappingParagraphs(paragraphs []internals.ParagraphWithContext) []internals.ParagraphWithContext {
	if len(paragraphs) == 0 {
		return paragraphs
	}

	sort.Slice(paragraphs, func(i, j int) bool {
		return paragraphs[i].StartLineNumber < paragraphs[j].StartLineNumber
	})

	var merged []internals.ParagraphWithContext
	current := paragraphs[0]

	for i := 1; i < len(paragraphs); i++ {
		next := paragraphs[i]

		if next.StartLineNumber <= current.EndLineNumber+1 {
			current.Content += "\n" + next.Content
			if next.EndLineNumber > current.EndLineNumber {
				current.EndLineNumber = next.EndLineNumber
			}
		} else {
			merged = append(merged, current)
			current = next
		}
	}

	merged = append(merged, current)

	return merged
}
