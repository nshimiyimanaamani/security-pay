package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (transactions.Repository) = (*txRepository)(nil)

type txRepository struct {
	*sql.DB
}

// NewTransactionRepository instanctiates a new transactions.Repository interface
func NewTransactionRepository(db *sql.DB) transactions.Repository {
	return &txRepository{db}
}

func (repo *txRepository) Save(ctx context.Context, tx transactions.Transaction) (string, error) {
	const op errors.Op = "store/postgres/transactionsRepository.Save"

	q := `
		INSERT INTO transactions (
			id, 
			madefor, 
			madeby, 
			amount, 
			method,
			invoice
		) VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at
	`

	err := repo.QueryRow(q, tx.ID, tx.MadeFor, tx.OwnerID, tx.Amount, tx.Method).Scan(&tx.DateRecorded)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return "", errors.E(op, err, "transaction already exists", errors.KindAlreadyExists)
			case errInvalid, errTruncation:
				return "", errors.E(op, err, "invalid transaction", errors.KindBadRequest)
			}
		}
		return "", errors.E(op, err, errors.KindUnexpected)
	}
	return tx.ID, nil
}

//seletect tx[id, amount, method, recorded]; properties[sector, cell, village] owner[fname, lname]
func (repo *txRepository) RetrieveByID(ctx context.Context, id string) (transactions.Transaction, error) {
	const op errors.Op = "store/postgres/transactionsRepository.RetrieveByID"

	q := `
		SELECT 
			transactions.id, transactions.amount, transactions.method, transactions.madefor,
			transactions.created_at, properties.sector, properties.cell, 
			properties.village, owners.id, owners.fname, owners.lname
		FROM 
			transactions
		INNER JOIN 
			properties ON transactions.madefor=properties.id
		INNER JOIN 
			owners ON transactions.madeby=owners.id
		WHERE transactions.id = $1
	`

	var tx = transactions.Transaction{}

	err := repo.QueryRow(q, id).Scan(
		&tx.ID, &tx.Amount, &tx.Method, &tx.MadeFor, &tx.DateRecorded,
		&tx.Sector, &tx.Cell, &tx.Village, &tx.OwnerID, &tx.OwneFname, &tx.OwnerLname,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)

		empty := transactions.Transaction{}

		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, err, "transaction not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	return tx, nil
}

func (repo *txRepository) RetrieveAll(ctx context.Context, offset uint64, limit uint64) (transactions.TransactionPage, error) {
	const op errors.Op = "store/postgres/transactionsRepository.RetrieveAll"

	q := `
	SELECT 
		transactions.id, transactions.amount, transactions.method, transactions.madefor,
		transactions.created_at, properties.sector, properties.cell, 
		properties.village, owners.id, owners.fname, owners.lname
	FROM 
		transactions
	INNER JOIN 
		properties ON transactions.madefor=properties.id
	INNER JOIN 
		owners ON transactions.madeby=owners.id 
	ORDER BY transactions.id LIMIT $1 OFFSET $2
	`
	empty := transactions.TransactionPage{}

	var items = []transactions.Transaction{}

	rows, err := repo.Query(q, limit, offset)
	if err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		c := transactions.Transaction{}

		if err := rows.Scan(
			&c.ID, &c.Amount, &c.Method, &c.MadeFor, &c.DateRecorded,
			&c.Sector, &c.Cell, &c.Village, &c.OwnerID, &c.OwneFname, &c.OwnerLname,
		); err != nil {
			return empty, errors.E(op, err, errors.KindUnexpected)
		}

		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM transactions`

	var total uint64
	if err := repo.QueryRow(q).Scan(&total); err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	page := transactions.TransactionPage{
		Transactions: items,
		PageMetadata: transactions.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}

func (repo *txRepository) RetrieveByProperty(ctx context.Context, property string, offset, limit uint64) (transactions.TransactionPage, error) {
	const op errors.Op = "store/postgres/transactionsRepository.RetrieveByProperty"

	q := `
	SELECT 
		transactions.id, transactions.amount, transactions.method, transactions.madefor, 
		transactions.invoice, transactions.created_at, properties.sector, properties.cell, 
		properties.village, owners.id, owners.fname, owners.lname
	FROM 
		transactions
	INNER JOIN 
		properties ON transactions.madefor=properties.id
	INNER JOIN 
		owners ON transactions.madeby=owners.id 
	WHERE 
		transactions.madefor = $1 
	ORDER BY transactions.id LIMIT $2 OFFSET $3
	`

	empty := transactions.TransactionPage{}

	var items = []transactions.Transaction{}

	rows, err := repo.Query(q, property, limit, offset)
	if err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		c := transactions.Transaction{}

		if err := rows.Scan(
			&c.ID, &c.Amount, &c.Method, &c.MadeFor, &c.Invoice, &c.DateRecorded,
			&c.Sector, &c.Cell, &c.Village, &c.OwnerID, &c.OwneFname, &c.OwnerLname,
		); err != nil {
			return empty, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM transactions WHERE madefor = $1`

	var total uint64
	if err := repo.QueryRow(q, property).Scan(&total); err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	page := transactions.TransactionPage{
		Transactions: items,
		PageMetadata: transactions.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, err
}

func (repo *txRepository) RetrieveByMethod(ctx context.Context, method string, offset, limit uint64) (transactions.TransactionPage, error) {
	const op errors.Op = "store/postgres/transactionsRepository.RetrieveByMethod"

	q := `
		SELECT 
			transactions.id,transactions.amount, transactions.method, transactions.madefor,
			transactions.created_at, properties.sector, properties.cell, 
			properties.village, owners.id, owners.fname, owners.lname
		FROM 
			transactions
		INNER JOIN 
			properties ON transactions.madefor=properties.id
		INNER JOIN 
			owners ON transactions.madeby=owners.id
		WHERE 
			transactions.method = $1 
		ORDER BY transactions.id LIMIT $2 OFFSET $3
	`
	empty := transactions.TransactionPage{}

	var items = []transactions.Transaction{}

	rows, err := repo.Query(q, method, limit, offset)
	if err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		c := transactions.Transaction{}

		if err := rows.Scan(
			&c.ID, &c.Amount, &c.Method, &c.MadeFor, &c.DateRecorded,
			&c.Sector, &c.Cell, &c.Village, &c.OwnerID, &c.OwneFname, &c.OwnerLname,
		); err != nil {
			return empty, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM transactions WHERE method = $1`

	var total uint64
	if err := repo.QueryRow(q, method).Scan(&total); err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	page := transactions.TransactionPage{
		Transactions: items,
		PageMetadata: transactions.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}
