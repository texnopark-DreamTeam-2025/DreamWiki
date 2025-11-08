package repository

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func embeddingToYDBList(embedding internals.Embedding) types.Value {
	embeddingValues := make([]types.Value, len(embedding))
	for i := range embedding {
		embeddingValues[i] = types.FloatValue(embedding[i])
	}
	return types.ListValue(embeddingValues...)
}
