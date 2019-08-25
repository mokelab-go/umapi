package service

import (
	"net/http"
	"testing"

	"github.com/mokelab-go/server/entity"
	"github.com/mokelab-go/umapi/mock"
	"github.com/mokelab-go/umapi/service"
)

type session struct {
	accountID string
}

func TestService_Create_ChangePassword_OK(t *testing.T) {
	v := mock.NewMockValidator()
	repo := mock.NewMockAccountRepository()
	s := service.New(v, repo)
	resp := s.CreateAccount("test@example.com", "pass0011", nil)
	// assert
	if !checkStatus(t, resp, http.StatusCreated) {
		return
	}
	id, _ := resp.Body["id"]
	if id != "id_0" {
		t.Errorf("Account ID must be %s but %s", "id_0", id)
		return
	}

	// change password
	resp = s.ChangePassword(&session{
		accountID: "id_0",
	}, "pass0011", "pass2233")
	// assert
	if !checkStatus(t, resp, http.StatusNoContent) {
		return
	}
}

func TestService_ChangePassword_WrongOld(t *testing.T) {
	v := mock.NewMockValidator()
	repo := mock.NewMockAccountRepository()
	s := service.New(v, repo)
	resp := s.CreateAccount("test@example.com", "pass0011", nil)
	// assert
	if !checkStatus(t, resp, http.StatusCreated) {
		return
	}
	id, _ := resp.Body["id"]
	if id != "id_0" {
		t.Errorf("Account ID must be %s but %s", "id_0", id)
		return
	}

	// change password - wrong
	resp = s.ChangePassword(&session{
		accountID: "id_0",
	}, "pass2233", "pass3344")
	// assert
	if !checkStatus(t, resp, http.StatusBadRequest) {
		return
	}
	code, _ := resp.Body["code"]
	if code != "input_error" {
		t.Errorf("code must be %s but %s", "input_error", code)
		return
	}
}

func TestService_ResetPassword_OK(t *testing.T) {
	v := mock.NewMockValidator()
	repo := mock.NewMockAccountRepository()
	s := service.New(v, repo)
	resp := s.CreateAccount("test@example.com", "pass0011", nil)
	// assert
	if !checkStatus(t, resp, http.StatusCreated) {
		return
	}
	id, _ := resp.Body["id"]
	if id != "id_0" {
		t.Errorf("Account ID must be %s but %s", "id_0", id)
		return
	}

	// reset password request
	resp = s.ResetPasswordRequest("test@example.com")
	// assert
	if !checkStatus(t, resp, http.StatusNoContent) {
		return
	}

	// get reset password request
	resp = s.GetResetPasswordRequest("token_1")
	// assert
	if !checkStatus(t, resp, http.StatusOK) {
		return
	}

	// reset
	resp = s.ResetPassword("token_1", "pass1122")
	if !checkStatus(t, resp, http.StatusNoContent) {
		return
	}

	// retry - deleted
	resp = s.GetResetPasswordRequest("token_1")
	// assert
	if !checkStatus(t, resp, http.StatusNotFound) {
		return
	}

	resp = s.ResetPassword("token_1", "pass4455")
	// assert
	if !checkStatus(t, resp, http.StatusNotFound) {
		return
	}
}

func (s *session) AccountID() string {
	return s.accountID
}

func checkStatus(t *testing.T, resp entity.Response, expected int) bool {
	if resp.Status == expected {
		return true
	}
	t.Errorf("Status must be %d but %d", expected, resp.Status)
	return false

}
