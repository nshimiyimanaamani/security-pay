package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/nshimiyimanaamani/paypack-backend/core/auth"
	"github.com/nshimiyimanaamani/paypack-backend/core/transactions"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

var _ (transactions.Repository) = (*transactionsStore)(nil)

type transactionsStore struct {
	*sql.DB
}

// NewTransactionRepository instanctiates a new transactions.Repository interface
func NewTransactionRepository(db *sql.DB) transactions.Repository {
	return &transactionsStore{db}
}

func (repo *transactionsStore) Save(ctx context.Context, tx transactions.Transaction) (string, error) {
	const op errors.Op = "store/postgres/transactionsRepository.Save"

	q := `
		INSERT INTO transactions (
			id, 
			madefor, 
			madeby, 
			amount, 
			method,
			invoice,
			namespace
		) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING created_at
	`

	err := repo.QueryRow(q,
		tx.ID,
		tx.MadeFor,
		tx.OwnerID,
		tx.Amount,
		tx.Method,
		tx.Invoice,
		tx.Namespace,
	).Scan(&tx.DateRecorded)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return "", errors.E(op, "transaction already exists", errors.KindAlreadyExists)
			case errInvalid, errTruncation:
				return "", errors.E(op, err, "invalid transaction", errors.KindBadRequest)
			}
		}
		return "", errors.E(op, err, errors.KindUnexpected)
	}
	return tx.ID, nil
}

