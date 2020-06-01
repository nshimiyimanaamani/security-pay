package postgres_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/accounts"
	"github.com/rugwirobaker/paypack-backend/core/invoices"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/transactions"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/stretchr/testify/require"
)

func saveOwner(t *testing.T, db *sql.DB, owner properties.Owner) properties.Owner {
	t.Helper()

	q := `
		INSERT INTO owners (
			id, 
			fname, 
			lname, 
			phone,
			namespace
		) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	_, err := db.Exec(q,
		&owner.ID,
		&owner.Fname,
		&owner.Lname,
		&owner.Phone,
		&owner.Namespace,
	)

	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return owner
}

func saveTx(t *testing.T, db *sql.DB, tx transactions.Transaction) transactions.Transaction {
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

	err := db.QueryRow(q, tx.ID, tx.MadeFor, tx.OwnerID, tx.Amount, tx.Method, tx.Invoice).Scan(&tx.DateRecorded)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return tx
}

func saveAgent(t *testing.T, db *sql.DB, agent users.Agent) users.Agent {
	t.Helper()

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
	_, err := db.Exec(q, agent.Telephone, agent.FirstName, agent.Role,
		agent.Account, agent.CreatedAt, agent.UpdatedAt)

	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return agent
}

func saveAccount(t *testing.T, db *sql.DB, acc accounts.Account) accounts.Account {
	t.Helper()

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

	_, err := db.Exec(q, acc.ID, acc.Name, acc.Type, acc.NumberOfSeats, acc.CreatedAt, acc.UpdatedAt)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return acc
}

func saveProperty(t *testing.T, db *sql.DB, pp properties.Property) properties.Property {
	t.Helper()

	q := `
		INSERT INTO properties (
			id, 
			owner, 
			due,
			sector, 
			cell, 
			village, 
			recorded_by, 
			occupied,
			namespace
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING  created_at, updated_at`

	err := db.QueryRow(q, pp.ID, pp.Owner.ID, pp.Due, pp.Address.Sector,
		pp.Address.Cell, pp.Address.Village, pp.RecordedBy, pp.Occupied, pp.Namespace).Scan(&pp.CreatedAt, &pp.UpdatedAt)

	if err != nil {
		t.Fatalf("err:%v", err)
	}
	return pp
}

func saveInvoice(t *testing.T, db *sql.DB, inv invoices.Invoice) invoices.Invoice {
	t.Helper()

	q := `
		INSERT INTO invoices (
			amount, 
			property, 
			status,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`

	err := db.QueryRow(q, inv.Amount, inv.Property, inv.Status, inv.CreatedAt, inv.UpdatedAt).Scan(&inv.ID)

	if err != nil {
		t.Fatalf("err: %v", err)
	}
	return inv
}

func savePropertyOn(t *testing.T, db *sql.DB, pro properties.Property) properties.Property {
	t.Helper()

	q := `
		INSERT INTO properties (
			id, 
			owner, 
			due,
			sector, 
			cell, 
			village, 
			recorded_by, 
			occupied,
			namespace,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := db.Exec(q,
		pro.ID,
		pro.Owner.ID,
		pro.Due,
		pro.Address.Sector,
		pro.Address.Cell,
		pro.Address.Village,
		pro.RecordedBy,
		pro.Occupied,
		pro.Namespace,
		pro.CreatedAt,
		pro.UpdatedAt,
	)

	if err != nil {
		t.Fatalf("err: %v", err)
	}
	return pro
}

func retrieveInvoice(t *testing.T, db *sql.DB, id string) invoices.Invoice {
	var val invoices.Invoice

	q := `
		SELECT 
			id,
			property,
			amount,
			status,
			created_at,
			updated_at
		FROM 
			invoices
		WHERE property=$1
	`
	err := db.QueryRow(q, id).Scan(
		&val.ID, &val.Property, &val.Amount, &val.Status, &val.CreatedAt, &val.UpdatedAt)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	return val
}

func CleanDB(t *testing.T, db *sql.DB) {
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
			accounts,
			villages,
			cells,
			sectors
	`

	_, err := db.Exec(q)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
}
