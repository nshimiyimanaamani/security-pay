package http

//Error ...
type Error struct {
	Message string `json:"message"`
}

// NewError creates a new http error
func NewError(err error) Error {
	message := err.Error()
	return Error{
		Message: message,
	}
}
