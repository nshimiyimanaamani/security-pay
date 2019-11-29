package feedback

import (
	"time"

	"github.com/ttacon/libphonenumber"
)

// Message ...
type Message struct {
	ID        string    `json:"id"`
	Title     string    `json:"title,omitempty"`
	Body      string    `json:"body,omitempty"`
	Creator   string    `json:"creator,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"update_at,omitempty"`
}

// Validate validates a Message instance
func (msg *Message) Validate() error {
	if msg.Title == "" || msg.Body == "" || msg.Creator == "" {
		return ErrInvalidEntity
	}
	num, _ := libphonenumber.Parse(msg.Creator, "RW")
	if !libphonenumber.IsValidNumberForRegion(num, "RW") {
		return ErrInvalidEntity
	}
	return nil
}
