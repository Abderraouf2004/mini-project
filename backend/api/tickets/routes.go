package tickets

import (
	"net/http"

	appErrors "mini-project/backend/apperr"
	coreTickets "mini-project/backend/core/tickets"
	"mini-project/backend/middleware"
)

func RegisterRoutes(mux *http.ServeMux) {
	controller := coreTickets.NewController()

	mux.Handle(
		"POST /tickets",
		middleware.Auth(
			appErrors.ValidateRequestBody[coreTickets.CreateTicketDTO](
				http.HandlerFunc(controller.Create),
			),
		),
	)
	mux.Handle(
		"GET /tickets",
		middleware.Auth(
			http.HandlerFunc(controller.Get),
		),
	)
	mux.Handle(
		"GET /tickets/{id}",
		middleware.Auth(
			http.HandlerFunc(controller.Getbyid),
		),
	)
	mux.Handle(
		"PUT /tickets/{id}",
		middleware.Auth(
			appErrors.ValidateRequestBody[coreTickets.UpdateTicketDTO](
				http.HandlerFunc(controller.Update),
			),
		),
	)
	mux.Handle(
		"DELETE /tickets/{id}",
		middleware.Auth(
			http.HandlerFunc(controller.Delete),
		),
	)
}
