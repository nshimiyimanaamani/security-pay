package postgres

import (
	"context"
	"database/sql"
	"time"

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
	const op errors.Op = "postgres.paymentRepo.Save"

	q := `INSERT INTO transactions (id, madefor, amount, method, date_recorded) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	tx.RecordedAt = time.Now()

	_, err := repo.Exec(q, tx.ID, tx.Code, tx.Amount, tx.Method, tx.RecordedAt)
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

func (repo *paymentRepo) RetrieveCode(ctx context.Context, code string) (string, error) {
	const op errors.Op = "postgres.paymentRepo.RetrieveCode"

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
