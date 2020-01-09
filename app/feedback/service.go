package feedback

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/app/identity"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Service ...
type Service interface {
	Record(ctx context.Context, msg *Message) (*Message, error)
	Retrieve(ctx context.Context, id string) (Message, error)
	List(ctx context.Context, offset, limit uint64) (MessagePage, error)
	Update(ctx context.Context, msg Message) error
	Delete(ctx context.Context, id string) error
}

type service struct {
	idp  identity.Provider
	repo Repository
}

// Options ...
type Options struct {
	Idp  identity.Provider
	Repo Repository
}

// New ...
func New(opts *Options) Service {
	return &service{
		idp:  opts.Idp,
		repo: opts.Repo,
	}
}

func (svc *service) Record(ctx context.Context, msg *Message) (*Message, error) {
	const op errors.Op = "app/feedback/service.Record"

	if err := msg.Validate(); err != nil {
		return nil, errors.E(op, err)
	}
	msg = svc.newMsg(msg)

	return svc.repo.Save(ctx, msg)
}
func (svc *service) Retrieve(ctx context.Context, id string) (Message, error) {
	const op errors.Op = "app/feedback/service.Retrieve"

	msg, err := svc.repo.Retrieve(ctx, id)
	if err != nil {
		return Message{}, errors.E(op, err)
	}

	return msg, nil
}
func (svc *service) Update(ctx context.Context, msg Message) error {
	const op errors.Op = "app/feedback/service.Update"

	if err := msg.Validate(); err != nil {
		return errors.E(op, err)
	}

	if err := svc.repo.Update(ctx, msg); err != nil {
		return errors.E(op, err)
	}
	return svc.repo.Update(ctx, msg)
}
func (svc *service) Delete(ctx context.Context, id string) error {
	const op errors.Op = "app/feedback/service.Delete"

	if err := svc.repo.Delete(ctx, id); err != nil {
		return errors.E(op, err)
	}

	return nil
}

func (svc *service) newMsg(msg *Message) *Message {
	msg.ID = svc.idp.ID()
	return msg
}

func (svc *service) List(ctx context.Context, offset, limit uint64) (MessagePage, error) {
	const op errors.Op = "app/feedback/service.List"

	page, err := svc.List(ctx, offset, limit)
	if err != nil {
		return MessagePage{}, errors.E(op, err)
	}
	return page, nil
}
