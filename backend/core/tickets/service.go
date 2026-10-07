package tickets

import "mini-project/backend/models"

type Service struct {
	repo *Repository
}

func NewService() *Service {
	repo := NewRepository()

	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateTicket(data CreateTicketDTO, userID string) (*models.Ticket, error) {
	ticket := &models.Ticket{
		Title:       data.Title,
		Description: data.Description,
		UserID:      userID,
	}

	return s.repo.CreateTicket(ticket)
}
