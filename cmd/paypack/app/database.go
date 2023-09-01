package app

import (
	"database/sql"

	"github.com/go-redis/redis/v7"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/config"
	"github.com/nshimiyimanaamani/paypack-backend/pkg/errors"
	"github.com/nshimiyimanaamani/paypack-backend/store/postgres"
)

// PostgresConnect returns a sql.DB connection to postgres
func PostgresConnect(config *config.PostgresConfig) (*sql.DB, error) {
	const op errors.Op = "app.PostgresConnect"

	db, err := postgres.Connect(config.URL)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}
	return db, nil
}

// RedisConnect returns a redis client
func RedisConnect(config *config.RedisConfig) (*redis.Client, error) {
	const op errors.Op = "app.RedisConnect"

	opts, err := redis.ParseURL(config.URL)
	if err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}

	client := redis.NewClient(opts)
	if _, err := client.Ping().Result(); err != nil {
		return nil, errors.E(op, err, errors.KindUnexpected)
	}
	return client, nil
}
