package api

import (
	"net/http"

	"mini-project/backend/api/auth"
	"mini-project/backend/api/tickets"
)

func RegisterRoutes(mux *http.ServeMux) {
	tickets.RegisterRoutes(mux)
	auth.RegisterRoutes(mux)
}
