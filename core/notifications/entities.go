package notifications

import "github.com/rugwirobaker/paypack-backend/pkg/errors"

// Payload represents sms details
type Payload struct {
	ID         string   `json:"_"`
	Message    string   `json:"message"  validate:"required"`
	Recipients []string `json:"recipients" validate:"required"`
}

func (p *Payload) Validate() error {
	const op errors.Op = "core/notifications/Payload.Validate"
	if p.Message == "" {
		return errors.E(op, "invalid payload: message is required", errors.KindBadRequest)
	}
	if p.Recipients == nil || len(p.Recipients) == 0 {
		return errors.E(op, "invalid payload: recipients must be a non empty array", errors.KindBadRequest)
	}
	return nil
}
