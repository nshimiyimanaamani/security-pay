package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // required driver for postgres access
	migrate "github.com/rubenv/sql-migrate"
)

const (
	errDuplicate  = "unique_violation"
	errFK         = "foreign_key_violation"
	errInvalid    = "invalid_text_representation"
	errTruncation = "string_data_right_truncation"
)

// Config defines the options that are used when connecting to a PostgreSQL instance
type Config struct {
	Host        string
	Port        string
	User        string
	Pass        string
	Name        string
	SSLMode     string
	SSLCert     string
	SSLKey      string
	SSLRootCert string
}

//Connect creates and returns a connection to a PostgreSQl instance.
//and returns a non-nil error if there is a failure.
func Connect(cfg Config) (*sql.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s sslcert=%s sslkey=%s sslrootcert=%s", cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSLMode, cfg.SSLCert, cfg.SSLKey, cfg.SSLRootCert)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	if err := migrateDB(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDB(db *sql.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "paypack_1",

				Up: []string{
					`CREATE TABLE IF NOT EXISTS users (
						id       	UUID,
						email    	VARCHAR(254) UNIQUE,
						password 	CHAR(60)	 NOT NULL,
						PRIMARY  	KEY (id)
					)`,
					`CREATE TABLE IF NOT EXISTS transactions (
						id 			UUID,
						property 	UUID,
						amount    	VARCHAR(254),
						method  	VARCHAR(254),
						PRIMARY KEY	(id)
					)`,
				},

				Down: []string{
					"DROP TABLE users",
					"DROP TABLE transactions",
				},
			},
			{

				Id: "paypack_2",

				Up: []string{
					`CREATE TABLE IF NOT EXISTS owners (
						id	   		UUID,
						fname  		VARCHAR(1024) NOT NULL,
						lname  		VARCHAR(1024) NOT NULL,
						phone  		VARCHAR(15)   NOT NULL,
						PRIMARY 	KEY(id)
					)`,

					`CREATE TABLE IF NOT EXISTS properties (
						id			UUID,
						owner		UUID		 NOT NULL,
						sector		VARCHAR(254) NOT NULL,
						cell		VARCHAR(254) NOT NULL,
						village		VARCHAR(254) NOT NULL,
						FOREIGN KEY(owner) references owners(id) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY 	KEY(id)
					)`,
				},

				Down: []string{
					"DROP TABLE properties",
				},
			},
			{
				Id: "paypack_3",

				Up: []string{
					`ALTER TABLE properties ADD COLUMN due NUMERIC (9, 2) NOT NULL DEFAULT (0);`,
				},

				Down: []string{
					`ALTER TABLE properties DROP COLUMN  monthlty_due`,
				},
			},
			{
				Id: "paypack_4",

				Up: []string{
					`ALTER TABLE transactions ADD COLUMN date_recorded TIMESTAMP;`,
				},

				Down: []string{
					`ALTER TABLE transactions DROP COLUMN  date_recorded`,
				},
			},
		},
	}
	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	return err
}
