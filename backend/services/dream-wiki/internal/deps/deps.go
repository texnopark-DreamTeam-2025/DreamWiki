package deps

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/inference_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/ywiki_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

type Deps struct {
	DB              table.Client
	Config          *config.Config
	Logger          logger.Logger
	InferenceClient inference_client.InferenceClient
	YWikiClient     ywiki_client.YWikiClient
}
