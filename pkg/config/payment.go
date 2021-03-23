package config

import validate "github.com/go-playground/validator/v10"

// PaymentConfig ...
type PaymentConfig struct {
	PaymentURL string `validate:"required" envconfig:"PAYPACK_PAYMENT_APP_URL"`
	Secret     string `validate:"required" envconfig:"PAYPACK_PAYMENT_APP_SECRET"`
	AppID      string `validate:"required" envconfig:"PAYPACK_PAYMENT_APP_ID"`
	DCallback  string `validate:"required" envconfig:"PAYPACK_PAYMENT_DEBIT_CALLBACK"`
	CCallback  string `validate:"required" envconfig:"PAYPACK_PAYMENT_CREDIT_CALLBACK"`
}

// Validate PaymentConfig
func (conf *PaymentConfig) Validate() error {
	validator := validate.New()
	return validator.Struct(conf)
}
