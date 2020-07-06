package payment

import (
	"bytes"
	"context"
	"fmt"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/core/invoices"
	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/core/owners"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/transactions"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

const header = "Murakoze kwishyura umusanzu wanyu na Paypack\n"

// Service is the api interface to the payment module
type Service interface {
	// Pull initializes payment from an external account
	Pull(ctx context.Context, tx Payment) (Response, error)

	// Push  initializes payment to an external account
	Push(ctx context.Context, tx Payment) (Response, error)

	// ConfirmPull processes debit callback
	ConfirmPull(ctx context.Context, res Callback) error

	// ConfirmPush processes credit callback
	ConfirmPush(ctx context.Context, res Callback) error
}

// Options simplifies New func signature
type Options struct {
	Idp          identity.Provider
	Backend      Client
	Queue        Queue
	SMS          notifs.Service
	Owners       owners.Repository
	Properties   properties.Repository
	Invoices     invoices.Repository
	Transactions transactions.Repository
}
type service struct {
	backend      Client
	idp          identity.Provider
	queue        Queue
	sms          notifs.Service
	owners       owners.Repository
	properties   properties.Repository
	transactions transactions.Repository
	invoices     invoices.Repository
}

// New initializes the payment service
func New(opts *Options) Service {
	return &service{
		queue:        opts.Queue,
		idp:          opts.Idp,
		backend:      opts.Backend,
		owners:       opts.Owners,
		properties:   opts.Properties,
		sms:          opts.SMS,
		invoices:     opts.Invoices,
		transactions: opts.Transactions,
	}
}

func (svc service) Pull(ctx context.Context, py Payment) (Response, error) {
	const op errors.Op = "core/app/payment/service.Push"

	failed := Response{TxState: "failed"}
	if err := py.Validate(); err != nil {
		return failed, errors.E(op, err)
	}

	property, err := svc.properties.RetrieveByID(ctx, py.Code)
	if err != nil {
		return failed, errors.E(op, err)
	}

	invoice, err := svc.invoices.Earliest(ctx, property.ID)
	if err != nil {
		return failed, errors.E(op, err)
	}

	if err := invoice.Verify(py.Amount); err != nil {
		return failed, errors.E(op, err)
	}

	py.Invoice = invoice.ID
	py.ID = svc.idp.ID()

	status, err := svc.backend.Pull(ctx, py)
	if err != nil {
		return failed, errors.E(op, err)
	}
	if err := svc.queue.Set(ctx, py); err != nil {
		return failed, errors.E(op, err)
	}
	return status, nil
}

func (svc *service) Push(ctx context.Context, py Payment) (Response, error) {
	const op errors.Op = "core/app/payment/service.Push"

	failed := Response{TxState: "failed"}

	if err := py.HackyValidation(); err != nil {
		return failed, errors.E(op, err)
	}

	py.ID = svc.idp.ID()

	status, err := svc.backend.Push(ctx, py)
	if err != nil {
		return failed, errors.E(op, err)
	}
	if err := svc.queue.Set(ctx, py); err != nil {
		return failed, errors.E(op, err)
	}
	return status, nil
}

func (svc *service) ConfirmPull(ctx context.Context, cb Callback) error {
	const op errors.Op = "core/app/payment/service. ConfirmPull"

	if err := cb.Validate(); err != nil {
		return errors.E(op, err)
	}

	if cb.Data.State != Successful {
		return errors.E(op, "transaction failed unexpectedly", errors.KindUnexpected)
	}

	payment, err := svc.queue.Get(ctx, cb.Data.TrxRef)
	if err != nil {
		return errors.E(op, err)
	}
	property, err := svc.properties.RetrieveByID(ctx, payment.Code)
	if err != nil {
		return errors.E(op, err)
	}

	transaction := svc.PaymentToTransaction(payment)

	owner, err := svc.owners.Retrieve(ctx, property.Owner.ID)
	if err != nil {
		return errors.E(op, err)
	}
	transaction.OwnerID = owner.ID

	err = svc.Notify(ctx, payment)
	if err != nil {
		return errors.E(op, err)
	}

	if _, err := svc.transactions.Save(ctx, transaction); err != nil {
		return errors.E(op, err)
	}
	//remove tx from the cache
	if err := svc.queue.Remove(ctx, payment.ID); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) ConfirmPush(ctx context.Context, cb Callback) error {
	const op errors.Op = "core/app/payment/service. ConfirmPush"

	if err := cb.Validate(); err != nil {
		return errors.E(op, err)
	}

	if cb.Data.State != Successful {
		return errors.E(op, "transaction failed unexpectedly", errors.KindUnexpected)
	}

	tx, err := svc.queue.Get(ctx, cb.Data.TrxRef)
	if err != nil {
		return errors.E(op, err)
	}
	//remove tx from the cache
	if err := svc.queue.Remove(ctx, tx.ID); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) Notify(ctx context.Context, payment Payment) error {
	const op errors.Op = "core/app/payment/service.Notify"

	property, err := svc.properties.RetrieveByID(ctx, payment.Code)
	if err != nil {
		return errors.E(op, err)
	}

	owner, err := svc.owners.Retrieve(ctx, property.Owner.ID)
	if err != nil {
		return errors.E(op, err)
	}

	message := formatMessage(payment, owner, property)

	notification := notifs.Notification{
		Recipients: []string{owner.Phone, payment.Phone}, //owners
		Sender:     property.Namespace,                   //account
		Message:    message,
	}
	if _, err := svc.sms.Send(ctx, notification); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func formatMessage(payment Payment, own owners.Owner, pr properties.Property) string {
	var buf bytes.Buffer

	buf.WriteString(header)
	buf.WriteString(fmt.Sprintf("nimero yishuriweho': %s\n", payment.Phone))
	buf.WriteString(fmt.Sprintf("TxRef:%s\n", payment.ID))
	buf.WriteString(fmt.Sprintf("Nimero ya fagitire: %d\n", payment.Invoice))
	buf.WriteString(fmt.Sprintf("Umubare w' amafaranga:%f\n", payment.Amount))
	buf.WriteString(fmt.Sprintf("Mwishyuriye %sinzu yishyuriwe:", pr.ID))
	buf.WriteString(fmt.Sprintf("ya %s %s\n:", own.Fname, own.Lname))
	return buf.String()
}

func (svc *service) PaymentToTransaction(tx Payment) transactions.Transaction {
	return transactions.Transaction{
		ID:           tx.ID,
		MadeFor:      tx.Code,
		OwnerID:      tx.ID,
		Amount:       tx.Amount,
		Invoice:      tx.Invoice,
		Method:       tx.Method,
		Namespace:    tx.Namespace,
		DateRecorded: tx.RecordedAt,
	}
}
