package service

import "errors"

// AccountRepository provides account related API
type AccountRepository interface {
	// Create new account
	Create(identifier, password string, params map[string]interface{}) (Account, error)

	// GetWithIDAndPassword gets account with id and password
	GetWithIDAndPassword(id, password string) (Account, error)

	// GetWithID gets account with identifier
	GetWithIdentifier(identifier string) (Account, error)

	UpdatePassword(id, newPassword string) error

	CreateResetPasswordRequest(account Account) error

	GetResetPasswordRequestWithToken(token string) (ResetPasswordRequest, error)

	DeleteResetPasswordRequest(token string) error
}

// Account is app Account
type Account interface {
	ID() string
	ToJSON() map[string]interface{}
}

// ResetPasswordRequest is request
type ResetPasswordRequest interface {
	AccountID() string
	ToJSON() map[string]interface{}
}

var (
	// ErrAlreadyExists will be returned in Create()
	ErrAlreadyExists = errors.New("already exists")
	// ErrNotFound will be returned in GetResetPasswordRequestWithToken()
	ErrNotFound = errors.New("not found")
	// ErrWrongPassword will be returned in GetWithIDAndPassword()
	ErrWrongPassword = errors.New("wrong password")
)
