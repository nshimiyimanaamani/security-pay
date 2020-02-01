package invoices

import "time"

// Status of the invoices
type Status string

// possible invoice states
const (
	Pending Status = "pending"
	Payed   Status = "payed"
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

// PageMetadata ...
type PageMetadata struct {
	Total  uint
	Months uint
}

// InvoicePage ...
type InvoicePage struct {
	Invoices []Invoice
	PageMetadata
}
