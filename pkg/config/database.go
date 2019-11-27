package config

// DBConfig ...
type DBConfig struct {
	URL string `validate:"required" envconfig:"PAYPACK_DB_URL"`
}
