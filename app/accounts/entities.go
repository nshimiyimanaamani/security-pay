package accounts

import "github.com/rugwirobaker/paypack-backend/pkg/errors"

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
	ID            string
	Name          string
	Type          AccountType
	NumberOfSeats int
	Active        bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Validate validates an account at registration time
func (ac *Account) Validate() error {
	const op errors.Op = "accounts/account.Validate"

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
