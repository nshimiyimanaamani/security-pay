package payment_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rugwirobaker/paypack-backend/core/identity/uuid"
	"github.com/rugwirobaker/paypack-backend/core/invoices"
	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/core/owners"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/core/payment/mocks"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
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
		payment payment.Payment
		state   payment.State
		err     error
	}{
		{
			desc:    "initialize payment with valid data",
			payment: payment.Payment{Code: property.ID, Amount: invoice.Amount, Phone: "0784607135", Method: "mtn-momo-rw"},
			state:   "processing",
			err:     nil,
		},
		{
			desc:    "initialize payment with invalid data",
			payment: payment.Payment{Code: property.ID, Amount: invoice.Amount, Phone: "0784607135"},
			state:   "failed",
			err:     errors.E(op, "payment method must be specified"),
		},
		{
			desc:    "initialize payment with unsaved house property.ID",
			payment: payment.Payment{Code: uuid.New().ID(), Amount: invoice.Amount, Phone: "0784607135", Method: "mtn-momo-rw"},
			state:   "failed",
			err:     errors.E(op, "property not found"),
		},
		{
			desc:    "initialize payment with invalid amount(different from invoice)",
			payment: payment.Payment{Code: property.ID, Amount: 100, Phone: "0784607135", Method: "mtn-momo-rw"},
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

	tx := payment.Payment{
		ID:     uuid.New().ID(),
		Code:   property.ID,
		Amount: 1000, Phone: "0784607135",
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
				Status: "success",
				Data: payment.CallBackData{
					GwRef:  uuid.New().ID(),
					TrxRef: res.TxID,
					State:  payment.Successful,
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
					State:  payment.Successful,
				},
			},
			err: errors.E(op, "status field must not be empty"),
		},
	}

	for _, tc := range cases {
		ctx := context.Background()
		err := svc.ConfirmPull(ctx, tc.callback)
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
		payment payment.Payment
		state   payment.State
		err     error
	}{
		{
			desc:    "initialize payment with valid data",
			payment: payment.Payment{Code: property.ID, Amount: 5000, Phone: "0784607135", Method: "mtn-momo-rw"},
			state:   "processing",
			err:     nil,
		},
		{
			desc:    "initialize payment with invalid method",
			payment: payment.Payment{Code: property.ID, Amount: 5000, Phone: "0784607135"},
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

	tx := payment.Payment{
		ID:     uuid.New().ID(),
		Code:   property.ID,
		Amount: 1000, Phone: "0784607135",
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
				Status: "success",
				Data: payment.CallBackData{
					GwRef:  uuid.New().ID(),
					TrxRef: res.TxID,
					State:  payment.Successful,
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
					State:  payment.Successful,
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

func newService(ws owners.Repository, ps properties.Repository, vc invoices.Repository) payment.Service {
	var opts payment.Options
	opts.Owners = ws
	opts.Properties = ps
	opts.Invoices = vc
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
