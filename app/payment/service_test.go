package payment_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rugwirobaker/paypack-backend/app/nanoid"
	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/app/payment/mocks"
	"github.com/stretchr/testify/assert"
)

func newService() payment.Service {
	idp := mocks.NewIdentityProvider()
	backend := mocks.NewBackend()
	opts := &payment.Options{Idp: idp, Backend: backend}
	return payment.New(opts)
}

func TestInitialize(t *testing.T) {
	svc := newService()

	cases := []struct {
		desc    string
		payment payment.Transaction
		status  string
		errKind int
		err     error
	}{
		{
			desc:    "initialize payment with valid data",
			payment: payment.Transaction{Code: nanoid.New(nil).ID(), Amount: 1000, Phone: "0784607135"},
		},
		{
			desc:    "initialize payment with invalid data",
			payment: payment.Transaction{Amount: 1000, Phone: "0784607135"},
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err, status := svc.Initilize(ctx, tc.payment)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.err, err))
		assert.Equal(t, tc.status, status, fmt.Sprintf("%s: expected %s got %s\n", tc.desc, tc.status, status))
	}

}

func TestValidate(t *testing.T) {}
