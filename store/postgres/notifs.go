package postgres

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

var _ notifs.Repository = (*notifsRepo)(nil)

type notifsRepo struct {
	*sql.DB
}

// NewNotifsRepository creates a postgres backed notids.Repository
func NewNotifsRepository(db *sql.DB) notifs.Repository {
	return &notifsRepo{db}
}

func (repo *notifsRepo) Save(ctx context.Context, sms notifs.Notification) (notifs.Notification, error) {
	const op errors.Op = "store/postgres/notifsRepo.Save"

	q := `
		INSERT INTO sms_notifications (
			message,
			sender,
			recipients
		) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at;
	`

	var empty notifs.Notification

	err := repo.QueryRow(q,
		sms.Message,
		sms.Sender,
		pq.Array(sms.Recipients),
	).Scan(&sms.ID, &sms.CreatedAt, &sms.UpdatedAt)

	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok {
			switch pqErr.Code.Name() {
			case errDuplicate:
				return empty, errors.E(op, err, "sms id conflict", errors.KindAlreadyExists)
			}
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	return sms, nil
}

func (repo *notifsRepo) Find(ctx context.Context, id string) (notifs.Notification, error) {
	const op errors.Op = "store/postgres/notifsRepo.Find"

	q := `
		SELECT 
			id,
			message,
			sender,
			recipients,
			created_at,
			updated_at
		FROM 
			sms_notifications
		WHERE id=$1
	`
	var sms notifs.Notification

	var recipients []string

	err := repo.QueryRow(q, id).Scan(
		&sms.ID,
		&sms.Message,
		&sms.Sender,
		pq.Array(&recipients),
		&sms.CreatedAt,
		&sms.UpdatedAt,
	)
	if err != nil {
		var empty notifs.Notification

		pqErr, ok := err.(*pq.Error)
		if err == sql.ErrNoRows || ok && errInvalid == pqErr.Code.Name() {
			return empty, errors.E(op, err, "sms not found", errors.KindNotFound)
		}
		return empty, errors.E(op, err, errors.KindUnexpected)
	}
	sms.Recipients = recipients

	return sms, nil
}

func (repo *notifsRepo) List(ctx context.Context, nspace string, offset, limit uint64) (notifs.NoticationPage, error) {
	const op errors.Op = "store/postgres/notifsRepo.List"

	q := `
		SELECT
			id,
			message,
			sender,
			recipients,
			created_at,
			updated_at
		FROM
			sms_notifications
		WHERE sender=$1 OFFSET $2 LIMIT $3
	`

	var items = []notifs.Notification{}

	rows, err := repo.Query(q, nspace, offset, limit)
	if err != nil {
		return notifs.NoticationPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	defer rows.Close()

	for rows.Next() {
		row := notifs.Notification{}

		var recipients []string

		err := rows.Scan(&row.ID, &row.Message, &row.Sender, pq.Array(&recipients), &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return notifs.NoticationPage{}, errors.E(op, err, errors.KindUnexpected)
		}
		items = append(items, row)
	}

	q = `SELECT count(*) FROM sms_notifications WHERE sender=$1`

	var total uint64
	if err := repo.QueryRow(q, nspace).Scan(&total); err != nil {
		return notifs.NoticationPage{}, errors.E(op, err, errors.KindUnexpected)
	}

	page := notifs.NoticationPage{
		Notifications: items,
		PageMetadata: notifs.PageMetadata{
			Total:  total,
			Offset: offset,
			Limit:  limit,
		},
	}
	return page, nil
}

func (repo *notifsRepo) Count(ctx context.Context, nspace string) (uint64, error) {
	const op errors.Op = "store/postgres/notifsRepo.Save"

	q := `SELECT count(*) FROM sms_notifications WHERE sender=$1`

	var total uint64
	if err := repo.QueryRow(q, nspace).Scan(&total); err != nil {
		return 0, errors.E(op, err, errors.KindUnexpected)
	}
	return total, nil
}
