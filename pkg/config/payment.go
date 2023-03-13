package config

import validate "github.com/go-playground/validator/v10"

// PaymentConfig ...
type PaymentConfig struct {
	PaymentURL string `validate:"required" envconfig:"PAYMENT_BASE_URL"`
	Secret     string `validate:"required" envconfig:"PAYMENT_CLIENT_SECRET"`
	AppID      string `validate:"required" envconfig:"PAYMENT_CLIENT_ID"`
}

// Validate PaymentConfig
func (conf *PaymentConfig) Validate() error {
	validator := validate.New()
	return validator.Struct(conf)
}
