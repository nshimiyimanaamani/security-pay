package auditor

import (
	"context"

	batch "github.com/nshimiyimanaamani/paypack-backend/pkg/batcher"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

// Service ...
type Service interface {
	Schedule(ctx context.Context, count, bsize int) error
}

// Options ...
type Options struct {
	Executor Executor
}

type service struct {
	exec Executor
}

// New ...
func New(opts *Options) Service {
	return &service{exec: opts.Executor}
}

func (svc *service) Schedule(ctx context.Context, count, bsize int) error {
	const op errors.Op = "core/auditor/service.Schedule"

	return batch.All(count, bsize, func(start, end int) error {
		if err := ctx.Err(); err != nil {
			return errors.E(op, err)
		}
		if _, err := svc.exec.AuditFunc(ctx, start, bsize); err != nil {
			return errors.E(op, err)
		}
		return nil
	})
}
