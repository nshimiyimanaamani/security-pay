package app

import (
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

	db, err := ConnectToDB(conf.DB)
	if err != nil {
		err = fmt.Errorf("error connecting to database (%s)", err)
		return nil, err
	}

	logLvl, err := logrus.ParseLevel(conf.LogLevel)
	if err != nil {
		return nil, err
	}
	lggr := log.New(conf.CloudRuntime, logLvl)

	//create protocol

	services := Init(db, conf.Secret)

	handlerOpts := NewHandlerOptions(services, lggr)

	r := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(false)

	r.Use(mw.LogEntryMiddleware(lggr))

	Register(r, handlerOpts)

	if conf.GoEnv == "development" {
		r.Use(mw.RequestLogger)
	}

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	return cors(r), nil
}
