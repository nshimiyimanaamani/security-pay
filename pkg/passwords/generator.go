package passwords

import "context"

// Generator ...
type Generator interface {
	Generate(ctx context.Context) (string, error)
}
