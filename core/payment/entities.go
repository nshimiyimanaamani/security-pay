package payment

import (
	"time"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// TxExpiration is the time it takes for a non confirmed treansaction to expire
const TxExpiration = time.Minute * 10

// payment methods
const (
	MTN    = "mtn"
	AIRTEL = "airtel"
)

// State defines transaction states
type State string

//Possible transaction states
const (
	Pending    State = "processing"
	Failed     State = "failed"
	Successful State = "successful"
)

// Response ...
type Response struct {
	Status  string `json:"status,omitempty"`
	TxID    string `json:"transaction_id,omitempty"`
	Message string `json:"message,omitempty"`
	TxState State  `json:"transaction_state,omitempty"`
}

// Callback defines the response got from the callback
type Callback struct {
	Status string       `json:"status,omitempty"`
	Data   CallBackData `json:"data,omitempty"`
}

// CallBackData ...
type CallBackData struct {
	Message    string `json:"message,omitempty"`
	GwRef      string `json:"gwRef,omitempty"`
	TrxRef     string `json:"trxRef,omitempty"`
	ChannelRef string `json:"channelRef,omitempty"`
	State      State  `json:"state,omitempty"`
}

// Validate validats a callback
func (cb *Callback) Validate() error {
	const op errors.Op = "core/payment/Callback.Validate"

	if cb.Status == "" {
		return errors.E(op, "status field must not be empty", errors.KindBadRequest)
	}

	if cb.Data.TrxRef == "" {
		return errors.E(op, "trxRef field must not be empty", errors.KindBadRequest)
	}

	if cb.Data.GwRef == "" {
		return errors.E(op, "gwRef field must not be empty", errors.KindBadRequest)
	}

	if cb.Data.State == "" {
		return errors.E(op, "state field must not be empty", errors.KindBadRequest)
	}
	return nil
}

// Payment ...
type Payment struct {
	ID        string    `json:"id,omitempty"`
	Code      string    `json:"code,omitempty"`
	Amount    float64   `json:"amount,string,omitempty"`
	Phone     string    `json:"phone,omitempty"`
	Invoice   uint64    `json:"invoce_id,omitempty"`
	Method    string    `json:"payment_method,omitempty"`
	Namespace string    `json:"namespace,omitempty"`
	CreatedAt time.Time `json:"recorded_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Validate returns an error if the Transaction entity doesn't adhere to
// the requirements
func (py *Payment) Validate() error {
	const op errors.Op = "core/payment/Transaction.Validate"
	if py.Code == "" {
		return errors.E(op, "missing house code", errors.KindBadRequest)
	}

	if py.Phone == "" {
		return errors.E(op, "missing phone number", errors.KindBadRequest)
	}

	if py.Amount == float64(0) {
		return errors.E(op, "amount must be greater than zero", errors.KindBadRequest)
	}

	if py.Method == "" {
		return errors.E(op, "payment method must be specified", errors.KindBadRequest)
	}
	return nil
}

// HackyValidation is to satisfy current needs
func (py *Payment) HackyValidation() error {
	const op errors.Op = "core/payment/Transaction.HackyValidation"

	if py.Phone == "" {
		return errors.E(op, "missing phone number", errors.KindBadRequest)
	}

	if py.Amount == float64(0) {
		return errors.E(op, "amount must be greater than zero", errors.KindBadRequest)
	}

	if py.Method == "" {
		return errors.E(op, "payment method must be specified", errors.KindBadRequest)
	}
	return nil
}

// Invoice ...
type Invoice struct {
	ID     uint64
	Amount float64
	Status string
}

// Satisfy checkes wether the invoice satisfies requirements to be paid
func (inv *Invoice) Satisfy(amount float64) error {
	const op errors.Op = "core/payment/Invoice.Satisfy"

	if inv.Status == "payed" {
		return errors.E(op, "you already payed", errors.KindRateLimit)
	}
	if inv.Amount != amount {
		return errors.E(op, "amount doesn't match invoice", errors.KindBadRequest)
	}
	return nil
}
