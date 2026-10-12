package tickets

type CreateTicketDTO struct {
	Title       string `json:"title" validate:"required,notblank,max=200"`
	Description string `json:"description" validate:"required,notblank,max=500"`
}

type UpdateTicketDTO struct {
	Title       string `json:"title" validate:"omitempty,notblank,max=200"`
	Description string `json:"description" validate:"omitempty,notblank,max=500"`
	Status      string `json:"status" validate:"omitempty,oneof=OPEN IN_PROGRESS CLOSED"`
}
