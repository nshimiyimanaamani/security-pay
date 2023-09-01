package scheduler

import (
	"context"
	"fmt"

	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	audit    = "audit"
	reminder = "reminder"
	archive  = "archive"
)

// Task is schedulable unit of work
type Task func(context.Context, string) (map[string]interface{}, error)

// SelectTask selects task based on user input
func (svc *service) Run(ctx context.Context, name string) error {
	const op errors.Op = "core/scheduler/service.Run"

	switch name {
	case audit:
		return svc.AuditTask(ctx, name)
	case reminder:
		return svc.ReminderTask(ctx, name)
	case archive:
		return svc.ArchiveTask(ctx, name)
	default:
		return svc.UnknownTask(ctx, name)
	}
}

// AuditTask schedules the audit task
func (svc *service) AuditTask(ctx context.Context, name string) error {
	const op errors.Op = "core/scheduler/service.AuditTask"

	const batch = 50

	var args = make(map[string]interface{})

	properties, err := svc.auditble.Count(ctx)
	if err != nil {
		return errors.E(op, err)
	}

	logrus.Info("testing")

	args["properties"] = properties
	args["batch"] = batch

	err = svc.queue.Enqueue(ctx, name, args)
	if err != nil {
		return err
	}
	return nil
}

// UnknownTask is returned for undefined tasks
func (svc *service) UnknownTask(ctx context.Context, name string) error {
	const op errors.Op = "core/scheduler/service.UnknownTask"
	return errors.E(op, fmt.Errorf("task(%s) not found", name), errors.KindNotFound)
}

// ReminderTask schedules the reminder task
func (svc *service) ReminderTask(ctx context.Context, name string) error {
	const op errors.Op = "core/scheduler/service.ReminderTask"

	return errors.E(op, errors.KindNotImplemented)
}

func (svc *service) ArchiveTask(ctx context.Context, name string) error {
	const op errors.Op = "core/scheduler/service.ArchiveTask"

	const batch = 50

	var args = make(map[string]interface{})

	invoices, err := svc.invoices.Archivable(ctx)
	if err != nil {
		return errors.E(op, err)
	}

	args["invoices"] = invoices.Total
	args["batch"] = batch

	err = svc.queue.Enqueue(ctx, name, args)
	if err != nil {
		return errors.E(op, err)
	}

	return nil
}
