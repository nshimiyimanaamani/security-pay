package archiver

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/nshimiyimanaamani/paypack-backend/core/archiver"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"
)

// ArchiveHandler ...
func ArchiveHandler(lgger log.Entry, svc archiver.Service) asynq.Handler {
	const op errors.Op = "api/work/ArchiveHandler"

	f := func(ctx context.Context, task *asynq.Task) error {
		var payload = task.Payload

		invoices, err := payload.GetInt("invoices")
		if err != nil {
			err := errors.E(op, err, errors.KindBadRequest)
			lgger.SystemErr(err)
			return err
		}

		batch, err := payload.GetInt("batch")
		if err != nil {
			err := errors.E(op, err, errors.KindBadRequest)
			lgger.SystemErr(err)
			return err
		}

		if err := svc.Schedule(ctx, invoices, batch); err != nil {
			err := errors.E(op, err)
			lgger.SystemErr(err)
			return err
		}
		return nil

	}

	return asynq.HandlerFunc(f)
}
