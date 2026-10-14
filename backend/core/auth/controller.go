package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"mini-project/backend/apperr"
	"net/http"
	"strings"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	return &Controller{
		service: NewService(),
	}
}

func handleAuthError(w http.ResponseWriter, err error, operation string) {
	switch {
	case errors.Is(err, ErrEmailAlreadyInUse):
		apperr.WriteError(w, apperr.NewApiError(
			apperr.Conflict,
			"Email already registered",
			"",
		))

	case errors.Is(err, ErrInvalidCredentials):
		apperr.WriteError(w, apperr.NewApiError(
			apperr.Unauthorized,
			"Invalid email or password",
			"",
		))

	case errors.Is(err, ErrInvalidToken),
		errors.Is(err, ErrTokenRevoked):
		apperr.WriteError(w, apperr.NewApiError(
			apperr.Unauthorized,
			"Invalid, expired, or revoked token",
			"",
		))

	default:
		slog.Error("Authentication operation failed",
			"operation", operation,
			"error", err,
		)
		apperr.WriteError(w, apperr.NewApiError(
			apperr.InternalServerError,
			"Internal server error",
			"",
		))
	}
}

func (c *Controller) Signup(w http.ResponseWriter, r *http.Request) {
	data, ok := apperr.ValidatedBody[SignupDTO](r)
	if !ok {
		apperr.WriteError(w, apperr.NewApiError(
			apperr.BadRequest,
			"Invalid request body",
			"",
		))
		return
	}

	tokens, err := c.service.Signup(data)
	if err != nil {
		handleAuthError(w, err, "signup")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    tokens,
		"message": "Sign Up successfully",
	})
}

func (c *Controller) Signin(w http.ResponseWriter, r *http.Request) {
	data, ok := apperr.ValidatedBody[SigninDTO](r)
	if !ok {
		apperr.WriteError(w, apperr.NewApiError(
			apperr.BadRequest,
			"Invalid request body",
			"",
		))
		return
	}

	tokens, err := c.service.Signin(data)
	if err != nil {
		handleAuthError(w, err, "signin")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    tokens,
		"message": "Sign In successfully",
	})
}

func (c *Controller) VerifyToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token, ok := strings.CutPrefix(authHeader, "Bearer ")

	if !ok || token == "" || strings.Contains(token, " ") {
		apperr.WriteError(w, apperr.NewApiError(
			apperr.Unauthorized,
			"Invalid authorization header",
			"",
		))
		return
	}

	user, err := c.service.VerifyToken(token)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) ||
			errors.Is(err, ErrTokenRevoked) {
			apperr.WriteError(w, apperr.NewApiError(
				apperr.Unauthorized,
				"Invalid, expired, or revoked token",
				"",
			))
			return
		}

		slog.Error("Token verification failed", "error", err)
		apperr.WriteError(w, apperr.NewApiError(
			apperr.InternalServerError,
			"Internal server error",
			"",
		))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": map[string]interface{}{
			"user": map[string]interface{}{
				"id": user.ID,
			},
		},
		"message": "Token is valid",
	})
}

func (c *Controller) Refresh(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token, ok := strings.CutPrefix(authHeader, "Bearer ")

	if !ok || token == "" || strings.Contains(token, " ") {
		apperr.WriteError(w, apperr.NewApiError(
			apperr.Unauthorized,
			"Invalid authorization header",
			"",
		))
		return
	}

	accessToken, err := c.service.RefreshToken(token)
	if err != nil {
		if errors.Is(err, ErrInvalidToken) ||
			errors.Is(err, ErrTokenRevoked) {
			apperr.WriteError(w, apperr.NewApiError(
				apperr.Unauthorized,
				"Invalid, expired, or revoked token",
				"",
			))
			return
		}

		slog.Error("Token refresh failed", "error", err)
		apperr.WriteError(w, apperr.NewApiError(
			apperr.InternalServerError,
			"Internal server error",
			"",
		))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"data": map[string]string{
			"accessToken": accessToken,
		},
		"message": "Token refreshed successfully",
	})
}

func (c *Controller) SignOut(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	token, ok := strings.CutPrefix(authHeader, "Bearer ")

	if !ok || token == "" || strings.Contains(token, " ") {
		apperr.WriteError(w, apperr.NewApiError(
			apperr.Unauthorized,
			"Invalid authorization header",
			"",
		))
		return
	}

	if err := c.service.SignOut(token); err != nil {
		if errors.Is(err, ErrInvalidToken) ||
			errors.Is(err, ErrTokenRevoked) {
			apperr.WriteError(w, apperr.NewApiError(
				apperr.Unauthorized,
				"Invalid, expired, or revoked token",
				"",
			))
			return
		}

		slog.Error("Sign out failed", "error", err)
		apperr.WriteError(w, apperr.NewApiError(
			apperr.InternalServerError,
			"Internal server error",
			"",
		))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Signed out successfully",
	})
}
