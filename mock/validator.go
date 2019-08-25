package mock

import (
	"github.com/mokelab-go/umapi/service"
)

type validator struct{}

// NewMockValidator creates empty validator
func NewMockValidator() service.Validator {
	return &validator{}
}

func (v *validator) ValidateIdentifier(identifier string) error {
	return nil
}

func (v *validator) ValidatePassword(password string) error {
	return nil
}
