package ussd

import (
	"context"
	"time"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/core/ussd/actions"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/ussd"
	"github.com/rugwirobaker/paypack-backend/pkg/ussd/command"
	"github.com/rugwirobaker/paypack-backend/pkg/ussd/executor"
)

// Service is the ussd service interface definition.
type Service interface {
	Process(ctx context.Context, r Request) (Response, error)
}

// Options aggregates service creation settings
type Options struct {
	Idp identity.Provider
}

type service struct {
	executor executor.Executor
	idp      identity.Provider
}

// New initialises the ussd service.
func New(opts *Options) Service {
	base := actions.Action0()
	return &service{
		idp:      opts.Idp,
		executor: executor.NewSimpleExecutor(base),
	}
}

func (svc *service) Process(ctx context.Context, req Request) (Response, error) {
	const op errors.Op = "core/ussd/service.Processs"

	if err := req.Validate(); err != nil {
		return Response{}, errors.E(op, err)
	}

	cmd, err := command.Parse(req.UserInput)
	if err != nil {
		return Response{}, errors.E(op, err)
	}

	res, err := svc.executor.Execute(cmd)
	if err != nil {
		return Response{}, errors.E(op, err)
	}

	return svc.Respond(res, req), nil
}

func (svc *service) Respond(res ussd.Result, req Request) Response {
	return Response{
		SessionID:  req.SessionID,
		GatewayRef: req.GwRef,
		AppRef:     svc.idp.ID(),
		GwTstamp:   time.Now(),
		Text:       res.String(),
		End:        res.End(),
	}
}
