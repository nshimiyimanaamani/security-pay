package app

import (
	"database/sql"

	"github.com/go-redis/redis/v7"
	"github.com/rugwirobaker/paypack-backend/core/accounts"
	"github.com/rugwirobaker/paypack-backend/core/auth"
	"github.com/rugwirobaker/paypack-backend/core/feedback"
	"github.com/rugwirobaker/paypack-backend/core/invoices"
	"github.com/rugwirobaker/paypack-backend/core/metrics"
	"github.com/rugwirobaker/paypack-backend/core/nanoid"
	"github.com/rugwirobaker/paypack-backend/core/notifications"
	"github.com/rugwirobaker/paypack-backend/core/owners"
	"github.com/rugwirobaker/paypack-backend/core/payment"
	"github.com/rugwirobaker/paypack-backend/core/properties"
	"github.com/rugwirobaker/paypack-backend/core/scheduler"
	"github.com/rugwirobaker/paypack-backend/core/transactions"
	"github.com/rugwirobaker/paypack-backend/core/users"
	"github.com/rugwirobaker/paypack-backend/core/ussd"
	"github.com/rugwirobaker/paypack-backend/core/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/encrypt"
	"github.com/rugwirobaker/paypack-backend/pkg/passwords/bcrypt"
	"github.com/rugwirobaker/paypack-backend/pkg/passwords/rand"
	"github.com/rugwirobaker/paypack-backend/pkg/tasks/queue"
	"github.com/rugwirobaker/paypack-backend/pkg/tokens/jwt"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	rstore "github.com/rugwirobaker/paypack-backend/store/redis"
)

// Services aggrates all the services
type Services struct {
	Accounts      accounts.Service
	Auth          auth.Service
	Feedback      feedback.Service
	Notifications notifications.Service
	Owners        owners.Service
	Payment       payment.Service
	Properties    properties.Service
	Transactions  transactions.Service
	Users         users.Service
	Invoices      invoices.Service
	Stats         metrics.Service
	USSD          ussd.Service
	Scheduler     scheduler.Service
}

// Init initialises all services
func Init(
	db *sql.DB,
	rclient *redis.Client,
	queue *queue.Queue,
	p payment.Backend,
	s notifications.Backend,
	secret string,
	namespace string,
	prefix string,
) *Services {
	services := &Services{
		Accounts:      bootAccountsService(db),
		Feedback:      bootFeedbackService(db),
		Notifications: bootNotifService(s),
		Owners:        bootOwnersService(db),
		Payment:       bootPaymentService(db, rclient, p),
		Properties:    bootPropertiesService(db),
		Transactions:  bootTransactionsService(db),
		Users:         bootUserService(db, secret),
		Auth:          bootAuthService(db, secret),
		Invoices:      bootInvoiceService(db),
		Stats:         bootStatsService(db),
		Scheduler:     bootScheduler(db, queue),
		USSD:          bootUSSDService(prefix, db, rclient, p),
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

func bootPaymentService(db *sql.DB, rclient *redis.Client, bc payment.Backend) payment.Service {
	idp := uuid.New()
	repo := postgres.NewPaymentRepo(db)
	queue := rstore.NewQueue(rclient)
	opts := &payment.Options{Idp: idp, Backend: bc, Repo: repo, Queue: queue}
	return payment.New(opts)
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

func bootNotifService(sms notifications.Backend) notifications.Service {
	idp := uuid.New()
	opts := &notifications.Options{Backend: sms, IDP: idp}
	return notifications.New(opts)
}

func bootUSSDService(prefix string, db *sql.DB, rclient *redis.Client, bc payment.Backend) ussd.Service {
	idp := uuid.New()
	properties := postgres.NewPropertyStore(db)
	owners := postgres.NewOwnerRepo(db)
	inv := postgres.NewPaymentRepo(db)
	payment := bootPaymentService(db, rclient, bc)
	opts := &ussd.Options{
		Prefix:     prefix,
		IDP:        idp,
		Owners:     owners,
		Properties: properties,
		Invoices:   inv,
		Payment:    payment,
	}
	return ussd.New(opts)
}

func bootScheduler(db *sql.DB, queue *queue.Queue) scheduler.Service {
	counter := postgres.NewCounter(db)
	opts := &scheduler.Options{Queue: queue, Auditable: counter}
	return scheduler.New(opts)
}
