package app

import (
	"github.com/rugwirobaker/paypack-backend/pkg/config"
	"github.com/rugwirobaker/paypack-backend/pkg/tasks/queue"
)

// InitQueue ...
func InitQueue(config *config.RedisConfig) (*queue.Queue, error) {
	opts, err := queue.ParseOptions("redis://localhost:6379")
	if err != nil {
		return nil, err
	}
	return queue.New(opts), nil
}
