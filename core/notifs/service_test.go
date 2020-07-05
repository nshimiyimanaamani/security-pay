package notifs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/core/notifs/mocks"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	svc := newService()

	const op = "core/sms/service.Send"

	cases := []struct {
		desc       string
		sms        string
		recipients []string
		sender     string
		err        error
	}{
		{
			desc:       "sending valid sms",
			sms:        "hello",
			recipients: []string{"0789657200"},
			sender:     "Remera",
			err:        nil,
		},
		{
			desc:       "sending sms with no recipients",
			sms:        "hello",
			recipients: []string{},
			sender:     "Remera",
			err:        errors.E(op, "the number of recipients must be greater than zero"),
		},
		{
			desc:       "sending sms with empty sms body",
			sms:        "",
			recipients: []string{"0789657200"},
			sender:     "Remera",
			err:        errors.E(op, "sms body must not be empty"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		sms := notifs.Notification{Message: tc.sms, Recipients: tc.recipients, Sender: tc.sender}
		_, err := svc.Send(ctx, sms)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestFind(t *testing.T) {
	svc := newService()

	sms := notifs.Notification{
		Message:    "Hello",
		Recipients: []string{"0789657200"},
		Sender:     "Remera",
	}

	sms, err := svc.Send(context.Background(), sms)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op = "core/sms/service.Find"

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
			err:  errors.E(op, "sms not found"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := svc.Find(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}

}

func TestList(t *testing.T) {
	svc := newService()

	namespace := "Remera"

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		sms := notifs.Notification{
			Message:    "Hello",
			Recipients: []string{"0789657200"},
			Sender:     namespace,
		}

		sms, err := svc.Send(context.Background(), sms)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	cases := []struct {
		desc   string
		nspace string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc:   "list all accounts",
			nspace: namespace,
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the accounts",
			nspace: namespace,
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			nspace: namespace,
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			nspace: namespace,
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		page, err := svc.List(ctx, tc.nspace, tc.offset, tc.limit)
		size := uint64(len(page.Notifications))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}

}

func TestCount(t *testing.T) {
	svc := newService()

	namespace := "Remera"

	expected := uint64(10)
	for i := uint64(0); i < expected; i++ {
		sms := notifs.Notification{
			Message:    "Hello",
			Recipients: []string{"0789657200"},
			Sender:     namespace,
		}

		sms, err := svc.Send(context.Background(), sms)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	got, err := svc.Count(context.Background(), namespace)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	assert.Equal(t, expected, got, fmt.Sprintf("expected %d got %d", expected, got))
}

func newService() notifs.Service {
	opts := &notifs.Options{
		Backend: mocks.NewBackend(),
		IDP:     mocks.NewIdentityProvider(),
		Store:   mocks.NewRepository(),
	}
	return notifs.New(opts)
}
