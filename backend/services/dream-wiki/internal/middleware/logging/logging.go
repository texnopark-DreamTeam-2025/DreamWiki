package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"go.uber.org/zap"
)

func LoggingMiddleware(log *logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug("incoming request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("remote_addr", r.RemoteAddr),
			)
			next.ServeHTTP(w, r)
		})
	}
}
