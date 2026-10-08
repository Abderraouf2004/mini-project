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

func (s *Service) GetTicketByID(id string, userID string) (*models.Ticket, error) {
	ticket, err := s.repo.GetTicketByID(id, userID)
	if err != nil {
		return nil, err
	}
	return ticket, nil
}

func (s *Service) UpdateTicket(id string, data UpdateTicketDTO, userID string) (*models.Ticket, error) {
	ticket, err := s.repo.GetTicketByID(id, userID)
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

func (s *Service) DeleteTicket(id string, userID string) error {
	ticket, err := s.repo.GetTicketByID(id, userID)
	if err != nil {
		return err
	}

	return s.repo.DeleteTicket(ticket)
}
