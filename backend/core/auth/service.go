package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"

	"mini-project/backend/database"
	"mini-project/backend/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

var (
	ErrEmailAlreadyInUse      = errors.New("email already in use")
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrInvalidToken           = errors.New("invalid or expired token")
	ErrTokenRevoked           = errors.New("token has been revoked")
	ErrAuthServiceUnavailable = errors.New("authentication service unavailable")
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
	return s.verifyToken(tokenString)
}

func (s *Service) verifyToken(tokenString string) (*models.User, error) {
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

	if database.Redis == nil {
		return nil, ErrAuthServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := database.Redis.Get(
		ctx,
		blacklistKey(tokenString),
	).Result()

	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("%w: Redis blacklist check failed: %v", ErrAuthServiceUnavailable, err)
	}

	if result != "" {
		return nil, ErrTokenRevoked
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

func (s *Service) RefreshToken(Token string) (string, error) {
	user, err := s.verifyToken(Token)
	if err != nil {
		return "", err
	}

	accessToken, err := GenerateToken(
		user.ID,
		15*time.Minute,
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
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
		if isUniqueConstraintViolation(err) {
			return nil, ErrEmailAlreadyInUse
		}
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
	accessToken, err := GenerateToken(userID, 15*time.Minute)
	if err != nil {
		return nil, err
	}

	refreshToken, err := GenerateToken(userID, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func isUniqueConstraintViolation(err error) bool {
	var sqlStateError interface{ SQLState() string }
	return errors.As(err, &sqlStateError) && sqlStateError.SQLState() == "23505"
}

func blacklistKey(token string) string {
	hash := sha256.Sum256([]byte(token))
	return "blacklist:" + hex.EncodeToString(hash[:])
}

func (s *Service) revokeToken(tokenString string) error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return fmt.Errorf("JWT_SECRET is not configured")
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
		return ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ErrInvalidToken
	}

	exp, err := claims.GetExpirationTime()
	if err != nil || exp == nil {
		return ErrInvalidToken
	}

	ttl := time.Until(exp.Time)
	if ttl <= 0 {
		return nil
	}

	if database.Redis == nil {
		return ErrAuthServiceUnavailable
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := database.Redis.Set(
		ctx,
		blacklistKey(tokenString),
		"blacklisted",
		ttl,
	).Err(); err != nil {
		return fmt.Errorf("%w: Redis token revocation failed: %v", ErrAuthServiceUnavailable, err)
	}

	return nil
}

func (s *Service) SignOut(tokenString string) error {
	return s.revokeToken(tokenString)
}
