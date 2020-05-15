package ussd

import (
	"context"
	"time"

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

// Settings aggregates service creation settings
type Settings struct {
	AppID string
}

type service struct {
	AppID    string
	executor executor.Executor
}

// New initialises the ussd service.
func New(sets *Settings) Service {
	base := actions.Action0()
	return &service{
		AppID:    sets.AppID,
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
		Session:    req.Session,
		GatewayRef: req.GatewayRef,
		AppRef:     svc.AppID,
		Timestamp:  time.Now(),
		Text:       res.String(),
		End:        res.End(),
	}
}
