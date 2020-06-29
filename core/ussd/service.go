package ussd

import (
	"context"
	"strings"
	"time"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/core/owners"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/platypus"
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
}

type service struct {
	prefix     string
	idp        identity.Provider
	properties properties.Repository
	owners     owners.Repository
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
	}
}

// register handlers
func register(prefix string, svc *service, mux *platypus.Mux) *platypus.Mux {
	prefix = strings.TrimSuffix(prefix, "#")
	mux.Handle(prefix+"*1", platypus.HandlerFunc(svc.Action1), platypus.TrimTrailHash)
	mux.Handle(prefix+"*1*:id*1#", platypus.HandlerFunc(svc.Action1_1_1), nil)
	mux.Handle(prefix+"*1*:id", platypus.HandlerFunc(svc.Action1_1), platypus.TrimTrailHash)
	mux.Handle(prefix+"*2", platypus.HandlerFunc(svc.action2), platypus.TrimTrailHash)
	mux.Handle(prefix+"*2*:phone#", platypus.HandlerFunc(svc.action2_1), nil)
	return mux
}

func (svc *service) Process(ctx context.Context, req *Request) (Response, error) {
	const op errors.Op = "core/ussd/service.Process"

	if err := req.Validate(); err != nil {
		return Response{}, errors.E(op, err)
	}

	cmd := &platypus.Command{Pattern: req.UserInput}

	result, err := svc.mux.Process(ctx, cmd)

	if err != nil {
		return Response{}, errors.E(op, err, errors.KindUnexpected)
	}
	return respond(svc.idp.ID(), result, req), nil
}

func (svc *service) MakePayment(ctx context.Context, property properties.Property) error {
	tx := payment.Transaction{
		Code:   property.ID,
		Method: payment.MTN,
	}
	_, err := svc.payment.Initilize(ctx, tx)
	if err != nil {
		return err
	}
	return nil
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
