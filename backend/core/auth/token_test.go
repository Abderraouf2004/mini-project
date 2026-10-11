package auth

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-for-unit-tests")

	tokenString, err := GenerateToken("user-123")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
	)

	if err != nil {
		t.Fatalf("expected valid token, got: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("expected JWT claims")
	}

	if claims["id"] != "user-123" {
		t.Fatalf("expected user ID user-123, got: %v", claims["id"])
	}
}

func TestGenerateTokenWithoutSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	_, err := GenerateToken("user-123")
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is empty")
	}
}
