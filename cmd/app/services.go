package app

import (
	"database/sql"

	"github.com/rugwirobaker/paypack-backend/app/feedback"
	"github.com/rugwirobaker/paypack-backend/app/nanoid"
	"github.com/rugwirobaker/paypack-backend/app/owners"
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
	"github.com/rugwirobaker/paypack-backend/app/users"
	"github.com/rugwirobaker/paypack-backend/app/users/bcrypt"
	"github.com/rugwirobaker/paypack-backend/app/users/jwt"
	"github.com/rugwirobaker/paypack-backend/app/uuid"
	"github.com/rugwirobaker/paypack-backend/store/postgres"
)

// Services aggrates all the services
type Services struct {
	Feedback     feedback.Service
	Owners       owners.Service
	Properties   properties.Service
	Transactions transactions.Service
	Users        users.Service
}

// Init initialises all services
func Init(db *sql.DB, secret string) *Services {
	services := &Services{
		Feedback:     bootFeedbackService(db),
		Owners:       bootOwnersService(db),
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
