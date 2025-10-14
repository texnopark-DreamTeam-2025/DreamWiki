package ywiki_client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/models"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/ywiki_client_gen"
)

type (
	yWikiClientImpl struct {
		client              *ywiki_client_gen.ClientWithResponses
		authorizationHeader string
		yandexCloudOrgID    string
	}

	YWikiClient interface {
		GetPage(ctx context.Context, pageSlug string) (*ywiki_client_gen.V1PageResponse, error)
	}
)

func NewYWikiClient(config *config.Config) (YWikiClient, error) {
	client, err := ywiki_client_gen.NewClientWithResponses("https://api.wiki.yandex.net")
	if err != nil {
		return nil, err
	}
	return &yWikiClientImpl{
		client:              client,
		authorizationHeader: fmt.Sprintf("OAuth %s", config.YWikiToken),
		yandexCloudOrgID:    config.YandexCloudOrgID,
	}, nil
}

func (c *yWikiClientImpl) GetPage(ctx context.Context, pageSlug string) (*ywiki_client_gen.V1PageResponse, error) {
	fields := "content"
	response, err := c.client.GetPageWithResponse(ctx, &ywiki_client_gen.GetPageParams{
		Slug:          pageSlug,
		Authorization: c.authorizationHeader,
		XCloudOrgId:   c.yandexCloudOrgID,
		Fields:        &fields,
	})
	if err != nil {
		return nil, err
	}

	switch response.HTTPResponse.StatusCode {
	case http.StatusNotFound:
		return nil, fmt.Errorf("YWiki GetPage: %w", models.ErrNotFound)
	case http.StatusOK:
		if response.JSON200 == nil {
			return nil, fmt.Errorf("200 response is nil")
		}
	default:
		return nil, fmt.Errorf("unexpected code: %d", response.HTTPResponse.StatusCode)
	}

	return response.JSON200, nil
}
