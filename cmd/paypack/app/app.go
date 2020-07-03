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

const namespace = "paypack"

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

	//init payment backend
	pb := InitPaymentClient(ctx, conf.Payment)

	//init sms backend
	sms, err := InitSMSBackend(ctx, conf.SMS)
	if err != nil {
		lggr.Errorf("error connecting to sms backend (%s)", err)
		return nil, err
	}

	queue, err := InitQueue(conf.Redis)
	if err != nil {
		lggr.Errorf("error connecting to redis queue (%s)", err)
		return nil, err
	}

	//create protocol
	services := Init(db, rclient, queue, pb, sms, conf.Secret, namespace)

	handlerOpts := NewHandlerOptions(services, lggr)

	r := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(false)

	r.Use(mw.LogEntryMiddleware(lggr))
	r.Use(mw.Recover())

	Register(r, handlerOpts)

	if conf.GoEnv == "development" {
		r.Use(mw.RequestLogger)
	}

	//recover := handlers.RecoveryHandler()(r)

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	return cors(r), nil
}
