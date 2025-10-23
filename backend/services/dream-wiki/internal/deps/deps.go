package deps

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/github_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/inference_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/ycloud_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/ywiki_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

type Deps struct {
	YDBDriver       *ydb.Driver
	Config          *config.Config
	Logger          logger.Logger
	InferenceClient inference_client.InferenceClient
	YWikiClient     ywiki_client.YWikiClient
	GitHubClient    github_client.GitHubClient
	YCloudClient    ycloud_client.YCloudClient
}
