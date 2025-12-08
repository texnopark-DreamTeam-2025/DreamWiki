package deps

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/github_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/inference_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/ycloud_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/ywiki_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/db_adapter"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
)

type Deps struct {
	YDBDriver       db_adapter.DBAdapter
	Config          *config.Config
	Logger          logger.Logger
	InferenceClient inference_client.InferenceClient
	YWikiClient     ywiki_client.YWikiClient
	GitHubClient    github_client.GitHubClient
	YCloudClient    ycloud_client.YCloudClient
}

type RepositoryDeps struct {
	TX   db_adapter.Transaction
	Deps *Deps
}
