package archiver

import "context"

// Executor ...
type Executor interface {
	ArchiveFunc(ctx context.Context, offset, limit int) (int, error)
}
