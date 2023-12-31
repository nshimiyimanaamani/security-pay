package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/nshimiyimanaamani/paypack-backend/core/notifs"
	"github.com/nshimiyimanaamani/paypack-backend/core/payment"
	"github.com/nshimiyimanaamani/paypack-backend/core/properties"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/clock"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"
	"golang.org/x/sync/errgroup"
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
			p.id, 
			p.ref,
			p.amount, 
			p.msisdn, 
			p.status,
			p.method, 
			p.invoice, 
			p.property,
			p.confirmed,
			p.created_at,
			p.updated_at,
			i.created_at
		FROM 
			payments p INNER JOIN invoices i ON p.invoice=i.id
		WHERE p.ref=$1
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
			&trns.PayedDate,
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
		g := errgroup.Group{}
		g.Go(func() error {
			_, err := repo.sms.Send(ctx,
				notifs.Notification{
					Sender:     property.Namespace,
					Recipients: []string{property.Owner.Phone},
					Message:    formMessage(payments, property)},
			)

			if err != nil {
				return err
			}

			log.NoOpLogger().Infof("sms sent to %s successful", property.Owner.Phone)

			return nil
		})

		if err := g.Wait(); err != nil {
			return errors.E(op, err, errors.KindUnexpected)
		}

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
			o.id,
			o.fname, 
			o.lname,
			o.phone,
			i.property,
			i.amount
		FROM 
			owners o
		JOIN properties p 
			ON p.owner = o.id
		JOIN invoices i ON 
			i.property = p.id
		WHERE 1 = 1
	`
	// get creds

	if flts.Status != nil {
		selectQuery += fmt.Sprintf(" AND i.status = '%s'", *flts.Status)
	}
	// check on from date
	if flts.From != nil {
		selectQuery += fmt.Sprintf(" AND i.created_at >= '%s'", *flts.From)
	}
	// check on to date
	if flts.To != nil {
		selectQuery += fmt.Sprintf(" AND i.created_at <= '%s'", *flts.To)
	}
	if flts.Sector != nil {
		selectQuery += fmt.Sprintf(" AND p.sector = '%s'", *flts.Sector)
	}

	if flts.Cell != nil {
		selectQuery += fmt.Sprintf(" AND p.cell = '%s'", *flts.Cell)
	}

	if flts.Village != nil {
		selectQuery += fmt.Sprintf(" AND p.village = '%s'", *flts.Village)
	}

	if flts.Namespace != nil {
		selectQuery += fmt.Sprintf(" AND p.namespace = '%s'", *flts.Namespace)
	}

	selectQuery += fmt.Sprintf(" ORDER BY i.created_at DESC OFFSET %d LIMIT %d", *flts.Offset, *flts.Limit)
	rows, err := tx.QueryContext(ctx, selectQuery)
	if err != nil {
		return payment.PaymentResponse{}, errors.E(op, err, errors.KindUnexpected)
	}
	defer rows.Close()

	var payments = []payment.Payment{}

	for rows.Next() {
		pmt := payment.Payment{}
		err := rows.Scan(
			&pmt.ID,
			&pmt.Fname,
			&pmt.Lname,
			&pmt.Phone,
			&pmt.PropertyID,
			&pmt.Amount,
		)
		if err != nil {
			return payment.PaymentResponse{}, errors.E(op, err, errors.KindUnexpected)
		}
		payments = append(payments, pmt)
	}

	countQuery := `SELECT COUNT(*), COALESCE(SUM(i.amount), 0.0) FROM invoices i JOIN properties p ON i.property = p.id`
	countQuery += " WHERE 1 = 1"

	if flts.Status != nil {
		countQuery += fmt.Sprintf(" AND i.status = '%s'", *flts.Status)
	}

	// check on from date
	if flts.From != nil {
		countQuery += fmt.Sprintf(" AND i.created_at >= '%s'", *flts.From)
	}
	// check on to date
	if flts.To != nil {
		countQuery += fmt.Sprintf(" AND i.created_at <= '%s'", *flts.To)
	}

	if flts.Sector != nil {
		countQuery += fmt.Sprintf(" AND p.sector = '%s'", *flts.Sector)
	}

	if flts.Cell != nil {
		countQuery += fmt.Sprintf(" AND p.cell = '%s'", *flts.Cell)
	}

	if flts.Village != nil {
		countQuery += fmt.Sprintf(" AND p.village = '%s'", *flts.Village)
	}

	// return data by the namespace
	if flts.Namespace != nil {
		selectQuery += fmt.Sprintf(" AND p.namespace = '%s'", *flts.Namespace)
	}

	var (
		total  uint64
		amount float64
	)

	if err := repo.QueryRow(countQuery).Scan(&total, &amount); err != nil {
		return payment.PaymentResponse{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := payment.PaymentResponse{
		Payments: payments,
		PageMetadata: payment.PageMetadata{
			Total:  total,
			Offset: *flts.Offset,
			Limit:  *flts.Limit,
			Amount: amount,
		},
	}
	return page, nil
}

func (repo *paymentStore) TodayTransaction(ctx context.Context, flts *payment.MetricFilters) (payment.Transaction, error) {
	const op errors.Op = "store/postgres/paymentStore.ListTodaysTransactions"

	tx, err := repo.BeginTx(ctx, nil)
	if err != nil {
		return payment.Transaction{}, errors.E(op, err)
	}
	defer tx.Rollback()

	selectQuery := `
	SELECT 
		COUNT(*) AS total,
		COALESCE(SUM(t.amount), 0) as amount
	FROM 
		transactions t
	JOIN 
		properties p ON t.madefor = p.id
	WHERE  
		DATE(t.created_at) = DATE(now())`

	if flts.Sector != nil {
		selectQuery += fmt.Sprintf(" AND p.sector = '%s'", *flts.Sector)
	}
	if flts.Village != nil {
		selectQuery += fmt.Sprintf(" AND p.village = '%s'", *flts.Village)
	}
	if flts.Cell != nil {
		selectQuery += fmt.Sprintf(" AND p.cell = '%s'", *flts.Cell)
	}
	if flts.Namespace != nil {
		selectQuery += fmt.Sprintf(" AND t.namespace = '%s'", *flts.Namespace)
	}

	selectQuery += ` GROUP BY DATE(t.created_at)`

	rows, err := tx.Query(selectQuery)
	if err != nil {
		return payment.Transaction{}, errors.E(op, err)
	}
	defer rows.Close()

	if rows.Next() {
		var transaction payment.Transaction
		err = rows.Scan(&transaction.Transactions, &transaction.Amount)
		if err != nil {
			return payment.Transaction{}, errors.E(op, err)
		}

		return transaction, tx.Commit()
	}

	return payment.Transaction{}, nil
}

// Implement the ListDailyTransactions
func (repo *paymentStore) ListDailyTransactions(ctx context.Context, flts *payment.MetricFilters) (payment.Transactions, error) {
	const op errors.Op = "store/postgres/paymentStore.ListDailyTransactions"

	tx, err := repo.BeginTx(ctx, nil)
	if err != nil {
		return payment.Transactions{}, errors.E(op, err)
	}
	defer tx.Rollback()

	selectQuery := `
	SELECT 
		COUNT(p.id) AS transactions,
		COALESCE(SUM(t.amount), 0),
		DATE(t.created_at) as date
	FROM 
		transactions t
	JOIN 
		properties p ON t.madefor = p.id
	WHERE 1=1 `

	if flts.Sector != nil {
		selectQuery += fmt.Sprintf(" AND p.sector = '%s'", *flts.Sector)
	}
	if flts.Village != nil {
		selectQuery += fmt.Sprintf(" AND p.village = '%s'", *flts.Village)
	}
	if flts.Cell != nil {
		selectQuery += fmt.Sprintf(" AND p.cell = '%s'", *flts.Cell)
	}
	if flts.From != nil {
		selectQuery += fmt.Sprintf(" AND DATE(t.created_at) >= '%s'", *flts.From)
	}
	if flts.To != nil {
		selectQuery += fmt.Sprintf(" AND DATE(t.created_at) <= '%s'", *flts.To)
	}

	if flts.Namespace != nil {
		selectQuery += fmt.Sprintf(" AND t.namespace = '%s'", *flts.Namespace)
	}

	selectQuery += ` GROUP BY  DATE(t.created_at)`

	selectQuery += ` ORDER BY DATE(t.created_at) DESC`

	selectQuery += fmt.Sprintf(" OFFSET %d LIMIT %d", *flts.Offset, *flts.Limit)

	rows, err := tx.QueryContext(ctx, selectQuery)
	if err != nil {
		return payment.Transactions{}, errors.E(op, err)
	}
	defer rows.Close()

	out := []payment.Transaction{}
	for rows.Next() {
		var transaction payment.Transaction
		err = rows.Scan(&transaction.Transactions, &transaction.Amount, &transaction.Date)
		if err != nil {
			return payment.Transactions{}, errors.E(op, err)
		}

		out = append(out, transaction)

	}

	if err = rows.Err(); err != nil {
		return payment.Transactions{}, errors.E(op, err)
	}

	// calculate total
	countQuery := `
		SELECT
		COUNT(*) OVER() AS total
	FROM
		transactions t
	JOIN
		properties p ON t.madefor = p.id
	WHERE 1=1 `

	if flts.Sector != nil {
		countQuery += fmt.Sprintf(" AND p.sector = '%s'", *flts.Sector)
	}
	if flts.Village != nil {
		countQuery += fmt.Sprintf(" AND p.village = '%s'", *flts.Village)
	}
	if flts.Cell != nil {
		countQuery += fmt.Sprintf(" AND p.cell = '%s'", *flts.Cell)
	}

	if flts.From != nil {
		countQuery += fmt.Sprintf(" AND DATE(t.created_at) >= '%s'", *flts.From)
	}
	if flts.To != nil {
		countQuery += fmt.Sprintf(" AND DATE(t.created_at) <= '%s'", *flts.To)
	}

	if flts.Namespace != nil {
		countQuery += fmt.Sprintf(" AND t.namespace = '%s'", *flts.Namespace)
	}
	countQuery += ` GROUP BY  DATE(t.created_at)`

	var total uint64
	if err := tx.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return payment.Transactions{}, errors.E(op, err, errors.KindUnexpected)
	}
	// return the transactionpage
	resp := payment.Transactions{
		Transactions: out,
		PageMetadata: payment.PageMetadata{
			Total:  total,
			Offset: *flts.Offset,
			Limit:  *flts.Limit,
		},
	}

	return resp, tx.Commit()
}

func (repo *paymentStore) TodaySummary(ctx context.Context, flts *payment.MetricFilters) (payment.Summaries, error) {

	const op errors.Op = "store/postgres/paymentStore.TodaySummary"

	tx, err := repo.BeginTx(ctx, nil)
	if err != nil {
		return payment.Summaries{}, errors.E(op, err)
	}
	defer tx.Rollback()

	selectQuery := `
			SELECT
		COUNT(properties.id) AS total_houses,
		SUM(amount) AS amount,
			cell,
			village,
		DATE(transactions.created_at) AS date
		FROM transactions
		INNER JOIN properties
		ON transactions.madefor = properties.id
		WHERE 1=1
		`
	if flts.Date != nil {
		selectQuery += fmt.Sprintf(" AND DATE(transactions.created_at) = '%s'", *flts.Date)
	} else {
		selectQuery += "AND DATE(transactions.created_at) = DATE(now())"
	}
	if flts.Sector != nil {
		selectQuery += fmt.Sprintf(" AND sector = '%s'", *flts.Sector)
	}
	if flts.Village != nil {
		selectQuery += fmt.Sprintf(" AND village = '%s'", *flts.Village)
	}
	if flts.Cell != nil {
		selectQuery += fmt.Sprintf(" AND cell = '%s'", *flts.Cell)
	}

	if flts.Namespace != nil {
		selectQuery += fmt.Sprintf(" AND transactions.namespace = '%s'", *flts.Namespace)
	}

	selectQuery += ` GROUP BY cell,village,DATE(transactions.created_at)`

	rows, err := tx.Query(selectQuery)
	if err != nil {
		return payment.Summaries{}, errors.E(op, err)
	}
	defer rows.Close()

	items := []payment.Summary{}
	for rows.Next() {
		var transaction payment.Summary
		err = rows.Scan(&transaction.Houses, &transaction.Amount, &transaction.Cell, &transaction.Village, &transaction.Created_at)
		if err != nil {
			return payment.Summaries{}, errors.E(op, err)
		}

		items = append(items, transaction)
	}
	out := payment.Summaries{
		Summaries: items,
	}

	return out, tx.Commit()

}

// UnpaidHouses
func (repo *paymentStore) UnpaidHouses(ctx context.Context, flts *payment.MetricFilters) (payment.PaymentResponse, error) {
	const op errors.Op = "store/postgres/paymentStore.UnpaidHouses"
	tx, err := repo.BeginTx(ctx, nil)
	if err != nil {
		return payment.PaymentResponse{}, errors.E(op, err)
	}
	defer tx.Rollback()

	selectQuery := `SELECT 
			o.id,
			o.fname, 
			o.lname,
			o.phone,
			i.property,
			i.amount,
			p.sector,
			p.village,
			p.cell
			
		FROM 
			owners o
		JOIN properties p 
			ON p.owner = o.id
		JOIN invoices i ON 
			i.property = p.id
		WHERE i.status = 'pending'
	`
	if flts.Username != nil {
		selectQuery += fmt.Sprintf(" AND p.recorded_by = '%s'", *flts.Username)
	}
	if flts.Namespace != nil {
		selectQuery += fmt.Sprintf(" AND p.namespace = '%s'", *flts.Namespace)
	}

	if flts.Month != nil {
		selectQuery += fmt.Sprintf(" AND EXTRACT(MONTH FROM i.created_at) = %d", *flts.Month)
	}

	selectQuery += fmt.Sprintf(" ORDER BY i.created_at DESC OFFSET %d LIMIT %d", *flts.Offset, *flts.Limit)
	rows, err := tx.QueryContext(ctx, selectQuery)
	if err != nil {
		return payment.PaymentResponse{}, errors.E(op, err)
	}
	defer rows.Close()
	var payments = []payment.Payment{}

	for rows.Next() {
		pmt := payment.Payment{}
		err := rows.Scan(
			&pmt.ID,
			&pmt.Fname,
			&pmt.Lname,
			&pmt.Phone,
			&pmt.PropertyID,
			&pmt.Amount,
			&pmt.Sector,
			&pmt.Village,
			&pmt.Cell,
		)
		if err != nil {
			return payment.PaymentResponse{}, errors.E(op, err, errors.KindUnexpected)
		}
		payments = append(payments, pmt)
	}
	countQuery := `SELECT COUNT(*), COALESCE(SUM(i.amount), 0.0) FROM invoices i JOIN properties p ON i.property = p.id`
	countQuery += " WHERE i.status = 'pending'"

	// check on from date
	if flts.Username != nil {
		countQuery += fmt.Sprintf(" AND p.recorded_by = '%s'", *flts.Username)
	}
	if flts.Namespace != nil {
		countQuery += fmt.Sprintf(" AND p.namespace = '%s'", *flts.Namespace)
	}
	if flts.Month != nil {
		countQuery += fmt.Sprintf(" AND EXTRACT(MONTH FROM i.created_at) = %d", *flts.Month)
	}

	var (
		total  uint64
		amount float64
	)

	if err := repo.QueryRow(countQuery).Scan(&total, &amount); err != nil {
		return payment.PaymentResponse{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := payment.PaymentResponse{
		Payments: payments,
		PageMetadata: payment.PageMetadata{
			Total:  total,
			Offset: *flts.Offset,
			Limit:  *flts.Limit,
			Amount: amount,
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
		month = fmt.Sprintf("Wishyuriye Ukwezi kwa: %d\n", int(tx[0].PayedDate.Month()))
	}

	var buf bytes.Buffer

	buf.WriteString(header)
	// buf.WriteString(selectActivity(pr.Address.Sector))
	buf.WriteString(" mu murenge wa ")
	buf.WriteString(fmt.Sprintf("%s.\n", prop.Address.Sector))
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
	now := time.Now()
	at := clock.TimeIn(now.Local().UTC(), clock.EAST)
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
