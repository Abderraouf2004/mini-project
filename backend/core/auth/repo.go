package auth

import (
	"mini-project/backend/database"
	"mini-project/backend/models"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) ReadByEmail(email string) (*models.User, error) {
	var user models.User

	err := database.DB.
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Signup(user *models.User) (*models.User, error) {
	err := database.DB.Create(user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}
