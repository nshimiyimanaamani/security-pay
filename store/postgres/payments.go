package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

type paymentRepo struct {
	*sql.DB
}

// NewPaymentRepo ...
func NewPaymentRepo(db *sql.DB) payment.Repository {
	return &paymentRepo{db}
}

func (repo *paymentRepo) Save(ctx context.Context, tx payment.Transaction) error {
	const op errors.Op = "store/postgres/paymentRepo.Save"

	q := `
		INSERT INTO transactions (
			id, 
			madefor, 
			amount, 
			method, 
			invoice
		) VALUES ($1, $2, $3, $4, $5) RETURNING created_at`

	err := repo.QueryRow(q, tx.ID, tx.Code, tx.Amount, tx.Method, tx.Invoice).Scan(&tx.RecordedAt)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return errors.E(op, "duplicate transaction", errors.KindAlreadyExists)
			case errInvalid, errTruncation:
				return errors.E(op, "invalid transaction entity", errors.KindBadRequest)
			}
		}
		return errors.E(op, err, errors.KindUnexpected)
	}
	return nil
}

func (repo *paymentRepo) RetrieveProperty(ctx context.Context, code string) (string, error) {
	const op errors.Op = "store/postgres/paymentRepo.RetrieveProperty"

	q := `SELECT id FROM properties WHERE id=$1`

	var property string

	if err := repo.QueryRow(q, code).Scan(&property); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return "", errors.E(op, err, "property not found", errors.KindNotFound)
		}
		return "", errors.E(op, err, errors.KindUnexpected)
	}

	return property, nil
}

func (repo *paymentRepo) OldestInvoice(ctx context.Context, property string) (uint64, error) {
	const op errors.Op = "store/postgres/paymentRepo.OldestInvoice"

	q := `SELECT id FROM invoices WHERE created_at = (SELECT MIN(created_at) FROM invoices WHERE property=$1 AND status='pending');`

	var invoice uint64

	if err := repo.QueryRow(q, property).Scan(&invoice); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return 0, errors.E(op, err, "no invoice found", errors.KindNotFound)
		}
		return 0, errors.E(op, err, errors.KindUnexpected)
	}

	return invoice, nil
}
