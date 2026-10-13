package tickets_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	ticketAPI "mini-project/backend/api/tickets"
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
