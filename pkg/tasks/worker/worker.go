package worker

import "github.com/hibiken/asynq"

// number of concurent workers
const concurrency = 10

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

// New ...
func New(opts *Options) *asynq.Background {
	bg := asynq.NewBackground(opts.RedisOpts, &asynq.Config{
		Concurrency: concurrency,
	})
	return bg
}
