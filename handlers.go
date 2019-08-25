package umapi

import (
	"net/http"

	"github.com/mokelab-go/hop"
)

func createAccount(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := hop.BodyJSON(r.Context())
		identifier, _ := params["identifier"].(string)
		password, _ := params["password"].(string)

		resp := s.CreateAccount(identifier, password, params)
		resp.Write(w)
	}
}

func changePassword(s Service, getSession func(r *http.Request) Session) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := hop.BodyJSON(r.Context())
		session := getSession(r)

		oldPassword, _ := params["old_password"].(string)
		newPassword, _ := params["new_password"].(string)

		resp := s.ChangePassword(session, oldPassword, newPassword)
		resp.Write(w)
	}
}

func resetPasswordRequest(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := hop.BodyJSON(r.Context())

		identifier, _ := params["identifier"].(string)

		resp := s.ResetPasswordRequest(identifier)
		resp.Write(w)
	}
}

func getResetPasswordRequest(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		token := hop.PathString(c, "token")

		resp := s.GetResetPasswordRequest(token)
		resp.Write(w)
	}
}

func resetPassword(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		params := hop.BodyJSON(c)
		token := hop.PathString(c, "token")

		newPassword, _ := params["new_password"].(string)

		resp := s.ResetPassword(token, newPassword)
		resp.Write(w)
	}
}
