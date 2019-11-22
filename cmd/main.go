package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rugwirobaker/paypack-backend"
	feedbackEndpoints "github.com/rugwirobaker/paypack-backend/api/http/feedback"
	"github.com/rugwirobaker/paypack-backend/api/http/health"
	ownersEndpoints "github.com/rugwirobaker/paypack-backend/api/http/owners"
	paymentEndpoints "github.com/rugwirobaker/paypack-backend/api/http/payment"
	prtEndpoints "github.com/rugwirobaker/paypack-backend/api/http/properties"
	trxEndpoints "github.com/rugwirobaker/paypack-backend/api/http/transactions"
	usersEndpoints "github.com/rugwirobaker/paypack-backend/api/http/users"
	"github.com/rugwirobaker/paypack-backend/api/http/version"
	"github.com/rugwirobaker/paypack-backend/app/feedback"
	"github.com/rugwirobaker/paypack-backend/app/nanoid"
	"github.com/rugwirobaker/paypack-backend/app/owners"
	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/app/users/bcrypt"
	"github.com/rugwirobaker/paypack-backend/app/users/jwt"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/build"
	"github.com/rugwirobaker/paypack-backend/logger"
	"github.com/rugwirobaker/paypack-backend/nova"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
)

const (
	defLogLevel        = "error"
	defDBHost          = "localhost"
	defDBPort          = "5432"
	defDBUser          = "paypack"
	defDBPass          = "paypack"
	defDBName          = "users"
	defDBSSLMode       = "disable"
	defDBSSLCert       = ""
	defDBSSLKey        = ""
	defDBSSLRootCert   = ""
	defHTTPPort        = "8080"
	defSecret          = "users"
	defServerCert      = ""
	defServerKey       = ""
	defPaymentEndpoint = ""
	defPaymentToken    = ""
	envLogLevel        = "PAYPACK_LOG_LEVEL"
	envDBHost          = "PAYPACK_DB_HOST"
	envDBPort          = "PAYPACK_DB_PORT"
	envDBUser          = "PAYPACK_DB_USER"
	envDBPass          = "PAYPACK_DB_PASS"
	envDBName          = "PAYPACK_DB"
	envDBSSLMode       = "PAYPACK_DB_SSL_MODE"
	envDBSSLCert       = "PAYPACK_DB_SSL_CERT"
	envDBSSLKey        = "PAYPACK_DB_SSL_KEY"
	envDBSSLRootCert   = "PAYPACK_DB_SSL_ROOT_CERT"
	envHTTPPort        = "PORT"
	envSecret          = "PAYPACK_SECRET"
	envServerCert      = "PAYPACK_SERVER_CERT"
	envServerKey       = "PAYPACK_SERVER_KEY"
	envPaymentEndpoint = "PAYPACK_PAYMENT_ENDPOINT"
	envPaymentToken    = "PAYPACK_PAYMENT_TOKEN"
)

var vers = flag.Bool("version", false, "Print version information and exit")

type config struct {
	logLevel       string
	dbConfig       postgres.Config
	paymentEndoint string
	paymentToken   string
	httpPort       string
	secret         string
}

func main() {
	flag.Parse()
	if *vers {
		fmt.Println(build.String())
		os.Exit(0)
	}
	ctx, cancel := context.WithCancel(context.Background())

	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	db := connectToDB(cfg.dbConfig, logger)
	defer db.Close()

	novaCfg := &nova.Config{
		Endpoint: cfg.paymentEndoint,
		Token:    cfg.paymentToken,
	}
	pGateway := nova.New(novaCfg)

	users := newUserService(db, cfg.secret)
	transactions := newTransactionService(db, users)
	properties := newPropertyService(db, users)
	payment := newPaymentService(db, pGateway)
	feedback := newFeedbackService(db)
	owners := newOwnersService(db)

	payOpts := paymentEndpoints.HandlerOpts{
		Service: payment,
		Logger:  logger,
	}
	feedOpts := feedbackEndpoints.HandlerOpts{
		Service: feedback,
		Logger:  logger,
	}
	proOpts := prtEndpoints.HandlerOpts{
		Service: properties,
		Logger:  logger,
	}

	ownersOpts := ownersEndpoints.HandlerOpts{
		Service: owners,
		Logger:  logger,
	}

	usersOpts := usersEndpoints.HandlerOpts{
		Service: users,
		Logger:  logger,
	}

	srvOptions := serverOptions{
		transactions:  transactions,
		payOptions:    &payOpts,
		feedOptions:   &feedOpts,
		proOptions:    &proOpts,
		ownersOptions: &ownersOpts,
		usersOptions:  &usersOpts,
		port:          cfg.httpPort,
		logger:        logger,
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		logger.Info("signal caught. shutting down...")
		cancel()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		startHTTPServer(ctx, srvOptions)
	}()

	wg.Wait()
}

