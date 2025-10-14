package panic

import (
	"net/http"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
	"go.uber.org/zap"
)

// PanicMiddleware catches panics and logs them, returning a 500 error to the client
func PanicMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("panic recovered",
						zap.Any("error", err),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("remote_addr", r.RemoteAddr),
					)

					// Return a 500 error to the client
					http.Error(w, `{"message": "Internal server error"}`, http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
