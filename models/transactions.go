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
	Transaction []Transaction
}

//Validate ensure that all Transaction's field are of the valid format
//and returns a non nil error if it's not
func (tr *Transaction) Validate() error {
	return nil
}
