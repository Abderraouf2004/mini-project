package tickets

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	appErrors "mini-project/backend/apperr"
	"mini-project/backend/middleware"

	"gorm.io/gorm"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

func ticketError(w http.ResponseWriter, err error, operation string) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.NotFound,
			"Ticket not found",
			"",
		))
		return
	}

	slog.Error("Ticket operation failed",
		"operation", operation,
		"error", err,
	)

	appErrors.WriteError(w, appErrors.NewApiError(
		appErrors.InternalServerError,
		"Internal server error",
		"",
	))
}

func authenticatedUserID(w http.ResponseWriter, r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.Unauthorized,
			"User not authenticated",
			"",
		))
		return "", false
	}

	return userID, true
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

	userID, ok := authenticatedUserID(w, r)
	if !ok {
		return
	}

	ticket, err := c.service.CreateTicket(data, userID)
	if err != nil {
		ticketError(w, err, "create")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"data":    ticket,
		"message": "Ticket created successfully",
	})
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserID(w, r)
	if !ok {
		return
	}

	tickets, err := c.service.GetTickets(userID)
	if err != nil {
		ticketError(w, err, "list")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":    tickets,
		"message": "Tickets retrieved successfully",
	})
}

func (c *Controller) Getbyid(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserID(w, r)
	if !ok {
		return
	}

	ticket, err := c.service.GetTicketByID(r.PathValue("id"), userID)
	if err != nil {
		ticketError(w, err, "get_by_id")
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

	userID, ok := authenticatedUserID(w, r)
	if !ok {
		return
	}

	ticket, err := c.service.UpdateTicket(r.PathValue("id"), data, userID)
	if err != nil {
		ticketError(w, err, "update")
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":    ticket,
		"message": "Ticket updated successfully",
	})
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticatedUserID(w, r)
	if !ok {
		return
	}

	err := c.service.DeleteTicket(r.PathValue("id"), userID)
	if err != nil {
		ticketError(w, err, "delete")
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
			"Internal server error",
			"",
		))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(append(body, '\n'))
}
