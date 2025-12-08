package dreamwikihttpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/app/delivery"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/components/component"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/auth"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/cors"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/logging"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/middleware/panic"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/pkg/api"
	"go.uber.org/zap"
)

type DreamWikiHTTPAPI struct {
	deps *deps.Deps
}

func NewDreamWikiHTTPAPI(deps *deps.Deps) *DreamWikiHTTPAPI {
	return &DreamWikiHTTPAPI{
		deps: deps,
	}
}

var _ component.Component = &DreamWikiHTTPAPI{}

func (d *DreamWikiHTTPAPI) Run(ctx context.Context) error {
	appDelivery := delivery.NewAppDelivery(d.deps)

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

	router.Use(panic.PanicMiddleware(d.deps.Logger))
	router.Use(auth.AuthMiddleware(d.deps))
	router.Use(logging.LoggingMiddleware(d.deps.Logger))
	routerWithCORS := cors.CORSMiddleware(router)

	port := ":" + d.deps.Config.ServerPort
	server := &http.Server{
		Addr:              port,
		ReadHeaderTimeout: 1 * time.Second,
		Handler:           routerWithCORS,
	}

	serverErr := make(chan error, 1)

	go func() {
		d.deps.Logger.Info("starting HTTP server. Address: ", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			d.deps.Logger.Error("server error", zap.Error(err))
			serverErr <- err
		}
	}()

	<-ctx.Done()
	d.deps.Logger.Info("server is shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		d.deps.Logger.Error("server shutdown error", zap.Error(err))
		return err
	}

	d.deps.Logger.Info("server stopped")
	return nil
}

func (d *DreamWikiHTTPAPI) Name() string {
	return "DreamWikiHTTPAPI"
}
