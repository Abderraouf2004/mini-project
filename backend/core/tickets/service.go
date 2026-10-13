package tickets

import (
	"errors"
	"strings"
	"uuid"

	"mini-project/backend/models"
)

var ErrInvalidTicketID = errors.New("invalid ticket ID")

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{
		repo: NewRepository(),
	}
}

func (s *Service) CreateTicket(data CreateTicketDTO, userID string) (*models.Ticket, error) {
	ticket := &models.Ticket{
		Title:       strings.TrimSpace(data.Title),
		Description: strings.TrimSpace(data.Description),
		UserID:      userID,
	}

	return s.repo.CreateTicket(ticket)
}

func (s *Service) GetTickets(userID string) ([]*models.Ticket, error) {
	return s.repo.GetTickets(userID)
}

func (s *Service) GetTicketByID(id string, userID string) (*models.Ticket, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidTicketID
	}
	return s.repo.GetTicketByID(id, userID)
}

func (s *Service) UpdateTicket(id string, data UpdateTicketDTO, userID string) (*models.Ticket, error) {
	if strings.TrimSpace(data.Title) == "" &&
		strings.TrimSpace(data.Description) == "" &&
		strings.TrimSpace(data.Status) == "" {
		return nil, errors.New("at least one field must be provided")
	}

	ticket, err := s.repo.GetTicketByID(id, userID)
	if err != nil {
		return nil, err
	}

	if data.Title != "" {
		ticket.Title = strings.TrimSpace(data.Title)
	}

	if data.Description != "" {
		ticket.Description = strings.TrimSpace(data.Description)
	}

	if data.Status != "" {
		ticket.Status = data.Status
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
