package accounts

import "github.com/nshimiyimanaamani/paypack-backend/pkg/errors"

import "time"

// AccountType represents the type of an account and the privileges it has
type AccountType string

const (
	// Devs represents developer account
	Devs AccountType = "dev"
	//Bens is beneficiary account or sectors
	Bens AccountType = "ben"
)

// Account represents an account entity
type Account struct {
	ID            string      `json:"id,omitempty"`
	Name          string      `json:"name,omitempty"`
	Type          AccountType `json:"type,omitempty"`
	NumberOfSeats int         `json:"seats,omitempty"`
	Active        bool        `json:"active,omitempty"`
	CreatedAt     time.Time   `json:"created_at,omitempty"`
	UpdatedAt     time.Time   `json:"updated_at,omitempty"`
}

// Validate validates an account at registration time
func (ac *Account) Validate() error {
	const op errors.Op = "app/accounts/account.Validate"

	if ac.ID == "" {
		return errors.E(op, "invalid account: missing id", errors.KindBadRequest)
	}

	if ac.Name == "" {
		return errors.E(op, "invalid account: missing name", errors.KindBadRequest)
	}

	if ac.Type == "" {
		return errors.E(op, "invalid account: missing type", errors.KindBadRequest)
	}
	return nil
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}

// AccountPage is a collection of accounts plus some metadata
type AccountPage struct {
	Accounts []Account
	PageMetadata
}
