package postgres

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/properties"
)

//OwnerStore store is a postgres implementation of the properties.OwnerStore
type ownerStore struct {
	db *sql.DB
}

var _ (properties.OwnerStore) = (*ownerStore)(nil)

//NewOwnerStore creates an instance of OwnerStore.
func NewOwnerStore(db *sql.DB) properties.OwnerStore {
	return &ownerStore{db}
}

func (str *ownerStore) Save(owner properties.Owner) (string, error) {
	q := `INSERT INTO owners (id, fname, lname, phone, password) VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	_, err := str.db.Exec(q, &owner.ID, &owner.Fname, &owner.Lname, &owner.Phone, &owner.Password)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return "", properties.ErrConflict
			case errInvalid, errTruncation:
				return "", properties.ErrInvalidEntity
			}
		}
		return "", err
	}
	return owner.ID, nil
}

func (str *ownerStore) Update(owner properties.Owner) error {
	q := `UPDATE owners SET fname=$1, lname=$2, phone=$3, password=$4 WHERE id=$5;`

	res, err := str.db.Exec(q, owner.Fname, owner.Lname, owner.Phone, owner.Password, owner.ID)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return properties.ErrInvalidEntity
			}
		}
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return properties.ErrNotFound
	}
	return nil
}

func (str *ownerStore) Retrieve(id string) (properties.Owner, error) {
	q := `SELECT id, fname, lname, phone, password FROM owners WHERE id = $1`

	var owner properties.Owner

	if err := str.db.QueryRow(q, id).Scan(&owner.ID, &owner.Fname, &owner.Lname, &owner.Phone, &owner.Password); err != nil {
		empty := properties.Owner{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, properties.ErrNotFound
		}
		return empty, err
	}
	return owner, nil
}

func (str *ownerStore) FindOwner(fname, lname, phone string) (properties.Owner, error) {
	q := `SELECT id, fname, lname, phone FROM owners WHERE fname=$1 AND lname=$2 AND phone=$3;`

	var owner properties.Owner

	err := str.db.QueryRow(q, fname, lname, phone).Scan(&owner.ID, &owner.Fname, &owner.Lname, &owner.Phone)
	if err != nil {
		empty := properties.Owner{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, properties.ErrNotFound
		}
		return empty, err
	}
	return owner, nil
}

func (str *ownerStore) RetrieveAll(offset, limit uint64) (properties.OwnerPage, error) {
	q := `SELECT id, fname, lname, phone FROM owners ORDER BY id LIMIT $1 OFFSET $2;`

	var items = []properties.Owner{}

	rows, err := str.db.Query(q, limit, offset)
	if err != nil {
		return properties.OwnerPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Owner{}

		if err := rows.Scan(&c.ID, &c.Fname, &c.Lname, &c.Phone); err != nil {
			return properties.OwnerPage{}, err
		}
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM owners;`

	var total uint64
	if err := str.db.QueryRow(q).Scan(&total); err != nil {
		return properties.OwnerPage{}, err
	}

	page := properties.OwnerPage{
		Owners: items,
		PageMetadata: properties.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}

func (str *ownerStore) RetrieveByPhone(phone string) (properties.Owner, error) {
	q := `SELECT id, fname, lname, phone, password FROM owners WHERE phone = $1`
	var owner properties.Owner

	if err := str.db.QueryRow(q, phone).Scan(&owner.ID, &owner.Fname, &owner.Lname, &owner.Phone, &owner.Password); err != nil {
		empty := properties.Owner{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, properties.ErrNotFound
		}
		return empty, err
	}
	return owner, nil
}
