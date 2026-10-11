package tickets

type CreateTicketDTO struct {
	Title       string `json:"title" validate:"required,max=200"`
	Description string `json:"description" validate:"required,max=500"`
}

type UpdateTicketDTO struct {
	Title       string `json:"title" validate:"omitempty,max=200"`
	Description string `json:"description" validate:"omitempty,max=500"`
	Status      string `json:"status" validate:"omitempty,oneof=OPEN IN_PROGRESS CLOSED"`
}
