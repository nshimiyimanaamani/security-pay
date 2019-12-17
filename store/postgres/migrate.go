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

					`CREATE TABLE IF NOT EXISTS users (
						username    VARCHAR(254),
						password 	CHAR(60)	 NOT NULL,
						role	 	SMALLINT NOT NULL DEFAULT 1 CHECK(role in (1, 2, 3, 4)),
						account		UUID NOT NULL,
						created_at 	TIMESTAMP,
						updated_at  TIMESTAMP,
						FOREIGN KEY(account) references accounts(id) ON DELETE CASCADE ON UPDATE CASCADE,
						UNIQUE(username, role),
						PRIMARY KEY (username)
					);`,

					`CREATE TABLE IF NOT EXISTS developers (
						email 	VARCHAR(254) PRIMARY KEY REFERENCES users(username),
						role	SMALLINT NOT NULL DEFAULT(1) check (role = 1),
						FOREIGN KEY(email, role) REFERENCES users(username, role)
					);`,

					`CREATE TABLE IF NOT EXISTS admins (
						email 		VARCHAR(254) PRIMARY KEY REFERENCES users(username),
						role	 	SMALLINT NOT NULL DEFAULT(2) check (role = 2),
						FOREIGN KEY(email, role) REFERENCES users(username, role)
					);`,

					`CREATE TABLE IF NOT EXISTS managers (
						email 		VARCHAR(254) PRIMARY KEY REFERENCES users(username),
						cell 		TEXT NOT NULL DEFAULT 'not set',
						role	 	SMALLINT NOT NULL DEFAULT(2) check (role = 3),
						FOREIGN KEY(email, role) REFERENCES users(username, role)
					);`,

					`CREATE table IF NOT EXISTS agents (
						telephone 	VARCHAR(254) PRIMARY KEY REFERENCES users(username),
						first_name 	TEXT NOT NULL DEFAULT 'not set', 
						last_name 	TEXT NOT NULL DEFAULT 'not set',
						cell 		VARCHAR(254) NOT NULL DEFAULT 'not set',
						sector 		VARCHAR(254) NOT NULL DEFAULT 'not set',
						village 	VARCHAR(254) NOT NULL DEFAULT 'not set',
						role	 	SMALLINT NOT NULL DEFAULT(3)  check (role = 4),
						FOREIGN KEY(telephone, role) REFERENCES users(username, role)
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
						recorded_by VARCHAR(15) NOT NULL,
						FOREIGN KEY(recorded_by) references users(username) ON DELETE CASCADE ON UPDATE CASCADE,
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
					"DROP TABLE accounts",
					"DROP TABLE agents",
					"DROP TABLE managers",
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
