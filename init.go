package umapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mokelab-go/hop"
)

// InitRouter is an entry point
func InitRouter(router *mux.Router, s Service,
	getSessionOperation func(next http.HandlerFunc) http.HandlerFunc,
	getSession func(r *http.Request) Session) {
	router.Methods(http.MethodPost).
		Path("/accounts").
		Handler(hop.Operations(
			hop.GetBodyAsJSON,
		)(createAccount(s)))
	// Change password
	router.Methods(http.MethodPut).
		Path("/password").
		Handler(hop.Operations(
			getSessionOperation,
			hop.GetBodyAsJSON,
		)(changePassword(s, getSession)))
	// Reset password request
	router.Methods(http.MethodPost).
		Path("/password/reset").
		Handler(hop.Operations(
			hop.GetBodyAsJSON,
		)(resetPasswordRequest(s)))
	// Get Reset password request
	router.Methods(http.MethodGet).
		Path("/password/reset/{token}").
		Handler(hop.Operations(
			hop.GetPathParams,
			hop.GetPathString("token"),
		)(getResetPasswordRequest(s)))
	// Reset password
	router.Methods(http.MethodPost).
		Path("/password/reset/{token}").
		Handler(hop.Operations(
			hop.GetPathParams,
			hop.GetPathString("token"),
			hop.GetBodyAsJSON,
		)(resetPassword(s)))
}
