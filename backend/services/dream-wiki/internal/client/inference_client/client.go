package inference_client

import (
	"context"
	"fmt"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/inference_client_gen"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/internals"
)

type (
	inferenceClientImpl struct {
		client *inference_client_gen.ClientWithResponses
	}

	InferenceClient interface {
		GenerateEmbedding(ctx context.Context, text string) (internals.Embedding, error)
		GenerateEmbeddings(ctx context.Context, texts []string) ([]internals.Embedding, error)
	}
)

func NewInferenceClient(config *config.Config) (InferenceClient, error) {
	client, err := inference_client_gen.NewClientWithResponses(config.InferenceAPIURL)
	if err != nil {
		return nil, err
	}
	return &inferenceClientImpl{client}, nil
}

func (c *inferenceClientImpl) GenerateEmbedding(ctx context.Context, text string) (internals.Embedding, error) {
	resp, err := c.client.GenerateEmbeddingWithResponse(ctx, inference_client_gen.GenerateEmbeddingJSONRequestBody{Texts: []string{text}})
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to generate embedding")
	}
	if len(resp.JSON200.Embeddings) != 1 {
		return nil, fmt.Errorf("no embeddings returned")
	}
	return internals.Embedding(resp.JSON200.Embeddings[0]), nil
}

func (c *inferenceClientImpl) GenerateEmbeddings(ctx context.Context, texts []string) ([]internals.Embedding, error) {
	resp, err := c.client.GenerateEmbeddingWithResponse(ctx, inference_client_gen.GenerateEmbeddingJSONRequestBody{Texts: texts})
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, fmt.Errorf("failed to generate embeddings")
	}

	embeddings := make([]internals.Embedding, len(resp.JSON200.Embeddings))
	for i, embedding := range resp.JSON200.Embeddings {
		embeddings[i] = internals.Embedding(embedding)
	}
	return embeddings, nil
}
