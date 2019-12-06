package payment

import "github.com/rugwirobaker/paypack-backend/pkg/errors"

import "time"

// State indicates the transaction state(failed, pending, success)
type State uint

// transaction statuses
const (
	Successful State = iota + 1
	Pending
	Failed
)

func (s State) String() string {
	statuses := []string{"success, pending, failed"}
	switch s {
	case 0:
		return statuses[0]
	case 1:
		return statuses[1]
	case 2:
		return statuses[1]
	default:
		return "unknown"
	}
}

// Transaction ...
type Transaction struct {
	ID         string  `json:"id,omitempty"`
	Code       string  `json:"code,omitempty"`
	Amount     float64 `json:"amount,string,omitempty"`
	Phone      string  `json:"phone,omitempty"`
	Method     string  `json:"payment_method,omitempty"`
	RecordedAt time.Time
}

// Validate returns an error if the Transaction entity doesn't adhere to
// the requirements
func (py *Transaction) Validate() error {
	const op errors.Op = "payment.Validate"
	if py.Code == "" {
		return errors.E(op, "code field must not be empty", errors.KindBadRequest)
	}

	if py.Phone == "" {
		return errors.E(op, "phone field must not be empty", errors.KindBadRequest)
	}

	if py.Amount == float64(0) {
		return errors.E(op, "payment amount must be greate than zero", errors.KindBadRequest)
	}

	if py.Method == "" {
		return errors.E(op, "payment method must be specified", errors.KindBadRequest)
	}
	return nil
}
