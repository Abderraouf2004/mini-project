package auth

import (
	"errors"

	"mini-project/backend/database"
	"mini-project/backend/models"

	"gorm.io/gorm"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Signup(user *models.User) (*models.User, error) {
	err := database.DB.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) ReadByEmail(email string) (*models.User, error) {
	var user models.User

	err := database.DB.
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repository) ReadByID(id string) (*models.User, error) {
	var user models.User

	err := database.DB.
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}
