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

func (s *Service) GetTickets() ([]*models.Ticket, error) {
	return s.repo.GetTickets()
}

func (s *Service) GetTicketByID(id string) (*models.Ticket, error) {
	return s.repo.GetTicketByID(id)
}

func (s *Service) UpdateTicket(id string, data UpdateTicketDTO) (*models.Ticket, error) {
	ticket, err := s.repo.GetTicketByID(id)
	if err != nil {
		return nil, err
	}

	if data.Title != "" {
		ticket.Title = data.Title
	}
	if data.Description != "" {
		ticket.Description = data.Description
	}

	return s.repo.UpdateTicket(ticket)
}

func (s *Service) DeleteTicket(id string) error {
	ticket, err := s.repo.GetTicketByID(id)
	if err != nil {
		return err
	}

	return s.repo.DeleteTicket(ticket)
}
