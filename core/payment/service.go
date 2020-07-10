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

const header = "Murakoze kwishyura umusanzu w' umutekano mu murenge wa "

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

func (svc service) Pull(ctx context.Context, payment Payment) (Response, error) {
	const op errors.Op = "core/payment/service.Pull"

	failed := Response{TxState: "failed"}

	// check the bare minimum
	if err := payment.HasCode(); err != nil {
		return failed, errors.E(op, err)
	}

	property, err := svc.properties.RetrieveByID(ctx, payment.Code)
	if err != nil {
		return failed, errors.E(op, err)
	}

	owner, err := svc.owners.Retrieve(ctx, property.Owner.ID)
	if err != nil {
		return failed, errors.E(op, err)
	}

	invoice, err := svc.invoices.Earliest(ctx, property.ID)
	if err != nil {
		return failed, errors.E(op, err)
	}

	// verify invoice
	if err := invoice.Verify(payment.Amount); err != nil {
		return failed, errors.E(op, err)
	}
	payment.ID = svc.idp.ID()

	// validate payment
	if err := payment.Ready(); err != nil {
		return failed, errors.E(op, err)
	}

	status, err := svc.backend.Pull(ctx, payment)
	if err != nil {
		return failed, errors.E(op, err)
	}

	// validate transaction
	transaction := svc.NewTransaction(payment, property, owner, invoice)
	if err := transaction.Validate(); err != nil {
		return failed, errors.E(op, err)
	}

	if _, err = svc.transactions.Save(ctx, transaction); err != nil {
		return failed, errors.E(op, err)
	}
	return status, nil
}

func (svc *service) Push(ctx context.Context, payment Payment) (Response, error) {
	const op errors.Op = "core/payment/service.Push"

	failed := Response{TxState: "failed"}

	if err := payment.Ready(); err != nil {
		return failed, errors.E(op, err)
	}

	payment.ID = svc.idp.ID()

	status, err := svc.backend.Push(ctx, payment)
	if err != nil {
		return failed, errors.E(op, err)
	}
	if err := svc.queue.Set(ctx, payment); err != nil {
		return failed, errors.E(op, err)
	}
	return status, nil
}

func (svc *service) ConfirmPull(ctx context.Context, cb Callback) error {
	const op errors.Op = "core/payment/service.ConfirmPull"

	if err := cb.Validate(); err != nil {
		return errors.E(op, err)
	}

	if cb.Data.State != Successful {
		return errors.E(op, "transaction failed unexpectedly", errors.KindUnexpected)
	}

	transaction, err := svc.transactions.RetrieveByID(ctx, cb.Data.TrxRef)
	if err != nil {
		return errors.E(op, err)
	}
	transaction.Confirm()

	if err := svc.transactions.Update(ctx, transaction); err != nil {
		return errors.E(op, err)
	}

	err = svc.Notify(ctx, transaction)
	if err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) ConfirmPush(ctx context.Context, cb Callback) error {
	const op errors.Op = "core/payment/service.ConfirmPush"

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

func (svc *service) Notify(ctx context.Context, tx transactions.Transaction) error {
	const op errors.Op = "core/app/payment/service.Notify"

	property, err := svc.properties.RetrieveByID(ctx, tx.MadeFor)
	if err != nil {
		return errors.E(op, err)
	}

	owner, err := svc.owners.Retrieve(ctx, tx.OwnerID)
	if err != nil {
		return errors.E(op, err)
	}

	message := formatMessage(tx, owner, property)

	notification := notifs.Notification{
		Recipients: []string{owner.Phone, tx.MSISDN}, //owners
		Sender:     property.Namespace,               //account
		Message:    message,
	}
	if _, err := svc.sms.Send(ctx, notification); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func formatMessage(tx transactions.Transaction, own owners.Owner, pr properties.Property) string {
	var buf bytes.Buffer

	buf.WriteString(header)
	buf.WriteString(fmt.Sprintf("%s.\n\n", pr.Address.Sector))
	buf.WriteString(fmt.Sprintf("Nimero yishyuriweho: %s\n", tx.MSISDN))
	buf.WriteString(fmt.Sprintf("Nimero ya fagitire: %d\n", tx.Invoice))
	buf.WriteString(fmt.Sprintf("Umubare w' amafaranga: %dRWF\n", int(tx.Amount)))
	buf.WriteString(fmt.Sprintf("Inzu yishyuriwe ni iya %s %s\n", own.Fname, own.Lname))
	buf.WriteString(fmt.Sprintf("Code y' inzu ni: %s\n\n", tx.MadeFor))
	buf.WriteString("Binyuze muri Paypack")
	return buf.String()
}

func (svc *service) NewTransaction(
	payment Payment,
	property properties.Property,
	owner owners.Owner,
	invoice invoices.Invoice,

) transactions.Transaction {
	return transactions.Transaction{
		ID:        payment.ID,
		MadeFor:   property.ID,
		OwnerID:   owner.ID,
		MSISDN:    payment.Phone,
		Amount:    payment.Amount,
		Invoice:   invoice.ID,
		Method:    payment.Method,
		Namespace: property.Namespace,
	}
}
