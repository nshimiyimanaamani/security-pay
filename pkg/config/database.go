package config

import validate "github.com/go-playground/validator/v10"

// PostgresConfig ...
type PostgresConfig struct {
	URL string `validate:"required" envconfig:"DATABASE_URL"`
}

// Validate database configuration
func (conf *PostgresConfig) Validate() error {
	validator := validate.New()
	return validator.Struct(conf)
}
