package tickets

import "time"

type CreateTicketDTO struct {
	Title       string `json:"title" validate:"required,max=200"`
	Description string `json:"description" validate:"required,max=500"`
}

type UpdateTicketDTO struct {
	Title       string `json:"title" validate:"optional,max=200"`
	Description string `json:"description" validate:"optional,max=500"`
}
type DTO struct {
	ID          string    `json:"id" validate:"required,uuid"`
	Title       string    `json:"title" validate:"required,max=200"`
	Description string    `json:"description" validate:"required,max=500"`
	Status      string    `json:"status" validate:"required,oneof:'open' 'in_progress' 'closed'"`
	CreatedAt   time.Time `json:"created_at" validate:"required, datetime"`
	UpdatedAt   time.Time `json:"updated_at" validate:"required, datetime"`
	UserID      string    `json:"user_id" validate:"required,uuid"`
}
