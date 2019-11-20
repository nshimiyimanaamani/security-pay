package feedback

import (
	"time"

	"github.com/ttacon/libphonenumber"
)

// Message ...
type Message struct {
	ID        string
	Title     string
	Body      string
	Creator string
	CreatedAt time.Time
	UpdatedAt time.Time
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
