package logging

import (
	"net/http"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("request started", map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
		})

		next.ServeHTTP(w, r)

		logger.Info("request completed", map[string]any{
			"method": r.Method,
			"path":   r.URL.Path,
		})
	})
}
