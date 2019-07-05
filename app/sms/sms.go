package sms

// Message defines an SMS message
type Message struct {
	Body        string
	Destination string
}

// Validate validates the struct of a message
func (msg *Message) Validate() error {
	if msg.Body == "" || msg.Destination == "" {
		return ErrInvalidEntity
	}
	return nil
}
