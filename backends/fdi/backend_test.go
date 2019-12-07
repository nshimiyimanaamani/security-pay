package fdi_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/backends/fdi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const timeOut = 30 * time.Second

// test environment
var (
	url       = "https://private-15d6f5-fdipaymentsapi.apiary-mock.com/v2"
	accountID = "92234DCC-FE88-4F2E-941B-E44F06F2B12D"
	appSecret = os.Getenv("PAYPACK_PAYMENT_SECRET")
	callback  = "https://codechef-inlets.herokuapp.com"
)

func newBackend() payment.Backend {
	opts := &fdi.ClientOptions{
		URL:       url,
		AppID:     accountID,
		AppSecret: appSecret,
		Callback:  callback,
	}
	bc := fdi.NewBackend(opts)
	return bc
}

func TestStatus(t *testing.T) {
	bck := newBackend()

	ctx := context.Background()

	status, err := bck.Status(ctx)

	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))
	require.Equal(t, http.StatusOK, status, fmt.Sprintf("expected %d got %d", http.StatusOK, status))
}

func TestAuth(t *testing.T) {
	bck := newBackend()

	ctx := context.Background()

	status, _ := bck.Status(ctx)
	require.Equal(t, http.StatusOK, status, fmt.Sprintf("expected %d got %d", http.StatusOK, status))

	err := bck.Auth(ctx)
	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))
}

func TestPull(t *testing.T) {
	bck := newBackend()

	ctx := context.Background()

	status, _ := bck.Status(ctx)
	require.Equal(t, http.StatusOK, status, fmt.Sprintf("expected %d got %d", http.StatusOK, status))

	err := bck.Auth(ctx)
	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

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
