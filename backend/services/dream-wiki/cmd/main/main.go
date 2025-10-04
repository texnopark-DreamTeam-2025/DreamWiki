package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/config"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/logging"
	recovery "github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/panic"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
)

func main() {
	logger.Init()
	defer logger.Close()

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.Fatal("failed to load config", err, nil)
	}

	logger.Info("configuration loaded successfully", nil)

	// layers

	handlerChain := func(h http.Handler) http.Handler {
		return recovery.RecoveryMiddleware(logging.LoggingMiddleware(h))
	}

	mux := http.NewServeMux()

	// handlers

	port := ":" + cfg.ServerPort
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	serverErr := make(chan error, 1)

	go func() {
		logger.Info("starting HTTP server", map[string]any{
			"address": "http://localhost" + port,
		})
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", err, nil)
			serverErr <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		logger.Error("failed to start server", err, nil)
		os.Exit(1)
	case sig := <-quit:
		logger.Info("server is shutting down", map[string]any{
			"signal": sig.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Error("server shutdown error", err, nil)
			os.Exit(1)
		}

		logger.Info("server stopped", nil)
	}

}
