package scheduler

import (
	"context"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
)

// Service exposes the scheduler use casea
type Service interface {
	Schedule(ctx context.Context, task string) error
}

// Queue abstracts the task queue
type Queue interface {
	Enqueue(ctx context.Context, name string, args map[string]interface{}) error
}

// Auditable  a given collection's items in a datasource
type Auditable interface {
	Count(context.Context) (int, error)
}

// Options to configure a new scheduler service
type Options struct {
	Queue
	Auditable
}

type service struct {
	queue    Queue
	auditble Auditable
	tasks    map[string]Task
}

// New ...
func New(opts *Options) Service {
	svc := &service{
		queue:    opts.Queue,
		auditble: opts.Auditable,
	}
	return svc
}

func (svc *service) Schedule(ctx context.Context, task string) error {
	const op errors.Op = "core/scheduler/service.Schedule"

	err := svc.Run(ctx, task)
	if err != nil {
		return errors.E(op, err)
	}
	return nil
}
