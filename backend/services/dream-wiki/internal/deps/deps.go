package deps

import (
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	inference_client "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/inference"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
)

type Deps struct {
	DB              table.Client
	Config          *config.Config
	Logger          logger.Logger
	InferenceClient *inference_client.ClientWithResponses
}
