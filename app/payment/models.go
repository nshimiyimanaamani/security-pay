package payment

import "time"

// Status indicates the transaction status
type Status uint

// transaction statuses
const (
	Successful Status = iota + 1
	Pending
	Failed
)

func (s Status) String() string {
	return ""
}

// Payment ...
type Payment struct {
	ID     string  `json:"id,omitempty"`
	Code   string  `json:"code,omitempty"`
	Amount float64 `json:"amount,omitempty"`
	Phone  string  `json:"phone,omitempty"`
}

// Message is a message to the client
type Message map[string]interface{}

// Validation is the request we recieve via the callback
type Validation struct {
	Token                  string  `json:"token,omitempty"`
	Amount                 float64 `json:"amount,omitempty"`
	MSISDN                 string  `json:"msisdn,omitempty"`
	ExternalTransactionsID string  `json:"external_transaction_id,omitempty"`
	Status                 int     `json:"status,omitempty"`
	Ref                    string  `json:"reference_id,omitempty"`
}

// Property ...
type Property struct {
	ID    string
	Owner string
}

// Transaction ...
type Transaction struct {
	ID           string
	MadeFor      string
	MadeBy       string
	Amount       float64
	Method       string
	DateRecorded time.Time
}
