package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
)

var _ (transactions.Store) = (*transactionStore)(nil)

type transactionStore struct {
	db *sql.DB
}

//NewTransactionStore instanctiates a new transactiob store interface
func NewTransactionStore(db *sql.DB) transactions.Store {
	return &transactionStore{db}
}

func (str *transactionStore) Save(trx transactions.Transaction) (string, error) {
	q := `
		INSERT INTO transactions (id, madefor, madeby, 
		amount, method, date_recorded) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`

	_, err := str.db.Exec(q, trx.ID, trx.MadeFor, trx.MadeBy, trx.Amount, trx.Method, trx.DateRecorded)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return "", transactions.ErrConflict
			case errInvalid, errTruncation:
				return "", transactions.ErrInvalidEntity
			}
		}
		return "", err
	}
	return trx.ID, nil
}

//seletect trx[id, amount, method, recorded]; properties[sector, cell, village] owner[fname, lname]
func (str *transactionStore) RetrieveByID(id string) (transactions.Transaction, error) {
	q := `
		SELECT 
			transactions.id, transactions.amount, transactions.method, 
			transactions.date_recorded, properties.sector, properties.cell, 
			properties.village, owners.fname, owners.lname
		FROM 
			transactions
		INNER JOIN 
			properties ON transactions.madefor=properties.id
		INNER JOIN 
			owners ON transactions.madeby=owners.id
		WHERE transactions.id = $1
	`

	var trx = transactions.Transaction{
		Address: make(map[string]string),
	}

	var sector, cell, village string

	var fname, lname string

	err := str.db.QueryRow(q, id).Scan(
		&trx.ID, &trx.Amount, &trx.Method, &trx.DateRecorded,
		&sector, &cell, &village, &fname, &lname,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)

		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return trx, transactions.ErrNotFound
		}
		return trx, err
	}
	if sector != "" && cell != "" && village != "" {
		trx.Address["sector"] = sector
		trx.Address["cell"] = cell
		trx.Address["village"] = village
	}
	if fname != "" && lname != "" {
		trx.MadeBy = fmt.Sprintf("%s %s", fname, lname)
	}
	return trx, nil
}

func (str *transactionStore) RetrieveAll(offset uint64, limit uint64) (transactions.TransactionPage, error) {
	q := `
	SELECT 
		transactions.id, transactions.amount, transactions.method, 
		transactions.date_recorded, properties.sector, properties.cell, 
		properties.village, owners.fname, owners.lname
	FROM 
		transactions
	INNER JOIN 
		properties ON transactions.madefor=properties.id
	INNER JOIN 
		owners ON transactions.madeby=owners.id 
	ORDER BY transactions.id LIMIT $1 OFFSET $2
	`

	var items = []transactions.Transaction{}

	rows, err := str.db.Query(q, limit, offset)
	if err != nil {
		return transactions.TransactionPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := transactions.Transaction{
			Address: make(map[string]string),
		}

		var sector, cell, village string

		var fname, lname string

		if err := rows.Scan(
			&c.ID, &c.Amount, &c.Method, &c.DateRecorded,
			&sector, &cell, &village, &fname, &lname,
		); err != nil {
			return transactions.TransactionPage{}, err
		}
		if sector != "" && cell != "" && village != "" {
			c.Address["sector"] = sector
			c.Address["cell"] = cell
			c.Address["village"] = village
		}
		if fname != "" && lname != "" {
			c.MadeBy = fmt.Sprintf("%s %s", fname, lname)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM transactions`

	var total uint64
	if err := str.db.QueryRow(q).Scan(&total); err != nil {
		return transactions.TransactionPage{}, err
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

func (str *transactionStore) RetrieveByProperty(property string, offset, limit uint64) (transactions.TransactionPage, error) {
	q := `
	SELECT 
		transactions.id, transactions.amount, transactions.method, 
		transactions.date_recorded, properties.sector, properties.cell, 
		properties.village, owners.fname, owners.lname
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

	var items = []transactions.Transaction{}

	rows, err := str.db.Query(q, property, limit, offset)
	if err != nil {
		return transactions.TransactionPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := transactions.Transaction{
			Address: make(map[string]string),
		}
		var sector, cell, village string

		var fname, lname string

		if err := rows.Scan(
			&c.ID, &c.Amount, &c.Method, &c.DateRecorded,
			&sector, &cell, &village, &fname, &lname,
		); err != nil {
			return transactions.TransactionPage{}, err
		}
		if sector != "" && cell != "" && village != "" {
			c.Address["sector"] = sector
			c.Address["cell"] = cell
			c.Address["village"] = village
		}
		if fname != "" && lname != "" {
			c.MadeBy = fmt.Sprintf("%s %s", fname, lname)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM transactions WHERE madefor = $1`

	var total uint64
	if err := str.db.QueryRow(q, property).Scan(&total); err != nil {
		return transactions.TransactionPage{}, err
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

func (str *transactionStore) RetrieveByMethod(method string, offset, limit uint64) (transactions.TransactionPage, error) {
	q := `
		SELECT 
			transactions.id,transactions.amount, transactions.method, 
			transactions.date_recorded, properties.sector, properties.cell, 
			properties.village, owners.fname, owners.lname
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

	var items = []transactions.Transaction{}

	rows, err := str.db.Query(q, method, limit, offset)
	if err != nil {
		return transactions.TransactionPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := transactions.Transaction{
			Address: make(map[string]string),
		}
		var sector, cell, village string

		var fname, lname string

		if err := rows.Scan(
			&c.ID, &c.Amount, &c.Method, &c.DateRecorded,
			&sector, &cell, &village, &fname, &lname,
		); err != nil {
			return transactions.TransactionPage{}, err
		}
		if sector != "" && cell != "" && village != "" {
			c.Address["sector"] = sector
			c.Address["cell"] = cell
			c.Address["village"] = village
		}
		if fname != "" && lname != "" {
			c.MadeBy = fmt.Sprintf("%s %s", fname, lname)
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM transactions WHERE method = $1`

	var total uint64
	if err := str.db.QueryRow(q, method).Scan(&total); err != nil {
		return transactions.TransactionPage{}, err
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

func (str *transactionStore) RetrieveByMonth(string, uint64, uint64) (transactions.TransactionPage, error) {
	return transactions.TransactionPage{}, nil
}

func (str *transactionStore) RetrieveByYear(string, uint64, uint64) (transactions.TransactionPage, error) {
	return transactions.TransactionPage{}, nil
}

func (str *transactionStore) UpdateTransaction(tx transactions.Transaction) error {
	q := `UPDATE transactions SET date_modified=$1, is_valid=TRUE WHERE id=$2;`

	res, err := str.db.Exec(q, tx.DateRecorded, tx.ID)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return transactions.ErrInvalidEntity
			}
		}
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return transactions.ErrNotFound
	}
	return nil
}
