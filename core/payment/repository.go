package payment

import (
	"context"
)

// Repository saves validated Transactions to the underlying datastore
type Repository interface {
	//Save a new pending payment to the database
	Save(ctx context.Context, payment *TxRequest) error
	//Find payment by id
	Find(ctx context.Context, id string) ([]*TxRequest, error)

	//Update the state of an existing payment
	Update(context.Context, string, []*TxRequest) error

	//BulkSave saves multiple payments to the database
	BulkSave(context.Context, []*TxRequest) error

	// PaymentRequest generates all payments
	List(context.Context, *Filters) (PaymentResponse, error)

	//Returns Transactions Per Sector,cell,village
	TodayTransaction(context.Context, *MetricFilters) (Transaction, error)

	//Returns Transactions Per Sector,cell,village
	ListDailyTransactions(context.Context, *MetricFilters) ([]Transactions, error)
}
