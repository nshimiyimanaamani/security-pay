package auditor

import "context"

// Executor generates reports for a batch of properties
type Executor interface {
	AuditFunc(ctx context.Context, offset, limit int) (int, error)
}
