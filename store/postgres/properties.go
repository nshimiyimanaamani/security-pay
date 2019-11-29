package postgres

import (
	"context"
	"database/sql"

	//"github.com/lib/pq"
	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/properties"
)

var _ (properties.Repository) = (*propertiesStore)(nil)

type propertiesStore struct {
	db *sql.DB
}

// NewPropertyStore instanctiates a new transactiob store interface
func NewPropertyStore(db *sql.DB) properties.Repository {
	return &propertiesStore{db}
}

func (str *propertiesStore) Save(ctx context.Context, pro properties.Property) (properties.Property, error) {
	q := `INSERT INTO properties (id, owner, due, sector, cell, village) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	empty := properties.Property{}

	_, err := str.db.Exec(q, pro.ID, pro.Owner.ID, pro.Due, pro.Address.Sector, pro.Address.Cell, pro.Address.Village)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return empty, properties.ErrConflict
			case errFK:
				return empty, properties.ErrOwnerNotFound
			case errInvalid, errTruncation:
				return empty, properties.ErrInvalidEntity
			}
		}
		return empty, err
	}

	return pro, nil
}

func (str *propertiesStore) UpdateProperty(ctx context.Context, pro properties.Property) error {
	q := `UPDATE properties SET owner=$1, due=$2 WHERE id=$3;`

	res, err := str.db.Exec(q, pro.Owner.ID, pro.Due, pro.ID)
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
		return properties.ErrPropertyNotFound
	}
	return nil
}

func (str *propertiesStore) RetrieveByID(ctx context.Context, id string) (properties.Property, error) {
	q := `
		SELECT 
			properties.id, properties.sector, properties.cell,  
			properties.village, properties.due, 
			owners.id, owners.fname, owners.lname, owners.phone
		FROM 
			properties
		INNER JOIN 
			owners ON properties.owner=owners.id 
		WHERE properties.id = $1
	`

	var prt = properties.Property{}

	if err := str.db.QueryRow(q, id).Scan(
		&prt.ID, &prt.Address.Sector, &prt.Address.Cell, &prt.Address.Village,
		&prt.Due, &prt.Owner.ID, &prt.Owner.Fname, &prt.Owner.Lname, &prt.Owner.Phone,
	); err != nil {
		empty := properties.Property{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, properties.ErrPropertyNotFound
		}

		return empty, err
	}
	return prt, nil
}

func (str *propertiesStore) RetrieveByOwner(ctx context.Context, owner string, offset, limit uint64) (properties.PropertyPage, error) {
	q := `SELECT 
			properties.id, properties.sector, properties.cell, 
			properties.village, properties.due, 
			owners.id, owners.fname, owners.lname, owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id
		WHERE 
			properties.owner = $1 ORDER BY properties.id LIMIT $2 OFFSET $3
	`

	var items = []properties.Property{}

	rows, err := str.db.Query(q, owner, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Property{}

		if err := rows.Scan(
			&c.ID, &c.Address.Sector, &c.Address.Cell, &c.Address.Village,
			&c.Due, &c.Owner.ID, &c.Owner.Fname, &c.Owner.Lname, &c.Owner.Phone,
		); err != nil {
			return properties.PropertyPage{}, err
		}

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

func (str *propertiesStore) RetrieveBySector(ctx context.Context, sector string, offset, limit uint64) (properties.PropertyPage, error) {
	q := `
		SELECT 
			properties.id, properties.sector, properties.cell, 
			properties.village, properties.due, 
			owners.id, owners.fname, owners.lname, owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 
		WHERE 
			properties.sector = $1 ORDER BY  properties.id LIMIT $2 OFFSET $3
	`

	var items = []properties.Property{}

	rows, err := str.db.Query(q, sector, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Property{}

		if err := rows.Scan(
			&c.ID, &c.Address.Sector, &c.Address.Cell, &c.Address.Village,
			&c.Due, &c.Owner.ID, &c.Owner.Fname, &c.Owner.Lname, &c.Owner.Phone,
		); err != nil {
			return properties.PropertyPage{}, err
		}

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

func (str *propertiesStore) RetrieveByCell(ctx context.Context, cell string, offset, limit uint64) (properties.PropertyPage, error) {
	q := `
		SELECT 
			properties.id, properties.sector, properties.cell, 
			properties.village, properties.due, 
			owners.id, owners.fname, owners.lname, owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 	
		WHERE properties.cell = $1 ORDER BY properties.id LIMIT $2 OFFSET $3
	`

	var items = []properties.Property{}

	rows, err := str.db.Query(q, cell, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Property{}

		if err := rows.Scan(
			&c.ID, &c.Address.Sector, &c.Address.Cell, &c.Address.Village,
			&c.Due, &c.Owner.ID, &c.Owner.Fname, &c.Owner.Lname, &c.Owner.Phone,
		); err != nil {
			return properties.PropertyPage{}, err
		}

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

func (str *propertiesStore) RetrieveByVillage(ctx context.Context, village string, offset, limit uint64) (properties.PropertyPage, error) {
	q := `
		SELECT 
			properties.id, properties.sector, properties.cell, 
			properties.village, properties.due, 
			owners.id, owners.fname, owners.lname, owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 
		WHERE properties.village = $1 ORDER BY properties.id LIMIT $2 OFFSET $3
	`

	var items = []properties.Property{}

	rows, err := str.db.Query(q, village, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, err
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Property{}

		if err := rows.Scan(
			&c.ID, &c.Address.Sector, &c.Address.Cell, &c.Address.Village,
			&c.Due, &c.Owner.ID, &c.Owner.Lname, &c.Owner.Lname, &c.Owner.Phone,
		); err != nil {
			return properties.PropertyPage{}, err
		}

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
