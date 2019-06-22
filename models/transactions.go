package models

//Transaction defines a payment made for a property(i.e house).
type Transaction struct {
	ID       string
	Amount   string
	Method   string
	Property string
	Date     string
}

// TransactionPage represents a list of transaction.
type TransactionPage struct {
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
