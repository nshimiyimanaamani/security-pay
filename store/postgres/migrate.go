package postgres

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

func migrateDB(db *sql.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "v1.0.0",

				Up: []string{
					`CREATE TABLE IF NOT EXISTS users (
						id       	UUID,
						username    VARCHAR(254) UNIQUE,
						cell 		VARCHAR(254) NOT NULL DEFAULT 'not set',
						sector 		VARCHAR(254) NOT NULL DEFAULT 'not set',
						village 	VARCHAR(254) NOT NULL DEFAULT 'not set',
						password 	CHAR(60)	 NOT NULL,
						PRIMARY  	KEY (id)
					);`,

					`CREATE table IF NOT EXISTS accounts (
						id 				UUID,
						name 			TEXT NOT NULL,
						type 			VARCHAR(3) NOT NULL,
						active			BOOLEAN DEFAULT true,
						seats 			INTEGER,
						created_at 	TIMESTAMP,
						updated_at  TIMESTAMP,
						PRIMARY KEY(id)
					);`,

					`CREATE table IF NOT EXISTS developers (
						email 		TEXT NOT NULL,
						password 	VARCHAR(60) NOT NULL,
						account		UUID NOT NULL,
						created_at 	TIMESTAMP,
						updated_at  TIMESTAMP,
						FOREIGN KEY(account) references accounts(id) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY KEY(email)
					);`,

					`CREATE table IF NOT EXISTS managers (
						email 		TEXT NOT NULL,
						cell 		TEXT NOT NULL,
						password 	VARCHAR(60) NOT NULL,
						account		UUID NOT NULL,
						created_at 	TIMESTAMP,
						updated_at  TIMESTAMP,
						FOREIGN KEY(account) references accounts(id) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY KEY(email)
					);`,

					`CREATE table IF NOT EXISTS agents (
						telephone 	VARCHAR(15) NOT NULL,
						first_name 	TEXT NOT NULL DEFAULT 'not set', 
						last_name 	TEXT NOT NULL DEFAULT 'not set', 
						password 	VARCHAR(60) NOT NULL,
						account		UUID NOT NULL,
						created_at 	TIMESTAMP,
						updated_at  TIMESTAMP,
						FOREIGN KEY (account) references accounts(id) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY KEY(telephone)
					);`,

					`CREATE TABLE IF NOT EXISTS owners (
						id	   		UUID,
						fname  		VARCHAR(1024) NOT NULL,
						lname  		VARCHAR(1024) NOT NULL,
						phone  		VARCHAR(15)   NOT NULL UNIQUE,
						PRIMARY 	KEY(id)
					);`,

					`CREATE TABLE IF NOT EXISTS properties (
						id			TEXT,
						owner		UUID,
						due 		NUMERIC (9, 2) NOT NULL DEFAULT (0),
						sector		VARCHAR(254) NOT NULL,
						cell		VARCHAR(254) NOT NULL,
						village		VARCHAR(254) NOT NULL,
						for_rent 	BOOLEAN DEFAULT FALSE,
						occupied 	BOOLEAN DEFAULT TRUE,
						recorded_by UUID NOT NULL,
						FOREIGN KEY(recorded_by) references users(id) ON DELETE CASCADE ON UPDATE CASCADE,
						FOREIGN KEY(owner) references owners(id) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY 	KEY(id)
					);`,

					`CREATE TABLE IF NOT EXISTS transactions (
						id 				UUID,
						madeby 			UUID,
						madefor			TEXT,
						amount    		VARCHAR(254),
						method  		VARCHAR(254),
						date_recorded 	TIMESTAMP,
						is_valid    	BOOLEAN DEFAULT false,
						FOREIGN KEY(madefor) references properties(id) ON DELETE CASCADE ON UPDATE CASCADE,
						FOREIGN KEY(madeby) references owners(id) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY KEY	(id)
					);`,

					`CREATE TABLE IF NOT EXISTS messages (
						id 			UUID,
						title 		TEXT NOT NULL,
						body  		TEXT,
						hidden 		BOOLEAN DEFAULT false,
						creator		VARCHAR(15) NOT NULL,
						created_at 	TIMESTAMP,
						updated_at  TIMESTAMP,
						PRIMARY KEY(id)
					);`,
				},

				Down: []string{
					"DROP TABLE users",
					"DROP TABLE owners",
					"DROP TABLE properties",
					"DROP TABLE transactions",
					"DROP TABLE messages",
				},
			},
			// {
			// 	Id: "paypack_8",

			// 	Up: []string{
			// 		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`,

			// 		`ALTER EXTENSION "uuid-ossp" SET SCHEMA public;`,

			// 		`ALTER TABLE transactions ALTER COLUMN id TYPE uuid USING (uuid_generate_v4());`,

			// 		`ALTER TABLE owners DROP COLUMN  password;`,
			// 	},
			// 	Down: []string{
			// 		`ALTER TABLE transactions ALTER COLUMN id TYPE TEXT;`,

			// 		`ALTER TABLE owners ADD COLUMN password VARCHAR(60) NOT NULL;`,
			// 	},
			// },
		},
	}
	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	return err
}
