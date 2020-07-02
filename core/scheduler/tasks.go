package scheduler

import (
	"context"
	"fmt"

	"github.com/rugwirobaker/paypack-backend/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	audit    = "audit"
	reminder = "reminder"
)

// Task is schedulable unit of work
type Task func(context.Context, string) (map[string]interface{}, error)

// SelectTask selects task based on user input
func (svc *service) Run(ctx context.Context, name string) error {
	const op errors.Op = "core/scheduler/service.SelectTask"

	switch name {
	case audit:
		return svc.AuditTask(ctx, name)
	case reminder:
		return svc.ReminderTask(ctx, name)
	default:
		return svc.UnknownTask(ctx, name)
	}
}

// AuditTask schedules the audit task
func (svc *service) AuditTask(ctx context.Context, name string) error {
	const op errors.Op = "core/scheduler/service.AuditTask"

	const batch = 20

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
