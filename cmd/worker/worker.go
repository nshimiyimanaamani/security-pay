package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/rugwirobaker/paypack-backend/cmd/worker/app"
	"github.com/rugwirobaker/paypack-backend/pkg/build"
	"github.com/rugwirobaker/paypack-backend/pkg/config"
	"github.com/rugwirobaker/paypack-backend/pkg/tasks/worker"
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

	mux, err := app.Bootstrap(conf)
	if err != nil {
		log.Fatal(err)
	}

	opts, err := worker.ParseOptions(conf.Redis.URL)
	if err != nil {
		log.Fatal(err)
	}

	worker := worker.New(opts)

	worker.Run(mux)

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

}
