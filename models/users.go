package models

//User defines a system user
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"Password"`
}

//NewUser instasiate a new Operator
func NewUser(email, password string) User {
	return User{
		ID:       "1",
		Email:    email,
		Password: password,
	}
}

//Validate ensure that all User's field are of the valid format
//and returns a non nil error if it's not.
func (user *User) Validate() error { return nil }
