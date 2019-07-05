package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	paypack "github.com/rugwirobaker/paypack-backend"
	prtEndpoints "github.com/rugwirobaker/paypack-backend/api/http/properties"
	trxAdapters "github.com/rugwirobaker/paypack-backend/api/http/transactions"
	usersAdapters "github.com/rugwirobaker/paypack-backend/api/http/users"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/app/users/bcrypt"
	"github.com/rugwirobaker/paypack-backend/app/users/jwt"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/logger"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
)

const (
	defLogLevel      = "error"
	defDBHost        = "localhost"
	defDBPort        = "5432"
	defDBUser        = "paypack"
	defDBPass        = "paypack"
	defDBName        = "users"
	defDBSSLMode     = "disable"
	defDBSSLCert     = ""
	defDBSSLKey      = ""
	defDBSSLRootCert = ""
	defHTTPPort      = "8080"
	defSecret        = "users"
	defServerCert    = ""
	defServerKey     = ""
	envLogLevel      = "PAYPACK_LOG_LEVEL"
	envDBHost        = "PAYPACK_DB_HOST"
	envDBPort        = "PAYPACK_DB_PORT"
	envDBUser        = "PAYPACK_DB_USER"
	envDBPass        = "PAYPACK_DB_PASS"
	envDBName        = "PAYPACK_DB"
	envDBSSLMode     = "PAYPACK_DB_SSL_MODE"
	envDBSSLCert     = "PAYPACK_DB_SSL_CERT"
	envDBSSLKey      = "PAYPACK_DB_SSL_KEY"
	envDBSSLRootCert = "PAYPACK_DB_SSL_ROOT_CERT"
	envHTTPPort      = "PORT"
	envSecret        = "PAYPACK_SECRET"
	envServerCert    = "PAYPACK_SERVER_CERT"
	envServerKey     = "PAYPACK_SERVER_KEY"
)

type config struct {
	logLevel string
	dbConfig postgres.Config
	httpPort string
	secret   string
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	db := connectToDB(cfg.dbConfig, logger)
	defer db.Close()

	users := newUserService(db, cfg.secret)
	transactions := newTransactionService(db)
	properties := newPropertyService(db)

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
		startHTTPServer(ctx, users, transactions, properties, cfg.httpPort, logger)
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
		logLevel: paypack.Env(envLogLevel, defLogLevel),
		dbConfig: dbConfig,
		httpPort: paypack.Env(envHTTPPort, defHTTPPort),
		secret:   paypack.Env(envSecret, defSecret),
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

func newTransactionService(db *sql.DB) transactions.Service {
	idp := uuid.New()
	store := postgres.NewTransactionStore(db)
	return transactions.New(idp, store)
}

func newPropertyService(db *sql.DB) properties.Service {
	idp := uuid.New()
	store := postgres.NewPropertyStore(db)
	return properties.New(idp, store)
}

func startHTTPServer(ctx context.Context,
	users users.Service,
	trx transactions.Service,
	prt properties.Service,
	port string, logger logger.Logger,
) {
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	router := mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(true)

	userRoutes := router.PathPrefix("/users").Subrouter()
	usersAdapters.MakeAdapter(userRoutes)(users)

	trxRoutes := router.PathPrefix("/transactions").Subrouter()
	trxAdapters.MakeAdapter(trxRoutes)(trx)

	prtRoutes := router.PathPrefix("/properties").Subrouter()
	prtEndpoints.MakeEndpoint(prtRoutes)(prt)

	s := &http.Server{
		Addr:        fmt.Sprintf(":%s", port),
		Handler:     cors(router),
		ReadTimeout: 2 * time.Minute,
	}

	done := make(chan struct{})
	go func() {
		<-ctx.Done()
		if err := s.Shutdown(context.Background()); err != nil {
			logger.Error(fmt.Sprintf("paypack backend service stopped with error %v", err))
		}
		close(done)
	}()

	logger.Info(fmt.Sprintf("serving api at http://127.0.0.1:%s", port))
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error(fmt.Sprintf("paypack backend service stopped with error %v", err))
	}
	<-done
}
