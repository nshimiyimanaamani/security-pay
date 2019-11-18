package feedback

import "time"

// Message ...
type Message struct {
	ID        string
	Title     string
	Body      string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Validate validates a Message instance
func (msg *Message) Validate() error {
	if msg.Title == "" || msg.Body == "" {
		return ErrInvalidEntity
	}
	return nil
}
