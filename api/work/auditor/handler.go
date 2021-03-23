package auditor

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rugwirobaker/paypack-backend/core/auditor"
	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
	"github.com/sirupsen/logrus"
)

// AuditHandler ...
func AuditHandler(lgger log.Entry, svc auditor.Service) asynq.Handler {
	const op errors.Op = "api/work/AuditHandler"

	f := func(ctx context.Context, task *asynq.Task) error {
		var payload = task.Payload

		properties, err := payload.GetInt("properties")
		if err != nil {
			err := errors.E(op, err, errors.KindBadRequest)
			logrus.Error(err)
			return err
		}

		batch, err := payload.GetInt("batch")
		if err != nil {
			err := errors.E(op, err, errors.KindBadRequest)
			logrus.Error(err)
			return err
		}

		if err := svc.Schedule(ctx, properties, batch); err != nil {
			err := errors.E(op, err)
			logrus.Error(err)
			return err
		}
		return nil

	}

	return asynq.HandlerFunc(f)
}
