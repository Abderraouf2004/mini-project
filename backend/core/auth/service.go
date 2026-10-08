package auth

import (
	"fmt"
	"mini-project/backend/models"
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	repo := NewRepository()

	return &Service{
		repo: repo,
	}
}

func (s *Service) Signup(data SignupDTO) (*models.User, error) {
	existingUser, err := s.repo.ReadByEmail(data.Email)

	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already in use")
	}

	hashedPassword, err := HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        data.Email,
		PasswordHash: hashedPassword,
	}

	createdUser, err := s.repo.Signup(user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// func (s *Service) Signin(data SigninDTO) (*models.User, error) {
// 	user, err := s.repo.ReadByEmail(data.Email)

// 	if err != nil || user == nil {
// 		return nil, fmt.Errorf("invalid email or password")
// 	}

// 	isValid := CheckPassword(data.Password, user.PasswordHash)

// 	if !isValid {
// 		return nil, fmt.Errorf("invalid email or password")
// 	}

//		return user, nil
//	}
func (s *Service) Signin(data SigninDTO) (string, error) {
	user, err := s.repo.ReadByEmail(data.Email)

	if err != nil || user == nil {
		return "", fmt.Errorf("invalid email or password")
	}

	isValid := CheckPassword(data.Password, user.PasswordHash)

	if !isValid {
		return "", fmt.Errorf("invalid email or password")
	}

	accessToken, err := GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
