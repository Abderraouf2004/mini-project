package api

import (
	"context"
	"net/http"
	"time"

	"mini-project/backend/api/auth"
	"mini-project/backend/api/tickets"

	"gorm.io/gorm"
)

func RegisterRoutes(mux *http.ServeMux, db *gorm.DB) {
	tickets.RegisterRoutes(mux)
	auth.RegisterRoutes(mux)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err != nil {
			http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			http.Error(w, "Database unavailable", http.StatusServiceUnavailable)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})
}