func loadConfig() config {
	dbConfig := postgres.Config{
		Host:        paypack.Env(envDBHost, defDBHost),
		Port:        paypack.Env(envDBPort, defDBPort),
		User:        paypack.Env(envDBUser, defDBUser),
		Pass:        paypack.Env(envDBPass, defDBPass),
		Name:        paypack.Env(envDBName, defDBName),
		SSLMode:     paypack.Env(envDBSSLMode, defDBSSLMode),
		SSLCert:     paypack.Env(envDBSSLCert, defDBSSLCert),
		SSLKey:      paypack.Env(envDBSSLKey, defDBSSLKey),
		SSLRootCert: paypack.Env(envDBSSLRootCert, defDBSSLRootCert),
	}
	return config{
		logLevel:       paypack.Env(envLogLevel, defLogLevel),
		dbConfig:       dbConfig,
		paymentEndoint: paypack.Env(envPaymentEndpoint, defPaymentEndpoint),
		paymentToken:   paypack.Env(envPaymentToken, defPaymentToken),
		httpPort:       paypack.Env(envHTTPPort, defHTTPPort),
		secret:         paypack.Env(envSecret, defSecret),
	}
}

func connectToDB(config postgres.Config, logger logger.Logger) *sql.DB {
	db, err := postgres.Connect(config)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to connect to postgres: %s", err))
		os.Exit(1)
	}
	return db
}

func newUserService(db *sql.DB, secret string) users.Service {
	hasher := bcrypt.New()
	tempid := jwt.New(secret)
	idp := uuid.New()
	store := postgres.NewUserStore(db)
	return users.New(hasher, tempid, idp, store)
}

func newTransactionService(db *sql.DB, users users.Service) transactions.Service {
	cfg := &nanoid.Config{
		Length: transactions.Length, Alphabet: transactions.Alphabet,
	}
	idp := nanoid.New(cfg)
	store := postgres.NewTransactionStore(db)
	auth := transactions.NewAuthBackend(users)
	return transactions.New(idp, store, auth)
}

func newPropertyService(db *sql.DB, users users.Service) properties.Service {
	cfg := &nanoid.Config{
		Length: properties.Length, Alphabet: properties.Alphabet,
	}
	idp := nanoid.New(cfg)
	props := postgres.NewPropertyStore(db)
	auth := properties.NewAuthBackend(users)
	return properties.New(idp, props, auth)
}

func newOwnersService(db *sql.DB) owners.Service {
	repo := postgres.NewOwnerRepo(db)
	idp := uuid.New()
	opts := &owners.Options{
		Repo: repo,
		Idp:  idp,
	}
	return owners.New(opts)
}

func newPaymentService(db *sql.DB, gw payment.Gateway) payment.Service {
	transactions := postgres.NewTransactionStore(db)
	properties := postgres.NewPropertyStore(db)

	repoOptions := &payment.RepoOptions{
		Transactions: transactions,
		Properties:   properties,
	}
	repo := payment.NewRepo(repoOptions)

	cfg := &nanoid.Config{
		Length: payment.Length, Alphabet: payment.Alphabet,
	}

	idp := nanoid.New(cfg)

	opts := &payment.ServiceOptions{
		Gateway:    gw,
		IDP:        idp,
		Repository: repo,
	}
	return payment.New(opts)
}

func newFeedbackService(db *sql.DB) feedback.Service {
	repo := postgres.NewMessageStore(db)
	idp := uuid.New()
	opts := &feedback.Options{
		Repo: repo,
		Idp:  idp,
	}
	return feedback.New(opts)
}

type serverOptions struct {
	// old API
	//users        users.Service
	transactions transactions.Service
	//properties   properties.Service

	// new API
	payOptions    *paymentEndpoints.HandlerOpts
	feedOptions   *feedbackEndpoints.HandlerOpts
	ownersOptions *ownersEndpoints.HandlerOpts
	proOptions    *prtEndpoints.HandlerOpts
	usersOptions  *usersEndpoints.HandlerOpts

	logger logger.Logger
	port   string
}

func startHTTPServer(ctx context.Context, opts serverOptions) {

	if opts.logger == nil {
		panic("can't use a nil logger")
	}

	if opts.payOptions == nil || opts.feedOptions == nil || opts.ownersOptions == nil || opts.proOptions == nil {
		panic("absolutely unacceptable start server opts")
	}

	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	router := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(false)

	router.HandleFunc("/healthz", health.Health).Methods(http.MethodGet)

	router.HandleFunc("/version", version.Build).Methods(http.MethodGet)

	feedbackEndpoints.RegisterHandlers(router, opts.feedOptions)

	ownersEndpoints.RegisterHandlers(router, opts.ownersOptions)

	paymentEndpoints.RegisterHandlers(router, opts.payOptions)

	prtEndpoints.RegisterHandlers(router, opts.proOptions)

	usersEndpoints.RegisterHandlers(router, opts.usersOptions)

	trxRoutes := router.PathPrefix("/transactions").Subrouter()
	trxEndpoints.MakeAdapter(trxRoutes)(opts.transactions)

	s := &http.Server{
		Addr:        fmt.Sprintf(":%s", opts.port),
		Handler:     cors(router),
		ReadTimeout: 2 * time.Minute,
	}

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			opts.logger.Error(fmt.Sprintf("paypack backend service stopped with error %v", err))
		}
		close(done)
	}()

	opts.logger.Info(fmt.Sprintf("serving api at http://127.0.0.1:%s", opts.port))
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		opts.logger.Error(fmt.Sprintf("paypack backend service stopped with error %v", err))
	}
	<-done
}
