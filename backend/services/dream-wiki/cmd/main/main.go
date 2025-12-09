package main

import (
	"context"
	"fmt"
	"os"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/github_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/inference_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/ycloud_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/ywiki_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/components/component"
	dreamwikihttpapi "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/components/dreamwiki_http_api"
	dreamwikitaskactionresultstopicreader "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/components/dreamwiki_task_action_results_topic_reader"
	dreamwikitaskactionstopicreader "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/components/dreamwiki_task_actions_topic_reader"
	staletaskfailer "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/components/stale_task_failer"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/db_adapter"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/db"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
)

func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	logger, err := logger.New(appConfig.LogMode)
	if err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		_ = logger.Sync()
	}()

	logger.Info("configuration loaded successfully")
	config.LogConfig(appConfig, logger)

	logger.Debug("connecting to ydb")
	ydbDriver, err := db.ConnectToYDB(appConfig, logger)
	if err != nil {
		logger.Fatalf("failed to connect to ydb: %v", err)
	}
	defer func() {
		_ = ydbDriver.Close(context.Background())
	}()

	inferenceClient, err := inference_client.NewInferenceClient(appConfig)
	if err != nil {
		logger.Fatalf("failed to initialize inference client: %v", err)
	}

	yWikiClient, err := ywiki_client.NewYWikiClient(appConfig)
	if err != nil {
		logger.Fatalf("failed to initialize ywiki client: %v", err)
	}

	yCloudClient, err := ycloud_client.NewYCloudClient(appConfig)
	if err != nil {
		logger.Fatalf("failed to initialize ycloud client: %v", err)
	}

	gitHubClient, err := github_client.NewGitHubClient(appConfig)
	if err != nil {
		logger.Fatalf("failed to initialize github client: %v", err)
	}
	dbAdapter := db_adapter.NewDBAdapter(appConfig, logger)
	defer dbAdapter.Close()

	deps := deps.Deps{
		YDBDriver:       dbAdapter,
		Config:          appConfig,
		Logger:          logger,
		InferenceClient: inferenceClient,
		YWikiClient:     yWikiClient,
		GitHubClient:    gitHubClient,
		YCloudClient:    yCloudClient,
	}

	taskActionsTopicReader := dreamwikitaskactionstopicreader.NewDreamWikiTaskActionsTopicReader(&deps)
	taskActionResultsTopicReader := dreamwikitaskactionresultstopicreader.NewDreamWikiTaskActionResultsTopicReader(&deps)
	httpAPI := dreamwikihttpapi.NewDreamWikiHTTPAPI(&deps)
	staleTaskFailer := staletaskfailer.NewStaleTaskFailer(&deps)

	err = component.RunComponents(
		taskActionsTopicReader,
		taskActionResultsTopicReader,
		httpAPI,
		staleTaskFailer,
	)
	if err != nil {
		logger.Error("one or more components shutted down with error: %v", err)
	}
}
