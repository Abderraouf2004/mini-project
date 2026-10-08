package tickets

import (
	"encoding/json"
	"mini-project/backend/middleware"
	"net/http"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	service := NewService()

	return &Controller{
		service: service,
	}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var data CreateTicketDTO

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// userID := "8ff90057-48ee-4a28-8f1f-4f9529d20ffd"
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)

	if !ok {
		http.Error(w, "user not authenticated", http.StatusUnauthorized)
		return
	}

	ticket, err := c.service.CreateTicket(data, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    ticket,
		"message": "Ticket created successfully",
	})
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	tickets, err := c.service.GetTickets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    tickets,
		"message": "Tickets retrieved successfully",
	})
}

func (c *Controller) Getbyid(w http.ResponseWriter, r *http.Request) {
	// Extract the ticket ID from the URL path
	ticketID := r.PathValue("id")

	ticket, err := c.service.GetTicketByID(ticketID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    ticket,
		"message": "Ticket retrieved successfully",
	})
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	var data UpdateTicketDTO
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	ticketID := r.PathValue("id")
	ticket, err := c.service.UpdateTicket(ticketID, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    ticket,
		"message": "Ticket updated successfully",
	})
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	ticketID := r.PathValue("id")
	err := c.service.DeleteTicket(ticketID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Ticket deleted successfully",
	})
}
