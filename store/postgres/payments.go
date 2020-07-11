package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (payment.Repository) = (*paymentStore)(nil)

type paymentStore struct {
	*sql.DB
}

// NewPaymentRepository creates a new postgres backed payment.Repository
func NewPaymentRepository(db *sql.DB) payment.Repository {
	return &paymentStore{db}
}

func (repo *paymentStore) Save(ctx context.Context, payment payment.Payment) error {
	const op errors.Op = "store/postgres/paymentStore.Save"

	q := `INSERT INTO payments(
			id,
			amount,
			msisdn,
			method,
			invoice,
			property,
			confirmed
		) VALUES($1, $2, $3, $4, $5, $6, $7);
	`
	_, err := repo.ExecContext(ctx, q,
		payment.ID,
		payment.Amount,
		payment.MSISDN,
		payment.Method,
		payment.Invoice,
		payment.Code,
		payment.Confirmed,
	)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return errors.E(op, "duplicate payment id", errors.KindAlreadyExists)
			case errInvalid, errTruncation:
				return errors.E(op, "invalid payment entity", errors.KindBadRequest)
			}
		}
		return errors.E(op, err, errors.KindUnexpected)
	}
	return nil
}

func (repo *paymentStore) Find(ctx context.Context, id string) (payment.Payment, error) {
	const op errors.Op = "store/postgres/paymentStore.Find"

	q := `SELECT 
			id, 
			amount, 
			msisdn, 
			method, 
			invoice, 
			property,
			confirmed,
			created_at,
			updated_at
		FROM payments WHERE id=$1
	`
	var pmt payment.Payment

	err := repo.QueryRowContext(ctx, q, id).Scan(
		&pmt.ID,
		&pmt.Amount,
		&pmt.MSISDN,
		&pmt.Method,
		&pmt.Invoice,
		&pmt.Code,
		&pmt.Confirmed,
		&pmt.CreatedAt,
		&pmt.UpdatedAt,
	)
	if err != nil {
		var empty payment.Payment
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, err, "payment not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	return pmt, nil
}

func (repo *paymentStore) Update(ctx context.Context, payment payment.Payment) error {
	const op errors.Op = "store/postgres/paymentStore.Update"

	q := `UPDATE payments SET confirmed=$1 WHERE id=$2`

	res, err := repo.ExecContext(ctx, q, payment.Confirmed, payment.ID)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return errors.E(op, err, "payment not found", errors.KindNotFound)
			}
		}
		return errors.E(op, err, errors.KindUnexpected)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	if cnt == 0 {
		return errors.E(op, err, "payment not found", errors.KindNotFound)
	}
	return nil
}
