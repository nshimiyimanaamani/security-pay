package transactions

import (
	"time"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Transaction defines a payment made for a property(i.e house).
type Transaction struct {
	ID           string    `json:"id,omitempty"`
	MadeFor      string    `json:"madefor,omitempty"`
	OwnerID      string    `json:"owner_id,omitempty"`
	OwneFname    string    `json:"owner_firstname,omitempty"`
	OwnerLname   string    `json:"owner_lastname,omitempty"`
	Cell         string    `json:"cell,omitempty"`
	Sector       string    `json:"sector,omitempty"`
	Village      string    `json:"village,omitempty"`
	Amount       float64   `json:"amount,omitempty"`
	Method       string    `json:"method,omitempty"`
	Invoice      uint64    `json:"invoice,omitempty"`
	Namespace    string    `json:"namespace,omitempty"`
	DateRecorded time.Time `json:"date_recorded,omitempty"`
	Paypack_fee  float64   `json:"paypack_fee,omitempty"`
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}

// TransactionPage represents a list of transaction.
type TransactionPage struct {
	PageMetadata
	Transactions []Transaction
}

// Validate ensure that all Transaction's field are of the valid format
// and returns a non nil error if it's not
func (tr *Transaction) Validate() error {
	const op errors.Op = "app/transactions/transaction.Validate"

	if tr.Amount == 0 || tr.Method == "" || tr.MadeFor == "" || tr.OwnerID == "" {
		return errors.E(op, "invalid transaction", errors.KindBadRequest)
	}
	return nil
}
