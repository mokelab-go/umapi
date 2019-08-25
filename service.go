package umapi

import (
	"github.com/mokelab-go/server/entity"
)

// Service is service
type Service interface {
	// CreateAccount is called in POST /accounts
	CreateAccount(identifier, password string, params map[string]interface{}) entity.Response

	// ChangePassword is called in PUT /password
	ChangePassword(session Session, oldPassword, newPassword string) entity.Response

	// ResetPasswordRequest is called in POST /password/reset
	ResetPasswordRequest(identifier string) entity.Response

	// ResetPasswordRequest is called in GET /password/reset/{token}
	GetResetPasswordRequest(token string) entity.Response

	// ResetPasswordRequest is called in POST /password/reset/{token}
	ResetPassword(token, newPassword string) entity.Response
}
