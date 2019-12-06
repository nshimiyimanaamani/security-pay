package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

type paymentRepo struct {
	db *sql.DB
}

// NewPaymentRepo ...
func NewPaymentRepo(db *sql.DB) payment.Repository {
	return &paymentRepo{}
}

func (repo *paymentRepo) Save(ctx context.Context, tx payment.Transaction) error {
	const op errors.Op = "paymentRepo.Save"

	q := `INSERT INTO transactions (id, madefor, amount, method, date_recorded) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	_, err := repo.db.Exec(q, tx.ID, tx.Code, tx.Amount, tx.Method, tx.RecordedAt)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return errors.E(op, "duplicate transaction", errors.KindAlreadyExists)
			case errInvalid, errTruncation:
				return errors.E(op, "invalid entity transaction", errors.KindBadRequest)
			}
		}
		return err
	}

	return errors.E(op, errors.KindNotImplemented)
}
