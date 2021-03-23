package notifs

import (
	"time"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Notification represents sms details
type Notification struct {
	ID         string    `json:"_"`
	Message    string    `json:"message"  validate:"required"`
	Sender     string    `json:"namespace" validate:"required"` // corresponds to a Namespace
	Recipients []string  `json:"recipients" validate:"required"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// NoticationPage is a collection of notifcations plus some metadata
type NoticationPage struct {
	Notifications []Notification
	PageMetadata
}

// PageMetadata adds context which helps in navigation
type PageMetadata struct {
	Total  uint64
	Offset uint64
	Limit  uint64
}

// Validate message
func (p *Notification) Validate() error {
	const op errors.Op = "core/notifs/Notification.Validate"
	if p.Message == "" {
		return errors.E(op, "invalid payload: message is required", errors.KindBadRequest)
	}
	if p.Recipients == nil || len(p.Recipients) == 0 {
		return errors.E(op, "invalid payload: recipients must be a non empty array", errors.KindBadRequest)
	}
	return nil
}
