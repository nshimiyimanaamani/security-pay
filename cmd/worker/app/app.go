package app

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rugwirobaker/paypack-backend/api/work/auditor"
	"github.com/rugwirobaker/paypack-backend/pkg/config"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
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

	opts := &auditor.HandlerOpts{
		Logger:  lggr,
		Service: bootAuditor(db),
	}

	mux := asynq.NewServeMux()

	auditor.RegisterHandlers(mux, opts)

	return mux, nil
}
