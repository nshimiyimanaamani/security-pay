package config

// RedisConfig ...
type RedisConfig struct {
	URL      string `validate:"required" envconfig:"PAYPACK_REDIS_URL"`
	Password string `envconfig:"PAYPACK_REDIS_PASS"`
	DB       string `envconfig:"PAYPACK_REDIS_DB"`
}
