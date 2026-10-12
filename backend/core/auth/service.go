// package auth

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	"mini-project/backend/models"

// 	"github.com/golang-jwt/jwt/v5"
// )

// type Service struct {
// 	repo *Repository
// }

// func NewService() *Service {
// 	return &Service{
// 		repo: NewRepository(),
// 	}
// }

// type TokenPair struct {
// 	AccessToken  string `json:"accessToken"`
// 	RefreshToken string `json:"refreshToken"`
// }

// func (s *Service) VerifyToken(tokenString string) (*models.User, error) {
// 	return s.verifyToken(tokenString, "access")
// }

// func (s *Service) verifyToken(tokenString, expectedType string) (*models.User, error) {
// 	secret := os.Getenv("JWT_SECRET")
// 	if secret == "" {
// 		return nil, fmt.Errorf("JWT_SECRET is not configured")
// 	}

// 	token, err := jwt.Parse(
// 		tokenString,
// 		func(token *jwt.Token) (interface{}, error) {
// 			if token.Method != jwt.SigningMethodHS256 {
// 				return nil, fmt.Errorf("invalid token")
// 			}
// 			return []byte(secret), nil
// 		},
// 		jwt.WithValidMethods([]string{"HS256"}),
// 		jwt.WithExpirationRequired(),
// 	)
// 	if err != nil || token == nil || !token.Valid {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	userID, ok := claims["id"].(string)
// 	if !ok || userID == "" {
// 		return nil, fmt.Errorf("invalid token")
// 	}

// 	tokenType, ok := claims["typ"].(string)
// 	if !ok || tokenType != expectedType {
// 		return nil, fmt.Errorf("invalid token type")
// 	}

// 	user, err := s.repo.ReadByID(userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if user == nil {
// 		return nil, fmt.Errorf("user not found")
// 	}

// 	return user, nil
// }

// func (s *Service) RefreshToken(refreshToken string) (string, error) {
// 	user, err := s.verifyToken(refreshToken, "refresh")
// 	if err != nil {
// 		return "", err
// 	}

// 	return GenerateToken(user.ID, "access", 15*time.Minute)
// }

// func (s *Service) Signup(data SignupDTO) (*TokenPair, error) {
// 	existingUser, err := s.repo.ReadByEmail(data.Email)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if existingUser != nil {
// 		return nil, fmt.Errorf("email already in use")
// 	}

// 	hashedPassword, err := HashPassword(data.Password)
// 	if err != nil {
// 		return nil, err
// 	}

// 	user := &models.User{
// 		Email:        data.Email,
// 		PasswordHash: hashedPassword,
// 	}

// 	createdUser, err := s.repo.Signup(user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return s.generateTokens(createdUser.ID)
// }

// func (s *Service) Signin(data SigninDTO) (*TokenPair, error) {
// 	user, err := s.repo.ReadByEmail(data.Email)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if user == nil || !CheckPassword(data.Password, user.PasswordHash) {
// 		return nil, fmt.Errorf("invalid email or password")
// 	}

// 	return s.generateTokens(user.ID)
// }

// func (s *Service) generateTokens(userID string) (*TokenPair, error) {
// 	accessToken, err := GenerateToken(userID, "access", 15*time.Minute)
// 	if err != nil {
// 		return nil, err
// 	}

// 	refreshToken, err := GenerateToken(userID, "refresh", 7*24*time.Hour)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &TokenPair{
// 		AccessToken:  accessToken,
// 		RefreshToken: refreshToken,
// 	}, nil
// }

package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"mini-project/backend/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEmailAlreadyInUse  = errors.New("email already in use")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

type Service struct {
	repo *Repository
}

func NewService() *Service {
	return &Service{
		repo: NewRepository(),
	}
}

type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (s *Service) VerifyToken(tokenString string) (*models.User, error) {
	return s.verifyToken(tokenString, "access")
}

func (s *Service) verifyToken(tokenString, expectedType string) (*models.User, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is not configured")
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, ErrInvalidToken
			}
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil || token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	userID, ok := claims["id"].(string)
	if !ok || userID == "" {
		return nil, ErrInvalidToken
	}

	tokenType, ok := claims["typ"].(string)
	if !ok || tokenType != expectedType {
		return nil, ErrInvalidToken
	}

	user, err := s.repo.ReadByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidToken
	}

	return user, nil
}

func (s *Service) RefreshToken(refreshToken string) (string, error) {
	user, err := s.verifyToken(refreshToken, "refresh")
	if err != nil {
		return "", err
	}

	return GenerateToken(user.ID, "access", 15*time.Minute)
}

func (s *Service) Signup(data SignupDTO) (*TokenPair, error) {
	existingUser, err := s.repo.ReadByEmail(data.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyInUse
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

	return s.generateTokens(createdUser.ID)
}

func (s *Service) Signin(data SigninDTO) (*TokenPair, error) {
	user, err := s.repo.ReadByEmail(data.Email)
	if err != nil {
		return nil, err
	}
	if user == nil || !CheckPassword(data.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokens(user.ID)
}

func (s *Service) generateTokens(userID string) (*TokenPair, error) {
	accessToken, err := GenerateToken(userID, "access", 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateToken(userID, "refresh", 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
