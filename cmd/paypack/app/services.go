package app

import (
	"database/sql"

	"github.com/go-redis/redis/v7"
	"github.com/rugwirobaker/paypack-backend/app/accounts"
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
	"github.com/rugwirobaker/paypack-backend/store/postgres"
	rstore "github.com/rugwirobaker/paypack-backend/store/redis"
)

// Services aggrates all the services
type Services struct {
	Accounts     accounts.Service
	Feedback     feedback.Service
	Owners       owners.Service
	Payment      payment.Service
	Properties   properties.Service
	Transactions transactions.Service
	Users        users.Service
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
		Users:        bootUserService(db, secret),
	}
	return services
}

// bootUserService configures the users service
func bootUserService(db *sql.DB, secret string) users.Service {
	hasher := bcrypt.New()
	tempid := jwt.New(secret)
	idp := uuid.New()
	store := postgres.NewUserStore(db)
	return users.New(hasher, tempid, idp, store)
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
