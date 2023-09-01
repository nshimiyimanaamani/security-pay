package auth

import "github.com/nshimiyimanaamani/paypack-backend/pkg/errors"

// Role represents user access level
type Role string

const (
	// Dev has access oveer all accounts
	Dev = "dev"
	// Admin has access level only to the sector they manage
	Admin = "admin"
	// Basic has access to only a single cell.
	Basic = "basic"

	// Min is the minimun privilage level
	Min = "min"
)

// Credentials is user login info
type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
	Account  string `json:"account,omitempty"`
}

// Validate credentials
func (creds *Credentials) Validate() error {
	const op errors.Op = "app/auth/credentials.Validate"
	if creds.Username == "" {
		return errors.E(op, "invalid credentials: missing username", errors.KindNotFound)
	}
	if creds.Password == "" {
		return errors.E(op, "invalid credentials: missing password", errors.KindNotFound)
	}
	return nil
}
