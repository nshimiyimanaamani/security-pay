package properties

import (
	"github.com/rugwirobaker/paypack-backend/app/users"
)

// AuthBackend calls the user service to identify a user
type AuthBackend interface {
	// Identity takes a token and returns a user'id if and a nil error
	// if the user is succefully identified.
	Identity(token string) (string, error)
}

var _ (AuthBackend) = (*authotization)(nil)

type authotization struct {
	svc users.Service
}

// NewAuthBackend creates a new authorization backend instance
func NewAuthBackend(svc users.Service) AuthBackend {
	return &authotization{svc}
}

func (auth *authotization) Identity(token string) (string, error) {
	return auth.svc.Identify(token)
}
