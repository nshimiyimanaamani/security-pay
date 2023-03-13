package fdi_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/quarksgroup/paypack-go/paypack/api"
	"github.com/rugwirobaker/paypack-backend/backends/fdi"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/stretchr/testify/assert"
)

const timeOut = 30 * time.Second

// test environment
var (
	url       = "https://private-15d6f5-fdipaymentsapi.apiary-mock.com/v2"
	appID     = "92234DCC-FE88-4F2E-941B-E44F06F2B12D"
	appSecret = os.Getenv("PAYPACK_PAYMENT_SECRET")
	env       = "development"
)

func newBackend() (payment.Client, error) {

	cli, err := api.New(url, nil)
	if err != nil {
		return nil, err
	}

	return fdi.New(cli, appID, appSecret, env)
}

func TestPull(t *testing.T) {
	t.Skip("Skipping testing in CI environment")

	bck, err := newBackend()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		desc        string
		transaction *payment.TxRequest
		state       string
		err         error
	}{
		{
			desc:        "pull payment request with valid data",
			transaction: &payment.TxRequest{ID: uuid.New().ID(), Amount: 5, MSISDN: "0785447001"},
			state:       "processing",
			err:         nil,
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		status, err := bck.Pull(ctx, tc.transaction)
		assert.Equal(t, tc.err, err, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.err, err))
		assert.Equal(t, tc.state, status.TxState, fmt.Sprintf("%s: expected %s got %s", tc.desc, tc.state, status.TxState))
	}

}
