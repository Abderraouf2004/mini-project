package auth

import (
	"encoding/json"
	"net/http"
)

type Controller struct {
	service *Service
}

func NewController() *Controller {
	service := NewService()

	return &Controller{
		service: service,
	}
}

func (c *Controller) Signup(w http.ResponseWriter, r *http.Request) {
	var data SignupDTO

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tokens, err := c.service.Signup(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    tokens,
		"message": "Sign Up successfully",
	})
}

func (c *Controller) Signin(w http.ResponseWriter, r *http.Request) {
	var data SigninDTO

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	accessToken, err := c.service.Signin(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": map[string]string{
			"accessToken": accessToken,
		},
		"message": "Sign In successfully",
	})
}
