package archiver

import (
	"context"

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

// New ...
func New(opts *Options) Service {
	return &service{exec: opts.Executor}
}

type service struct {
	exec Executor
}

func (svc *service) Schedule(ctx context.Context, count, bsize int) error {
	const op errors.Op = "core/archiver/service.Schedule"

	var offset = 0

	for {
		affected, err := svc.exec.ArchiveFunc(ctx, offset, bsize)
		if err != nil {
			return errors.E(op, err)
		}

		if affected == 0 {
			break
		}
	}
	return nil
}
