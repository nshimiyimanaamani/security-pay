package postgres

import (
	"database/sql"

	"github.com/rugwirobaker/paypack-backend/models"
	"github.com/rugwirobaker/paypack-backend/store/transactions"
)

var _ (transactions.Store) = (*transactionStore)(nil)

type transactionStore struct {
	db *sql.DB
}

//New instanctiates a new transactiob store interface
func New(db *sql.DB) transactions.Store {
	return &transactionStore{db}
}

func (str *transactionStore) Save(models.Transaction) (string, error) {
	return "", nil
}

func (str *transactionStore) RetrieveByID(string) (models.Transaction, error) {
	return models.Transaction{}, nil
}

func (str *transactionStore) RetrieveAll(uint64, uint64) models.TransactionPage {
	return models.TransactionPage{}
}

func (str *transactionStore) RetrieveByProperty(string, uint64, uint64) models.TransactionPage {
	return models.TransactionPage{}
}

func (str *transactionStore) RetrieveByMethod(string, uint64, uint64) models.TransactionPage {
	return models.TransactionPage{}
}

func (str *transactionStore) RetrieveByMonth(string, uint64, uint64) models.TransactionPage {
	return models.TransactionPage{}
}

func (str *transactionStore) RetrieveByYear(string, uint64, uint64) models.TransactionPage {
	return models.TransactionPage{}
}
