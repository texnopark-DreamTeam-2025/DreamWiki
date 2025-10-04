package recovery

import (
	"fmt"
	"net/http"

	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/utils/logger"
)

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered", fmt.Errorf("%v", err), map[string]any{
					"method": r.Method,
					"path":   r.URL.Path,
				})
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
