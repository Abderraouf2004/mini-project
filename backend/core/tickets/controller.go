package tickets

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

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
	if errors.Is(err, ErrInvalidTicketID) {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.BadRequest,
			"Invalid ticket ID",
			"",
		))
		return
	}

	if errors.Is(err, ErrEmptyTicketUpdate) || errors.Is(err, ErrInvalidTicketUpdate) {
		appErrors.WriteError(w, appErrors.NewApiError(
			appErrors.ValidationError,
			"Invalid ticket update",
			"Provide at least one valid field to update",
		))
		return
	}

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
	data, ok := appErrors.ValidatedBody[CreateTicketDTO](r)
	if !ok {
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
	page := 1
	pageSize := 10

	if value := r.URL.Query().Get("page"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 1_000_000 {
			http.Error(w, "invalid page", http.StatusBadRequest)
			return
		}
		page = parsed
	}

	if value := r.URL.Query().Get("page_size"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 100 {
			http.Error(w, "invalid page_size", http.StatusBadRequest)
			return
		}
		pageSize = parsed
	}
	tickets, total, err := c.service.GetTickets(userID, page, pageSize)
	if err != nil {
		ticketError(w, err, "list")
		return
	}
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	response := map[string]interface{}{
		"data": tickets,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": totalPages,
		},
		"message": "Tickets retrieved successfully",
	}
	writeJSON(w, http.StatusOK, response)
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
	data, ok := appErrors.ValidatedBody[UpdateTicketDTO](r)
	if !ok {
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
