package app

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/config"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/log"
	"github.com/sirupsen/logrus"
)

const namespace = "paypack"

// Bootstrap worker pool
func Bootstrap(conf *config.Config) (*asynq.ServeMux, error) {
	db, err := PostgresConnect(conf.Postgres)
	if err != nil {
		err = fmt.Errorf("error connecting to postgres (%s)", err)
		return nil, err
	}

	logLvl, err := logrus.ParseLevel(conf.LogLevel)
	if err != nil {
		return nil, err
	}
	lggr := log.New(conf.CloudRuntime, logLvl)

	services := ProvideServices(db)

	handlerOpts := ProvideHandlerOptions(services, lggr)

	mux := asynq.NewServeMux()

	Register(mux, handlerOpts)

	return mux, nil
}
