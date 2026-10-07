package api

import (
	"net/http"

	"mini-project/backend/api/tickets"
)

func RegisterRoutes(mux *http.ServeMux) {
	tickets.RegisterRoutes(mux)
}