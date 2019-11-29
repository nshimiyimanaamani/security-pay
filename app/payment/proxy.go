package payment

import (
	"github.com/rugwirobaker/paypack-backend/app/properties"
	"github.com/rugwirobaker/paypack-backend/app/transactions"
)

type repository struct {
	transactions transactions.Repository
	properties   properties.Repository
}

// RepoOptions ...
type RepoOptions struct {
	Transactions transactions.Repository
	Properties   properties.Repository
}

// NewRepo creates a repo proxy
func NewRepo(opts *RepoOptions) Repository {
	return &repository{
		transactions: opts.Transactions,
		properties:   opts.Properties,
	}
}

func (repo *repository) SaveTransaction(tx Transaction) (string, error) {
	return "", nil
}

func (repo *repository) UpdateTransaction(tx Transaction) error {
	return nil
}

func (repo *repository) RetrieveProperty(id string) (Property, error) {
	empty := Property{}
	pro, err := repo.properties.RetrieveByID(id)
	if err != nil {
		return empty, nil
	}

	property := Property{
		ID:      pro.ID,
		OwnerID: pro.Owner.ID,
	}
	return property, nil
}
