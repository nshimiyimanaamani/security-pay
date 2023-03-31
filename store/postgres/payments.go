package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/clock"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ (payment.Repository) = (*paymentStore)(nil)

type paymentStore struct {
	*sql.DB
	sms notifs.Service
}

// NewPaymentRepository creates a new postgres backed payment.Repository
func NewPaymentRepository(db *sql.DB, sms notifs.Service) payment.Repository {
	return &paymentStore{db, sms}
}

func (repo *paymentStore) Save(ctx context.Context, payment *payment.TxRequest) error {
	const op errors.Op = "store/postgres/paymentStore.Save"

	q := `INSERT INTO payments(
			id,
			ref,
			amount,
			msisdn,
			method,
			invoice,
			property,
			confirmed
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8);
	`
	_, err := repo.ExecContext(ctx, q,
		payment.ID,
		payment.Ref,
		payment.Amount,
		payment.MSISDN,
		payment.Method,
		payment.Invoice,
		payment.Code,
		payment.Confirmed,
	)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return errors.E(op, "duplicate payment id", errors.KindAlreadyExists)
			case errInvalid, errTruncation:
				return errors.E(op, "invalid payment entity", errors.KindBadRequest)
			}
		}
		return errors.E(op, err, errors.KindUnexpected)
	}
	return nil
}

func (repo *paymentStore) Find(ctx context.Context, id string) ([]*payment.TxRequest, error) {

	const op errors.Op = "store/postgres/paymentStore.Find"

	tx, err := repo.BeginTx(ctx, nil)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}
	defer tx.Rollback()

	query := `
		SELECT 
			id, 
			ref,
			amount, 
			msisdn, 
			status,
			method, 
			invoice, 
			property,
			confirmed,
			created_at,
			updated_at
		FROM 
			payments
		WHERE ref=$1
	`
	out := make([]*payment.TxRequest, 0)

	rows, err := tx.QueryContext(ctx, query, id)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	for rows.Next() {
		trns := new(payment.TxRequest)
		err = rows.Scan(
			&trns.ID,
			&trns.Ref,
			&trns.Amount,
			&trns.MSISDN,
			&trns.Status,
			&trns.Method,
			&trns.Invoice,
			&trns.Code,
			&trns.Confirmed,
			&trns.CreatedAt,
			&trns.UpdatedAt,
		)
		if err != nil {
			pqErr, ok := err.(*pq.Error)
			if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
				return nil, errors.E(op, err, "payment not found", errors.KindNotFound)
			}
			return nil, errors.E(op, err, errors.KindUnexpected)
		}
		out = append(out, trns)
	}

	return out, tx.Commit()
}

func (repo *paymentStore) Update(ctx context.Context, status string, payments []*payment.TxRequest) error {
	const op errors.Op = "store/postgres/paymentStore.Update"

	tx, err := repo.BeginTx(ctx, nil)
	if err != nil {
		return errors.E(op, err, errors.KindUnexpected)
	}
	defer tx.Rollback()

	property := new(properties.Property)

	// get property by property code
	q := `
		SELECT
			owner,
			namespace,
			sector,
			fname,
			lname,
			phone
		FROM
			properties INNER JOIN owners ON properties.owner = owners.id
		WHERE properties.id=$1
	`
	if err = tx.QueryRowContext(ctx, q, payments[0].Code).Scan(
		&property.Owner.ID,
		&property.Namespace,
		&property.Address.Sector,
		&property.Owner.Fname,
		&property.Owner.Lname,
		&property.Owner.Phone,
	); err != nil {
		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return errors.E(op, err, "property not found", errors.KindNotFound)
		}
		return errors.E(op, err, errors.KindUnexpected)
	}

	pos, args := []string{}, []interface{}{}

	i := 0
	for _, txns := range payments {
		if status == "successful" {
			pos = append(pos,
				fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*9+1, i*9+2, i*9+3, i*9+4, i*9+5, i*9+6, i*9+7, i*9+8, i*9+9),
			)
			args = append(
				args,
				txns.ID,
				txns.Ref,
				status,
				txns.Code,
				property.Owner.ID,
				txns.Amount,
				txns.Method,
				txns.Invoice,
				property.Namespace,
			)
			i++
		}
	}

	if len(pos) > 0 {
		var query = insertTxQuery + strings.Join(pos, ",")

		_, err = tx.ExecContext(
			ctx,
			query,
			args...,
		)
		if err != nil {
			pqErr, ok := err.(*pq.Error)
			if ok {
				switch pqErr.Code.Name() {
				case errDuplicate:
					return errors.E(op, "transactions already existed", errors.KindAlreadyExists)
				}
			}
			return errors.E(op, err, errors.KindUnexpected)
		}
	}

	_, err = tx.ExecContext(ctx, updatePayQuery, true, status, payments[0].Ref)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errInvalid, errTruncation:
				return errors.E(op, err, "payment not found", errors.KindNotFound)
			}
		}
		return errors.E(op, err, errors.KindUnexpected)
	}

	if status == "successful" {
		go func() error {
			_, err := repo.sms.Send(ctx,
				notifs.Notification{
					Sender:     property.Namespace,
					Recipients: []string{property.Owner.Phone},
					Message:    formMessage(payments, property)},
			)
			if err != nil {
				return errors.E(op, err, errors.KindUnexpected)
			}
			return nil
		}()
	}

	return tx.Commit()
}

