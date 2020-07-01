package payment

import (
	"time"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// TxExpiration is the time it takes for a non confirmed treansaction to expire
const TxExpiration = time.Minute * 10

// Method payment method
type Method string

// supported payment methods
const (
	MTN    Method = "momo-mtn-rw"
	AIRTEL Method = "momo-airtel-rw"
)

// State defines transaction states
type State string

//Possible transaction states
const (
	Pending    State = "processing"
	Failed     State = "failed"
	Successful State = "successful"
)

// Status ...
type Status struct {
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
	const op errors.Op = "payment.Callback.Validate"
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

// Transaction ...
type Transaction struct {
	ID         string    `json:"id,omitempty"`
	Code       string    `json:"code,omitempty"`
	Amount     float64   `json:"amount,string,omitempty"`
	Phone      string    `json:"phone,omitempty"`
	Invoice    uint64    `json:"invoce_id,omitempty"`
	Method     Method    `json:"payment_method,omitempty"`
	Namespace  string    `json:"namespace,omitempty"`
	RecordedAt time.Time `json:"recorded_at,omitempty"`
}

// Validate returns an error if the Transaction entity doesn't adhere to
// the requirements
func (py *Transaction) Validate() error {
	const op errors.Op = "payment.Transaction.Validate"
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

// Invoice ...
type Invoice struct {
	ID     uint64
	Amount float64
}
