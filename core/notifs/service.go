package notifs

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/core/identity"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Service provides sms facilities to end user.
type Service interface {
	// Send sms message to a collection of recipients
	Send(ctx context.Context, n Notification) (Notification, error)

	//Retrieve message
	Find(ctx context.Context, id string) (Notification, error)

	//List messages per namespace
	List(ctx context.Context, nspace string, offset, limit uint64) (NoticationPage, error)

	// Count the number of messages sent by an account
	Count(ctx context.Context, nspace string) (uint64, error)
}

// Options to New function
type Options struct {
	Backend Backend
	IDP     identity.Provider
	Store   Repository
}

type service struct {
	backend Backend
	idp     identity.Provider
	store   Repository
}

// New instantiaces the sms service
func New(opts *Options) Service {
	return &service{
		backend: opts.Backend,
		idp:     opts.IDP,
		store:   opts.Store,
	}
}

func (svc *service) Send(ctx context.Context, sms Notification) (Notification, error) {
	const op errors.Op = "core/sms/service.Send"

	var empty Notification

	sms.ID = svc.idp.ID()

	if err := svc.backend.Send(ctx, sms.ID, sms.Message, sms.Recipients); err != nil {
		return empty, errors.E(op, err)
	}

	sms, err := svc.store.Save(ctx, sms)
	if err != nil {
		return empty, errors.E(op, err)
	}
	return sms, nil
}

func (svc *service) Find(ctx context.Context, id string) (Notification, error) {
	const op errors.Op = "core/sms/service.Find"

	var empty Notification

	n, err := svc.store.Find(ctx, id)
	if err != nil {
		return empty, errors.E(op, err)
	}
	return n, nil
}

func (svc *service) List(ctx context.Context, nspace string, offset, limit uint64) (NoticationPage, error) {
	const op errors.Op = "core/sms/service.List"

	page, err := svc.store.List(ctx, nspace, offset, limit)
	if err != nil {
		return page, errors.E(op, err)
	}

	return page, nil
}

func (svc *service) Count(ctx context.Context, nspace string) (uint64, error) {
	const op errors.Op = "core/sms/service.Count"

	count, err := svc.store.Count(ctx, nspace)
	if err != nil {
		return 0, errors.E(op, err)
	}

	return count, nil

}
