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

func (r *Repository) GetTickets() ([]*models.Ticket, error) {
	var tickets []*models.Ticket
	err := database.DB.Find(&tickets).Error
	if err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *Repository) GetTicketByID(id string, userID string) (*models.Ticket, error) {
	var ticket models.Ticket
	err := database.DB.First(&ticket, "id = ? AND user_id = ?", id, userID).Error
	if err != nil {
		return nil, err
	}

	return &ticket, nil
}

func (r *Repository) UpdateTicket(ticket *models.Ticket) (*models.Ticket, error) {
	err := database.DB.Save(ticket).Error
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *Repository) DeleteTicket(ticket *models.Ticket) error {
	return database.DB.Delete(ticket).Error
}
