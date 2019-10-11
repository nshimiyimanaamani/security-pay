package payment

import (
	"time"

	"github.com/rugwirobaker/paypack-backend/app/uuid"
)

// Service is the api interface to the payment module
type Service interface {
	Initilize(r Payment) (Message, error)

	Validate(r Validation) (Validation, error)
}

// ServiceOptions simplifies New func signature
type ServiceOptions struct {
	Gateway    Gateway
	Repository Repository
}
type service struct {
	gateway    Gateway
	repository Repository
}

// New initializes the payment service
func New(opts *ServiceOptions) Service {
	return &service{
		gateway:    opts.Gateway,
		repository: opts.Repository,
	}
}

func (svc service) Initilize(r Payment) (Message, error) {
	empty := Message{}
	property, err := svc.repository.RetrieveProperty(r.Code)
	if err != nil {
		return empty, err
	}

	transaction := Transaction{
		ID:           uuid.New().ID(),
		MadeFor:      property.ID,
		MadeBy:       property.OwnerID,
		Amount:       r.Amount,
		Method:       "mtn",
		DateRecorded: time.Now(),
	}
	id, err := svc.repository.SaveTransaction(transaction)
	if err != nil {
		return empty, err
	}
	r.ID = id
	return svc.gateway.Initiate(r)
}
func (svc service) Validate(r Validation) (Validation, error) {
	tx := Transaction{ID: r.ExternalTransactionsID, DateRecorded: time.Now()}

	res := svc.gateway.Validate(r)

	status := Status(r.Status)

	switch status {
	case Successful:
		return res, svc.repository.UpdateTransaction(tx)
	case Pending:
		return res, nil
	case Failed:
		return res, nil
	}
	return res, nil
}
