package tickets

import (
	"net/http"

	coreTickets "mini-project/backend/core/tickets"
	appErrors "mini-project/backend/errors"
)

func RegisterRoutes(mux *http.ServeMux) {
	controller := coreTickets.NewController()

	mux.Handle(
		"POST /tickets",
		appErrors.ValidateRequestBody[coreTickets.CreateTicketDTO](
			http.HandlerFunc(controller.Create),
		),
	)
}
