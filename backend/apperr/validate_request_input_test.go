package apperr

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateRequestBody(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "valid JSON",
			body:       `{"title":"Learn Go"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing required field",
			body:       `{}`,
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "invalid JSON",
			body:       `{"title":`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "unknown field",
			body:       `{"title":"Learn Go","admin":true}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "multiple JSON values",
			body:       `{"title":"Learn Go"} {"title":"Another"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	handler := ValidateRequestBody[testDTO](
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"success"}`))
		}),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodPost,
				"/test",
				strings.NewReader(tt.body),
			)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf(
					"got status %d, want %d; body: %s",
					rec.Code,
					tt.wantStatus,
					rec.Body.String(),
				)
			}

			if tt.wantStatus != http.StatusOK {
				var response map[string]any
				if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
					t.Fatalf("expected JSON error response, got: %s", rec.Body.String())
				}

				if response["code"] == nil {
					t.Fatalf("expected error code in response, got: %s", rec.Body.String())
				}
			}
		})
	}
}
