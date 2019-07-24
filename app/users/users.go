package users

import "github.com/asaskevich/govalidator"

//User defines a system user
type User struct {
	ID       string `json:"id"`
	Cell     string `json:"cell"`
	Email    string `json:"email"`
	Password string `json:"Password"`
}

//Validate ensure that all User's field are of the valid format
//and returns a non nil error if it's not.
func (user *User) Validate() error {
	if user.Email == "" || user.Password == "" {
		return ErrInvalidEntity
	}

	if !govalidator.IsEmail(user.Email) {
		return ErrInvalidEntity
	}

	return nil
}

// CheckCell verifies whehter the Cell field is populated
func (user *User) CheckCell() bool {
	if user.Cell != "" {
		return true
	}
	return false
}
