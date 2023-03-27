package ussd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/core/invoices"
	"github.com/rugwirobaker/paypack-backend/core/owners"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/platypus"
)

// Service is the ussd service interface definition.
type Service interface {
	Process(ctx context.Context, req *Request) (Response, error)
}

// Options aggregates service creation settings
type Options struct {
	Prefix     string
	IDP        identity.Provider
	Properties properties.Repository
	Owners     owners.Repository
	Payment    payment.Service
	Invoices   invoices.Repository
	Agents     users.AgentsRepository
}

type service struct {
	prefix     string
	idp        identity.Provider
	properties properties.Repository
	owners     owners.Repository
	agents     users.AgentsRepository
	invoice    invoices.Repository
	payment    payment.Service
	mux        *platypus.Mux
}

// New initialises the ussd service.
func New(opts *Options) Service {
	svc := new(opts)
	mux := platypus.New(opts.Prefix, platypus.HandlerFunc(svc.Action0))
	mux = register(opts.Prefix, svc, mux)
	svc.mux = mux
	return svc
}

func new(opts *Options) *service {
	return &service{
		idp:        opts.IDP,
		properties: opts.Properties,
		payment:    opts.Payment,
		owners:     opts.Owners,
		agents:     opts.Agents,
		invoice:    opts.Invoices,
	}
}

// register handlers
func register(prefix string, svc *service, mux *platypus.Mux) *platypus.Mux {
	prefix = strings.TrimSuffix(prefix, "#")
	mux.Handle(prefix+"*1", platypus.HandlerFunc(svc.Action1), platypus.TrimTrailHash)
	mux.Handle(prefix+"*1*:id*1*1#", platypus.HandlerFunc(svc.Action1_1_1_1), nil)
	mux.Handle(prefix+"*1*:id*2*1#", platypus.HandlerFunc(svc.Action1_1_1_2), nil)
	mux.Handle(prefix+"*1*:id*1", platypus.HandlerFunc(svc.ActionPreview), platypus.TrimTrailHash)
	mux.Handle(prefix+"*1*:id*2", platypus.HandlerFunc(svc.ActionPreview), platypus.TrimTrailHash)
	mux.Handle(prefix+"*1*:id", platypus.HandlerFunc(svc.Action1_1), platypus.TrimTrailHash)
	mux.Handle(prefix+"*2", platypus.HandlerFunc(svc.action2), platypus.TrimTrailHash)
	mux.Handle(prefix+"*2*:phone#", platypus.HandlerFunc(svc.action2_1), nil)
	mux.Handle(prefix+"*3", platypus.HandlerFunc(svc.action3), platypus.TrimTrailHash)

	//kwishura amezi menshi ukoresheje nimero inzu yanditseho
	mux.Handle(prefix+"*1*:id*1*2*:input*1#", platypus.HandlerFunc(svc.Action1_1_1_1_2_Input_1), nil)
	mux.Handle(prefix+"*1*:id*1*2*:input", platypus.HandlerFunc(svc.Action1_1_1_1_2_Input), platypus.TrimTrailHash)
	mux.Handle(prefix+"*1*:id*1*2", platypus.HandlerFunc(svc.Action1_1_1_1_2), platypus.TrimTrailHash)

	//kwishura amezi menshi ukoresheje nimero uri gukoresha
	mux.Handle(prefix+"*1*:id*2*2", platypus.HandlerFunc(svc.Action1_1_1_1_2), platypus.TrimTrailHash)
	mux.Handle(prefix+"*1*:id*2*2*:input", platypus.HandlerFunc(svc.Action1_1_1_1_2_Input), platypus.TrimTrailHash)
	mux.Handle(prefix+"*1*:id*2*2*:input*1#", platypus.HandlerFunc(svc.Action1_1_1_2_2_Input_1), nil)

	//Kwishyura ikirarane ukoresheje nimero inzu yanditseho
	mux.Handle(prefix+"*1*:id*1*3", platypus.HandlerFunc(svc.Action1_1_1_1_3), platypus.TrimTrailHash)
	mux.Handle(prefix+"*1*:id*1*3*1#", platypus.HandlerFunc(svc.Action1_1_1_1_3_1), nil)

	//Kwishyura ikirarane ukoresheje nimero uri gukoresha
	mux.Handle(prefix+"*1*:id*2*3", platypus.HandlerFunc(svc.Action1_1_1_1_3), platypus.TrimTrailHash)
	mux.Handle(prefix+"*1*:id*2*3*1#", platypus.HandlerFunc(svc.Action1_1_1_2_3_1), nil)
	return mux
}

