package payment_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/nanoid"
	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestPaymentValidate(t *testing.T) {
	cases := []struct {
		desc    string
		payment payment.Transaction
		errMsg  string
		kind    int
		err     error
	}{
		{
			desc:    "validate valid payment",
			payment: payment.Transaction{Code: nanoid.New(nil).ID(), Amount: 1000, Phone: "0784607135", Method: "mtn-momo-rw"},
		},
		{
			desc:    "validate with missing house code",
			payment: payment.Transaction{Amount: 1000, Phone: "0784607135"},
			errMsg:  "",
			kind:    http.StatusBadRequest,
		},
		{
			desc:    "validate with zero amount",
			payment: payment.Transaction{Code: nanoid.New(nil).ID(), Phone: "0784607135"},
			errMsg:  "",
			kind:    http.StatusBadRequest,
		},
		{
			desc:    "validate with missing phone payment",
			payment: payment.Transaction{Code: nanoid.New(nil).ID(), Amount: 1000},
			errMsg:  "",
			kind:    http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		err := tc.payment.Validate()
		if err != nil {
			assert.Equal(t, tc.kind, errors.Kind(err), fmt.Sprintf("%s: expected error of kind (%d) got (%d)\n", tc.desc, tc.kind, errors.Kind(err)))
		}
	}
}
