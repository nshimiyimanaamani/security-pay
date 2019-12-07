package redis_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/ory/dockertest"
)

var redisClient *redis.Client

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	container, err := pool.Run("redis", "latest", nil)
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	if err := pool.Retry(func() error {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("localhost:%s", container.GetPort("6379/tcp")),
			Password: "",
			DB:       0,
		})
		return redisClient.Ping().Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	if err := pool.Purge(container); err != nil {
		log.Fatalf("Could not purge container: %s", err)
	}

	os.Exit(code)
}