func (svc *service) Process(ctx context.Context, req *Request) (Response, error) {
	const op errors.Op = "core/ussd/service.Process"

	// if err := req.Validate(); err != nil {
	// 	return Response{}, errors.E(op, err)
	// }
	cmd := platypus.NewCommand(req.Msisdn, req.UserInput)

	result, err := svc.mux.Process(ctx, cmd)
	if err != nil {
		return respond(svc.idp.ID(), result, req), errors.E(op, err)
	}

	return respond(svc.idp.ID(), result, req), nil
}

func (svc *service) Pay(ctx context.Context, p properties.Property, phone string) (string, error) {

	//remove country code prefix so that payment can be able to pass through on the payment gateway(paypack-payments)
	phone = strings.TrimPrefix(phone, "25")
	phone = strings.TrimPrefix(phone, "+25")

	tx := &payment.TxRequest{
		Code:   p.ID,
		MSISDN: phone,
		Amount: p.Due,
		Method: SelectMethod(phone),
	}
	status, err := svc.payment.Pull(ctx, tx)
	if err != nil {
		return status.Message, err
	}
	return status.Message, nil
}

func (svc *service) BulkPay(ctx context.Context, p properties.Property, phone string, month int) (string, error) {

	phone = strings.TrimPrefix(phone, "25")
	phone = strings.TrimPrefix(phone, "+25")

	tx := &payment.TxRequest{
		Code:   p.ID,
		MSISDN: phone,
		Amount: p.Due * float64(month),
		Method: SelectMethod(phone),
	}

	status, err := svc.payment.BulkPull(ctx, tx, month)
	if err != nil {
		return status.Message, err
	}
	return status.Message, nil
}

// credit payment initiation
func (svc *service) CreditPay(ctx context.Context, p properties.Property, phone string, invoices []invoices.Invoice) (string, error) {

	phone = strings.TrimPrefix(phone, "25")
	phone = strings.TrimPrefix(phone, "+25")

	tx := &payment.TxRequest{
		Code:   p.ID,
		MSISDN: phone,
		Amount: p.Due,
		Method: SelectMethod(phone),
	}

	status, err := svc.payment.CreditPull(ctx, tx, invoices)
	if err != nil {
		return status.Message, err
	}
	return status.Message, nil
}

func sequence(res platypus.Result) int {
	if res.Tail() {
		return 0
	}
	return 1
}

func respond(ref string, result platypus.Result, req *Request) Response {
	stamp := time.Now()

	return Response{
		SessionID: req.SessionID,
		GwRef:     req.GwRef,
		AppRef:    ref,
		GwTstamp:  stamp,
		Text:      result.Out,
		End:       sequence(result),
	}
}

// SelectMethod selects payment method based on
func SelectMethod(phone string) payment.Method {
	phone = NormalizePhoneNumber(phone)

	if strings.HasPrefix(phone, "25073") || strings.HasPrefix(phone, "25072") {
		return payment.AIRTEL
	}
	if strings.HasPrefix(phone, "25078") || strings.HasPrefix(phone, "25079") {
		return payment.MTN
	}
	return ""
}

// NormalizePhoneNumber ...
func NormalizePhoneNumber(phone string) string {
	const op errors.Op = "core/ussd/NormalizePhone"
	if len(phone) < 10 || len(phone) > 12 {
		return ""
	}
	if len(phone) == 10 {
		return fmt.Sprintf("%s%s", "25", phone)
	}
	return phone
}
