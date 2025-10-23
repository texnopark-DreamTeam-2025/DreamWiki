package ycloud_client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/ycloud_client_gen"
)

type (
	yCloudClientImpl struct {
		operationClient *ycloud_client_gen.ClientWithResponses
		llmClient       *ycloud_client_gen.ClientWithResponses
	}

	YCloudClient interface {
		StartAsyncLLMRequest(ctx context.Context, messages []ycloud_client_gen.Message) (*ycloud_client_gen.OperationID, error)
		GetLLMResponse(ctx context.Context, operationID ycloud_client_gen.OperationID) (*ycloud_client_gen.Operation, error)
	}
)

var (
	_ YCloudClient = &yCloudClientImpl{}
)

func NewYCloudClient(config *config.Config) (YCloudClient, error) {
	authorizationRequestEditor := func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.YandexCloudToken))
		return nil
	}

	llmClient, err := ycloud_client_gen.NewClientWithResponses(
		"https://llm.api.cloud.yandex.net",
		ycloud_client_gen.WithRequestEditorFn(authorizationRequestEditor),
	)
	if err != nil {
		return nil, err
	}

	operationClient, err := ycloud_client_gen.NewClientWithResponses(
		"https://operation.api.cloud.yandex.net",
		ycloud_client_gen.WithRequestEditorFn(authorizationRequestEditor),
	)
	if err != nil {
		return nil, err
	}

	return &yCloudClientImpl{
		llmClient:       llmClient,
		operationClient: operationClient,
	}, nil
}

func (c *yCloudClientImpl) StartAsyncLLMRequest(ctx context.Context, messages []ycloud_client_gen.Message) (*ycloud_client_gen.OperationID, error) {
	response, err := c.llmClient.PostFoundationModelsV1CompletionAsyncWithResponse(ctx, ycloud_client_gen.FoundationModelsV1CompletionAsyncRequest{
		CompletionOptions: ycloud_client_gen.CompletionOptions{
			Stream:      false,
			Temperature: 0.1,
			MaxTokens:   "2000",
			ReasoningOptions: ycloud_client_gen.ReasoningOptions{
				Mode: "DISABLED",
			},
		},
		Messages: messages,
		ModelUri: "gpt://b1gji9k43bb3qbc31oim/yandexgpt-lite/rc",
	})
	if err != nil {
		return nil, err
	}

	operationID := response.JSON200.Id
	return &operationID, nil
}

func (c *yCloudClientImpl) GetLLMResponse(ctx context.Context, operationID ycloud_client_gen.OperationID) (*ycloud_client_gen.Operation, error) {
	response, err := c.operationClient.GetOperationsOperationIdWithResponse(ctx, operationID)
	if err != nil {
		return nil, err
	}

	return response.JSON200, nil
}