// seletect tx[id, amount, method, recorded]; properties[sector, cell, village] owner[fname, lname]
func (repo *transactionsStore) RetrieveByID(ctx context.Context, id string) (transactions.Transaction, error) {
	const op errors.Op = "store/postgres/transactionsRepository.RetrieveByID"

	q := `
		SELECT 
			transactions.id, 
			transactions.amount, 
			transactions.method, 
			transactions.madefor,
			transactions.invoice,
			transactions.created_at, 
			properties.sector, 
			properties.cell, 
			properties.village, 
			owners.id, 
			owners.fname, 
			owners.lname
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
		&tx.ID,
		&tx.Amount,
		&tx.Method,
		&tx.MadeFor,
		&tx.Invoice,
		// &tx.Namespace,
		&tx.DateRecorded,
		&tx.Sector,
		&tx.Cell,
		&tx.Village,
		&tx.OwnerID,
		&tx.OwneFname,
		&tx.OwnerLname,
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
func (repo *transactionsStore) RetrieveAll(ctx context.Context, offset uint64, limit uint64) (transactions.TransactionPage, error) {
	const op errors.Op = "store/postgres/transactionsRepository.RetrieveAll"

	txn, err := repo.DB.BeginTx(ctx, nil)
	if err != nil {
		return transactions.TransactionPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer txn.Rollback()

	selectQuery := `
	SELECT 
		transactions.id, transactions.amount,transactions.method, transactions.madefor,
		transactions.invoice, transactions.created_at, properties.sector, properties.cell, 
		properties.village, owners.id, owners.fname, owners.lname
	FROM 
		transactions
	INNER JOIN 
		properties ON transactions.madefor=properties.id
	INNER JOIN 
		owners ON transactions.madeby=owners.id 
	WHERE
		1 = 1
	`
	creds := auth.CredentialsFromContext(ctx)

	selectQuery += fmt.Sprintf(" AND transactions.namespace='%s'", creds.Account)

	if creds.Role == "min" {
		var cell, sector, village string
		err := txn.QueryRow(`
			SELECT
				cell, sector, village
			FROM
				agents
			WHERE
			   telephone=$1
		`, creds.Username).Scan(&cell, &sector, &village)
		if err != nil {
			return transactions.TransactionPage{}, errors.E(op, err, errors.KindUnexpected)
		}
		selectQuery += fmt.Sprintf(" AND properties.sector='%s' AND properties.cell='%s' AND properties.village='%s'", sector, cell, village)
	}

	selectQuery += fmt.Sprintf(" ORDER BY transactions.created_at DESC LIMIT %d OFFSET %d", limit, offset)

	rows, err := txn.QueryContext(ctx, selectQuery)
	if err != nil {
		return transactions.TransactionPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	out := []transactions.Transaction{}

	for rows.Next() {
		tx := transactions.Transaction{}
		err := rows.Scan(
			&tx.ID,
			&tx.Amount,
			&tx.Method,
			&tx.MadeFor,
			&tx.Invoice,
			&tx.DateRecorded,
			&tx.Sector,
			&tx.Cell,
			&tx.Village,
			&tx.OwnerID,
			&tx.OwneFname,
			&tx.OwnerLname,
		)
		if err != nil {
			return transactions.TransactionPage{}, errors.E(op, err, errors.KindUnexpected)
		}
		tx.Fees = tx.Amount * transactions.RateFee
		out = append(out, tx)
	}

	countQuery := `SELECT 
		count(transactions.id) as total
	FROM 
		transactions
	INNER JOIN 
		properties ON transactions.madefor=properties.id
	INNER JOIN 
		owners ON transactions.madeby=owners.id 
	WHERE
		1 = 1
		`

	countQuery += fmt.Sprintf(" AND transactions.namespace='%s'", creds.Account)

	if creds.Role == "min" {
		var cell, sector, village string
		err := txn.QueryRow(`
			SELECT
				cell, sector, village
			FROM
				agents
			WHERE
			   telephone=$1
		`, creds.Username).Scan(&cell, &sector, &village)
		if err != nil {
			return transactions.TransactionPage{}, errors.E(op, err, errors.KindUnexpected)
		}
		countQuery += fmt.Sprintf(" AND properties.sector='%s' AND properties.cell='%s' AND properties.village='%s'", sector, cell, village)
	}

	var total uint64
	if err := txn.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return transactions.TransactionPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := transactions.TransactionPage{
		Transactions: out,
		PageMetadata: transactions.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, txn.Commit()
}

func (repo *transactionsStore) RetrieveByProperty(ctx context.Context, property string, offset, limit uint64) (transactions.TransactionPage, error) {
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
	AND 
		transactions.namespace=$2
	ORDER BY 
		transactions.id LIMIT $3 OFFSET $4
	`

	empty := transactions.TransactionPage{}

	var items = []transactions.Transaction{}

	creds := auth.CredentialsFromContext(ctx)

	rows, err := repo.Query(q, property, creds.Account, limit, offset)
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

	q = `SELECT COUNT(*) FROM transactions WHERE madefor = $1 AND namespace=$2`

	var total uint64
	if err := repo.QueryRow(q, property, creds.Account).Scan(&total); err != nil {
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

func (repo *transactionsStore) RetrieveByMethod(ctx context.Context, method string, offset, limit uint64) (transactions.TransactionPage, error) {
	const op errors.Op = "store/postgres/transactionsRepository.RetrieveByMethod"

	q := `
		SELECT 
			transactions.id,
			transactions.amount, 
			transactions.method, 
			transactions.madefor,
			transactions.created_at, 
			properties.sector, 
			properties.cell, 
			properties.village, 
			owners.id, 
			owners.fname, 
			owners.lname
		FROM 
			transactions
		INNER JOIN 
			properties ON transactions.madefor=properties.id
		INNER JOIN 
			owners ON transactions.madeby=owners.id
		WHERE 
			transactions.method = $1 
		AND 
			transactions.namespace
		ORDER BY 
			transactions.id LIMIT $2 OFFSET $3
	`
	empty := transactions.TransactionPage{}

	var items = []transactions.Transaction{}

	creds := auth.CredentialsFromContext(ctx)

	rows, err := repo.Query(q, method, creds.Account, limit, offset)
	if err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		c := transactions.Transaction{}

		if err := rows.Scan(
			&c.ID,
			&c.Amount,
			&c.Method,
			&c.MadeFor,
			&c.DateRecorded,
			&c.Sector,
			&c.Cell,
			&c.Village,
			&c.OwnerID,
			&c.OwneFname,
			&c.OwnerLname,
		); err != nil {
			return empty, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM transactions WHERE method = $1 AND namespace=$2`

	var total uint64

	if err := repo.QueryRow(q, method, creds.Account).Scan(&total); err != nil {
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

func (repo *transactionsStore) RetrieveByPropertyR(ctx context.Context, property string) (transactions.TransactionPage, error) {
	const op errors.Op = "store/postgres/transactionsRepository.RetrieveByProperty"

	q := `
	SELECT 
		transactions.id, 
		transactions.amount, 
		transactions.method, 
		transactions.madefor, 
		transactions.invoice, 
		transactions.created_at, 
		properties.sector, 
		properties.cell, 
		properties.village, 
		owners.id, 
		owners.fname,
		owners.lname
	FROM 
		transactions
	INNER JOIN 
		properties ON transactions.madefor=properties.id
	INNER JOIN 
		owners ON transactions.madeby=owners.id 
	WHERE 
		transactions.madefor=$1
	ORDER BY transactions.id
`

	empty := transactions.TransactionPage{}

	var items = []transactions.Transaction{}

	// creds := auth.CredentialsFromContext(ctx)

	rows, err := repo.Query(q, property)
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

	q = `SELECT COUNT(*) FROM transactions WHERE madefor=$1`

	var total uint64
	if err := repo.QueryRow(q, property).Scan(&total); err != nil {
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	page := transactions.TransactionPage{
		Transactions: items,
		PageMetadata: transactions.PageMetadata{
			Total: total,
		},
	}
	return page, nil
}
