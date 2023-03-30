package payment

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/core/invoices"
	"github.com/rugwirobaker/paypack-backend/core/notifs"
	"github.com/rugwirobaker/paypack-backend/core/owners"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/transactions"
	"github.com/rugwirobaker/paypack-backend/pkg/clock"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Service is the api interface to the payment module
type Service interface {
	// Pull initializes payment from an external account
	Pull(ctx context.Context, tx *TxRequest) (*TxResponse, error)

	// Push  initializes payment to an external account
	Push(ctx context.Context, tx *TxRequest) (*TxResponse, error)

	// ProcessHook processes debit callback
	ProcessHook(ctx context.Context, res Callback) error

	// ConfirmPush processes credit callback
	ConfirmPush(ctx context.Context, res Callback) error

	// PaymentRequest generates all payments
	PaymentReports(ctx context.Context, status, sector, cell, village string, limit, offset uint64) (PaymentResponse, error)
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
	Repository   Repository
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
	repository   Repository
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
		repository:   opts.Repository,
	}
}

// Pull is for initiating payment for last invoice
func (svc service) Pull(ctx context.Context, payment *TxRequest) (*TxResponse, error) {
	const op errors.Op = "core/payment/service.Pull"

	failed := &TxResponse{TxState: "failed"}

	// check the bare minimum
	if err := payment.HasCode(); err != nil {
		failed.Message = err.Error()
		return failed, errors.E(op, err)
	}

	property, err := svc.properties.RetrieveByID(ctx, payment.Code)
	if err != nil {
		failed.Message = err.Error()
		return failed, errors.E(op, err)
	}

	invoice, err := svc.invoices.Earliest(ctx, property.ID)
	if err != nil {
		failed.Message = err.Error()
		return failed, errors.E(op, err)
	}

	// verify invoice
	if err := invoice.Verify(payment.Amount); err != nil {
		failed.Message = err.Error()
		return failed, errors.E(op, err)
	}
	payment.ID = svc.idp.ID()

	payment.Invoice = invoice.ID
	if err := payment.HasInvoice(); err != nil {
		failed.Message = err.Error()
		return failed, errors.E(op, err)
	}

	// validate payment
	if err := payment.Ready(); err != nil {
		failed.Message = err.Error()
		return failed, errors.E(op, err)
	}

	res, err := svc.backend.Pull(ctx, payment)
	if err != nil {
		failed.Message = err.Error()
		return failed, errors.E(op, err)
	}

	payment.ID = uuid.NewV4().String()
	payment.Ref = res.TxID
	payment.Confirmed = false

	if err := svc.repository.Save(ctx, payment); err != nil {
		failed.Message = err.Error()
		return failed, errors.E(op, err)
	}
	res.TxState = "success"
	return res, nil
}

// BulkPull initiate payment for many invoices
func (svc service) BulkPull(ctx context.Context, payment *TxRequest, month int) (*TxResponse, error) {
	const op errors.Op = "core/payment/service.BulkPull"

	failed := &TxResponse{TxState: "failed"}

	// check the bare minimum
	if err := payment.HasCode(); err != nil {
		failed.Message = err.Error()
		return failed, errors.E(op, err)
	}

	// validate payment
	if err := payment.Ready(); err != nil {
		failed.Message = err.Error()
		return failed, errors.E(op, err)
	}

	payments := make([]*TxRequest, 0)

	invoices, err := svc.invoices.Generate(ctx, payment.Code, uint(payment.Amount), uint(month))
	if err != nil {
		failed.Message = "Mwihangane habaye ikibazo muri sisiteme mwongere mukanya"
		return failed, errors.E(op, err)
	}

	if len(invoices) == 0 {
		failed.Message = "Mwihangane habaye ikibazo muri sisiteme mwongere mukanya"
		return failed, errors.E(op, failed.Message, errors.KindUnexpected)
	}

	res, err := svc.backend.Pull(ctx, payment)
	if err != nil {
		failed.Message = err.Error()
		return failed, errors.E(op, err)
	}

	for _, invoice := range invoices {
		payment := &TxRequest{
			ID:        svc.idp.ID(),
			Amount:    invoice.Amount,
			Invoice:   invoice.ID,
			Method:    payment.Method,
			MSISDN:    payment.MSISDN,
			Code:      payment.Code,
			Ref:       res.TxID,
			Confirmed: false,
		}
		payments = append(payments, payment)
	}

	if err := svc.repository.BulkSave(ctx, payments); err != nil {
		failed.Message = "Mwihangane habaye ikibazo muri sisiteme mwongere mukanya"
		return failed, errors.E(op, err)
	}

	return res, nil
}

