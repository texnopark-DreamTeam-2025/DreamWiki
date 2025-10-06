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
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/cors"
	middleware "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/logging"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"go.uber.org/zap"

	"github.com/ydb-platform/ydb-go-sdk/v3"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("failed to load config: %v\n", err)
		os.Exit(1)
	}

	err = logger.Init(config.LogMode)
	if err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("configuration loaded successfully")
	logger.Debug("debug mode enabled",
		zap.String("log_mode", config.LogMode),
		zap.String("server_port", config.ServerPort),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Println("connecting to ydb")
	db, err := ydb.Open(ctx, config.YDBDSN,
		ydb.WithDialTimeout(1*time.Second),
	)
	fmt.Println("connected to ydb")
	if err != nil {
		logger.Fatalf("failed to connect ydb: ", err.Error())
	}
	logger.Info("connected to ydb")
	defer db.Close(ctx)

	appRepo := repository.NewAppRepository()
	appUsecase := usecase.NewAppUsecase(appRepo)
	appDelivery := delivery.NewAppDelivery(appUsecase)

	router := mux.NewRouter()

	router.Use(cors.CORS)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/v1/search", optionsHandler).Methods("OPTIONS")
	apiRouter.HandleFunc("/v1/diagnostic-info/get", optionsHandler).Methods("OPTIONS")
	apiRouter.HandleFunc("/v1/indexate-page", optionsHandler).Methods("OPTIONS")

	strictHandler := api.NewStrictHandler(appDelivery, []api.StrictMiddlewareFunc{})
	api.HandlerWithOptions(strictHandler, api.GorillaServerOptions{
		BaseRouter: apiRouter,
	})

	router.Use(middleware.LoggingMiddleware)

	port := ":" + config.ServerPort
	server := &http.Server{
		Addr:    port,
		Handler: router,
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

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown error", zap.Error(err))
			os.Exit(1)
		}

		logger.Info("server stopped")
	}
}

func optionsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
