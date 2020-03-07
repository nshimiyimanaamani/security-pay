package app

import (
	"github.com/rugwirobaker/paypack-backend/pkg/config"
	"github.com/rugwirobaker/paypack-backend/pkg/tasks/queue"
)

// InitQueue ...
func InitQueue(config *config.RedisConfig) (*queue.Queue, error) {
	opts, err := queue.ParseOptions(config.URL)
	if err != nil {
		return nil, err
	}
	return queue.New(opts), nil
}
