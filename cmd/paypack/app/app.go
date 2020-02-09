package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	mw "github.com/rugwirobaker/paypack-backend/api/http/middleware"
	"github.com/rugwirobaker/paypack-backend/pkg/config"
	"github.com/rugwirobaker/paypack-backend/pkg/log"
	"github.com/sirupsen/logrus"
)

// Bootstrap is where all routes and middleware for the server
func Bootstrap(conf *config.Config) (http.Handler, error) {

	ctx := context.Background()

	db, err := PostgresConnect(conf.Postgres)
	if err != nil {
		err = fmt.Errorf("error connecting to postgres (%s)", err)
		return nil, err
	}

	rclient, err := RedisConnect(conf.Redis)
	if err != nil {
		err = fmt.Errorf("error connecting to redis (%s)", err)
		return nil, err
	}

	logLvl, err := logrus.ParseLevel(conf.LogLevel)
	if err != nil {
		return nil, err
	}
	lggr := log.New(conf.CloudRuntime, logLvl)

	pb, err := InitPBackend(ctx, conf.Payment)
	if err != nil {
		lggr.Errorf("error connecting to payment backend (%s)", err)
	}

	//create protocol

	services := Init(db, rclient, pb, conf.Secret)

	handlerOpts := NewHandlerOptions(services, lggr)

	r := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(false)

	r.Use(mw.LogEntryMiddleware(lggr))

	Register(r, handlerOpts)

	if conf.GoEnv == "development" {
		r.Use(mw.RequestLogger)
	}

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	return cors(r), nil
}
