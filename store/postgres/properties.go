package postgres

import (
	"context"
	"database/sql"

	//"github.com/lib/pq"
	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (properties.Repository) = (*propertiesStore)(nil)

type propertiesStore struct {
	*sql.DB
}

// NewPropertyStore instanctiates a new transactiob store interface
func NewPropertyStore(db *sql.DB) properties.Repository {
	return &propertiesStore{db}
}

func (repo *propertiesStore) Save(ctx context.Context, pro properties.Property) (properties.Property, error) {
	const op errors.Op = "store/postgres/propertiesStore.Save"

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
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING created_at, updated_at`

	empty := properties.Property{}

	err := repo.QueryRow(q, pro.ID, pro.Owner.ID, pro.Due, pro.Address.Sector,
		pro.Address.Cell, pro.Address.Village, pro.RecordedBy, pro.Occupied).Scan(&pro.CreatedAt, &pro.UpdatedAt)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return empty, errors.E(op, err, "property already exists", errors.KindAlreadyExists)
			case errFK, errTruncation, errInvalid:
				return empty, errors.E(op, err, "owner not found", errors.KindNotFound)
			}
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}

	return pro, nil
}

func (repo *propertiesStore) UpdateProperty(ctx context.Context, pro properties.Property) error {
	const op errors.Op = "store/postgres/propertiesStore.UpdateProperty"

	q := `UPDATE properties SET owner=$1, due=$2, sector=$3, cell=$4, village=$5, occupied=$6, for_rent=$7 WHERE id=$8;`

	res, err := repo.Exec(q, pro.Owner.ID, pro.Due,
		pro.Address.Sector, pro.Address.Cell, pro.Address.Village, pro.ForRent, pro.Occupied, pro.ID,
	)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return errors.E(op, err, "invalid property", errors.KindBadRequest)
			}
		}
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if cnt == 0 {
		return errors.E(op, "property not found", errors.KindNotFound)
	}
	return nil
}

func (repo *propertiesStore) RetrieveByID(ctx context.Context, id string) (properties.Property, error) {
	const op errors.Op = "store/postgres/propertiesStore.RetrieveByID"

	q := `
		SELECT 
			properties.id, properties.sector, properties.cell,  
			properties.village, properties.due, properties.recorded_by,
			properties.occupied, properties.for_rent, properties.created_at, 
			properties.updated_at, owners.id, owners.fname, owners.lname, owners.phone
		FROM 
			properties
		INNER JOIN 
			owners ON properties.owner=owners.id 
		WHERE properties.id = $1
	`

	var prt = properties.Property{}

	if err := repo.QueryRow(q, id).Scan(
		&prt.ID, &prt.Address.Sector, &prt.Address.Cell, &prt.Address.Village,
		&prt.Due, &prt.RecordedBy, &prt.Occupied, &prt.ForRent, &prt.CreatedAt, &prt.UpdatedAt,
		&prt.Owner.ID, &prt.Owner.Fname, &prt.Owner.Lname, &prt.Owner.Phone,
	); err != nil {
		empty := properties.Property{}

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, err, "property not found", errors.KindNotFound)
		}

		return empty, err
	}
	return prt, nil
}

func (repo *propertiesStore) RetrieveByOwner(ctx context.Context, owner string, offset, limit uint64) (properties.PropertyPage, error) {
	const op errors.Op = "store/postgres/propertiesStore.RetrieveByOwner"

	q := `SELECT 
			properties.id, properties.sector, properties.cell, 
			properties.village, properties.due, properties.recorded_by,
			properties.occupied, properties.for_rent, properties.created_at,
			properties.updated_at, owners.id, owners.fname, owners.lname, owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id
		WHERE 
			properties.owner = $1 ORDER BY properties.id LIMIT $2 OFFSET $3
	`

	var items = []properties.Property{}

	rows, err := repo.Query(q, owner, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Property{}

		if err := rows.Scan(
			&c.ID, &c.Address.Sector, &c.Address.Cell, &c.Address.Village,
			&c.Due, &c.RecordedBy, &c.Occupied, &c.ForRent, &c.CreatedAt, &c.UpdatedAt,
			&c.Owner.ID, &c.Owner.Fname, &c.Owner.Lname, &c.Owner.Phone,
		); err != nil {
			return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
		}

		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM properties WHERE owner = $1`

	var total uint64
	if err := repo.QueryRow(q, owner).Scan(&total); err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
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

func (repo *propertiesStore) RetrieveBySector(ctx context.Context, sector string, offset, limit uint64) (properties.PropertyPage, error) {
	const op errors.Op = "store/postgres/propertiesStore.RetrieveBySector"

	q := `
		SELECT 
			properties.id, properties.sector, properties.cell, 
			properties.village, properties.due, properties.recorded_by, 
			properties.occupied, properties.for_rent, properties.created_at,
			properties.updated_at, owners.id, owners.fname, owners.lname, owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 
		WHERE 
			properties.sector = $1 ORDER BY  properties.id LIMIT $2 OFFSET $3
	`

	var items = []properties.Property{}

	rows, err := repo.Query(q, sector, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Property{}

		if err := rows.Scan(
			&c.ID, &c.Address.Sector, &c.Address.Cell, &c.Address.Village,
			&c.Due, &c.RecordedBy, &c.Occupied, &c.ForRent, &c.CreatedAt, &c.UpdatedAt,
			&c.Owner.ID, &c.Owner.Fname, &c.Owner.Lname, &c.Owner.Phone,
		); err != nil {
			return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
		}

		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM properties WHERE sector = $1`

	var total uint64
	if err := repo.QueryRow(q, sector).Scan(&total); err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
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

func (repo *propertiesStore) RetrieveByCell(ctx context.Context, cell string, offset, limit uint64) (properties.PropertyPage, error) {
	const op errors.Op = "store/postgres/propertiesStore.RetrieveByCell"

	q := `
		SELECT 
			properties.id, properties.sector, properties.cell, 
			properties.village, properties.due, properties.recorded_by, 
			properties.occupied, properties.for_rent, properties.created_at,
			properties.updated_at, owners.id, owners.fname, owners.lname, owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 	
		WHERE properties.cell = $1 ORDER BY properties.id LIMIT $2 OFFSET $3
	`

	var items = []properties.Property{}

	rows, err := repo.Query(q, cell, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Property{}

		if err := rows.Scan(
			&c.ID, &c.Address.Sector, &c.Address.Cell, &c.Address.Village,
			&c.Due, &c.RecordedBy, &c.Occupied, &c.ForRent, &c.CreatedAt, &c.UpdatedAt,
			&c.Owner.ID, &c.Owner.Fname, &c.Owner.Lname, &c.Owner.Phone,
		); err != nil {
			return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
		}

		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM properties WHERE cell = $1`

	var total uint64
	if err := repo.QueryRow(q, cell).Scan(&total); err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
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

func (repo *propertiesStore) RetrieveByVillage(ctx context.Context, village string, offset, limit uint64) (properties.PropertyPage, error) {
	const op errors.Op = "store/postgres/propertiesStore.RetrieveByVillage"

	q := `
		SELECT 
			properties.id, properties.sector, properties.cell, 
			properties.village, properties.due, properties.recorded_by, 
			properties.occupied, properties.for_rent, properties.created_at,
			properties.updated_at, owners.id, owners.fname, owners.lname, owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 
		WHERE properties.village = $1 ORDER BY properties.id LIMIT $2 OFFSET $3
	`

	var items = []properties.Property{}

	rows, err := repo.Query(q, village, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		c := properties.Property{}

		if err := rows.Scan(
			&c.ID, &c.Address.Sector, &c.Address.Cell, &c.Address.Village,
			&c.Due, &c.RecordedBy, &c.Occupied, &c.ForRent, &c.CreatedAt, &c.UpdatedAt,
			&c.Owner.ID, &c.Owner.Fname, &c.Owner.Lname, &c.Owner.Phone,
		); err != nil {
			return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
		}

		items = append(items, c)
	}

	q = `SELECT COUNT(*) FROM properties WHERE village = $1`

	var total uint64
	if err := repo.QueryRow(q, village).Scan(&total); err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
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
