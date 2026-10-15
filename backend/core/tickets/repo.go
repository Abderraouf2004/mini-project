package tickets

import (
	"mini-project/backend/database"
	"mini-project/backend/models"

	"gorm.io/gorm"
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

func (r *Repository) GetTickets(
	userID string,
	page int,
	pageSize int,
) ([]*models.Ticket, int64, error) {
	var tickets []*models.Ticket
	var total int64

	query := database.DB.Model(&models.Ticket{}).
		Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	err := query.
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&tickets).Error

	if err != nil {
		return nil, 0, err
	}

	return tickets, total, nil
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
	result := database.DB.
		Model(&models.Ticket{}).
		Where("id = ? AND user_id = ?", ticket.ID, ticket.UserID).
		Select("Title", "Description", "Status", "UpdatedAt").
		Updates(ticket)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return ticket, nil
}

func (r *Repository) DeleteTicket(ticket *models.Ticket) error {
	result := database.DB.
		Where("id = ? AND user_id = ?", ticket.ID, ticket.UserID).
		Delete(&models.Ticket{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
