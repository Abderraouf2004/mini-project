package tickets

type CreateTicketDTO struct {
	Title       string `json:"title" validate:"required,max=200"`
	Description string `json:"description" validate:"required,max=500"`
}
