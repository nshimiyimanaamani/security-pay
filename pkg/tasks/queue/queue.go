package queue

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/nshimiyimanaamani/paypack-backend/core/scheduler"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
)

var _ (scheduler.Queue) = (*Queue)(nil)

// Queue ...
type Queue struct {
	cli *asynq.Client
}

// Options ...
type Options struct {
	RedisOpts *asynq.RedisClientOpt
}

// ParseOptions from uri
func ParseOptions(uri string) (*Options, error) {
	if uri == "" {
		return defaultOpts(), nil
	}
	return ParseURL(uri)
}

func defaultOpts() *Options {
	opts := &asynq.RedisClientOpt{
		Addr: "localhost:6379",
		DB:   0,
	}
	return &Options{opts}
}

// New Queue
func New(opts *Options) *Queue {
	cli := asynq.NewClient(opts.RedisOpts)
	return &Queue{cli}
}

// Enqueue new task
func (queue *Queue) Enqueue(ctx context.Context, name string, args map[string]interface{}) error {
	const op errors.Op = "pkg/tasks/queue/Queue.Enqueue"

	task := asynq.NewTask(name, args)

	if _, err := queue.cli.Enqueue(task, asynq.MaxRetry(-1), asynq.Timeout(40*time.Minute)); err != nil {
		return errors.E(op, err)
	}
	return nil
}
