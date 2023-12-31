package payment_test

import (
	"fmt"
	"testing"

	"github.com/nshimiyimanaamani/paypack-backend/core/identity/uuid"
	"github.com/nshimiyimanaamani/paypack-backend/core/nanoid"
	"github.com/nshimiyimanaamani/paypack-backend/core/payment"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestPaymentHasCode(t *testing.T) {
	const op errors.Op = "core/payment/Payment.HasCode"

	cases := []struct {
		desc    string
		payment payment.TxRequest
		errMsg  string
		err     error
	}{
		{
			desc:    "validate valid payment",
			payment: payment.TxRequest{Code: nanoid.New(nil).ID(), Amount: 1000, MSISDN: "0784607135", Method: "mtn-momo-rw"},
			err:     nil,
		},
		{
			desc:    "validate with missing house code",
			payment: payment.TxRequest{Amount: 1000, MSISDN: "0784607135", Method: "mtn-momo-rw"},
			err:     errors.E(op, "missing house code", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		err := tc.payment.HasCode()
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))

	}
}

func TestPaymentReady(t *testing.T) {
	const op errors.Op = "core/payment/Payment.Ready"

	cases := []struct {
		desc    string
		payment payment.TxRequest
		errMsg  string
		err     error
	}{
		{
			desc:    "validate valid payment",
			payment: payment.TxRequest{Code: nanoid.New(nil).ID(), Amount: 1000, MSISDN: "0784607135", Method: "mtn-momo-rw"},
			err:     nil,
		},
		{
			desc:    "validate with zero amount",
			payment: payment.TxRequest{Code: nanoid.New(nil).ID(), MSISDN: "0784607135", Method: "mtn-momo-rw"},
			err:     errors.E(op, "amount must be greater than zero", errors.KindBadRequest),
		},
		{
			desc:    "validate with missing phone payment",
			payment: payment.TxRequest{Code: nanoid.New(nil).ID(), Amount: 1000, Method: payment.MTN},
			err:     errors.E(op, "missing phone number", errors.KindBadRequest),
		},
		{
			desc:    "validate with missing payment method",
			payment: payment.TxRequest{Code: nanoid.New(nil).ID(), Amount: 1000, MSISDN: "0784607135"},
			err:     errors.E(op, "payment method must be specified", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		err := tc.payment.Ready()
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))

	}
}

func TestValidateCallback(t *testing.T) {
	const op errors.Op = "core/payment/Callback.Validate"

	cases := []struct {
		desc     string
		callback payment.Callback
		err      error
	}{
		{
			desc: "validate valid callback",
			callback: payment.Callback{
				Kind: "transaction:processed",
				Data: payment.Data{
					Ref:    uuid.New().ID(),
					Status: string(payment.Successful),
				},
			},
			err: nil,
		},
		{
			desc: "validate callback without status",
			callback: payment.Callback{
				Kind: "transaction:processed",
				Data: payment.Data{
					Ref:    uuid.New().ID(),
					Status: string(payment.Failed),
				},
			},
			err: errors.E(op, "status field must not be empty", errors.KindBadRequest),
		},
		{
			desc: "validate without gwRef",
			callback: payment.Callback{
				Kind: "transaction:processed",
				Data: payment.Data{
					Ref:    uuid.New().ID(),
					Status: string(payment.Failed),
				},
			},
			err: errors.E(op, "gwRef field must not be empty", errors.KindBadRequest),
		},
		{
			desc: "validate without transaction reference",
			callback: payment.Callback{
				Kind: "transaction:processed",
				Data: payment.Data{
					Ref:    uuid.New().ID(),
					Status: string(payment.Failed),
				},
			},
			err: errors.E(op, "trxRef field must not be empty", errors.KindBadRequest),
		},
		{
			desc: "validate without state reference",
			callback: payment.Callback{
				Kind: "transaction:processed",
				Data: payment.Data{
					Ref:    uuid.New().ID(),
					Status: string(payment.Failed),
				},
			},
			err: errors.E(op, "state field must not be empty", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		err := tc.callback.Validate()
		if err != nil {
			assert.True(t, errors.ErrEqual(tc.err, err), fmt.Sprintf("%s: expected %s got '%s'\n", tc.desc, tc.err, err))
		}
	}

}
