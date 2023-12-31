package app

import (
	"database/sql"

	"github.com/go-redis/redis/v7"
	"github.com/nshimiyimanaamani/paypack-backend/core/accounts"
	"github.com/nshimiyimanaamani/paypack-backend/core/auth"
	"github.com/nshimiyimanaamani/paypack-backend/core/feedback"
	"github.com/nshimiyimanaamani/paypack-backend/core/invoices"
	"github.com/nshimiyimanaamani/paypack-backend/core/metrics"
	"github.com/nshimiyimanaamani/paypack-backend/core/nanoid"
	"github.com/nshimiyimanaamani/paypack-backend/core/notifs"
	"github.com/nshimiyimanaamani/paypack-backend/core/owners"
	"github.com/nshimiyimanaamani/paypack-backend/core/payment"
	"github.com/nshimiyimanaamani/paypack-backend/core/properties"
	"github.com/nshimiyimanaamani/paypack-backend/core/scheduler"
	"github.com/nshimiyimanaamani/paypack-backend/core/transactions"
	"github.com/nshimiyimanaamani/paypack-backend/core/users"
	"github.com/nshimiyimanaamani/paypack-backend/core/ussd"
	"github.com/nshimiyimanaamani/paypack-backend/core/uuid"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/encrypt"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/passwords/bcrypt"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/passwords/rand"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/tasks/queue"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/tokens/jwt"
	"github.com/nshimiyimanaamani/paypack-backend/store/postgres"
	rstore "github.com/nshimiyimanaamani/paypack-backend/store/redis"
)

// Services aggrates all the services
type Services struct {
	Accounts      accounts.Service
	Auth          auth.Service
	Feedback      feedback.Service
	Notifications notifs.Service
	Owners        owners.Service
	Payment       payment.Service
	Properties    properties.Service
	Transactions  transactions.Service
	Users         users.Service
	Invoices      invoices.Service
	Stats         metrics.Service
	USSD          ussd.Service
	Scheduler     scheduler.Service
	PaymentStore  payment.Repository
}

// Init initialises all services
func Init(
	db *sql.DB,
	rclient *redis.Client,
	queue *queue.Queue,
	pclient payment.Client,
	sms notifs.Backend,
	secret string,
	namespace string,
	prefix string,
) *Services {
	notifs := bootNotifService(db, sms)
	services := &Services{
		Accounts:      bootAccountsService(db),
		Feedback:      bootFeedbackService(db),
		Notifications: notifs,
		Owners:        bootOwnersService(db),
		Payment:       bootPaymentService(db, rclient, sms, pclient),
		Properties:    bootPropertiesService(db),
		Transactions:  bootTransactionsService(db),
		Users:         bootUserService(db, secret),
		Auth:          bootAuthService(db, secret),
		Invoices:      bootInvoiceService(db),
		Stats:         bootStatsService(db),
		Scheduler:     bootScheduler(db, queue),
		USSD:          bootUSSDService(prefix, db, rclient, sms, pclient),
		PaymentStore:  postgres.NewPaymentRepository(db, notifs),
	}
	return services
}

func bootAuthService(db *sql.DB, secret string) auth.Service {
	hasher := bcrypt.New()
	repo := postgres.NewAuthRepository(db)
	jwt := jwt.New(secret)
	encrypter, _ := encrypt.New(secret)
	opts := &auth.Options{Hasher: hasher, Repo: repo, JWT: jwt, Encrypter: encrypter}
	return auth.New(opts)
}

// bootUserService configures the users service
func bootUserService(db *sql.DB, secret string) users.Service {
	hasher := bcrypt.New()
	generator := rand.New()
	repo := postgres.NewUserRepository(db)
	encrypter, _ := encrypt.New(secret)
	opts := &users.Options{Repo: repo, Hasher: hasher, PGen: generator, Encrypter: encrypter}
	return users.New(opts)
}

// bootPropertyService configures the properties service
func bootPropertiesService(db *sql.DB) properties.Service {
	cfg := &nanoid.Config{Length: properties.Length, Alphabet: properties.Alphabet}
	idp := nanoid.New(cfg)
	props := postgres.NewPropertyStore(db)
	return properties.New(idp, props)
}

// bootOwnersService configures the owners service
func bootOwnersService(db *sql.DB) owners.Service {
	repo := postgres.NewOwnerRepo(db)
	idp := uuid.New()
	opts := &owners.Options{Repo: repo, Idp: idp}
	return owners.New(opts)
}

// bootTransactionsService configures the transactions service
func bootTransactionsService(db *sql.DB) transactions.Service {
	repo := postgres.NewTransactionRepository(db)
	idp := uuid.New()
	opts := &transactions.Options{Repo: repo, Idp: idp}
	return transactions.New(opts)
}

// bootFeedbackService configures the feedback service
func bootFeedbackService(db *sql.DB) feedback.Service {
	repo := postgres.NewMessageRepo(db)
	idp := uuid.New()
	opts := &feedback.Options{Repo: repo, Idp: idp}
	return feedback.New(opts)
}

func bootPaymentService(db *sql.DB, rclient *redis.Client, nclient notifs.Backend, pclient payment.Client) payment.Service {
	var opts payment.Options
	opts.Backend = pclient
	opts.Idp = uuid.New()
	opts.SMS = bootNotifService(db, nclient)
	opts.Repository = postgres.NewPaymentRepository(db, opts.SMS)
	opts.Queue = rstore.NewQueue(rclient)
	opts.Properties = postgres.NewPropertyStore(db)
	opts.Owners = postgres.NewOwnerRepo(db)
	opts.Invoices = postgres.NewInvoiceRepository(db)
	opts.Transactions = postgres.NewTransactionRepository(db)
	return payment.New(&opts)
}

func bootAccountsService(db *sql.DB) accounts.Service {
	repo := postgres.NewAccountRepository(db)
	idp := uuid.New()
	opts := &accounts.Options{Repository: repo, IDP: idp}
	return accounts.New(opts)
}

func bootInvoiceService(db *sql.DB) invoices.Service {
	repo := postgres.NewInvoiceRepository(db)
	opts := &invoices.Options{Repo: repo}
	return invoices.New(opts)
}

func bootStatsService(db *sql.DB) metrics.Service {
	repo := postgres.NewStatsRepository(db)
	opts := &metrics.Options{Repo: repo}
	return metrics.New(opts)
}

func bootNotifService(db *sql.DB, client notifs.Backend) notifs.Service {
	var opts notifs.Options
	opts.IDP = uuid.New()
	opts.Backend = client
	opts.Store = postgres.NewNotifsRepository(db)
	return notifs.New(&opts)
}

func bootUSSDService(prefix string, db *sql.DB, rclient *redis.Client, sms notifs.Backend, pclient payment.Client) ussd.Service {
	idp := uuid.New()
	properties := postgres.NewPropertyStore(db)
	owners := postgres.NewOwnerRepo(db)
	payment := bootPaymentService(db, rclient, sms, pclient)
	agents := postgres.NewUserRepository(db)
	invoice := postgres.NewInvoiceRepository(db)
	opts := &ussd.Options{
		Prefix:     prefix,
		IDP:        idp,
		Owners:     owners,
		Properties: properties,
		Payment:    payment,
		Agents:     agents,
		Invoices:   invoice,
	}
	return ussd.New(opts)
}

func bootScheduler(db *sql.DB, queue *queue.Queue) scheduler.Service {
	var opts scheduler.Options
	opts.Queue = queue
	opts.Counter = postgres.NewAuditableCounter(db)
	opts.Invoices = postgres.NewInvoiceRepository(db)
	return scheduler.New(&opts)
}
