package fdi_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/rugwirobaker/paypack-backend/backends/fdi"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const timeOut = 30 * time.Second

// test environment
var (
	url       = "https://private-15d6f5-fdipaymentsapi.apiary-mock.com/v2"
	appID     = "92234DCC-FE88-4F2E-941B-E44F06F2B12D"
	appSecret = os.Getenv("PAYPACK_PAYMENT_SECRET")
	callback  = "https://codechef-inlets.herokuapp.com"
)

func newBackend() payment.Client {
	opts := &fdi.ClientOptions{
		URL:       url,
		AppID:     appID,
		AppSecret: appSecret,
		Callback:  callback,
	}
	return fdi.New(opts)
}

func TestStatus(t *testing.T) {
	t.Skip("Skipping testing with external deps")
	cli := newBackend()

	ctx := context.Background()

	status, err := cli.Status(ctx)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))
	require.Equal(t, http.StatusOK, status, fmt.Sprintf("expected %d got %d", http.StatusOK, status))
}

func TestPull(t *testing.T) {
	t.Skip("Skipping testing in CI environment")
	bck := newBackend()

	ctx := context.Background()

	status, _ := bck.Status(ctx)
	require.Equal(t, http.StatusOK, status, fmt.Sprintf("expected %d got %d", http.StatusOK, status))

	cases := []struct {
		desc        string
		transaction payment.Transaction
		state       string
		err         error
	}{
		{
			desc:        "pull payment request with valid data",
			transaction: payment.Transaction{ID: uuid.New().ID(), Amount: 5, Phone: "0785447001"},
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
