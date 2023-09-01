package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/nshimiyimanaamani/paypack-backend/core/owners"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

// OwnerStore store is a postgres implementation of the owners.OwnerStore
type ownerRepo struct {
	db *sql.DB
}

var _ (owners.Repository) = (*ownerRepo)(nil)

// NewOwnerRepo creates an instance of OwnerStore.
func NewOwnerRepo(db *sql.DB) owners.Repository {
	return &ownerRepo{db}
}

func (str *ownerRepo) Save(ctx context.Context, owner owners.Owner) (owners.Owner, error) {
	q := `INSERT INTO owners (
			id, 
			fname, 
			lname, 
			phone
		) VALUES ($1, $2, $3, $4) RETURNING id;`

	_, err := str.db.Exec(q,
		&owner.ID,
		&owner.Fname,
		&owner.Lname,
		&owner.Phone,
	)

	if err != nil {
		empty := owners.Owner{}

		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return empty, owners.ErrConflict
			case errInvalid, errTruncation:
				return empty, owners.ErrInvalidEntity
			}
		}
		return empty, err
	}
	return owner, nil
}

func (str *ownerRepo) Update(ctx context.Context, owner owners.Owner) error {
	q := `UPDATE owners SET fname=$1, lname=$2, phone=$3 WHERE id=$4;`

	res, err := str.db.Exec(q, owner.Fname, owner.Lname, owner.Phone, owner.ID)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return owners.ErrInvalidEntity
			}
		}
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return owners.ErrNotFound
	}
	return nil
}

func (str *ownerRepo) Retrieve(ctx context.Context, id string) (owners.Owner, error) {
	q := `SELECT id, fname, lname, phone FROM owners WHERE id = $1`

	var owner owners.Owner

	if err := str.db.QueryRow(q, id).Scan(&owner.ID, &owner.Fname, &owner.Lname, &owner.Phone); err != nil {
		empty := owners.Owner{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, owners.ErrNotFound
		}
		return empty, err
	}
	return owner, nil
}

func (str *ownerRepo) Search(ctx context.Context, owner owners.Owner) (owners.Owner, error) {
	q := `SELECT id, fname, lname, phone FROM owners WHERE phone=$1;`

	err := str.db.QueryRow(q, owner.Phone).Scan(&owner.ID, &owner.Fname, &owner.Lname, &owner.Phone)
	if err != nil {
		empty := owners.Owner{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, owners.ErrNotFound
		}
		return empty, err
	}
	return owner, nil
}

func (str *ownerRepo) RetrieveAll(ctx context.Context, offset, limit uint64) (owners.OwnerPage, error) {
	q := `
		SELECT 
			id, 
			fname, 
			lname, 
			phone 
		FROM 
			owners 
		ORDER BY id LIMIT $1 OFFSET $2;`

	var items = []owners.Owner{}

	rows, err := str.db.Query(q, limit, offset)
	if err != nil {
		return owners.OwnerPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := owners.Owner{}

		if err := rows.Scan(&c.ID, &c.Fname, &c.Lname, &c.Phone); err != nil {
			return owners.OwnerPage{}, err
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM owners;`

	var total uint64
	if err := str.db.QueryRow(q).Scan(&total); err != nil {
		return owners.OwnerPage{}, err
	}

	page := owners.OwnerPage{
		Owners: items,
		PageMetadata: owners.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}

func (str *ownerRepo) RetrieveByPhone(ctx context.Context, phone string) (owners.Owner, error) {

	const op errors.Op = "store/postgres/RetrieveByPhone"

	q := `
		SELECT 
			id, 
			fname, 
			lname, 
			phone 
		FROM 
			owners 
		WHERE 
			phone = $1`

	var owner owners.Owner

	if err := str.db.QueryRow(q, phone).Scan(&owner.ID, &owner.Fname, &owner.Lname, &owner.Phone); err != nil {
		empty := owners.Owner{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, owners.ErrNotFound, errors.KindNotFound)
		}
		return empty, err
	}
	return owner, nil
}
