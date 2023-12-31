package postgres

import (
	"context"
	"database/sql"
	"fmt"

	//"github.com/lib/pq"
	"github.com/lib/pq"
	"github.com/nshimiyimanaamani/paypack-backend/core/auth"
	"github.com/nshimiyimanaamani/paypack-backend/core/properties"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
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
			occupied,
			namespace
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING created_at, updated_at`

	empty := properties.Property{}

	err := repo.QueryRow(q,
		pro.ID,
		pro.Owner.ID,
		pro.Due,
		pro.Address.Sector,
		pro.Address.Cell,
		pro.Address.Village,
		pro.RecordedBy,
		pro.Occupied,
		pro.Namespace,
	).Scan(&pro.CreatedAt, &pro.UpdatedAt)

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

func (repo *propertiesStore) Update(ctx context.Context, pro properties.Property) error {
	const op errors.Op = "store/postgres/propertiesStore.Update"

	q := `
		UPDATE properties SET 
			owner=$1, due=$2, sector=$3, 
			cell=$4, village=$5, occupied=$6, 
			for_rent=$7, namespace=$8
		WHERE id=$9;
	`

	res, err := repo.Exec(q,
		pro.Owner.ID,
		pro.Due,
		pro.Address.Sector,
		pro.Address.Cell,
		pro.Address.Village,
		pro.ForRent,
		pro.Occupied,
		pro.Namespace, pro.ID,
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

func (repo *propertiesStore) Delete(ctx context.Context, uid string) error {
	const op errors.Op = "store/postgres/propertiesStore.Delete"

	q := `DELETE FROM properties WHERE id=$1`

	res, err := repo.ExecContext(ctx, q, uid)
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	cnt, err := res.RowsAffected()
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
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
			properties.id, 
			properties.sector, 
			properties.cell,  
			properties.village, 
			properties.due, 
			properties.recorded_by,
			properties.occupied, 
			properties.for_rent, 
			properties.created_at, 
			properties.updated_at, 
			properties.namespace,
			owners.id, 
			owners.fname, 
			owners.lname, 
			owners.phone
		FROM 
			properties
		INNER JOIN 
			owners ON properties.owner=owners.id 
		WHERE properties.id = $1
	`

	var prt = properties.Property{}

	err := repo.QueryRow(q, id).Scan(
		&prt.ID,
		&prt.Address.Sector,
		&prt.Address.Cell,
		&prt.Address.Village,
		&prt.Due,
		&prt.RecordedBy,
		&prt.Occupied,
		&prt.ForRent,
		&prt.CreatedAt,
		&prt.UpdatedAt,
		&prt.Namespace,
		&prt.Owner.ID,
		&prt.Owner.Fname,
		&prt.Owner.Lname,
		&prt.Owner.Phone,
	)
	if err != nil {
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
			properties.id, 
			properties.sector, 
			properties.cell, 
			properties.village, 
			properties.due, 
			properties.recorded_by,
			properties.occupied, 
			properties.for_rent, 
			properties.created_at,
			properties.updated_at,
			properties.namespace,
			owners.id, 
			owners.fname, 
			owners.lname, 
			owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id
		WHERE 
			properties.owner = $1
		ORDER BY properties.id LIMIT $2 OFFSET $3
		
	`

	// creds := auth.CredentialsFromContext(ctx)

	var items = []properties.Property{}

	rows, err := repo.Query(q, owner, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		row := properties.Property{}

		err := rows.Scan(
			&row.ID,
			&row.Address.Sector,
			&row.Address.Cell,
			&row.Address.Village,
			&row.Due,
			&row.RecordedBy,
			&row.Occupied,
			&row.ForRent,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.Namespace,
			&row.Owner.ID,
			&row.Owner.Fname,
			&row.Owner.Lname,
			&row.Owner.Phone,
		)
		if err != nil {
			return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
		}

		items = append(items, row)
	}

	q = `SELECT COUNT(*),COALESCE(SUM(properties.due), 0)  FROM properties WHERE owner = $1`

	var total uint64
	var amount float64
	if err := repo.QueryRow(q, owner).Scan(&total, &amount); err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
			Amount: amount,
		},
	}
	return page, nil
}

func (repo *propertiesStore) RetrieveBySector(ctx context.Context, flts *properties.Filters) (properties.PropertyPage, error) {
	const op errors.Op = "store/postgres/propertiesStore.RetrieveBySector"

	tx, err := repo.BeginTx(ctx, nil)
	if err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer tx.Rollback()

	selectQuery := `
		SELECT 
			properties.id, 
			properties.sector, 
			properties.cell, 
			properties.village, 
			properties.due, 
			properties.recorded_by, 
			properties.occupied, 
			properties.for_rent, 
			properties.created_at,
			properties.updated_at, 
			properties.namespace,
			owners.id, 
			owners.fname, 
			owners.lname, 
			owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 
		WHERE  1=1
	`

	if flts.Namespace != nil {
		selectQuery += fmt.Sprintf(" AND properties.namespace='%s'", *flts.Namespace)
	}

	if flts.Sector != nil {
		selectQuery += fmt.Sprintf(" AND properties.sector='%s'", *flts.Sector)
	}

	if flts.Phone != nil {
		selectQuery += fmt.Sprintf(" AND owners.phone='%s'", *flts.Phone)
	}

	if flts.Names != nil {
		selectQuery += fmt.Sprintf(" AND (owners.fname LIKE '%%%s%%' OR owners.lname LIKE '%%%s%%') ", *flts.Names, *flts.Names)
	}

	selectQuery += fmt.Sprintf(" ORDER BY  properties.id LIMIT %d OFFSET %d", *flts.Limit, *flts.Offset)
	var items = []properties.Property{}

	rows, err := tx.QueryContext(ctx, selectQuery)
	if err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		row := properties.Property{}

		err := rows.Scan(
			&row.ID,
			&row.Address.Sector,
			&row.Address.Cell,
			&row.Address.Village,
			&row.Due,
			&row.RecordedBy,
			&row.Occupied,
			&row.ForRent,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.Namespace,
			&row.Owner.ID,
			&row.Owner.Fname,
			&row.Owner.Lname,
			&row.Owner.Phone,
		)

		if err != nil {
			return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
		}

		items = append(items, row)
	}

	countQuery := `SELECT 
		COUNT(*),
		COALESCE(SUM(properties.due), 0.0)
	FROM 
		properties
	INNER JOIN
		owners ON properties.owner=owners.id 
	WHERE 1=1 `

	if flts.Phone != nil {
		countQuery += fmt.Sprintf(" AND owners.phone='%s'", *flts.Phone)
	}

	if flts.Sector != nil {
		countQuery += fmt.Sprintf(" AND properties.sector='%s'", *flts.Sector)
	}

	if flts.Namespace != nil {
		countQuery += fmt.Sprintf(" AND properties.namespace='%s'", *flts.Namespace)
	}

	if flts.Names != nil {
		countQuery += fmt.Sprintf(" AND (owners.fname LIKE '%%%s%%' OR owners.lname LIKE '%%%s%%') ", *flts.Names, *flts.Names)
	}

	var (
		total  uint64
		amount float64
	)

	if err := tx.QueryRowContext(ctx, countQuery).Scan(&total, &amount); err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  total,
			Offset: *flts.Offset,
			Limit:  *flts.Limit,
			Amount: amount,
		},
	}
	return page, tx.Commit()
}

func (repo *propertiesStore) RetrieveByCell(ctx context.Context, cell string, offset, limit uint64, names string) (properties.PropertyPage, error) {
	const op errors.Op = "store/postgres/propertiesStore.RetrieveByCell"

	q := `
		SELECT 
			properties.id, 
			properties.sector, 
			properties.cell, 
			properties.village, 
			properties.due, 
			properties.recorded_by, 
			properties.occupied, 
			properties.for_rent, 
			properties.created_at,
			properties.updated_at, 
			properties.namespace,
			owners.id, 
			owners.fname, 
			owners.lname, 
			owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 	
		WHERE 
			properties.cell = $1 AND properties.namespace=$2
			AND (owners.fname LIKE '%' || $3 || '%' OR owners.lname  LIKE '%' || $3 || '%')
		ORDER BY properties.id LIMIT $4 OFFSET $5
	`

	var items = []properties.Property{}

	creds := auth.CredentialsFromContext(ctx)

	rows, err := repo.Query(q, cell, creds.Account, names, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		row := properties.Property{}

		if err := rows.Scan(
			&row.ID,
			&row.Address.Sector,
			&row.Address.Cell,
			&row.Address.Village,
			&row.Due,
			&row.RecordedBy,
			&row.Occupied,
			&row.ForRent,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.Namespace,
			&row.Owner.ID,
			&row.Owner.Fname, &row.Owner.Lname, &row.Owner.Phone,
		); err != nil {
			return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
		}

		items = append(items, row)
	}

	q = `SELECT 
			COUNT(*),
			COALESCE(SUM(properties.due), 0)
	    FROM 
		    properties 
		INNER JOIN
			owners ON properties.owner=owners.id 
	    WHERE 
			cell = $1 
	    AND namespace=$2 AND (owners.fname LIKE '%' || $3 || '%' OR owners.lname  LIKE '%' || $3 || '%')`

	var total uint64
	var amount float64
	if err := repo.QueryRow(q, cell, creds.Account, names).Scan(&total, &amount); err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
			Amount: amount,
		},
	}
	return page, nil
}

func (repo *propertiesStore) RetrieveByVillage(ctx context.Context, village string, offset, limit uint64, names string) (properties.PropertyPage, error) {
	const op errors.Op = "store/postgres/propertiesStore.RetrieveByVillage"

	q := `
		SELECT 
			properties.id, 
			properties.sector, 
			properties.cell, 
			properties.village, 
			properties.due, 
			properties.recorded_by, 
			properties.occupied, 
			properties.for_rent, 
			properties.created_at,
			properties.updated_at,
			properties.namespace, 
			owners.id, 
			owners.fname, 
			owners.lname, 
			owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 
		WHERE 
			properties.village = $1 
			AND properties.namespace=$2 
			AND (owners.fname LIKE '%' || $3 || '%' OR owners.lname  LIKE '%' || $3 || '%')
		ORDER BY properties.id 
		LIMIT $4 
		OFFSET $5
	`

	var items = []properties.Property{}

	creds := auth.CredentialsFromContext(ctx)

	rows, err := repo.Query(q, village, creds.Account, names, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		row := properties.Property{}

		err := rows.Scan(
			&row.ID,
			&row.Address.Sector,
			&row.Address.Cell,
			&row.Address.Village,
			&row.Due,
			&row.RecordedBy,
			&row.Occupied,
			&row.ForRent,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.Namespace,
			&row.Owner.ID,
			&row.Owner.Fname,
			&row.Owner.Lname,
			&row.Owner.Phone,
		)
		if err != nil {
			return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
		}

		items = append(items, row)
	}

	q = `SELECT 
		COUNT(*),
		COALESCE(SUM(properties.due), 0)
	FROM 
		properties 
	INNER JOIN
		owners ON properties.owner=owners.id 	
	WHERE 
		village = $1 
	AND namespace=$2 AND (owners.fname LIKE '%' || $3 || '%' OR owners.lname  LIKE '%' || $3 || '%')`

	var total uint64
	var amount float64
	if err := repo.QueryRow(q, village, creds.Account, names).Scan(&total, &amount); err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
			Amount: amount,
		},
	}
	return page, nil
}

func (repo *propertiesStore) RetrieveByRecorder(ctx context.Context, user string, offset, limit uint64) (properties.PropertyPage, error) {
	const op errors.Op = "store/postgres/propertiesStore.RetrieveByVillage"

	q := `
		SELECT 
			properties.id, 
			properties.sector, 
			properties.cell, 
			properties.village, 
			properties.due, 
			properties.recorded_by, 
			properties.occupied, 
			properties.for_rent, 
			properties.created_at,
			properties.updated_at,
			properties.namespace,
			owners.id, 
			owners.fname, 
			owners.lname, 
			owners.phone
		FROM 
			properties
		INNER JOIN
			owners ON properties.owner=owners.id 
		WHERE 
			properties.recorded_by = $1 AND properties.namespace=$2
		ORDER BY properties.id LIMIT $3 OFFSET $4
	`

	var items = []properties.Property{}

	creds := auth.CredentialsFromContext(ctx)

	rows, err := repo.Query(q, user, creds.Account, limit, offset)
	if err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		row := properties.Property{}

		err := rows.Scan(
			&row.ID,
			&row.Address.Sector,
			&row.Address.Cell,
			&row.Address.Village,
			&row.Due,
			&row.RecordedBy,
			&row.Occupied,
			&row.ForRent,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.Namespace,
			&row.Owner.ID,
			&row.Owner.Lname,
			&row.Owner.Lname,
			&row.Owner.Phone,
		)
		if err != nil {
			return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
		}

		items = append(items, row)
	}

	q = `SELECT COUNT(*), COALESCE(SUM(properties.due), 0)  FROM properties WHERE recorded_by = $1 AND namespace=$2`

	var total uint64
	var amount float64
	if err := repo.QueryRow(q, user, creds.Account).Scan(&total, &amount); err != nil {
		return properties.PropertyPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := properties.PropertyPage{
		Properties: items,
		PageMetadata: properties.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
			Amount: amount,
		},
	}
	return page, nil
}
