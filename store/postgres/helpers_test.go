package postgres_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/accounts"
	"github.com/rugwirobaker/paypack-backend/app/invoices"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/stretchr/testify/require"
)

func saveOwner(t *testing.T, db *sql.DB, owner properties.Owner) (properties.Owner, error) {
	t.Helper()

	q := `INSERT INTO owners (id, fname, lname, phone) VALUES ($1, $2, $3, $4) RETURNING id;`

	_, err := db.Exec(q, &owner.ID, &owner.Fname, &owner.Lname, &owner.Phone)
	if err != nil {
		return properties.Owner{}, err
	}

	return owner, nil
}

func saveTx(t *testing.T, db *sql.DB, tx transactions.Transaction) (transactions.Transaction, error) {
	t.Helper()

	q := `
		INSERT INTO transactions (
			id, 
			madefor, 
			madeby, 
			amount,
			method, 
			invoice
		) VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at;
	`

	var empty = transactions.Transaction{}

	err := db.QueryRow(q, tx.ID, tx.MadeFor, tx.OwnerID, tx.Amount, tx.Method, tx.Invoice).Scan(&tx.DateRecorded)
	if err != nil {
		return empty, err
	}
	return tx, nil
}

func saveAgent(t *testing.T, db *sql.DB, agent users.Agent) (users.Agent, error) {

	q := `
		INSERT into users (
			username, 
			password, 
			role, 
			account, 
			created_at, 
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6) RETURNING username
	`
	empty := users.Agent{}

	if _, err := db.Exec(q, agent.Telephone, agent.FirstName, agent.Role,
		agent.Account, agent.CreatedAt, agent.UpdatedAt); err != nil {
		return empty, err
	}
	return agent, nil
}

func saveAccount(t *testing.T, db *sql.DB, acc accounts.Account) (accounts.Account, error) {
	q := `
		INSERT INTO accounts (
			id, 
			name, 
			type, 
			seats, 
			created_at, 
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		) RETURNING id;`

	empty := accounts.Account{}

	_, err := db.Exec(q, acc.ID, acc.Name, acc.Type, acc.NumberOfSeats, acc.CreatedAt, acc.UpdatedAt)
	if err != nil {
		return empty, err
	}
	return acc, nil
}

func saveProperty(t *testing.T, db *sql.DB, pp properties.Property) (properties.Property, error) {
	q := `
		INSERT INTO properties (
			id, 
			owner, 
			due,
			sector, 
			cell, 
			village, 
			recorded_by, 
			occupied
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING  created_at, updated_at`

	err := db.QueryRow(q, pp.ID, pp.Owner.ID, pp.Due, pp.Address.Sector,
		pp.Address.Cell, pp.Address.Village, pp.RecordedBy, pp.Occupied).Scan(&pp.CreatedAt, &pp.UpdatedAt)

	if err != nil {
		return properties.Property{}, err
	}
	return pp, nil
}

func saveInvoice(t *testing.T, db *sql.DB, inv invoices.Invoice) (invoices.Invoice, error) {
	q := `
		INSERT INTO invoices (
			amount, 
			property, 
			status
		) VALUES (
			$1, $2, $3
	) RETURNING id, amount, property, status, created_at, updated_at;
	`

	err := db.QueryRow(q, inv.Amount, inv.Property, inv.Status).
		Scan(&inv.ID, &inv.Amount, &inv.Property, &inv.Status, &inv.CreatedAt, &inv.UpdatedAt)

	if err != nil {
		return invoices.Invoice{}, err
	}
	return inv, nil
}

func CleanDB(t *testing.T) {
	q := `
		TRUNCATE TABLE 
			messages, 
			transactions, 
			invoices, 
			properties,
			owners, 
			agents, 
			managers, 
			admins, 
			developers,
			users,
			accounts
	`

	_, err := db.Exec(q)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
}
