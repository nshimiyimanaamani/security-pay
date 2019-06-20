package models

//Operator defines a system user
type Operator struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"Password"`
}

//NewOperator instasiate a new Operator
func NewOperator(email, password string) Operator {
	return Operator{
		ID:       "1",
		Email:    email,
		Password: password,
	}
}

//Validate returns an nil error if all the operator fields are valid
func (opr *Operator) Validate() error { return nil }
