package auth

import "github.com/rugwirobaker/paypack-backend/pkg/errors"

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
