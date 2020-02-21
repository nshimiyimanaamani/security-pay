package config

import validate "github.com/go-playground/validator/v10"

// PostgresConfig ...
type PostgresConfig struct {
	URL string `validate:"required" envconfig:"PAYPACK_POSTGRES_URL"`
}

func (conf *PostgresConfig) Validate() error {
	validator := validate.New()
	return validator.Struct(conf)
}
