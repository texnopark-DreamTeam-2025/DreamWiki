package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/delivery"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/repository"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/usecase"
	middleware "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/logging"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"go.uber.org/zap"
)

func main() {
	// cfg, err := config.LoadConfig()
	// if err != nil {
	// 	fmt.Printf("error initializing config: %v\n", err)
	// 	os.Exit(1)
	// }

	err := logger.Init("dev")
	if err != nil {
		fmt.Printf("failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("configuration loaded successfully")
	logger.Debug("debug mode enabled",
		zap.String("log_mode", "dev"),
		zap.String("server_port", "8080"),
	)

	// driver, err := db.CreateYDBDriver(cfg)
	// if err != nil {
	// 	logger.Fatalf("failed to connect to DB %w", zap.Error(err))
	// }
	// defer driver.Close(context.Background())

	appRepo := repository.NewAppRepository()
	appUsecase := usecase.NewAppUsecase(appRepo)
	appDelivery := delivery.NewAppDelivery(appUsecase)

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

	router.Use(middleware.LoggingMiddleware)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173", "http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	port := ":" + "8080"
	server := &http.Server{
		Addr:    port,
		Handler: cors(router),
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
