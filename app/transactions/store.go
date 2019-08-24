package transactions

// Store defines the api to the transactions data store
type Store interface {
	// Save adds a new transactiob to the data store returns nil
	// if the operation is successful or otherwise an error.
	Save(Transaction) (string, error)

	// RetrieveByID retreives a transaction identified by the given id.
	RetrieveByID(string) (Transaction, error)

	// RetrieveAll retrieves the subset of transactions owned by the specified property.
	RetrieveAll(uint64, uint64) (TransactionPage, error)

	// RetrieveByMethod retrieves the subset of transactions that where made using the given method.
	RetrieveByProperty(string, uint64, uint64) (TransactionPage, error)

	// RetrieveByMethod retrieves the subset of transactions that where made using the given method.
	RetrieveByMethod(string, uint64, uint64) (TransactionPage, error)

	// RetrieveByMonth retrieves the subset of transactions that where made during the given month.
	RetrieveByMonth(string, uint64, uint64) (TransactionPage, error)

	// RetrieveByYear retrieves the subset of transactions that where made using the given year.
	RetrieveByYear(string, uint64, uint64) (TransactionPage, error)
}
