package main

import (
	"database/sql"

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
)

func newTransactionService(db *sql.DB, users users.Service) transactions.Service {
	idp := uuid.New()
	store := postgres.NewTransactionStore(db)
	return transactions.New(idp, store)
}

func newUserService(db *sql.DB, secret string) users.Service {
	hasher := bcrypt.New()
	tempid := jwt.New(secret)
	idp := uuid.New()
	store := postgres.NewUserStore(db)
	return users.New(hasher, tempid, idp, store)

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
