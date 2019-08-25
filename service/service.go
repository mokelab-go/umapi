package service

import (
	"errors"
	"net/http"

	"github.com/mokelab-go/server/entity"
	"github.com/mokelab-go/umapi"
)

var (
	errEmptyIdentifier = errors.New("empty identifier")
)

type service struct {
	validator   Validator
	accountRepo AccountRepository
}

// New creates instance
func New(validator Validator,
	accountRepo AccountRepository) umapi.Service {
	return &service{
		validator:   validator,
		accountRepo: accountRepo,
	}
}

func (s *service) CreateAccount(identifier, password string, params map[string]interface{}) entity.Response {
	err := s.validator.ValidateIdentifier(identifier)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "input_error", err)
	}
	err = s.validator.ValidatePassword(password)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "input_error", err)
	}

	account, err := s.accountRepo.Create(identifier, password, params)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "create_account_failed", err)
	}
	return entity.Response{
		Status: http.StatusCreated,
		Body:   account.ToJSON(),
	}
}

func (s *service) ChangePassword(session umapi.Session, oldPassword, newPassword string) entity.Response {
	err := s.validator.ValidatePassword(newPassword)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "input_error", err)
	}
	accountID := session.AccountID()

	_, err = s.accountRepo.GetWithIDAndPassword(accountID, oldPassword)
	if err != nil {
		if err == ErrWrongPassword {
			return errorResponse(http.StatusBadRequest, "input_error", err)
		}
		return errorResponse(http.StatusInternalServerError, "server_error", err)
	}

	err = s.accountRepo.UpdatePassword(accountID, newPassword)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "server_error", err)
	}

	return entity.Response{
		Status: http.StatusNoContent,
	}
}

func (s *service) ResetPasswordRequest(identifier string) entity.Response {
	if len(identifier) == 0 {
		return errorResponse(http.StatusBadRequest, "input_error", errEmptyIdentifier)
	}

	a, err := s.accountRepo.GetWithIdentifier(identifier)
	if err != nil {
		return entity.Response{
			Status: http.StatusNoContent,
		}
	}

	err = s.accountRepo.CreateResetPasswordRequest(a)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "server_error", err)
	}

	return entity.Response{
		Status: http.StatusNoContent,
	}
}

func (s *service) GetResetPasswordRequest(token string) entity.Response {
	r, err := s.accountRepo.GetResetPasswordRequestWithToken(token)
	if err != nil {
		if err == ErrNotFound {
			return errorResponse(http.StatusNotFound, "not_found", err)
		}
		return errorResponse(http.StatusInternalServerError, "server_error", err)
	}
	return entity.Response{
		Status: http.StatusOK,
		Body:   r.ToJSON(),
	}
}

func (s *service) ResetPassword(token, newPassword string) entity.Response {
	err := s.validator.ValidatePassword(newPassword)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "input_error", err)
	}
	r, err := s.accountRepo.GetResetPasswordRequestWithToken(token)
	if err != nil {
		if err == ErrNotFound {
			return errorResponse(http.StatusNotFound, "not_found", err)
		}
		return errorResponse(http.StatusInternalServerError, "server_error", err)
	}

	accountID := r.AccountID()
	err = s.accountRepo.UpdatePassword(accountID, newPassword)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "server_error", err)
	}

	err = s.accountRepo.DeleteResetPasswordRequest(token)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "server_error", err)
	}

	return entity.Response{
		Status: http.StatusNoContent,
	}
}

func errorResponse(status int, code string, err error) entity.Response {
	return entity.Response{
		Status: status,
		Body: map[string]interface{}{
			"code": code,
			"msg":  err.Error(),
		},
	}
}
