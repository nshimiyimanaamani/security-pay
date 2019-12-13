package users

//User defines a system user
type User struct {
	ID       string `json:"id,omitempty"`
	Cell     string `json:"cell,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Sector   string `json:"sector,omitempty"`
	Village  string `json:"village,omitemoty"`
}

//Validate ensure that all User's field are of the valid format
//and returns a non nil error if it's not.
func (user *User) Validate() error {
	if user.Password == "" {
		return ErrInvalidEntity
	}

	if user.Username == "" {
		return ErrInvalidEntity
	}

	if !user.isValidAddress() {
		return ErrInvalidEntity
	}
	return nil
}

// technical debt
func (user *User) validateUsername() error {
	return nil
}

// CheckCell verifies whehter the Cell field is populated
func (user *User) isValidAddress() bool {
	if user.Cell == "" {
		return false
	}
	if user.Sector == "" {
		return false
	}
	if user.Village == "" {
		return false
	}
	return true
}

// func (user *User) validateLogin() error {
// 	if user.Username == "" {
// 		return ErrInvalidEntity
// 	}
// 	if user.Password == "" {
// 		return ErrInvalidEntity
// 	}
// 	return nil
// }
