package ywiki_client

import (
	"context"
	"fmt"

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
		authorizationHeader: fmt.Sprintf("Oauth %s", config.YWikiToken),
		yandexCloudOrgID:    config.YandexCloudOrgID,
	}, nil
}

func (c *yWikiClientImpl) GetPage(ctx context.Context, pageSlug string) (*ywiki_client_gen.V1PageResponse, error) {
	response, err := c.client.GetPageWithResponse(ctx, &ywiki_client_gen.GetPageParams{
		Slug:          pageSlug,
		Authorization: c.authorizationHeader,
		XCloudOrgId:   c.yandexCloudOrgID,
	})
	if err != nil {
		return nil, err
	}
	return response.JSON200, nil
}
