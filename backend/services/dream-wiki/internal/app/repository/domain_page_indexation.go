package repository

import (
	"fmt"
	"strings"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func (r *appRepositoryImpl) RemovePageIndexation(pageID api.PageID) error {
	yql := `
		DELETE FROM Paragraph WHERE page_id=$pageID;
	`

	result, err := r.ydbClient.InTX().Execute(yql,
		table.ValueParam("$pageID", types.UuidValue(pageID)),
	)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}

func (r *appRepositoryImpl) AddIndexedParagraph(paragraph internals.ParagraphWithEmbedding) error {
	yql := `
		INSERT INTO Paragraph (page_id, line_number, content, embedding, anchor_link_slug, paragraph_index, headers, is_header)
		VALUES (
			$pageID,
			$lineNumber,
			$content,
			Untag(Knn::ToBinaryStringFloat($embedding), "FloatVector"),
			$anchorLinkSlug,
			$paragraphIndex,
			$headers,
			$isHeader
		);
	`

	anchorLinkSlug := types.TextValue("")
	if paragraph.AnchorSlug != nil {
		anchorLinkSlug = types.TextValue(*paragraph.AnchorSlug)
	}

	if len(paragraph.Embedding) == 0 {
		return fmt.Errorf("embedding is empty for paragraph with page_id: %s, line_number: %d", paragraph.PageId, paragraph.LineNumber)
	}

	result, err := r.ydbClient.InTX().Execute(yql,
		table.ValueParam("$pageID", types.UuidValue(paragraph.PageId)),
		table.ValueParam("$lineNumber", types.Int64Value(int64(paragraph.LineNumber))),
		table.ValueParam("$content", types.TextValue(paragraph.Content)),
		table.ValueParam("$embedding", embeddingToYDBList(paragraph.Embedding)),
		table.ValueParam("$anchorLinkSlug", anchorLinkSlug),
		table.ValueParam("$paragraphIndex", types.Int64Value(int64(paragraph.ParagraphIndex))),
		table.ValueParam("$headers", types.UTF8Value(strings.Join(paragraph.Headers, "\n"))),
		table.ValueParam("$isHeader", types.BoolValue(paragraph.IsHeader)),
	)
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}
