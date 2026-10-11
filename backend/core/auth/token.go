package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not configured")
	}

	if userID == "" {
		return "", fmt.Errorf("user ID is required")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	})

	return token.SignedString([]byte(secret))
}
