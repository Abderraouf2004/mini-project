package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

func WriteError(w http.ResponseWriter, err error) {
	var apiErr *ApiError

	if !errors.As(err, &apiErr) {
		apiErr = NewApiError(
			InternalServerError,
			"Internal server error",
			"",
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(apiErr.Status)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  apiErr.Status,
		"code":    apiErr.Code,
		"message": apiErr.Message,
		"details": apiErr.Details,
	})
}
