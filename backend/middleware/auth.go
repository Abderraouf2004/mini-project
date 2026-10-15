package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"mini-project/backend/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type contextKey string

const UserIDKey contextKey = "userID"

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Fields(r.Header.Get("Authorization"))

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			http.Error(w, "authentication is not configured", http.StatusInternalServerError)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(
			tokenString,
			func(token *jwt.Token) (interface{}, error) {
				if token.Method != jwt.SigningMethodHS256 {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			},
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithExpirationRequired(),
		)

		if err != nil || token == nil || !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["id"].(string)
		if !ok || userID == "" {
			http.Error(w, "invalid user id in token", http.StatusUnauthorized)
			return
		}

		if database.Redis == nil {
			http.Error(w, "authentication service unavailable", http.StatusServiceUnavailable)
			return
		}

		hash := sha256.Sum256([]byte(tokenString))
		key := "blacklist:" + hex.EncodeToString(hash[:])

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		result, err := database.Redis.Get(ctx, key).Result()

		if err != nil && !errors.Is(err, redis.Nil) {
			http.Error(w, "authentication service unavailable", http.StatusServiceUnavailable)
			return
		}

		if result != "" {
			http.Error(w, "token has been revoked", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
