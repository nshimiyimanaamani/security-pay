package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/accounts"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

type accountRepository struct {
	*sql.DB
}

// NewAccountRepository ...
func NewAccountRepository(db *sql.DB) accounts.Repository {
	return &accountRepository{db}
}

func (repo *accountRepository) Save(ctx context.Context, acc accounts.Account) (accounts.Account, error) {
	const op errors.Op = "store/postgres.accountRepository.Save"

	q := `
		INSERT INTO accounts (
			id, 
			name, 
			type, 
			seats
		) VALUES ($1, $2, $3, $4) RETURNING created_at, updated_at;`

	empty := accounts.Account{}

	err := repo.QueryRow(q, acc.ID, acc.Name, acc.Type, acc.NumberOfSeats).Scan(&acc.CreatedAt, &acc.UpdatedAt)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return empty, errors.E(op, err, "account already exists", errors.KindAlreadyExists)
			case errFK:
				return empty, errors.E(op, err, "invalid input data: sector not found", errors.KindNotFound)
			case errInvalid, errTruncation:
				return empty, errors.E(op, err, "invalid input data ", errors.KindBadRequest)
			}
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	return acc, nil
}

func (repo *accountRepository) Update(ctx context.Context, acc accounts.Account) error {
	const op errors.Op = "store/postgres.accountRepository.Update"

	q := `UPDATE accounts SET name=$1, type=$2, seats=$3, updated_at=$4 WHERE id=$5`

	res, err := repo.Exec(q, acc.Name, acc.Type, acc.NumberOfSeats, acc.UpdatedAt, acc.ID)

	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	if cnt == 0 {
		return errors.E(op, "account not found", errors.KindNotFound)
	}
	return nil
}

func (repo *accountRepository) Retrieve(ctx context.Context, id string) (accounts.Account, error) {
	const op errors.Op = "store/postgres.accountRepository.Retrieve"

	q := `SELECT id, name, type, seats, active, created_at, updated_at FROM accounts WHERE id=$1`

	var account = accounts.Account{}

	if err := repo.QueryRow(q, id).Scan(&id, &account.Name, &account.Type, &account.NumberOfSeats, &account.Active, &account.CreatedAt, &account.UpdatedAt); err != nil {
		empty := accounts.Account{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, "account not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	return account, nil
}

func (repo *accountRepository) List(ctx context.Context, offset, limit uint64) (accounts.AccountPage, error) {
	const op errors.Op = "store/postgres.accountRepository.List"

	q := `SELECT id, name, type, seats, active, created_at, updated_at FROM accounts ORDER BY id LIMIT $1 OFFSET $2;`

	var items = []accounts.Account{}

	rows, err := repo.Query(q, limit, offset)
	if err != nil {
		return accounts.AccountPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		c := accounts.Account{}

		if err := rows.Scan(&c.ID, &c.Name, &c.Type, &c.NumberOfSeats, &c.Active, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return accounts.AccountPage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM accounts;`

	var total uint64
	if err := repo.QueryRow(q).Scan(&total); err != nil {
		return accounts.AccountPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := accounts.AccountPage{
		Accounts: items,
		PageMetadata: accounts.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}
