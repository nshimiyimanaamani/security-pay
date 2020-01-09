package feedback_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/feedback"
	"github.com/rugwirobaker/paypack-backend/app/feedback/mocks"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var wrong = "wrong_value"

func newService() feedback.Service {
	opts := &feedback.Options{
		Idp:  mocks.NewIdentityProvider(),
		Repo: mocks.NewRepository(),
	}
	return feedback.New(opts)
}

func TestRecord(t *testing.T) {
	svc := newService()

	const op errors.Op = "app/feedback/service.Record"

	cases := []struct {
		desc string
		msg  feedback.Message
		err  error
	}{
		{
			desc: "record valid message",
			msg:  feedback.Message{Title: "title", Body: "body", Creator: "0784677882"},
			err:  nil,
		},
		{
			desc: "record message without body",
			msg:  feedback.Message{Title: "title"},
			err:  errors.E(op, "invalid message: missing body"),
		},
		{
			desc: "record message without title",
			msg:  feedback.Message{Body: "body"},
			err:  errors.E(op, "invalid message: missing title"),
		},
		{
			desc: "record message without creator",
			msg:  feedback.Message{Body: "body", Title: "title"},
			err:  errors.E(op, "invalid message: missing creator"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := svc.Record(ctx, &tc.msg)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestRetrieve(t *testing.T) {
	svc := newService()

	ctx := context.Background()
	msg := feedback.Message{Title: "title", Body: "body", Creator: "0784677882"}

	saved, err := svc.Record(ctx, &msg)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "app/feedback/service.Retrieve"

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "retrieve existant message",
			id:   saved.ID,
			err:  nil,
		},
		{
			desc: "retrieve non-existant message",
			id:   wrong,
			err:  errors.E(op, "message not found"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		_, err := svc.Retrieve(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}

}

func TestUpdate(t *testing.T) {
	svc := newService()

	ctx := context.Background()
	msg := feedback.Message{Title: "title", Body: "body", Creator: "0784677882"}

	saved, err := svc.Record(ctx, &msg)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "app/feedback/service.Update"

	cases := []struct {
		desc string
		msg  feedback.Message
		err  error
	}{
		{
			desc: "update existant message",
			msg:  *saved,
			err:  nil,
		},

		{
			desc: "update non existant message",
			msg:  feedback.Message{Title: "title", Body: "body", Creator: "0784677882"},
			err:  errors.E(op, "message not found"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := svc.Update(ctx, tc.msg)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestDelete(t *testing.T) {
	svc := newService()

	ctx := context.Background()
	msg := feedback.Message{Title: "title", Body: "body", Creator: "0784677882"}
	saved, err := svc.Record(ctx, &msg)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	const op errors.Op = "app/feedback/service.Delete"

	cases := []struct {
		desc string
		id   string
		err  error
	}{
		{
			desc: "delete existant message",
			id:   saved.ID,
			err:  nil,
		},
		{
			desc: "delete non-existant message",
			id:   wrong,
			err:  errors.E(op, "message not found"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := svc.Delete(ctx, tc.id)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}

func TestList(t *testing.T) {
	svc := newService()

	const op errors.Op = "app/accounts/service.List"

	msg := feedback.Message{Title: "title", Body: "body", Creator: "0784677882"}

	n := uint64(10)
	for i := uint64(0); i < n; i++ {
		ctx := context.Background()
		_, err := svc.Record(ctx, &msg)
		require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))
	}

	cases := []struct {
		desc   string
		offset uint64
		limit  uint64
		size   uint64
		err    error
	}{
		{
			desc:   "list all accounts",
			offset: 0,
			limit:  n,
			size:   n,
			err:    nil,
		},
		{
			desc:   "list half of the accounts",
			offset: n / 2,
			limit:  n,
			size:   n / 2,
			err:    nil,
		},
		{
			desc: "	list empty set",
			offset: n + 1,
			limit:  n,
			size:   0,
			err:    nil,
		},
		{
			desc:   "list with zero limit",
			offset: 1,
			limit:  0,
			size:   0,
			err:    nil,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		page, err := svc.List(ctx, tc.offset, tc.limit)
		size := uint64(len(page.Messages))
		assert.Equal(t, tc.size, size, fmt.Sprintf("%s: expected %d got %d\n", tc.desc, tc.size, size))
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
	}
}
