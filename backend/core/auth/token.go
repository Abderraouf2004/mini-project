// package auth

// import (
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/golang-jwt/jwt/v5"
// )

// func GenerateToken(userID string, duration time.Duration) (string, error) {
// 	secret := os.Getenv("JWT_SECRET")
// 	if secret == "" {
// 		return "", fmt.Errorf("JWT_SECRET is not configured")
// 	}

// 	if userID == "" {
// 		return "", fmt.Errorf("user ID is required")
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"id":  userID,
// 		"exp": time.Now().Add(duration).Unix(),
// 	})

// 	return token.SignedString([]byte(secret))
// }

package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID, tokenType string, duration time.Duration) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET is not configured")
	}

	if userID == "" {
		return "", fmt.Errorf("user ID is required")
	}

	if tokenType != "access" && tokenType != "refresh" {
		return "", fmt.Errorf("invalid token type")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userID,
		"typ": tokenType,
		"exp": time.Now().Add(duration).Unix(),
	})

	return token.SignedString([]byte(secret))
}
