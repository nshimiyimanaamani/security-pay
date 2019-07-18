package transactions

import (
	"time"
)

//Transaction defines a payment made for a property(i.e house).
type Transaction struct {
	ID           string
	Property     string
	Amount       string
	Method       string
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
	if tr.Amount == "" || tr.Method == "" || tr.Property == "" {
		return ErrInvalidEntity
	}
	return nil
}
