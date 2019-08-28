package postgres

import (
	"database/sql"
	"fmt"

	//"github.com/lib/pq"
	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/properties"
)

var _ (properties.PropertyStore) = (*propertiesStore)(nil)

type propertiesStore struct {
	db *sql.DB
}

// NewPropertyStore instanctiates a new transactiob store interface
func NewPropertyStore(db *sql.DB) properties.PropertyStore {
	return &propertiesStore{db}
}

func (str *propertiesStore) Save(pro properties.Property) (string, error) {
	q := `INSERT INTO properties (id, owner, due, sector, cell, village) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	_, err := str.db.Exec(q, pro.ID, pro.Owner, pro.Due, pro.Sector, pro.Cell, pro.Village)
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

	return pro.ID, nil
}

func (str *propertiesStore) UpdateProperty(pro properties.Property) error {
	q := `UPDATE properties SET owner=$1, due=$2 WHERE id=$3;`

	res, err := str.db.Exec(q, pro.Owner, pro.Due, pro.ID)
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

func (str *propertiesStore) RetrieveByID(id string) (properties.Property, error) {
	q := `
		SELECT 
			properties.id, properties.sector, properties.cell,  
			properties.village, properties.due, owners.fname, owners.lname 
		FROM 
			properties
		INNER JOIN 
			owners ON properties.owner=owners.id 
		WHERE properties.id = $1
	`

	var prt = properties.Property{}

	var fname, lname string

	if err := str.db.QueryRow(q, id).Scan(
		&prt.ID, &prt.Sector, &prt.Cell, &prt.Village, &prt.Due, &fname, &lname,
	); err != nil {
		empty := properties.Property{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, properties.ErrNotFound
		}

		return empty, err
	}
	prt.Owner = fmt.Sprintf("%s %s", fname, lname)
	return prt, nil
}

func (str *propertiesStore) RetrieveByOwner(owner string, offset, limit uint64) (properties.PropertyPage, error) {
	q := `SELECT 
			properties.id, properties.sector, properties.cell, 
			properties.village, properties.due, owners.fname, owners.lname 
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id
		WHERE 
			properties.owner = $1 ORDER BY id LIMIT $2 OFFSET $3
	`

	var items = []properties.Property{}

	rows, err := str.db.Query(q, owner, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Property{}

		var fname, lname string

		if err := rows.Scan(
			&c.ID, &c.Sector, &c.Cell, &c.Village, &c.Due, &fname, &lname,
		); err != nil {
			return properties.PropertyPage{}, err
		}

		c.Owner = fmt.Sprintf("%s %s", fname, lname)

		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM properties WHERE owner = $1`

	var total uint64
	if err := str.db.QueryRow(q, owner).Scan(&total); err != nil {
		return properties.PropertyPage{}, err
	}

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}

func (str *propertiesStore) RetrieveBySector(sector string, offset, limit uint64) (properties.PropertyPage, error) {
	q := `
		SELECT 
			properties.id, properties.sector, properties.cell, 
			properties.village, properties.due, owners.fname, owners.lname 
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 
		WHERE 
			properties.sector = $1 ORDER BY id LIMIT $2 OFFSET $3
	`

	var items = []properties.Property{}

	rows, err := str.db.Query(q, sector, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Property{}

		var fname, lname string

		if err := rows.Scan(
			&c.ID, &c.Sector, &c.Cell, &c.Village, &c.Due, &fname, &lname,
		); err != nil {
			return properties.PropertyPage{}, err
		}
		c.Owner = fmt.Sprintf("%s %s", fname, lname)
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM properties WHERE sector = $1`

	var total uint64
	if err := str.db.QueryRow(q, sector).Scan(&total); err != nil {
		return properties.PropertyPage{}, err
	}

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}

func (str *propertiesStore) RetrieveByCell(cell string, offset, limit uint64) (properties.PropertyPage, error) {
	q := `
		SELECT 
			properties.id, properties.sector, properties.cell, 
			properties.village, properties.due, owners.fname, owners.lname 
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 	
		WHERE properties.cell = $1 ORDER BY id LIMIT $2 OFFSET $3
	`

	var items = []properties.Property{}

	rows, err := str.db.Query(q, cell, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Property{}

		var fname, lname string

		if err := rows.Scan(
			&c.ID, &c.Sector, &c.Cell, &c.Village, &c.Due, &fname, &lname,
		); err != nil {
			return properties.PropertyPage{}, err
		}
		c.Owner = fmt.Sprintf("%s %s", fname, lname)
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM properties WHERE cell = $1`

	var total uint64
	if err := str.db.QueryRow(q, cell).Scan(&total); err != nil {
		return properties.PropertyPage{}, err
	}

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}

func (str *propertiesStore) RetrieveByVillage(village string, offset, limit uint64) (properties.PropertyPage, error) {
	q := `
		SELECT 
			properties.id, properties.sector, properties.cell, 
			properties.village, properties.due, owners.fname, owners.lname 
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 
		WHERE properties.village = $1 ORDER BY id LIMIT $2 OFFSET $3
	`

	var items = []properties.Property{}

	rows, err := str.db.Query(q, village, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Property{}
		var fname, lname string

		if err := rows.Scan(
			&c.ID, &c.Sector, &c.Cell, &c.Village, &c.Due, &fname, &lname,
		); err != nil {
			return properties.PropertyPage{}, err
		}
		c.Owner = fmt.Sprintf("%s %s", fname, lname)
		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM properties WHERE village = $1`

	var total uint64
	if err := str.db.QueryRow(q, village).Scan(&total); err != nil {
		return properties.PropertyPage{}, err
	}

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}
