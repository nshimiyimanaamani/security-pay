package users

// errorMessage ..,
type errorMessage struct {
	Message string `json:"message"`
}

// newErrorMessage creates a new http error
func newErrorMessage(message string) errorMessage {
	return errorMessage{
		Message: message,
	}
}
