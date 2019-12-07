package config

// PostgresConfig ...
type PostgresConfig struct {
	URL string `validate:"required" envconfig:"PAYPACK_POSTGRES_URL"`
}
