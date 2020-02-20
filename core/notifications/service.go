package notifications

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Service provides sms facilities to end user.
type Service interface {
	Send(ctx context.Context, message Payload) error
}

// Options to New function
type Options struct {
	Backend Backend
	IDP     identity.Provider
}

type service struct {
	backend Backend
	idp     identity.Provider
}

// New instantiaces the sms service
func New(opts *Options) Service {
	return &service{
		backend: opts.Backend,
		idp:     opts.IDP,
	}
}

func (svc *service) Send(ctx context.Context, payload Payload) error {
	const op errors.Op = "core/sms/service.Send"

	payload.ID = svc.idp.ID()

	if err := svc.backend.Send(ctx, payload.ID, payload.Message, payload.Recipients); err != nil {
		return errors.E(op, err)
	}
	return nil
}
