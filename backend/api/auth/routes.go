package auth

import (
	"net/http"

	coreAuth "mini-project/backend/core/auth"
	appErrors "mini-project/backend/errors"
)

func RegisterRoutes(mux *http.ServeMux) {
	controller := coreAuth.NewController()

	mux.Handle(
		"POST /auth/signup",
		appErrors.ValidateRequestBody[coreAuth.SignupDTO](
			http.HandlerFunc(controller.Signup),
		),
	)

	mux.Handle(
		"POST /auth/signin",
		appErrors.ValidateRequestBody[coreAuth.SigninDTO](
			http.HandlerFunc(controller.Signin),
		),
	)
}