// CreditPull initiate payment for credited invoices
func (svc service) CreditPull(ctx context.Context, payment *TxRequest, invoices []invoices.Invoice) (*TxResponse, error) {
	const op errors.Op = "core/payment/service.CreditPull"

	failed := &TxResponse{TxState: "failed"}

	// check the bare minimum
	if err := payment.HasCode(); err != nil {
		return failed, errors.E(op, err)
	}

	// validate payment
	if err := payment.Ready(); err != nil {
		return failed, errors.E(op, err)
	}

	var amount float64
	for _, invoice := range invoices {
		amount += invoice.Amount
	}
	payment.Amount = amount

	res, err := svc.backend.Pull(ctx, payment)
	if err != nil {
		return failed, errors.E(op, err)
	}

	payments := make([]*TxRequest, 0)

	for _, invoice := range invoices {
		payment := &TxRequest{
			ID:        svc.idp.ID(),
			Amount:    invoice.Amount,
			Invoice:   invoice.ID,
			Method:    payment.Method,
			MSISDN:    payment.MSISDN,
			Code:      payment.Code,
			Ref:       res.TxID,
			Confirmed: false,
		}
		payments = append(payments, payment)
	}

	if err := svc.repository.BulkSave(ctx, payments); err != nil {
		return failed, errors.E(op, err)
	}

	return res, nil
}

func (svc *service) Push(ctx context.Context, payment *TxRequest) (*TxResponse, error) {
	const op errors.Op = "core/payment/service.Push"

	failed := &TxResponse{TxState: "failed"}

	if err := payment.Ready(); err != nil {
		return failed, errors.E(op, err)
	}

	payment.ID = svc.idp.ID()

	res, err := svc.backend.Push(ctx, payment)
	if err != nil {
		return failed, errors.E(op, err)
	}

	payment.ID = res.TxID
	//save instead to payments
	if err := svc.queue.Set(ctx, payment); err != nil {
		return failed, errors.E(op, err)
	}
	return res, nil
}

func (svc *service) ProcessHook(ctx context.Context, cb Callback) error {
	const op errors.Op = "core/payment/service.ProcessHook"

	if err := cb.Validate(); err != nil {
		return errors.E(op, err)
	}

	payments, err := svc.repository.Find(ctx, cb.Data.Ref)
	if err != nil {
		return errors.E(op, err)
	}

	if len(payments) == 0 {
		return errors.E(op, fmt.Sprintf("no payments found for this ref %s", cb.Data.Ref), errors.KindUnexpected)
	}

	if payments[0].Status != "pending" {
		return errors.E(op, fmt.Sprintf("payment with ref %s is already processed", cb.Data.Ref), errors.KindUnexpected)
	}

	if err := svc.repository.Update(ctx, cb.Data.Status, payments); err != nil {
		return errors.E(op, err, errors.Kind(err))
	}

	return nil
}

