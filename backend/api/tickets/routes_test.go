package tickets_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	ticketAPI "mini-project/backend/api/tickets"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateTicketUnauthorized(t *testing.T) {
	mux := http.NewServeMux()
	ticketAPI.RegisterRoutes(mux)

	req := httptest.NewRequest(
		http.MethodPost,
		"/tickets",
		nil,
	)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status %d, got %d; body: %s",
			http.StatusUnauthorized,
			rec.Code,
			rec.Body.String(),
		)
	}
}

func TestCreateTicketValidation(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "missing token",
			body:       `{"title":"Learn Go","description":"Test"}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "missing title without token",
			body:       `{"description":"Test"}`,
			wantStatus: http.StatusUnauthorized,
		},
	}

	mux := http.NewServeMux()
	ticketAPI.RegisterRoutes(mux)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/tickets",
				strings.NewReader(tt.body),
			)

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d; body: %s",
					tt.wantStatus,
					rec.Code,
					rec.Body.String(),
				)
			}
		})
	}
}

func createTestAccessToken(t *testing.T) string {
	t.Helper()

	secret := "test-secret-key-at-least-32-characters-long"
	t.Setenv("JWT_SECRET", secret)

	claims := jwt.MapClaims{
		"id":  "test-user-id",
		"typ": "access",
		"exp": time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	return tokenString
}

func TestCreateTicketValidationWithToken(t *testing.T) {
	token := createTestAccessToken(t)

	mux := http.NewServeMux()
	ticketAPI.RegisterRoutes(mux)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "missing title",
			body:       `{"description":"Test description"}`,
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "whitespace title",
			body:       `{"title":"   ","description":"Test description"}`,
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "invalid JSON",
			body:       `{"title":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unknown field",
			body:       `{"title":"Learn Go","description":"Test","admin":true}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/tickets",
				strings.NewReader(tt.body),
			)
			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"expected status %d, got %d; body: %s",
					tt.wantStatus,
					rec.Code,
					rec.Body.String(),
				)
			}
		})
	}
}
