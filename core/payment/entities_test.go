package payment_test

import (
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/core/identity/uuid"
	"github.com/rugwirobaker/paypack-backend/core/nanoid"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestPaymentValidate(t *testing.T) {
	const op errors.Op = "payment.Transaction.Validate"

	cases := []struct {
		desc    string
		payment payment.Transaction
		errMsg  string
		err     error
	}{
		{
			desc:    "validate valid payment",
			payment: payment.Transaction{Code: nanoid.New(nil).ID(), Amount: 1000, Phone: "0784607135", Method: "mtn-momo-rw"},
			err:     nil,
		},
		{
			desc:    "validate with missing house code",
			payment: payment.Transaction{Amount: 1000, Phone: "0784607135", Method: "mtn-momo-rw"},
			err:     errors.E(op, "missing house code", errors.KindBadRequest),
		},
		{
			desc:    "validate with zero amount",
			payment: payment.Transaction{Code: nanoid.New(nil).ID(), Phone: "0784607135", Method: "mtn-momo-rw"},
			err:     errors.E(op, "amount must be greater than zero", errors.KindBadRequest),
		},
		{
			desc:    "validate with missing phone payment",
			payment: payment.Transaction{Code: nanoid.New(nil).ID(), Amount: 1000, Method: payment.MTN},
			err:     errors.E(op, "missing phone number", errors.KindBadRequest),
		},
		{
			desc:    "validate with missing payment method",
			payment: payment.Transaction{Code: nanoid.New(nil).ID(), Amount: 1000, Phone: "0784607135"},
			err:     errors.E(op, "payment method must be specified", errors.KindBadRequest),
		},
	}

	for _, tc := range cases {
		err := tc.payment.Validate()
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))

	}
}

func TestValidateCallback(t *testing.T) {
	const op errors.Op = ""

	cases := []struct {
		desc     string
		callback payment.Callback
		err      error
	}{
		{
			desc: "validate valid callback",
			callback: payment.Callback{
				Status: "success",
				Data: payment.CallBackData{
					GwRef:  uuid.New().ID(),
					TrxRef: uuid.New().ID(),
					State:  payment.Successful,
				},
			},
			err: nil,
		},
		{
			desc: "validate callback without status",
			callback: payment.Callback{
				Data: payment.CallBackData{
					GwRef:  uuid.New().ID(),
					TrxRef: uuid.New().ID(),
					State:  payment.Failed,
				},
			},
			err: errors.E(op, "status field must not be empty", errors.KindBadRequest),
		},
		{
			desc: "validate without gwRef",
			callback: payment.Callback{
				Status: "success",
				Data: payment.CallBackData{
					TrxRef: uuid.New().ID(),
					State:  payment.Failed,
				},
			},
			err: errors.E(op, "gwRef field must not be empty", errors.KindBadRequest),
		},
		{
			desc: "validate without transaction reference",
			callback: payment.Callback{
				Status: "success",
				Data: payment.CallBackData{
					GwRef: uuid.New().ID(),
					State: payment.Failed,
				},
			},
			err: errors.E(op, "trxRef field must not be empty", errors.KindBadRequest),
		},
		{
			desc: "validate without state reference",
			callback: payment.Callback{
				Status: "success",
				Data: payment.CallBackData{
					TrxRef: uuid.New().ID(),
					GwRef:  uuid.New().ID(),
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
