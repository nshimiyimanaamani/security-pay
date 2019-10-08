package nova

import (
	"github.com/rugwirobaker/paypack-backend/app/payment"
)

// messages
const (
	Message = "Paypack: kwishyura umutekano"
)

// Request to the initialization endpoint
type Request struct {
	Token                  string  `json:"token"`
	Amount                 float64 `json:"amount"`
	MSISDN                 string  `json:"msisdn"`
	ExternalTransactionsID string  `json:"external_transaction_id"`
	FromMsg                string  `json:"from_msg"`
	ToMsg                  string  `json:"to_msg"`
}

// PaymentToRequest converts payment to initialization request
func PaymentToRequest(token string, py payment.Payment) Request {
	return Request{
		Token:                  token,
		Amount:                 py.Amount,
		MSISDN:                 py.Phone,
		FromMsg:                Message,
		ToMsg:                  Message,
		ExternalTransactionsID: py.ID,
	}
}