func (repo *paymentStore) BulkSave(ctx context.Context, payments []*payment.TxRequest) error {
	const op errors.Op = "store/postgres/paymentStore.BulkSave"

	tx, err := repo.BeginTx(ctx, nil)
	if err != nil {
		return errors.E(op, err)
	}
	defer tx.Rollback()

	insertQuery := `
	INSERT INTO payments(
		id,
		ref,
		amount,
		msisdn,
		method,
		invoice,
		property,
		confirmed
	) VALUES 
	`

	pos, args := []string{}, []interface{}{}

	i := 0
	for _, item := range payments {
		pos = append(pos,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8),
		)
		args = append(
			args,
			item.ID,
			item.Ref,
			item.Amount,
			item.MSISDN,
			item.Method,
			item.Invoice,
			item.Code,
			item.Confirmed,
		)
		i++
	}

	var query = insertQuery + strings.Join(pos, ",")

	_, err = tx.ExecContext(
		ctx,
		query,
		args...,
	)
	if err != nil {
		fmt.Println("insert payment err", err)
		fmt.Println("insert payment input", payments)
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return errors.E(op, "duplicate payment id", errors.KindAlreadyExists)
			case errInvalid, errTruncation:
				return errors.E(op, "invalid payment entity", errors.KindBadRequest)
			}
		}
		return errors.E(op, err, errors.KindUnexpected)
	}
	return tx.Commit()
}

