package app

import (
	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"

	"github.com/hibiken/asynq"
	"github.com/nshimiyimanaamani/paypack-backend/api/work/archiver"
	"github.com/nshimiyimanaamani/paypack-backend/api/work/auditor"
)

// HandlerOptions ...
type HandlerOptions struct {
	ArchiveOptions *archiver.HandlerOpts
	AuditOptions   *auditor.HandlerOpts
}

// ProvideHandlerOptions ...
func ProvideHandlerOptions(services *Services, lggr *log.Logger) *HandlerOptions {
	audit := &auditor.HandlerOpts{
		Logger:  lggr,
		Service: services.Auditor,
	}
	archive := &archiver.HandlerOpts{
		Logger:  lggr,
		Service: services.Archiver,
	}

	return &HandlerOptions{
		ArchiveOptions: archive,
		AuditOptions:   audit,
	}
}

// Register registers all handlers
func Register(mux *asynq.ServeMux, opts *HandlerOptions) {
	if opts.AuditOptions == nil || opts.ArchiveOptions == nil {
		panic("absolutely unacceptable start server opts")
	}

	archiver.RegisterHandlers(mux, opts.ArchiveOptions)
	auditor.RegisterHandlers(mux, opts.AuditOptions)
}
