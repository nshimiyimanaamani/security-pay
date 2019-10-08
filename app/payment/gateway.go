package payment

// Gateway is the interface to the payment API
type Gateway interface {
	Initiate(Payment) (Message, error)
	Validate(Validation) Validation
}
