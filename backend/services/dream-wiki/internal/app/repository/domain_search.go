package repository

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func (r *appRepositoryImpl) SearchByEmbedding(query string, queryEmbedding internals.Embedding, limit int) ([]internals.SearchResultItem, error) {
	yql := `
	$targetEmbedding = Knn::ToBinaryStringFloat($queryEmbedding);

	SELECT
		par.page_id,
		page.ywiki_slug,
		par.paragraph_index,
		page.title,
		par.content,
		par.headers,
		par.anchor_link_slug,
		par.line_number,
		Unwrap(Knn::CosineDistance(Unwrap(par.embedding), $targetEmbedding)) As CosineDistance
	FROM Paragraph par
	JOIN Page page USING(page_id)
	WHERE NOT is_header
	ORDER BY Knn::CosineDistance(embedding, $targetEmbedding)
	LIMIT $limit;
`

	yqlEmbedding := embeddingToYDBList(queryEmbedding)

	result, err := r.tx.InTX().Execute(yql, table.ValueParam("$queryEmbedding", yqlEmbedding),
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
		var headers string
		var anchorLinkSlug string
		var lineNumber int64
		var distance float32
		err = result.FetchRow(
			&retrievedPageID,
			&pageSlug,
			&paragraphIndex,
			&title,
			&pageContent,
			&headers,
			&anchorLinkSlug,
			&lineNumber,
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
			ParagraphContent: pageContent,
			AnchorSlug:       &anchorLinkSlug,
			Headers:          strings.Split(headers, "\n"),
			ParagraphIndex:   int(paragraphIndex),
			LineIndex:        int(lineNumber),
		})
	}

	return searchResult, nil
}

