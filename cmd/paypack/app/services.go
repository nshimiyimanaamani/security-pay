package app

import (
	"database/sql"

	"github.com/go-redis/redis/v7"
	"github.com/rugwirobaker/paypack-backend/app/accounts"
	"github.com/rugwirobaker/paypack-backend/app/auth"
	"github.com/rugwirobaker/paypack-backend/app/feedback"
	"github.com/rugwirobaker/paypack-backend/app/invoices"
	"github.com/rugwirobaker/paypack-backend/app/nanoid"
	"github.com/rugwirobaker/paypack-backend/app/owners"
	"github.com/rugwirobaker/paypack-backend/app/payment"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/pkg/passwords/bcrypt"
	"github.com/rugwirobaker/paypack-backend/pkg/passwords/randgen"
	"github.com/rugwirobaker/paypack-backend/pkg/tokens/jwt"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	rstore "github.com/rugwirobaker/paypack-backend/store/redis"
)

// Services aggrates all the services
type Services struct {
	Accounts     accounts.Service
	Auth         auth.Service
	Feedback     feedback.Service
	Owners       owners.Service
	Payment      payment.Service
	Properties   properties.Service
	Transactions transactions.Service
	Users        users.Service
	Invoices     invoices.Service
}

// Init initialises all services
func Init(db *sql.DB, rclient *redis.Client, b payment.Backend, secret string) *Services {
	services := &Services{
		Accounts:     bootAccountsService(db),
		Feedback:     bootFeedbackService(db),
		Owners:       bootOwnersService(db),
		Payment:      bootPaymentService(db, rclient, b),
		Properties:   bootPropertiesService(db),
		Transactions: bootTransactionsService(db),
		Users:        bootUserService(db),
		Auth:         bootAuthService(db, secret),
		Invoices:     bootInvoiceService(db),
	}
	return services
}

func bootAuthService(db *sql.DB, secret string) auth.Service {
	hasher := bcrypt.New()
	repo := postgres.NewAuthRepository(db)
	jwt := jwt.New(secret)
	opts := &auth.Options{Hasher: hasher, Repo: repo, JWT: jwt}
	return auth.New(opts)
}

// bootUserService configures the users service
func bootUserService(db *sql.DB) users.Service {
	hasher := bcrypt.New()
	generator := randgen.New()
	repo := postgres.NewUserRepository(db)
	opts := &users.Options{Repo: repo, Hasher: hasher, PGen: generator}
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
	repo := postgres.NewMessageStore(db)
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
