package mock

import (
	"fmt"

	"github.com/mokelab-go/umapi/service"
)

type repo struct {
	idCount               int
	identifiers           map[string]*account
	ids                   map[string]*account
	resetPasswordRequests map[string]*resetPasswordRequest
}

type account struct {
	id         string
	identifier string
	password   string
}

type resetPasswordRequest struct {
	token     string
	accountID string
}

// NewMockAccountRepository creates mock repository
func NewMockAccountRepository() service.AccountRepository {
	return &repo{
		idCount:     0,
		identifiers: make(map[string]*account),
		ids:         make(map[string]*account),
		resetPasswordRequests: make(map[string]*resetPasswordRequest),
	}
}

// Create new account
func (r *repo) Create(identifier, password string, params map[string]interface{}) (service.Account, error) {
	if _, ok := r.identifiers[identifier]; ok {
		return nil, service.ErrAlreadyExists
	}
	id := fmt.Sprintf("id_%d", r.idCount)
	r.idCount++
	a := &account{
		id:         id,
		identifier: identifier,
		password:   password,
	}
	// put local DB
	r.identifiers[identifier] = a
	r.ids[id] = a
	return a, nil
}

func (r *repo) GetWithIDAndPassword(id, password string) (service.Account, error) {
	a, ok := r.ids[id]
	if !ok {
		return nil, service.ErrNotFound
	}
	if a.password != password {
		return nil, service.ErrWrongPassword
	}
	return a, nil
}

func (r *repo) GetWithIdentifier(identifier string) (service.Account, error) {
	a, ok := r.identifiers[identifier]
	if !ok {
		return nil, service.ErrNotFound
	}
	return a, nil
}

func (r *repo) UpdatePassword(id, newPassword string) error {
	a, ok := r.ids[id]
	if !ok {
		return service.ErrNotFound
	}
	a.password = newPassword
	return nil
}

func (r *repo) CreateResetPasswordRequest(account service.Account) error {
	token := fmt.Sprintf("token_%d", r.idCount)
	r.idCount++
	req := &resetPasswordRequest{
		token:     token,
		accountID: account.ID(),
	}
	r.resetPasswordRequests[token] = req
	return nil
}

func (r *repo) GetResetPasswordRequestWithToken(token string) (service.ResetPasswordRequest, error) {
	req, ok := r.resetPasswordRequests[token]
	if !ok {
		return nil, service.ErrNotFound
	}
	return req, nil
}

func (r *repo) DeleteResetPasswordRequest(token string) error {
	delete(r.resetPasswordRequests, token)
	return nil
}

func (a *account) ID() string {
	return a.id
}

func (a *account) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":         a.id,
		"identifier": a.identifier,
	}
}

func (r *resetPasswordRequest) AccountID() string {
	return r.accountID
}

func (r *resetPasswordRequest) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id": r.token,
	}
}
