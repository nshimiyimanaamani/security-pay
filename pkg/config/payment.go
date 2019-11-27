package config

// Payment ...
type Payment struct {
	PaymentEndoint string `validate:"required" envconfig:"PAYPACK_PAYMENT_ENDPOINT"`
	PaymentToken   string `validate:"required" envconfig:"PAYPACK_PAYMENT_TOKEN"`
}
