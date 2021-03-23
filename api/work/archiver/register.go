package archiver

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/rugwirobaker/paypack-backend/core/archiver"
	"github.com/rugwirobaker/paypack-backend/core/auditor"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
)

// LogEntryHandler pulls a log entry from the request context. Thanks to the
// LogEntryMiddleware, we should have a log entry stored in the context for each
// request with request-specific fields. This will grab the entry and pass it to
// the protocol handlers
func LogEntryHandler(ph ProtocolHandler, opts *HandlerOpts) asynq.Handler {
	f := func(ctx context.Context, task *asynq.Task) error {
		ent := log.EntryFromContext(ctx)
		handler := ph(ent, opts.Service)
		return handler.ProcessTask(ctx, task)
	}
	return asynq.HandlerFunc(f)
}

// ProtocolHandler adapts the feedback service into an  asynq..handler
type ProtocolHandler func(lgger log.Entry, svc archiver.Service) asynq.Handler

// HandlerOpts are the generic options
// for a ProtocolHandler
type HandlerOpts struct {
	Logger  *log.Logger
	Service auditor.Service
}

// RegisterHandlers ...
func RegisterHandlers(r *asynq.ServeMux, opts *HandlerOpts) {
	// If true, this would only panic at boot time, static nil checks anyone?
	if opts == nil || opts.Service == nil || opts.Logger == nil {
		panic("absolutely unacceptable handler opts")
	}
	r.Handle("archive", LogEntryHandler(ArchiveHandler, opts))
}
