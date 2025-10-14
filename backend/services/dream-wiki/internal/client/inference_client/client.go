package inference_client

import (
	"context"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/inference_client_gen"
)

type (
	inferenceClientImpl struct {
		client *inference_client_gen.ClientWithResponses
	}

	InferenceClient interface {
		GenerateEmbedding(ctx context.Context, text string) (models.Embedding, error)
	}
)

func NewInferenceClient(config *config.Config) (InferenceClient, error) {
	client, err := inference_client_gen.NewClientWithResponses(config.InferenceAPIURL)
	if err != nil {
		return nil, err
	}
	return &inferenceClientImpl{client}, nil
}

func (c *inferenceClientImpl) GenerateEmbedding(ctx context.Context, text string) (models.Embedding, error) {
	resp, err := c.client.GenerateEmbeddingWithResponse(ctx, inference_client_gen.GenerateEmbeddingJSONRequestBody{Text: text})
	if err != nil {
		return nil, err
	}
	return resp.JSON200.Embedding, nil
}
