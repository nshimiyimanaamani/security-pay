package feedback

import (
	"time"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/ttacon/libphonenumber"
)

// Message ...
type Message struct {
	ID          string    `json:"id"`
	Title       string    `json:"title,omitempty"`
	Body        string    `json:"body,omitempty"`
	Creator     string    `json:"creator,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"update_at,omitempty"`
}

// Validate validates a Message instance
func (msg *Message) Validate() error {
	const op errors.Op = "app/feedback/message.Validate"

	if msg.Title == "" {
		return errors.E(op, "invalid message: missing title", errors.KindBadRequest)
	}

	if msg.Body == "" {
		return errors.E(op, "invalid message: missing body", errors.KindBadRequest)
	}

	if msg.Creator == "" {
		return errors.E(op, "invalid message: missing creator", errors.KindBadRequest)
	}

	num, _ := libphonenumber.Parse(msg.Creator, "RW")
	if !libphonenumber.IsValidNumberForRegion(num, "RW") {
		return errors.E(op, "invalid message: invalid phone number", errors.KindBadRequest)
	}
	return nil
}

// MessagePage is a collection of messages
type MessagePage struct {
	Messages []Message
	PageMetadata
}

// PageMetadata contains page metadata that helps navigation.
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}
