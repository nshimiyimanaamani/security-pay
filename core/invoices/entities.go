package invoices

import (
	"time"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Status of the invoices
type Status string

// possible invoice states
const (
	Pending Status = "pending"
	Payed   Status = "payed"
	Expired Status = "expired"
)

// Invoice ...
type Invoice struct {
	ID        uint64    `json:"id"`
	Amount    float64   `json:"amount"`
	Property  string    `json:"property"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Verify checkes wether the invoice satisfies requirements to be paid
func (vc *Invoice) Verify(amount float64) error {
	const op errors.Op = "core/payment/Invoice.Satisfy"

	if vc.Status == Payed {
		return errors.E(op, "you already payed", errors.KindRateLimit)
	}
	if vc.Amount != amount {
		return errors.E(op, "amount doesn't match invoice", errors.KindBadRequest)
	}
	return nil
}

// PageMetadata ...
type PageMetadata struct {
	Total       uint
	Months      uint
	TotalAmount float64
}

// InvoicePage ...
type InvoicePage struct {
	Invoices []Invoice
	PageMetadata
}
