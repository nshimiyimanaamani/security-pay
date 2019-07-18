package postgres

import (
	"database/sql"
	"log"

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
		INSERT INTO transactions (id, property, 
		amount, method, date_recorded) VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	_, err := str.db.Exec(q, trx.ID, trx.Property, trx.Amount, trx.Method, trx.DateRecorded)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && errDuplicate == pqErr.Code.Name() {
			return "", transactions.ErrConflict
		}
		return "", err
	}
	return trx.ID, nil
}

func (str *transactionStore) RetrieveByID(id string) (transactions.Transaction, error) {
	q := `
		SELECT id, property, 
		amount, method, date_recorded FROM transactions WHERE id = $1
	`

	var trx = transactions.Transaction{}

	err := str.db.QueryRow(q, id).Scan(&trx.ID, &trx.Property, &trx.Amount, &trx.Method, &trx.Method)
	if err != nil {
		pqErr, ok := err.(*pq.Error)

		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return trx, transactions.ErrNotFound
		}
	}
	return trx, nil
}

func (str *transactionStore) RetrieveAll(offset uint64, limit uint64) transactions.TransactionPage {
	q := `
		SELECT id, property, 
		amount, method, date_recorded FROM transactions ORDER BY id LIMIT $1 OFFSET $2
	`

	var items = []transactions.Transaction{}

	rows, err := str.db.Query(q, limit, offset)
	if err != nil {
		//tr.log.Error(fmt.Sprintf("Failed to retrieve transactions due to %s", err))
		log.Printf("Failed to retrieve transactions due to %s", err)
		return transactions.TransactionPage{}
	}
	defer rows.Close()

	for rows.Next() {
		c := transactions.Transaction{}
		if err := rows.Scan(&c.ID, &c.Property, &c.Amount, &c.Method, &c.DateRecorded); err != nil {
			//tr.log.Error(fmt.Sprintf("Failed to retrieve transactions due to %s", err))
			log.Printf("Failed to retrieve transactions due to %s", err)
			return transactions.TransactionPage{}
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM transactions`

	var total uint64
	if err := str.db.QueryRow(q).Scan(&total); err != nil {
		log.Printf("Failed to retrieve transactions due to %s", err)
		return transactions.TransactionPage{}
	}

	page := transactions.TransactionPage{
		Transactions: items,
		PageMetadata: transactions.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page
}

func (str *transactionStore) RetrieveByProperty(property string, offset, limit uint64) transactions.TransactionPage {
	q := `
		SELECT id, property, amount, method, 
		date_recorded FROM transactions WHERE property = $1 ORDER BY id LIMIT $2 OFFSET $3
	`

	var items = []transactions.Transaction{}

	rows, err := str.db.Query(q, property, limit, offset)
	if err != nil {
		//tr.log.Error(fmt.Sprintf("Failed to retrieve transactions due to %s", err))
		log.Printf("Failed to retrieve transactions due to %s", err)
		return transactions.TransactionPage{}
	}
	defer rows.Close()

	for rows.Next() {
		c := transactions.Transaction{}
		if err := rows.Scan(&c.ID, &c.Property, &c.Amount, &c.Method, &c.DateRecorded); err != nil {
			//tr.log.Error(fmt.Sprintf("Failed to retrieve transactions due to %s", err))
			log.Printf("Failed to count transactions due to %s", err)
			return transactions.TransactionPage{}
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM transactions WHERE property = $1`

	var total uint64
	if err := str.db.QueryRow(q, property).Scan(&total); err != nil {
		log.Printf("Failed to count transactions due to %s", err)
		return transactions.TransactionPage{}
	}

	page := transactions.TransactionPage{
		Transactions: items,
		PageMetadata: transactions.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page
}

func (str *transactionStore) RetrieveByMethod(method string, offset, limit uint64) transactions.TransactionPage {
	q := `
		SELECT id, property, amount, method,
		date_recorded FROM transactions WHERE method = $1 ORDER BY id LIMIT $2 OFFSET $3
	`

	var items = []transactions.Transaction{}

	rows, err := str.db.Query(q, method, limit, offset)
	if err != nil {
		//tr.log.Error(fmt.Sprintf("Failed to retrieve transactions due to %s", err))
		log.Printf("Failed to retrieve transactions due to %s", err)
		return transactions.TransactionPage{}
	}
	defer rows.Close()

	for rows.Next() {
		c := transactions.Transaction{}
		if err := rows.Scan(&c.ID, &c.Property, &c.Amount, &c.Method, &c.DateRecorded); err != nil {
			//tr.log.Error(fmt.Sprintf("Failed to retrieve transactions due to %s", err))
			log.Printf("Failed to retrieve transactions due to %s", err)
			return transactions.TransactionPage{}
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM transactions WHERE method = $1`

	var total uint64
	if err := str.db.QueryRow(q, method).Scan(&total); err != nil {
		log.Printf("Failed to count transactions due to %s", err)
		return transactions.TransactionPage{}
	}

	page := transactions.TransactionPage{
		Transactions: items,
		PageMetadata: transactions.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page
}

func (str *transactionStore) RetrieveByMonth(string, uint64, uint64) transactions.TransactionPage {
	return transactions.TransactionPage{}
}

func (str *transactionStore) RetrieveByYear(string, uint64, uint64) transactions.TransactionPage {
	return transactions.TransactionPage{}
}
