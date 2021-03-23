package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/accounts"
	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveNotification(t *testing.T) {
	repo := postgres.NewNotifsRepository(db)
	defer CleanDB(t, db)

	account := accounts.Account{
		ID:            "gasabo.remera",
		Name:          "remera",
		NumberOfSeats: 10,
		Type:          accounts.Devs,
	}
	account = saveAccount(t, db, account)

	const op errors.Op = "store/postgres/notifsRepo.Save"

	cases := []struct {
		desc       string
		sms        string
		recipients []string
		sender     string
		err        error
	}{
		{
			desc:       "saving valid sms",
			sms:        "hello",
			recipients: []string{"0789657200"},
			sender:     account.ID,
			err:        nil,
		},
	}

	for _, tc := range cases {
		sms := notifs.Notification{Message: tc.sms, Sender: tc.sender, Recipients: tc.recipients}
		sms, err := repo.Save(context.Background(), sms)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%s' got err: '%s'", tc.desc, tc.err, err))
		_, err = uuid.FromString(sms.ID)
		assert.Nil(t, err, fmt.Sprintf("err:%v is not nil", err))
	}
}
func TestFindNotification(t *testing.T) {
	const op errors.Op = "store/postgres/notifsRepo.Find"

	repo := postgres.NewNotifsRepository(db)
	defer CleanDB(t, db)

	account := accounts.Account{
		ID:            "gasabo.remera",
		Name:          "remera",
		NumberOfSeats: 10,
		Type:          accounts.Devs,
	}
	account = saveAccount(t, db, account)

	sms := notifs.Notification{
		Message:    "Hello",
		Recipients: []string{"0789657200"},
		Sender:     account.ID,
	}
	sms, err := repo.Save(context.Background(), sms)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "find sent sms",
			id:   sms.ID,
			err:  nil,
		},
		{
			desc: "find a non existant message",
			id:   "invalid",
			err:  errors.E(op, "sms not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := repo.Find(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%s' got err: '%s'", tc.desc, tc.err, err))
	}
}
func TestListNotifications(t *testing.T) {
	const op errors.Op = "store/postgres/notifsRepo.List"

	repo := postgres.NewNotifsRepository(db)
	defer CleanDB(t, db)

	account := accounts.Account{
		ID:            "gasabo.remera",
		Name:          "remera",
		NumberOfSeats: 10,
		Type:          accounts.Devs,
	}
	account = saveAccount(t, db, account)

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		sms := notifs.Notification{
			Message:    "Hello",
			Recipients: []string{"0789657200"},
			Sender:     account.ID,
		}
		sms, err := repo.Save(context.Background(), sms)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	cases := map[string]struct {
		sender string
		offset uint64
		limit  uint64
		size   uint64
		total  uint64
	}{
		"list all accounts": {
			sender: account.ID,
			offset: 0,
			limit:  n,
			size:   n,
			total:  n,
		},
		"list half of the accounts": {
			sender: account.ID,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			total:  n,
		},
		"list empty set": {
			sender: account.ID,
			offset: n + 1,
			limit:  n,
			size:   0,
			total:  n,
		},
		"list with zero limit": {
			sender: account.ID,
			offset: 1,
			limit:  0,
			size:   0,
			total:  n,
		},
	}

	for desc, tc := range cases {
		ctx := context.Background()
		page, err := repo.List(ctx, tc.sender, tc.offset, tc.limit)
		size := uint64(len(page.Notifications))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.size, size))
		assert.Equal(t, tc.total, page.Total, fmt.Sprintf("%s: expected %d got %d\n", desc, tc.total, page.Total))
		assert.Nil(t, err, fmt.Sprintf("%s: expected no error got '%v'\n", desc, err))
	}

}
func TestCountNotifications(t *testing.T) {
	const op errors.Op = "store/postgres/notifsRepo.List"

	repo := postgres.NewNotifsRepository(db)
	defer CleanDB(t, db)

	account := accounts.Account{
		ID:            "gasabo.remera",
		Name:          "remera",
		NumberOfSeats: 10,
		Type:          accounts.Devs,
	}
	account = saveAccount(t, db, account)

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		sms := notifs.Notification{
			Message:    "Hello",
			Recipients: []string{"0789657200"},
			Sender:     account.ID,
		}
		sms, err := repo.Save(context.Background(), sms)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	cases := []struct {
		desc     string
		nspace   string
		expected uint64
		err      error
	}{
		{
			desc:     "count messages from valid namespace(sender)",
			nspace:   account.ID,
			expected: n,
			err:      nil,
		},
		{
			desc:     "count messages from invvalid namespace(sender)",
			nspace:   "invalid",
			expected: 0,
			err:      nil,
		},
	}

	for _, tc := range cases {
		got, err := repo.Count(context.Background(), tc.nspace)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%s' got err: '%s'", tc.desc, tc.err, err))
		assert.Equal(t, tc.expected, got, fmt.Sprintf("expected %d got %d", tc.expected, got))

	}
}
