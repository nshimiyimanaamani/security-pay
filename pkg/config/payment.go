package config

// PaymentConfig ...
type PaymentConfig struct {
	PaymentURL string `validate:"required" envconfig:"PAYPACK_PAYMENT_URL"`
	Secret     string `validate:"required" envconfig:"PAYPACK_PAYMENT_SECRET"`
	AppID      string `validate:"required" envconfig:"PAYPACK_PAYMENT_APP_ID"`
	// ChannelID  string `validate:"required" envconfig:"PAYPACK_PAYMENT_CHANNEL_ID"`
	Callback string `validate:"required" envconfig:"PAYPACK_PAYMENT_CALLBACK"`
}
