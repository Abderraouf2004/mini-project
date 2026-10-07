package tickets

import (
	"encoding/json"
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

	userID := "8ff90057-48ee-4a28-8f1f-4f9529d20ffd"

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
