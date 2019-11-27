package postgres_test

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest"
	// _ "github.com/lib/pq" // required driver for postgres access

	"github.com/rugwirobaker/paypack-backend/store/postgres"
)

var db *sql.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	cfg := []string{
		"POSTGRES_USER=test",
		"POSTGRES_PASSWORD=test",
		"POSTGRES_DB=test",
	}

	container, err := pool.Run("postgres", "alpine", cfg)
	if err != nil {
		log.Fatalf("Could not start container: %s", err)
	}

	port := container.GetPort("5432/tcp")
	if err := pool.Retry(func() error {
		url := fmt.Sprintf("host=localhost port=%s user=test dbname=test password=test sslmode=disable", port)
		db, err = sql.Open("postgres", url)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	config := struct {
		Host        string
		Port        string
		User        string
		Pass        string
		Name        string
		SSLMode     string
		SSLCert     string
		SSLKey      string
		SSLRootCert string
	}{
		Host:        "localhost",
		Port:        port,
		User:        "test",
		Pass:        "test",
		Name:        "test",
		SSLMode:     "disable",
		SSLCert:     "",
		SSLKey:      "",
		SSLRootCert: "",
	}

	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s",
		config.Host, config.Port, config.User, config.Name, config.Pass, config.SSLMode, config.SSLCert, config.SSLKey, config.SSLRootCert,
	)

	if db, err = postgres.Connect(connString); err != nil {
		log.Fatalf("failed to connect to test DB: %v", err)
	}
	defer db.Close()

	code := m.Run()

	if err = pool.Purge(container); err != nil {
		log.Fatalf("Could not purge container: %s", err)
	}

	os.Exit(code)
}
