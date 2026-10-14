package tickets

import (
	"errors"
	"fmt"
	"strings"
	"uuid"

	"mini-project/backend/apperr"
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

func (s *Service) CreateTicket(data CreateTicketDTO, userID string) (*DTO, error) {
	ticket := &models.Ticket{
		Title:       strings.TrimSpace(data.Title),
		Description: strings.TrimSpace(data.Description),
		UserID:      userID,
	}

	created, err := s.repo.CreateTicket(ticket)
	if err != nil {
		return nil, err
	}

	response, err := apperr.ValidateAndConvert[DTO](created)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (s *Service) GetTickets(userID string) ([]*DTO, error) {
	tickets, err := s.repo.GetTickets(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*DTO, 0, len(tickets))
	for _, ticket := range tickets {
		response, err := apperr.ValidateAndConvert[DTO](ticket)
		if err != nil {
			return nil, err
		}
		responses = append(responses, &response)
	}

	return responses, nil
}

func (s *Service) GetTicketByID(id string, userID string) (*DTO, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidTicketID
	}

	ticket, err := s.repo.GetTicketByID(id, userID)
	if err != nil {
		return nil, err
	}

	// response, err := apperr.ValidateAndConvert[DTO](ticket)
	// if err != nil {
	// 	return nil, err
	// }
	response, err := apperr.ValidateAndConvert[DTO](ticket)
	if err != nil {
		return nil, fmt.Errorf("validate ticket response: %w", err)
	}

	return &response, nil

}

func (s *Service) UpdateTicket(id string, data UpdateTicketDTO, userID string) (*DTO, error) {
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

	updated, err := s.repo.UpdateTicket(ticket)
	if err != nil {
		return nil, err
	}

	response, err := apperr.ValidateAndConvert[DTO](updated)
	if err != nil {
		return nil, err
	}

	return &response, nil

}

func (s *Service) DeleteTicket(id string, userID string) error {
	ticket, err := s.repo.GetTicketByID(id, userID)
	if err != nil {
		return err
	}

	return s.repo.DeleteTicket(ticket)
}
