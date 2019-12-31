package transactions

import (
	"time"
)

//Transaction defines a payment made for a property(i.e house).
type Transaction struct {
	ID           string
	MadeFor      string
	MadeBy       string
	Address      map[string]string
	Amount       float64
	Method       string
	Invoice      uint64
	DateRecorded time.Time
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
	if tr.Amount == 0 || tr.Method == "" || tr.MadeFor == "" || tr.MadeBy == "" {
		return ErrInvalidEntity
	}
	return nil
}
