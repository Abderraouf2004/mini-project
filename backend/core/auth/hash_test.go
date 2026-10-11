package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "MyPassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if hash == password {
		t.Fatal("password should not be stored as plain text")
	}

	if !CheckPassword(password, hash) {
		t.Fatal("correct password should be accepted")
	}

	if CheckPassword("WrongPassword123!", hash) {
		t.Fatal("wrong password should be rejected")
	}
}
