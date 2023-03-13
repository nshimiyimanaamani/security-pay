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

// Webhook defines the webhook mode
const WebhookMode = "production"

// State defines transaction states
type State string

// Possible transaction states
const (
	Pending    State = "processing"
	Failed     State = "failed"
	Successful State = "successful"
)

var ToTxState = map[string]State{
	"pending":    Pending,
	"failed":     Failed,
	"successful": Successful,
}

// PaymentTk...
type PaymentTk struct {
	Access  string `json:"access_token,omitempty"`
	Refresh string `json:"refresh_token,omitempty"`
	Expires int64  `json:"expires_at,omitempty"`
}

// TxResponse ...
type TxResponse struct {
	Status  string `json:"status,omitempty"`
	TxID    string `json:"transaction_id,omitempty"`
	Message string `json:"message,omitempty"`
	TxState State  `json:"transaction_state,omitempty"`
}

// Callback defines the response got from the callback
type Callback struct {
	Data Data   `json:"data"`
	Kind string `json:"kind"`
}

type Data struct {
	Ref       string     `json:"ref,omitempty"`
	Kind      string     `json:"kind,omitempty"`
	Fee       float64    `json:"fee,omitempty"`
	Client    string     `json:"client,omitempty"`
	Amount    float64    `json:"amount,omitempty"`
	Status    string     `json:"status,omitempty"`
	Created   *time.Time `json:"created_at,omitempty"`
	Processed *time.Time `json:"processed_at,omitempty"`
	Commited  *time.Time `json:"commited_at,omitempty"`
}

// Validate validats a callback
func (cb *Callback) Validate() error {
	const op errors.Op = "core/payment/Callback.Validate"

	if cb.Kind == "" {
		return errors.E(op, "Kind field must not be empty", errors.KindBadRequest)
	}

	if cb.Data.Ref == "" {
		return errors.E(op, "transaction ref field must not be empty", errors.KindBadRequest)
	}

	if cb.Data.Status == "" {
		return errors.E(op, "status field must not be empty", errors.KindBadRequest)
	}
	return nil
}

// TxRequest ...
type TxRequest struct {
	ID        string    `json:"id,omitempty"`
	Code      string    `json:"code,omitempty"`
	Amount    float64   `json:"amount,string,omitempty"`
	MSISDN    string    `json:"phone,omitempty"`
	Method    Method    `json:"payment_method,omitempty"`
	Invoice   uint64    `json:"invoce_id,omitempty"`
	Confirmed bool      `json:"confirmed,omitempty"`
	CreatedAt time.Time `json:"recorded_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Confirm payment
func (p *TxRequest) Confirm() {
	p.Confirmed = true
}

// HasCode checks whether a pull payment at least has the a property code
func (p *TxRequest) HasCode() error {
	const op errors.Op = "core/payment/Payment.HasCode"

	if p.Code == "" {
		return errors.E(op, "missing house code", errors.KindBadRequest)
	}
	return nil
}

// HasInvoice verfies invoice
func (p *TxRequest) HasInvoice() error {
	const op errors.Op = "core/payment/Payment.HasInvoice"

	if p.Invoice == 0 {
		return errors.E(op, "invoice id not set", errors.KindBadRequest)
	}
	return nil
}

// Ready to send be sent to the payment gateway
func (p *TxRequest) Ready() error {
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
