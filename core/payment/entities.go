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
	MSISDN    string    `json:"phone,omitempty"`
	Method    string    `json:"payment_method,omitempty"`
	Invoice   uint64    `json:"invoce_id,omitempty"`
	Confirmed bool      `json:"confirmed,omitempty"`
	CreatedAt time.Time `json:"recorded_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Confirm payment
func (p *Payment) Confirm() {
	p.Confirmed = true
}

// HasCode checks whether a pull payment at least has the a property code
func (p *Payment) HasCode() error {
	const op errors.Op = "core/payment/Payment.HasCode"

	if p.Code == "" {
		return errors.E(op, "missing house code", errors.KindBadRequest)
	}
	return nil
}

// HasInvoice verfies invoice
func (p *Payment) HasInvoice() error {
	const op errors.Op = "core/payment/Payment.HasInvoice"

	if p.Invoice == 0 {
		return errors.E(op, "invoice id not set", errors.KindBadRequest)
	}
	return nil
}

// Ready to send be sent to the payment gateway
func (p *Payment) Ready() error {
	const op errors.Op = "core/payment/Payment.Ready"

	if p.MSISDN == "" {
		return errors.E(op, "missing phone number", errors.KindBadRequest)
	}

	if p.Amount == float64(0) {
		return errors.E(op, "amount must be greater than zero", errors.KindBadRequest)
	}

	if p.Method == "" {
		return errors.E(op, "payment method must be specified", errors.KindBadRequest)
	}
	return nil
}