func (r *appRepositoryImpl) SearchByEmbeddingWithContext(query string, queryEmbedding internals.Embedding, contextSize int) ([]internals.ParagraphWithContext, error) {
	initialResults, err := r.SearchByEmbedding(query, queryEmbedding, 3)
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
			WHERE page_id = $page_id
				AND paragraph_index >= $start_index
				AND paragraph_index <= $end_index
				AND NOT is_header
			ORDER BY paragraph_index
		`

		result, err := r.tx.InTX().Execute(yql,
			table.ValueParam("$page_id", types.UuidValue(result.PageId)),
			table.ValueParam("$start_index", types.Int32Value(startIndex)),
			table.ValueParam("$end_index", types.Int32Value(endIndex)))
		if err != nil {
			return nil, err
		}
		defer result.Close()

		for result.NextRow() {
			var pageID api.PageID
			var lineNumber int64
			var content string
			var paragraphIndex int64

			err = result.FetchRow(&pageID, &lineNumber, &content, &paragraphIndex)
			if err != nil {
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
	}

	pageParagraphs := groupParagraphsByPages(allRetrievedParagraphs)

	for _, paragraphs := range pageParagraphs {
		merged := mergeOverlappingParagraphs(paragraphs)
		allParagraphs = append(allParagraphs, merged...)
	}

	return allParagraphs, nil
}

func (r *appRepositoryImpl) SearchByTerms(terms []string, limit int) ([]internals.SearchResultItem, error) {
	if len(terms) == 0 {
		return []internals.SearchResultItem{}, nil
	}

	totalDocsQuery := `SELECT Count(DISTINCT page_id) AS total FROM Paragraph`
	totalDocsResult, err := r.tx.InTX().Execute(totalDocsQuery)
	if err != nil {
		return nil, err
	}
	defer totalDocsResult.Close()

	var totalDocs uint64
	err = totalDocsResult.FetchExactlyOne(&totalDocs)
	if err != nil {
		return nil, err
	}

	termList := make([]types.Value, len(terms))
	for i, term := range terms {
		termList[i] = types.TextValue(term)
	}
	yqlTerms := types.ListValue(termList...)

	docFreqQuery := `
		SELECT
			term,
			Count(DISTINCT page_id) AS doc_freq
		FROM Term
		WHERE term IN $terms
		GROUP BY term
	`

	docFreqResult, err := r.tx.InTX().Execute(docFreqQuery, table.ValueParam("$terms", yqlTerms))
	if err != nil {
		return nil, err
	}
	defer docFreqResult.Close()

	termDocFreq := make(map[string]int64)
	for docFreqResult.NextRow() {
		var term string
		var docFreq uint64
		err = docFreqResult.FetchRow(&term, &docFreq)
		if err != nil {
			return nil, err
		}
		termDocFreq[term] = int64(docFreq)
	}

	paragraphsQuery := `
		SELECT
			t.page_id,
			t.paragraph_index,
			t.term,
			t.times_in,
			p.content,
			p.headers,
			page.ywiki_slug,
			page.title
		FROM Term t
		JOIN Paragraph p ON t.page_id = p.page_id AND t.paragraph_index = p.paragraph_index
		JOIN Page page ON t.page_id = page.page_id
		WHERE t.term IN $terms
	`

	paragraphsResult, err := r.tx.InTX().Execute(paragraphsQuery, table.ValueParam("$terms", yqlTerms))
	if err != nil {
		return nil, err
	}
	defer paragraphsResult.Close()

	type paragraphData struct {
		data            internals.SearchResultItem
		termFrequencies map[string]int64
	}

	paragraphs := make(map[string]*paragraphData)
	for paragraphsResult.NextRow() {
		var pageID api.PageID
		var paragraphIndex int64
		var term string
		var timesIn int64
		var content string
		var headers string
		var pageSlug string
		var title string

		err = paragraphsResult.FetchRow(&pageID, &paragraphIndex, &term, &timesIn, &content, &headers, &pageSlug, &title)
		if err != nil {
			return nil, err
		}

		key := fmt.Sprintf("%s_%d", pageID.String(), paragraphIndex)
		if _, exists := paragraphs[key]; !exists {
			paragraphs[key] = &paragraphData{
				data: internals.SearchResultItem{
					PageId:           pageID,
					ParagraphIndex:   int(paragraphIndex),
					ParagraphContent: content,
					Headers:          strings.Split(headers, "\n"),
					PageSlug:         pageSlug,
					PageTitle:        title,
				},
				termFrequencies: make(map[string]int64),
			}
		}
		paragraphs[key].termFrequencies[term] = timesIn
	}

	type scoredParagraph struct {
		data  internals.SearchResultItem
		score float64
	}

	var scoredParagraphs []scoredParagraph
	for _, paragraph := range paragraphs {
		score := 0.0
		totalTermsInParagraph := 0
		for _, freq := range paragraph.termFrequencies {
			totalTermsInParagraph += int(freq)
		}

		for term, freq := range paragraph.termFrequencies {
			if docFreq, exists := termDocFreq[term]; exists {
				tf := float64(freq) / float64(totalTermsInParagraph)
				idf := math.Log(float64(totalDocs) / float64(docFreq))
				score += tf * idf
			}
		}

		scoredParagraphs = append(scoredParagraphs, scoredParagraph{
			data:  paragraph.data,
			score: score,
		})
	}

	sort.Slice(scoredParagraphs, func(i, j int) bool {
		return scoredParagraphs[i].score > scoredParagraphs[j].score
	})

	resultLimit := min(len(scoredParagraphs), limit)

	searchResult := make([]internals.SearchResultItem, 0, resultLimit)
	for i := range resultLimit {
		paragraph := scoredParagraphs[i]
		searchResult = append(searchResult, paragraph.data)
	}

	return searchResult, nil
}

func (r *appRepositoryImpl) SearchByTermsWithContext(terms []string, contextSize int) ([]internals.ParagraphWithContext, error) {
	initialResults, err := r.SearchByTerms(terms, 3)
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

		result, err := r.tx.InTX().Execute(yql,
			table.ValueParam("$page_id", types.UuidValue(result.PageId)),
			table.ValueParam("$start_index", types.Int32Value(startIndex)),
			table.ValueParam("$end_index", types.Int32Value(endIndex)))
		if err != nil {
			return nil, err
		}
		defer result.Close()

		for result.NextRow() {
			var pageID api.PageID
			var lineNumber int64
			var content string
			var paragraphIndex int64

			err = result.FetchRow(&pageID, &lineNumber, &content, &paragraphIndex)
			if err != nil {
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
	}

	pageParagraphs := groupParagraphsByPages(allRetrievedParagraphs)

	for _, paragraphs := range pageParagraphs {
		merged := mergeOverlappingParagraphs(paragraphs)
		allParagraphs = append(allParagraphs, merged...)
	}

	return allParagraphs, nil
}
