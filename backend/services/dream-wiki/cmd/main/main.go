package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/delivery"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/github_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/inference_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/ycloud_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/ywiki_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/auth"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/cors"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/logging"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/panic"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/topic_reader"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/db"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"go.uber.org/zap"
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

	deps := deps.Deps{
		YDBDriver:       ydbDriver,
		Config:          appConfig,
		Logger:          logger,
		InferenceClient: inferenceClient,
		YWikiClient:     yWikiClient,
		GitHubClient:    gitHubClient,
		YCloudClient:    yCloudClient,
	}

	appDelivery := delivery.NewAppDelivery(&deps)

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods("GET")

	apiRouter := router.PathPrefix("/api").Subrouter()

	strictHandler := api.NewStrictHandler(appDelivery, []api.StrictMiddlewareFunc{})
	api.HandlerWithOptions(strictHandler, api.GorillaServerOptions{
		BaseRouter: apiRouter,
	})

	router.Use(panic.PanicMiddleware(logger))
	router.Use(auth.AuthMiddleware(&deps))
	router.Use(logging.LoggingMiddleware(logger))
	routerWithCORS := cors.CORSMiddleware(router)

	port := ":" + appConfig.ServerPort
	server := &http.Server{
		Addr:              port,
		ReadHeaderTimeout: 1 * time.Second,
		Handler:           routerWithCORS,
	}

	serverErr := make(chan error, 1)

	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	go func() {
		logger.Info("starting HTTP server",
			zap.String("address", server.Addr),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", zap.Error(err))
			serverErr <- err
		}
	}()

	rd, err := topic_reader.NewTopicReader(appCtx, &deps)
	if err != nil {
		logger.Error("failed to create topic reader", zap.Error(err))
		os.Exit(1)
	}
	rd.ReadMessages()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		logger.Error("failed to start server", zap.Error(err))
		os.Exit(1)
	case sig := <-quit:
		logger.Info("server is shutting down",
			zap.String("signal", sig.String()),
		)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("server shutdown error", zap.Error(err))
			cancel()
			os.Exit(1)
		}

		appCancel()

		if err := rd.Close(shutdownCtx); err != nil {
			logger.Error("topic readers shutdown error", zap.Error(err))
		}

		logger.Info("server stopped")
	}
}
