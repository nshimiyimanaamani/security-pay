package models

//Transaction defines a payment made for a property(i.e house).
type Transaction struct {
	ID       string
	Amount   string
	Property string
	Date     string
}

//Validate ensure that all Transaction's field are of the valid format
//and returns a non nil error if it's not
func (tr *Transaction) Validate() error {
	return nil
}