func (svc *service) ConfirmPush(ctx context.Context, cb Callback) error {
	const op errors.Op = "core/payment/service.ConfirmPush"

	if err := cb.Validate(); err != nil {
		return errors.E(op, err)
	}

	if State(cb.Data.Status) != Successful {
		return errors.E(op, cb.Data.Status, errors.KindUnexpected)
	}

	//retrieve from payments instead
	tx, err := svc.queue.Get(ctx, cb.Data.Ref)
	if err != nil {
		return errors.E(op, err)
	}
	//remove tx from the cache
	if err := svc.queue.Remove(ctx, tx.ID); err != nil {
		return errors.E(op, err)
	}
	return nil
}

func (svc *service) PaymentReports(ctx context.Context, status, sector, cell, village string, limit, offset uint64) (PaymentResponse, error) {

	const op errors.Op = "core/payment/service.PaymentReports"
	tx, err := svc.repository.List(ctx, status, sector, cell, village, limit, offset)
	if err != nil {
		return PaymentResponse{}, errors.E(op, err)
	}

	return tx, nil
}

func (svc *service) Notify(ctx context.Context, py TxRequest, tx transactions.Transaction) error {
	const op errors.Op = "core/app/payment/service.Notify"

	property, err := svc.properties.RetrieveByID(ctx, tx.MadeFor)
	if err != nil {
		return errors.E(op, err)
	}

	owner, err := svc.owners.Retrieve(ctx, tx.OwnerID)
	if err != nil {
		return errors.E(op, err)
	}

	invoice, err := svc.invoices.Find(ctx, py.Invoice)
	if err != nil {
		return errors.E(op, err)
	}

	message := FormatMessage(tx, invoice, py, owner, property, timestamp())

	notification := notifs.Notification{
		Recipients: []string{owner.Phone, py.MSISDN}, //owners
		Sender:     property.Namespace,               //account
		Message:    message,
	}

	if _, err := svc.sms.Send(ctx, notification); err != nil {
		return errors.E(op, err)
	}

	return nil
}

// FormatMessage creates sms message
func FormatMessage(
	tx transactions.Transaction,
	inv invoices.Invoice,
	py TxRequest,
	own owners.Owner,
	pr properties.Property,
	timestamp string,
) string {
	const header = "Murakoze kwishyura umusanzu w' isuku"

	var buf bytes.Buffer

	buf.WriteString(header)
	// buf.WriteString(selectActivity(pr.Address.Sector))
	buf.WriteString(" mu murenge wa ")
	buf.WriteString(fmt.Sprintf("%s.\n\n", pr.Address.Sector))
	buf.WriteString(fmt.Sprintf("Nimero yishyuriweho: %s\n", py.MSISDN))
	buf.WriteString(fmt.Sprintf("Itariki: %s\n", timestamp))
	buf.WriteString(fmt.Sprintf("Wishyuriye Ukwezi kwa: %d\n", inv.CreatedAt.Month()))
	buf.WriteString(fmt.Sprintf("Nimero ya fagitire: %d\n", tx.Invoice))
	buf.WriteString(fmt.Sprintf("Umubare w' amafaranga: %dRWF\n", int(tx.Amount)))
	buf.WriteString(fmt.Sprintf("Inzu yishyuriwe ni iya %s %s\n", own.Fname, own.Lname))
	buf.WriteString(fmt.Sprintf("Code y' inzu ni: %s", tx.MadeFor))
	return buf.String()
}

func timestamp() string {
	at := clock.TimeIn(time.Now(), clock.EAST)
	return clock.Format(at, clock.LayoutCustom)
}

// temp hack
func selectActivity(sector string) string {
	switch sector {
	case "Remera":
		return "umutekano"
	default:
		return "isuku"
	}
}

func (svc *service) NewTransaction(
	payment TxRequest,
	property properties.Property,
	owner owners.Owner,
) transactions.Transaction {
	return transactions.Transaction{
		ID:        payment.ID,
		MadeFor:   property.ID,
		OwnerID:   owner.ID,
		Amount:    payment.Amount,
		Invoice:   payment.Invoice,
		Method:    string(payment.Method),
		Namespace: property.Namespace,
	}
}
