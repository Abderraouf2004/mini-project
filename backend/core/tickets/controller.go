package tickets

import (
	"encoding/json"
	"net/http"

	appErrors "mini-project/backend/errors"
	"mini-project/backend/middleware"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var data CreateTicketDTO

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.BadRequest,
			"Invalid request body",
			"",
		))
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.Unauthorized,
			"User not authenticated",
			"",
		))
		return
	}

	ticket, err := c.service.CreateTicket(data, userID)
	if err != nil {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.BadRequest,
			"Failed to create ticket",
			err.Error(),
		))
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"data":    ticket,
		"message": "Ticket created successfully",
	})
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.Unauthorized,
			"User not authenticated",
			"",
		))
		return
	}

	tickets, err := c.service.GetTickets(userID)
	if err != nil {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.InternalServerError,
			"Failed to retrieve tickets",
			err.Error(),
		))
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":    tickets,
		"message": "Tickets retrieved successfully",
	})
}

func (c *Controller) Getbyid(w http.ResponseWriter, r *http.Request) {
	ticketID := r.PathValue("id")

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.Unauthorized,
			"User not authenticated",
			"",
		))
		return
	}

	ticket, err := c.service.GetTicketByID(ticketID, userID)
	if err != nil {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.InternalServerError,
			"Failed to retrieve ticket",
			err.Error(),
		))
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":    ticket,
		"message": "Ticket retrieved successfully",
	})
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	var data UpdateTicketDTO

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.BadRequest,
			"Invalid request body",
			"",
		))
		return
	}

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.Unauthorized,
			"User not authenticated",
			"",
		))
		return
	}

	ticketID := r.PathValue("id")

	ticket, err := c.service.UpdateTicket(ticketID, data, userID)
	if err != nil {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.BadRequest,
			"Failed to update ticket",
			err.Error(),
		))
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":    ticket,
		"message": "Ticket updated successfully",
	})
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	ticketID := r.PathValue("id")

	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.Unauthorized,
			"User not authenticated",
			"",
		))
		return
	}

	if err := c.service.DeleteTicket(ticketID, userID); err != nil {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.InternalServerError,
			"Failed to delete ticket",
			err.Error(),
		))
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Ticket deleted successfully",
	})
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	body, err := json.Marshal(data)
	if err != nil {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.InternalServerError,
			"Failed to encode response",
			"",
		))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(append(body, '\n'))
}
