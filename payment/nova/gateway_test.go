// +build integration

package nova_test

// import (
// 	"fmt"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/rugwirobaker/paypack-backend/app/payment"
// 	"github.com/rugwirobaker/paypack-backend/nova"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// var token = os.Getenv("PAYPACK_PAYMENT_TOKEN")

// func newGateway() payment.Gateway {
// 	cfg := &nova.Config{
// 		Endpoint: "https://novapay.rw/api/v1/novapay/initialize-payment",
// 		Token:    token,
// 		TimeOut:  30 * time.Second,
// 	}
// 	return nova.New(cfg)
// }

// func TestInitiate(t *testing.T) {
// 	expected := payment.Message{
// 		"message": "ubusabe bwawe bwo kwishyura bwakiriwe kanda *182*7# kwishyura kuri mobile money",
// 		"action":  200.0,
// 	}

// 	pgate := newGateway()

// 	payment := payment.Payment{
// 		Code:   "123434",
// 		Phone:  "+250785868145",
// 		Amount: 5,
// 	}

// 	actual, err := pgate.Initiate(payment)
// 	require.Nil(t, err, fmt.Sprintf("unexpected error %v", err))

// 	assert.Equal(t, expected, actual, fmt.Sprintf(("expected status: '%v' got '%v'"), expected, actual))
// }
