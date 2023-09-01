package payment_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nshimiyimanaamani/paypack-backend/core/identity/uuid"
	"github.com/nshimiyimanaamani/paypack-backend/core/invoices"
	"github.com/nshimiyimanaamani/paypack-backend/core/notifs"
	"github.com/nshimiyimanaamani/paypack-backend/core/owners"
	"github.com/nshimiyimanaamani/paypack-backend/core/payment"
	"github.com/nshimiyimanaamani/paypack-backend/core/payment/mocks"
	"github.com/nshimiyimanaamani/paypack-backend/core/properties"
	"github.com/nshimiyimanaamani/paypack-backend/core/transactions"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const namespace = "remera"

func TestPull(t *testing.T) {
	const op errors.Op = "core/payment/service.Pull"

	owners, owner := newOwnersStore()
	properties, property := newPropertiesStore(owner)
	invoices, invoice := newInvoiceStore(property)
	svc := newService(owners, properties, invoices)

	cases := []struct {
		desc    string
		payment *payment.TxRequest
		state   payment.State
		err     error
	}{
		{
			desc:    "initialize payment with valid data",
			payment: &payment.TxRequest{Code: property.ID, Amount: invoice.Amount, MSISDN: "0784607135", Method: "mtn-momo-rw"},
			state:   "processing",
			err:     nil,
		},
		{
			desc:    "initialize payment with invalid data",
			payment: &payment.TxRequest{Code: property.ID, Amount: invoice.Amount, MSISDN: "0784607135"},
			state:   "failed",
			err:     errors.E(op, "payment method must be specified"),
		},
		{
			desc:    "initialize payment with unsaved house property.ID",
			payment: &payment.TxRequest{Code: uuid.New().ID(), Amount: invoice.Amount, MSISDN: "0784607135", Method: "mtn-momo-rw"},
			state:   "failed",
			err:     errors.E(op, "property not found"),
		},
		{
			desc:    "initialize payment with invalid amount(different from invoice)",
			payment: &payment.TxRequest{Code: property.ID, Amount: 100, MSISDN: "0784607135", Method: "mtn-momo-rw"},
			state:   "failed",
			err:     errors.E(op, "amount doesn't match invoice"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		status, err := svc.Pull(ctx, tc.payment)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
		assert.Equal(t, tc.state, status.TxState, fmt.Sprintf("%s: expected %s got '%s'\n", tc.desc, tc.state, status.TxState))
	}

}

func TestConfirmPull(t *testing.T) {
	const op errors.Op = "core/payment/service.ConfirmPull"

	owners, owner := newOwnersStore()
	properties, property := newPropertiesStore(owner)
	invoices, _ := newInvoiceStore(property)
	svc := newService(owners, properties, invoices)

	tx := &payment.TxRequest{
		ID:     uuid.New().ID(),
		Code:   property.ID,
		Amount: 1000, MSISDN: "0784607135",
		Method: "mtn-momo-rw",
	}

	res, err := svc.Pull(context.Background(), tx)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc     string
		callback payment.Callback
		err      error
	}{
		{
			desc: "confirm valid callback",
			callback: payment.Callback{
				Kind: "transaction:processed",
				Data: payment.Data{
					Ref:    res.TxID,
					Status: string(payment.Successful),
				},
			},
			err: nil,
		},
		{
			desc: "confirm invalid callback",
			callback: payment.Callback{
				Kind: "transaction:processed",
				Data: payment.Data{
					Ref:    res.TxID,
					Status: string(payment.Successful),
				},
			},
			err: errors.E(op, "status field must not be empty"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := svc.ProcessHook(ctx, tc.callback)
		assert.True(t, errors.ErrEqual(tc.err, err), fmt.Sprintf("%s: expected %s got '%s'\n", tc.desc, tc.err, err))
	}
}

func TestPush(t *testing.T) {
	const op errors.Op = "core/payment/service.Push"

	owners, owner := newOwnersStore()
	properties, property := newPropertiesStore(owner)
	invoices, _ := newInvoiceStore(property)
	svc := newService(owners, properties, invoices)

	cases := []struct {
		desc    string
		payment *payment.TxRequest
		state   payment.State
		err     error
	}{
		{
			desc:    "initialize payment with valid data",
			payment: &payment.TxRequest{Code: property.ID, Amount: 5000, MSISDN: "0784607135", Method: "mtn-momo-rw"},
			state:   "processing",
			err:     nil,
		},
		{
			desc:    "initialize payment with invalid method",
			payment: &payment.TxRequest{Code: property.ID, Amount: 5000, MSISDN: "0784607135"},
			state:   "failed",
			err:     errors.E(op, "payment method must be specified"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		status, err := svc.Push(ctx, tc.payment)
		assert.True(t, errors.Match(tc.err, err), fmt.Sprintf("%s: expected err: '%v' got err: '%v'", tc.desc, tc.err, err))
		assert.Equal(t, tc.state, status.TxState, fmt.Sprintf("%s: expected %s got '%s'\n", tc.desc, tc.state, status.TxState))
	}

}

func TestConfirmPush(t *testing.T) {
	const op errors.Op = "core/payment/service.ConfirmPush"

	owners, owner := newOwnersStore()
	properties, property := newPropertiesStore(owner)
	invoices, _ := newInvoiceStore(property)
	svc := newService(owners, properties, invoices)

	tx := &payment.TxRequest{
		ID:     uuid.New().ID(),
		Code:   property.ID,
		Amount: 1000, MSISDN: "0784607135",
		Method: "mtn-momo-rw",
	}

	res, err := svc.Push(context.Background(), tx)
	require.Nil(t, err, fmt.Sprintf("unexpected error: '%v'", err))

	cases := []struct {
		desc     string
		callback payment.Callback
		err      error
	}{
		{
			desc: "confirm valid callback",
			callback: payment.Callback{
				Kind: "transaction:processed",
				Data: payment.Data{
					Ref:    res.TxID,
					Status: string(payment.Successful),
				},
			},
			err: nil,
		},
		{
			desc: "confirm invalid callback",
			callback: payment.Callback{
				Kind: "transaction:processed",
				Data: payment.Data{
					Ref:    res.TxID,
					Status: string(payment.Successful),
				},
			},
			err: errors.E(op, "status field must not be empty"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := svc.ConfirmPush(ctx, tc.callback)
		assert.True(t, errors.ErrEqual(tc.err, err), fmt.Sprintf("%s: expected %s got '%s'\n", tc.desc, tc.err, err))
	}
}

func TestFormatMessage(t *testing.T) {

	p := properties.Property{
		ID: "22C95179",
		Address: properties.Address{
			Sector: "Remera",
		},
		Due: 1000,
	}

	in := invoices.Invoice{
		ID:     234,
		Amount: p.Due,
	}

	tx := transactions.Transaction{
		ID:      "1",
		Amount:  p.Due,
		MadeFor: "22C95179",
		Invoice: in.ID,
	}

	py := payment.TxRequest{
		Amount: in.Amount,
		MSISDN: "0788123501",
	}

	owner := owners.Owner{
		Fname: "Todd",
		Lname: "Cantwell",
	}

	expected := `Murakoze kwishyura umusanzu w' umutekano mu murenge wa Remera.

Nimero yishyuriweho: 0788123501
Itariki: 1 Sep 2020 08:56
Wishyuriye Ukwezi kwa: 1
Nimero ya fagitire: 234
Umubare w' amafaranga: 1000RWF
Inzu yishyuriwe ni iya Todd Cantwell
Code y' inzu ni: 22C95179`

	got := payment.FormatMessage(tx, in, py, owner, p, "1 Sep 2020 08:56")

	assert.Equal(t, expected, got, fmt.Sprintf("expected '%s', got '%s'", expected, got))
}

func newService(ws owners.Repository, ps properties.Repository, vc invoices.Repository) payment.Service {
	var opts payment.Options
	opts.Owners = ws
	opts.Properties = ps
	opts.Invoices = vc
	opts.Repository = mocks.NewPaymentRepository()
	opts.SMS = newSMSService()
	opts.Idp = mocks.NewIdentityProvider()
	opts.Backend = mocks.NewBackend()
	opts.Queue = mocks.NewQueue()
	opts.Transactions = mocks.NewTransactionsRepository()
	return payment.New(&opts)
}

func newSMSService() notifs.Service {
	var opts notifs.Options
	opts.IDP = uuid.New()
	opts.Backend = mocks.NewSMSBackend()
	opts.Store = mocks.NewSMSRepository()
	return notifs.New(&opts)
}

func newPropertiesStore(owner owners.Owner) (properties.Repository, properties.Property) {
	var property properties.Property
	property.ID = uuid.New().ID()
	property.Due = 1000
	property.Owner = properties.Owner{ID: owner.ID}
	store := mocks.NewPropertyRepository()
	property, _ = store.Save(context.Background(), property)
	return store, property
}

func newOwnersStore() (owners.Repository, owners.Owner) {
	var owner owners.Owner
	owner.ID = uuid.New().ID()
	owner.Fname = "Jamie"
	owner.Lname = "Jones"
	owner.Phone = "0787205106"
	store := mocks.NewOwnersRepository()
	owner, _ = store.Save(context.Background(), owner)
	return store, owner
}
func newInvoiceStore(property properties.Property) (invoices.Repository, invoices.Invoice) {
	var invoice invoices.Invoice
	invoice.Status = invoices.Pending
	invoice.Amount = property.Due
	invoice.Property = property.ID
	creation := time.Now()
	var invs = map[string]invoices.Invoice{
		property.ID: {ID: 1, Amount: 1000, CreatedAt: creation, UpdatedAt: creation},
	}
	store := mocks.NewInvoiceRepository(invs)
	return store, invoice
}
