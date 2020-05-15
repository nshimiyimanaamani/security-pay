package config

// USSDConfig ...
type USSDConfig struct {
	ApplicationID string `envconfig:"PAYPACK_USSD_APP_REF"`
}
