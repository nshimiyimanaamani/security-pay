package transactions

import (
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"time"
)

//Transaction defines a payment made for a property(i.e house).
type Transaction struct {
	ID           string            `json:"id,omitempty"`
	MadeFor      string            `json:"madefor,omitempty"`
	MadeBy       string            `json:"madeby,omitempty"`
	Address      map[string]string `json:"address,omitempty"`
	Amount       float64           `json:"amount,omitempty"`
	Method       string            `json:"method,omitempty"`
	Invoice      uint64            `json:"invoice,omitempty"`
	DateRecorded time.Time         `json:"date_recorded,omitempty"`
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

//Validate ensure that all Transaction's field are of the valid format
//and returns a non nil error if it's not
func (tr *Transaction) Validate() error {
	const op errors.Op = "app/transactions/transaction.Validate"

	if tr.Amount == 0 || tr.Method == "" || tr.MadeFor == "" || tr.MadeBy == "" {
		return errors.E(op, "invalid transaction", errors.KindBadRequest)
	}
	return nil
}
