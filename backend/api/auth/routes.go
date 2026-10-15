package auth

import (
	"net/http"
	"time"

	appErrors "mini-project/backend/apperr"
	coreAuth "mini-project/backend/core/auth"
	"mini-project/backend/middleware"
)

func RegisterRoutes(mux *http.ServeMux) {
	controller := coreAuth.NewController()

	mux.Handle(
		"GET /auth/verify",
		middleware.RateLimit(20, time.Minute, "auth-verify")(
			http.HandlerFunc(controller.VerifyToken),
		),
	)

	mux.Handle(
		"POST /auth/refresh",
		middleware.RateLimit(20, time.Minute, "auth-refresh")(
			http.HandlerFunc(controller.Refresh),
		),
	)

	mux.Handle(
		"POST /auth/signup",
		middleware.RateLimit(5, time.Minute, "auth-signup")(
			appErrors.ValidateRequestBody[coreAuth.SignupDTO](
				http.HandlerFunc(controller.Signup),
			),
		),
	)

	mux.Handle(
		"POST /auth/signin",
		middleware.RateLimit(10, time.Minute, "auth-signin")(
			appErrors.ValidateRequestBody[coreAuth.SigninDTO](
				http.HandlerFunc(controller.Signin),
			),
		),
	)

	mux.Handle(
		"POST /auth/signout",
		middleware.RateLimit(20, time.Minute, "auth-signout")(
			http.HandlerFunc(controller.SignOut),
		),
	)
}
