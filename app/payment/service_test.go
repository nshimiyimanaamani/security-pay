package payment_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/identity/uuid"
	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/app/payment/mocks"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newService(properties []string) payment.Service {
	idp := mocks.NewIdentityProvider()
	backend := mocks.NewBackend()
	queue := mocks.NewQueue()
	repo := mocks.NewRepository(properties)
	opts := &payment.Options{Idp: idp, Backend: backend, Queue: queue, Repo: repo}
	return payment.New(opts)
}

func TestInitialize(t *testing.T) {
	code := uuid.New().ID()
	properties := []string{code}
	svc := newService(properties)

	const op errors.Op = "app.payment.Initialize"

	cases := []struct {
		desc    string
		payment payment.Transaction
		state   string
		errKind int
		err     error
	}{
		{
			desc:    "initialize payment with valid data",
			payment: payment.Transaction{Code: code, Amount: 1000, Phone: "0784607135", Method: "mtn-momo-rw"},
			state:   "processing",
			err:     nil,
		},
		{
			desc:    "initialize payment with invalid data",
			payment: payment.Transaction{Code: code, Amount: 1000, Phone: "0784607135"},
			err:     errors.E(op, "payment method must be specified", errors.KindBadRequest),
		},
		{
			desc:    "initialize payment with unsaved house code",
			payment: payment.Transaction{Code: uuid.New().ID(), Amount: 1000, Phone: "0784607135", Method: "mtn-momo-rw"},
			err:     errors.E(op, "property not found", errors.KindNotFound),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		status, err := svc.Initilize(ctx, tc.payment)
		assert.True(t, errors.ErrEqual(tc.err, err), fmt.Sprintf("%s: expected '%s' got '%s'\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.state, status.TxState, fmt.Sprintf("%s: expected %s got '%s'\n", tc.desc, tc.state, status.TxState))
	}

}

func TestConfirm(t *testing.T) {
	code := uuid.New().ID()
	properties := []string{code}
	svc := newService(properties)

	tx := payment.Transaction{
		ID:     uuid.New().ID(),
		Code:   code,
		Amount: 1000, Phone: "0784607135",
		Method: "mtn-momo-rw",
	}

	res, err := svc.Initilize(context.Background(), tx)
	require.Nil(t, err, fmt.Sprintf("unexpected error: %s", err))

	const op errors.Op = "payment.Confirm"

	cases := []struct {
		desc     string
		callback payment.Callback
		err      error
	}{
		{
			desc: "confirm valid callback",
			callback: payment.Callback{
				Status: "success",
				Data: payment.CallBackData{
					GwRef:  uuid.New().ID(),
					TrxRef: res.TxID,
					State:  "success",
				},
			},
			err: nil,
		},
		{
			desc: "confirm invalid callback",
			callback: payment.Callback{
				Data: payment.CallBackData{
					GwRef:  uuid.New().ID(),
					TrxRef: uuid.New().ID(),
					State:  "success",
				},
			},
			err: errors.E(op, "status field must not be empty", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := svc.Confirm(ctx, tc.callback)
		assert.True(t, errors.ErrEqual(tc.err, err), fmt.Sprintf("%s: expected %s got '%s'\n", tc.desc, tc.err, err))
	}
}
