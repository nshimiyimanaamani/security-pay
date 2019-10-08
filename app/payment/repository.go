package payment

// Repository ...
type Repository interface {
	SaveTransaction(Transaction) (string, error)

	UpdateTransaction(tx Transaction) error

	RetrieveProperty(id string) (Property, error)
}
