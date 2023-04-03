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
	"github.com/rugwirobaker/paypack-backend/web"
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

	logLvl, err := logrus.ParseLevel("info")
	if err != nil {
		return nil, err
	}
	lggr := log.New(conf.CloudRuntime, logLvl)

	//init payment backend
	pb, err := InitPaymentClient(ctx, conf)
	if err != nil {
		lggr.Errorf("error connecting to payment backend (%s)", err)
		return nil, err
	}

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
	services := Init(db, rclient, queue, pb, sms, conf.Secret, namespace, conf.USSD.Prefix)

	handlerOpts := NewHandlerOptions(services, lggr)

	router := mux.NewRouter()

	router.Use(mw.LogEntryMiddleware(lggr))

	if conf.GoEnv == "development" {
		router.Use(mw.RequestLogger)
	}

	api := router.PathPrefix("/api").Subrouter().StrictSlash(false)

	Register(api, handlerOpts)

	router.PathPrefix("/").Handler(web.Handler("/", "dist"))

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	return cors(router), nil
}
