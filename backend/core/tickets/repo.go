package tickets

import (
	"mini-project/backend/database"
	"mini-project/backend/models"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) CreateTicket(ticket *models.Ticket) (*models.Ticket, error) {
	err := database.DB.Create(ticket).Error

	if err != nil {
		return nil, err
	}

	return ticket, nil
}
