package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/rugwirobaker/paypack-backend/cmd/paypack/app"
	"github.com/rugwirobaker/paypack-backend/pkg/build"
	"github.com/rugwirobaker/paypack-backend/pkg/config"
)

var vers = flag.Bool("version", false, "Print version information and exit")

var prefix = "paypack"

func main() {
	flag.Parse()
	if *vers {
		fmt.Println(build.String())
		os.Exit(0)
	}

	//get configuration
	conf, err := config.Load(prefix)
	if err != nil {
		log.Fatalf("could not load configuration: %v", err)
	}

	if err := config.Validate(conf); err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	handler, err := app.Bootstrap(conf)
	if err != nil {
		log.Fatal(err)
	}

	// start the server
	srv := &http.Server{Addr: ":" + conf.Port, Handler: handler}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
		close(idleConnsClosed)
	}()

	log.Printf("Starting application at port %v", conf.Port)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-idleConnsClosed

}