func (repo *paymentStore) List(ctx context.Context, flts *payment.Filters) (payment.PaymentResponse, error) {
	const op errors.Op = "store/postgres/paymentStore.List"

	tx, err := repo.BeginTx(ctx, nil)
	if err != nil {
		return payment.PaymentResponse{}, errors.E(op, err)
	}
	defer tx.Rollback()

	selectQuery := `SELECT 
			i.id,
			o.fname, 
			o.lname,
			o.phone
		FROM 
			owners o
		JOIN properties p 
			ON p.owner = o.id
		JOIN invoices i ON 
			i.property = p.id
		WHERE 1 = 1
	`

	if *flts.Status != "" {
		selectQuery += fmt.Sprintf("\nAND i.status = '%s'", *flts.Status)
	}

	selectQuery += "\nAND DATE_TRUNC('month', i.created_at) = DATE_TRUNC('month', CURRENT_DATE)"

	if *flts.Sector != "" {
		selectQuery += fmt.Sprintf("\nAND p.sector = '%s'", *flts.Sector)
	}

	if *flts.Cell != "" {
		selectQuery += fmt.Sprintf("\nAND p.cell = '%s'", *flts.Cell)
	}

	if *flts.Village != "" {
		selectQuery += fmt.Sprintf("\nAND p.village = '%s'", *flts.Village)
	}

	selectQuery += "\nORDER BY i.created_at DESC"
	selectQuery += fmt.Sprintf("\nOFFSET %d LIMIT %d", *flts.Offset, *flts.Limit)

	rows, err := tx.QueryContext(ctx, selectQuery)
	if err != nil {
		return payment.PaymentResponse{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	var payments []payment.Payment
	for rows.Next() {
		var pmt payment.Payment
		err := rows.Scan(
			&pmt.ID,
			&pmt.Fname,
			&pmt.Lname,
			&pmt.Email,
		)
		if err != nil {
			return payment.PaymentResponse{}, errors.E(op, err, errors.KindUnexpected)
		}
		payments = append(payments, pmt)
	}

	selectQuery = `SELECT COUNT(*) FROM invoices i JOIN properties p ON i.property = p.id`
	selectQuery += "\nWHERE 1 = 1"

	if *flts.Status != "" {
		selectQuery += fmt.Sprintf("\nAND i.status = '%s'", *flts.Status)
	}

	selectQuery += "\nAND DATE_TRUNC('month', i.created_at) = DATE_TRUNC('month', CURRENT_DATE)"

	if *flts.Sector != "" {
		selectQuery += fmt.Sprintf("\nAND p.sector = '%s'", *flts.Sector)
	}

	if *flts.Cell != "" {
		selectQuery += fmt.Sprintf("\nAND p.cell = '%s'", *flts.Cell)
	}

	if flts.Village != nil {
		selectQuery += fmt.Sprintf("\nAND p.village = '%s'", *flts.Village)
	}

	var total uint64
	if err := repo.QueryRow(selectQuery).Scan(&total); err != nil {
		return payment.PaymentResponse{}, errors.E(op, err, errors.KindUnexpected)
	}
	//
	page := payment.PaymentResponse{
		Payments: payments,
		PageMetadata: payment.PageMetadata{
			Total:  total,
			Offset: *flts.Offset,
			Limit:  *flts.Limit,
		},
	}
	return page, nil
}

func formMessage(tx []*payment.TxRequest, prop *properties.Property) string {

	const header = "Murakoze kwishyura umusanzu w' isuku"
	var (
		amount          int
		invoices, month string
	)

	for _, item := range tx {
		amount += int(item.Amount)
		invoices += fmt.Sprintf("%d, ", item.Invoice)
	}

	if len(tx) > 1 {
		month = fmt.Sprintf("Wishyuriye Amezi %d\n", len(tx))
	} else {
		month = fmt.Sprintf("Wishyuriye Ukwezi kwa: %d\n", int(tx[0].CreatedAt.Month()))
	}

	var buf bytes.Buffer

	buf.WriteString(header)
	// buf.WriteString(selectActivity(pr.Address.Sector))
	buf.WriteString(" mu murenge wa ")
	buf.WriteString(fmt.Sprintf("%s.\n\n", prop.Address.Sector))
	buf.WriteString(fmt.Sprintf("Nimero yishyuriweho: %s\n", tx[0].MSISDN))
	buf.WriteString(fmt.Sprintf("Itariki: %s\n", timestamp()))
	buf.WriteString(month)
	buf.WriteString(fmt.Sprintf("Nimero ya fagitire: %s\n", invoices))
	buf.WriteString(fmt.Sprintf("Umubare w' amafaranga: %dRWF\n", amount))
	buf.WriteString(fmt.Sprintf("Inzu yishyuriwe ni iya %s %s\n", prop.Owner.Fname, prop.Owner.Lname))
	buf.WriteString(fmt.Sprintf("Code y' inzu ni: %s", tx[0].Code))
	return buf.String()
}

func timestamp() string {
	at := clock.TimeIn(time.Now(), clock.EAST)
	return clock.Format(at, clock.LayoutCustom)
}

// update transaction
var insertTxQuery = `
	INSERT INTO transactions (
		id,
		ref,
		status,
		madefor,
		madeby,
		amount,
		method,
		invoice,
		namespace
	) VALUES
`

// update payments table
var updatePayQuery = `
UPDATE
	payments 
SET 
	confirmed=$1, status=$2 
WHERE ref=$3
`
