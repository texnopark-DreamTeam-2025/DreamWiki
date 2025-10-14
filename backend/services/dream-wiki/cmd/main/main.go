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
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/inference_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/client/ywiki_client"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	auth_middleware "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/auth"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/cors"
	middleware "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/logging"
	panic_middleware "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/panic"
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
	defer logger.Sync()

	logger.Info("configuration loaded successfully")
	config.LogConfig(appConfig, logger)

	logger.Debug("connecting to ydb")
	db, err := db.ConnectToYDB(appConfig, logger)
	if err != nil {
		logger.Fatalf("failed to connect to ydb: %v", err)
	}
	logger.Info("connected to ydb")
	defer db.Close(context.Background())

	dbTable := db.Table()

	// Initialize inference client
	inferenceClient, err := inference_client.NewInferenceClient(appConfig)
	if err != nil {
		logger.Fatalf("failed to initialize inference client: %v", err)
	}

	yWikiClient, err := ywiki_client.NewYWikiClient(appConfig)
	if err != nil {
		logger.Fatalf("failed to initialize ywiki client: %v", err)
	}

	deps := deps.Deps{
		DB:              dbTable,
		Config:          appConfig,
		Logger:          logger,
		InferenceClient: inferenceClient,
		YWikiClient:     yWikiClient,
	}

	appDelivery := delivery.NewAppDelivery(&deps)

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	apiRouter := router.PathPrefix("/api").Subrouter()

	strictHandler := api.NewStrictHandler(appDelivery, []api.StrictMiddlewareFunc{})
	api.HandlerWithOptions(strictHandler, api.GorillaServerOptions{
		BaseRouter: apiRouter,
	})

	router.Use(panic_middleware.PanicMiddleware(logger))
	router.Use(auth_middleware.AuthMiddleware(&deps))
	router.Use(middleware.LoggingMiddleware(logger))
	routerWithCORS := cors.CORSMiddleware(router)

	port := ":" + appConfig.ServerPort
	server := &http.Server{
		Addr:    port,
		Handler: routerWithCORS,
	}

	serverErr := make(chan error, 1)

	go func() {
		logger.Info("starting HTTP server",
			zap.String("address", server.Addr),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", zap.Error(err))
			serverErr <- err
		}
	}()

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

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown error", zap.Error(err))
			os.Exit(1)
		}

		logger.Info("server stopped")
	}
}
