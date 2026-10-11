package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-for-unit-tests")

	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(UserIDKey).(string)
		if !ok || userID != "user-123" {
			t.Errorf("expected user ID in context, got: %v", userID)
		}
		w.WriteHeader(http.StatusOK)
	})

	handler := Auth(protectedHandler)

	tests := []struct {
		name   string
		token  string
		status int
	}{
		{
			name:   "missing token",
			token:  "",
			status: http.StatusUnauthorized,
		},
		{
			name:   "invalid token",
			token:  "Bearer abc123",
			status: http.StatusUnauthorized,
		},
		{
			name:   "valid token",
			token:  "Bearer " + createTestToken(t, time.Now().Add(time.Minute)),
			status: http.StatusOK,
		},
		{
			name:   "expired token",
			token:  "Bearer " + createTestToken(t, time.Now().Add(-time.Minute)),
			status: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/tickets", nil)

			if tt.token != "" {
				req.Header.Set("Authorization", tt.token)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.status {
				t.Errorf("expected status %d, got %d", tt.status, rec.Code)
			}
		})
	}
}

func createTestToken(t *testing.T, expiresAt time.Time) string {
	t.Helper()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  "user-123",
		"exp": expiresAt.Unix(),
	})

	tokenString, err := token.SignedString([]byte("test-secret-key-for-unit-tests"))
	if err != nil {
		t.Fatalf("failed to create test token: %v", err)
	}

	return tokenString
}
