package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/texnopark-DreamTeam-2025/DreamWiki/internal/deps"
)

// UserKey is the key for storing user information in the context
type UserKey string

const (
	// UserContextKey is the key used to store user information in the context
	UserContextKey UserKey = "user"
)

// User represents a user in the system
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(deps *deps.Deps) func(http.Handler) http.Handler {
	log := deps.Logger

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Skip authentication for login and other unprotected endpoints
			if strings.HasSuffix(r.URL.Path, "/v1/login") || strings.HasSuffix(r.URL.Path, "/health") {
				next.ServeHTTP(w, r)
				return
			}

			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Debug("Authorization header is missing")
				http.Error(w, `{"message": "Authorization header is required"}`, http.StatusUnauthorized)
				return
			}

			// Check if the header starts with "Bearer "
			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Debug("Authorization header does not start with Bearer")
				http.Error(w, `{"message": "Authorization header must start with Bearer"}`, http.StatusUnauthorized)
				return
			}

			// Extract the token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Parse and validate the token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Validate the signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(deps.Config.JWTSecretKey), nil
			})

			if err != nil {
				log.Debug("Error parsing JWT token", "error", err.Error())
				http.Error(w, `{"message": "Invalid token"}`, http.StatusUnauthorized)
				return
			}

			// Check if the token is valid
			if !token.Valid {
				log.Debug("JWT token is invalid")
				http.Error(w, `{"message": "Invalid token"}`, http.StatusUnauthorized)
				return
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Debug("Failed to extract claims from JWT token")
				http.Error(w, `{"message": "Invalid token claims"}`, http.StatusUnauthorized)
				return
			}

			// Create a user object from claims
			user := User{
				ID:       claims["id"].(string),
				Username: claims["username"].(string),
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), UserContextKey, user)

			// Call the next handler
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves the user from the request context
func GetUserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(UserContextKey).(User)
	if !ok {
		return nil, false
	}
	return &user, true
}
