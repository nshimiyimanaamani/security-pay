package config

// USSDConfig contains the setup configuration for ussd.
type USSDConfig struct {
	Prefix string `validate:"required" envconfig:"PAYPACK_USSD_PREFIX"`
}
