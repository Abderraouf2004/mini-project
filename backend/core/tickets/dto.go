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

type DTO struct {
	ID          string `json:"id" validate:"required,uuid"`
	Title       string `json:"title" validate:"required,notblank,max=200"`
	Description string `json:"description" validate:"required,notblank,max=500"`
	Status      string `json:"status" validate:"required,oneof=OPEN IN_PROGRESS CLOSED"`
	UserID      string `json:"user_id" validate:"required,uuid"`
	CreatedAt   string `json:"created_at" validate:"required"`
	UpdatedAt   string `json:"updated_at" validate:"required"`
}
