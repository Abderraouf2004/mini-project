// package apperr

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io"
// 	"net/http"
// )

// const maxRequestBodySize = 1 << 20 // 1 MB

// func ValidateRequestBody[T any](next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

// 		var data T

// 		decoder := json.NewDecoder(r.Body)
// 		decoder.DisallowUnknownFields()

// 		if err := decoder.Decode(&data); err != nil {
// 			WriteError(w, NewApiError(
// 				BadRequest,
// 				"Invalid request body",
// 				"",
// 			))
// 			return
// 		}

// 		var extra any
// 		if err := decoder.Decode(&extra); err != io.EOF {
// 			WriteError(w, NewApiError(
// 				BadRequest,
// 				"Request must contain a single JSON value",
// 				"",
// 			))
// 			return
// 		}

// 		if err := ValidateObject(data); err != nil {
// 			WriteError(w, NewApiError(
// 				ValidationError,
// 				"Request validation failed",
// 				err.Error(),
// 			))
// 			return
// 		}

// 		body, err := json.Marshal(data)
// 		if err != nil {
// 			WriteError(w, NewApiError(
// 				InternalServerError,
// 				"Internal server error",
// 				"",
// 			))
// 			return
// 		}

// 		r.Body = io.NopCloser(bytes.NewReader(body))
// 		next.ServeHTTP(w, r)
// 	})
// }

package apperr

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

const maxRequestBodySize = 1 << 20 // 1 MB

type validatedBodyKey struct{}

func ValidatedBody[T any](r *http.Request) (T, bool) {
	data, ok := r.Context().Value(validatedBodyKey{}).(T)
	return data, ok
}

func ValidateRequestBody[T any](next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodySize)

		var data T

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&data); err != nil {
			WriteError(w, NewApiError(
				BadRequest,
				"Invalid request body",
				"",
			))
			return
		}

		var extra any
		if err := decoder.Decode(&extra); err != io.EOF {
			WriteError(w, NewApiError(
				BadRequest,
				"Request must contain a single JSON value",
				"",
			))
			return
		}

		if err := ValidateObject(data); err != nil {
			WriteError(w, NewApiError(
				ValidationError,
				"Request validation failed",
				err.Error(),
			))
			return
		}

		ctx := context.WithValue(
			r.Context(),
			validatedBodyKey{},
			data,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
