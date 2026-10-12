// package auth

// import (
// 	"net/http"

// 	appErrors "mini-project/backend/apperr"
// 	coreAuth "mini-project/backend/core/auth"
// )

// func RegisterRoutes(mux *http.ServeMux) {
// 	controller := coreAuth.NewController()

// 	mux.Handle(
// 		"POST /auth/signup",
// 		appErrors.ValidateRequestBody[coreAuth.SignupDTO](
// 			http.HandlerFunc(controller.Signup),
// 		),
// 	)

// 	mux.Handle(
// 		"POST /auth/signin",
// 		appErrors.ValidateRequestBody[coreAuth.SigninDTO](
// 			http.HandlerFunc(controller.Signin),
// 		),
// 	)
// }

package auth

import (
	"net/http"

	appErrors "mini-project/backend/apperr"
	coreAuth "mini-project/backend/core/auth"
)

func RegisterRoutes(mux *http.ServeMux) {
	controller := coreAuth.NewController()

	mux.HandleFunc("GET /auth/verify", controller.VerifyToken)

	mux.HandleFunc("POST /auth/refresh", controller.Refresh)

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
