package local_model

import (
	"context"
	"fmt"
	"net/http"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	inference_client "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/inference"
)

type (
	Model struct {
		deps *deps.Deps
	}

	Embedding []float32
)

func NewModel(deps *deps.Deps) (*Model, error) {
	return &Model{
		deps: deps,
	}, nil
}

func (m *Model) TextToEmbedding(ctx context.Context, text string) (*Embedding, error) {
	return m.getTextEmbeddingFromService(ctx, text)
}

func (m *Model) getTextEmbeddingFromService(ctx context.Context, text string) (*Embedding, error) {
	requestBody := inference_client.GenerateEmbeddingJSONRequestBody{
		Text: text,
	}

	response, err := m.deps.InferenceClient.GenerateEmbeddingWithResponse(ctx, requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to call inference service: %w", err)
	}

	if response.StatusCode() != http.StatusOK {
		if response.JSON422 != nil {
			return nil, fmt.Errorf("validation error: %s", response.JSON422.Error)
		}
		if response.JSON500 != nil {
			return nil, fmt.Errorf("server error: %s", response.JSON500.Error)
		}
		return nil, fmt.Errorf("unexpected status code: %d", response.StatusCode())
	}

	// Convert the response to our embedding type
	if response.JSON200 == nil {
		return nil, fmt.Errorf("empty response from inference service")
	}

	embedding := make(Embedding, len(response.JSON200.Embedding))
	copy(embedding, response.JSON200.Embedding)

	return &embedding, nil
}
