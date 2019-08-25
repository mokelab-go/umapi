package service

// Validator validates input parameter
type Validator interface {
	ValidateIdentifier(identifier string) error

	ValidatePassword(password string) error
}
