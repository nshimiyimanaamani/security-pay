package payment

import (
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
)

type repository struct {
	transactions transactions.Store
	properties   properties.PropertyStore
}

// RepoOptions ...
type RepoOptions struct {
	Transactions transactions.Store
	Properties   properties.PropertyStore
}

// NewRepo creates a repo proxy
func NewRepo(opts *RepoOptions) Repository {
	return &repository{
		transactions: opts.Transactions,
		properties:   opts.Properties,
	}
}

func (repo *repository) SaveTransaction(tx Transaction) (string, error) {

	stx := transactions.Transaction{
		ID:           tx.ID,
		MadeFor:      tx.MadeFor,
		MadeBy:       tx.MadeBy,
		Amount:       tx.Amount,
		Method:       tx.Method,
		DateRecorded: tx.DateRecorded,
	}

	return repo.transactions.Save(stx)
}

func (repo *repository) UpdateTransaction(tx Transaction) error {
	stx := transactions.Transaction{
		ID:           tx.ID,
		MadeFor:      tx.MadeFor,
		Amount:       tx.Amount,
		Method:       tx.Method,
		DateRecorded: tx.DateRecorded,
	}
	return repo.transactions.UpdateTransaction(stx)
}

func (repo *repository) RetrieveProperty(id string) (Property, error) {
	empty := Property{}
	pro, err := repo.properties.RetrieveByID(id)
	if err != nil {
		return empty, nil
	}

	property := Property{
		ID:      pro.ID,
		OwnerID: pro.OwnerID,
	}
	return property, nil
}
