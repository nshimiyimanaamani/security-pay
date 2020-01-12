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
					`
					CREATE OR REPLACE FUNCTION trigger_set_timestamp()
					RETURNS TRIGGER AS $$
					BEGIN
					  NEW.updated_at = NOW();
					  RETURN NEW;
					END;
					$$ LANGUAGE plpgsql;`,

					`
					CREATE OR REPLACE FUNCTION trigger_set_invoice_status()
					RETURNS TRIGGER AS $$
					BEGIN
					  UPDATE invoices SET status='payed' WHERE status='pending' AND invoices.id=NEW.invoice;
					  RETURN NEW;
					END;
					$$ LANGUAGE plpgsql;`,

					`
					CREATE OR REPLACE FUNCTION trigger_initial_invoice()
					RETURNS TRIGGER AS $$
					BEGIN
						INSERT INTO invoices (amount, property) VALUES (NEW.due, NEW.id);
						RETURN NEW;
					END;
					$$ LANGUAGE plpgsql;`,

					`CREATE table IF NOT EXISTS sectors(
						sector			VARCHAR(256),
						created_at 		TIMESTAMP NOT NULL DEFAULT NOW(),
						updated_at 		TIMESTAMP NOT NULL DEFAULT NOW(),
						PRIMARY KEY(sector)
					);
					
					CREATE TRIGGER set_timestamp
					BEFORE UPDATE ON sectors
					FOR EACH ROW
					EXECUTE PROCEDURE trigger_set_timestamp();
					`,

					`CREATE table IF NOT EXISTS cells(
						cell			VARCHAR(256),
						sector			VARCHAR(256),
						created_at 		TIMESTAMP NOT NULL DEFAULT NOW(),
						updated_at 		TIMESTAMP NOT NULL DEFAULT NOW(),
						FOREIGN KEY (sector) references sectors(sector) ON DELETE CASCADE ON UPDATE CASCADE,
						UNIQUE(cell,sector),
						PRIMARY KEY(cell)
					);
					
					CREATE TRIGGER set_timestamp
					BEFORE UPDATE ON cells
					FOR EACH ROW
					EXECUTE PROCEDURE trigger_set_timestamp();
					`,

					`CREATE table IF NOT EXISTS villages(
						village 		VARCHAR(256),
						cell			VARCHAR(256),
						sector			VARCHAR(256),
						created_at 		TIMESTAMP NOT NULL DEFAULT NOW(),
						updated_at 		TIMESTAMP NOT NULL DEFAULT NOW(),
						FOREIGN KEY (cell, sector) references cells(cell, sector) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY KEY(village)
					);
					
					CREATE TRIGGER set_timestamp
					BEFORE UPDATE ON villages
					FOR EACH ROW
					EXECUTE PROCEDURE trigger_set_timestamp();
					`,

					`CREATE table IF NOT EXISTS accounts (
						id 				VARCHAR(256),
						name 			TEXT NOT NULL,
						type 			VARCHAR(3) NOT NULL CHECK(type in ('dev', 'ben')),
						active			BOOLEAN DEFAULT true,
						seats 			INTEGER,
						created_at 		TIMESTAMP NOT NULL DEFAULT NOW(),
						updated_at 		TIMESTAMP NOT NULL DEFAULT NOW(),
						FOREIGN KEY(id) references sectors(sector) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY KEY(id)
					);

					CREATE TRIGGER set_timestamp
					BEFORE UPDATE ON accounts
					FOR EACH ROW
					EXECUTE PROCEDURE trigger_set_timestamp();
					`,

					`CREATE TABLE IF NOT EXISTS users (
						username    VARCHAR(254),
						password 	CHAR(60)	 NOT NULL,
						role	 	VARCHAR(5) NOT NULL DEFAULT 'dev' CHECK(role in ('dev', 'admin', 'basic', 'min')),
						account		VARCHAR(256) NOT NULL,
						created_at 	TIMESTAMP NOT NULL DEFAULT NOW(),
						updated_at 	TIMESTAMP NOT NULL DEFAULT NOW(),
						FOREIGN KEY(account) references accounts(id) ON DELETE CASCADE ON UPDATE CASCADE,
						UNIQUE(username, role),
						PRIMARY KEY (username)
					);
					
					CREATE TRIGGER set_timestamp
					BEFORE UPDATE ON users
					FOR EACH ROW
					EXECUTE PROCEDURE trigger_set_timestamp();
					`,

					`CREATE TABLE IF NOT EXISTS developers (
						email 	VARCHAR(254) PRIMARY KEY,
						role	VARCHAR(5)  NOT NULL DEFAULT('dev') check (role = 'dev'),
						FOREIGN KEY(email, role) REFERENCES users(username, role) ON DELETE CASCADE ON UPDATE CASCADE
					);`,

					`CREATE TABLE IF NOT EXISTS admins (
						email 		VARCHAR(254) PRIMARY KEY,
						role	 	VARCHAR(5) NOT NULL DEFAULT('admin') check (role = 'admin'),
						FOREIGN KEY(email, role) REFERENCES users(username, role) ON DELETE CASCADE ON UPDATE CASCADE
					);`,

					`CREATE TABLE IF NOT EXISTS managers (
						email 		VARCHAR(254) PRIMARY KEY,
						cell 		TEXT NOT NULL DEFAULT 'not set',
						role	 	VARCHAR(5) NOT NULL DEFAULT('basic') check (role = 'basic'),
						FOREIGN KEY(email, role) REFERENCES users(username, role) ON DELETE CASCADE ON UPDATE CASCADE
					);`,

					`CREATE table IF NOT EXISTS agents (
						telephone 	VARCHAR(254) PRIMARY KEY,
						first_name 	TEXT NOT NULL DEFAULT 'not set', 
						last_name 	TEXT NOT NULL DEFAULT 'not set',
						cell 		VARCHAR(254) NOT NULL DEFAULT 'not set',
						sector 		VARCHAR(254) NOT NULL DEFAULT 'not set',
						village 	VARCHAR(254) NOT NULL DEFAULT 'not set',
						role	 	VARCHAR(5) NOT NULL DEFAULT('agent') check (role = 'min'),
						created_at 	TIMESTAMP NOT NULL DEFAULT NOW(),
						updated_at 	TIMESTAMP NOT NULL DEFAULT NOW(),
						FOREIGN KEY(telephone, role) REFERENCES users(username, role) ON DELETE CASCADE ON UPDATE CASCADE
					);
					
					CREATE TRIGGER set_timestamp
					BEFORE UPDATE ON agents
					FOR EACH ROW
					EXECUTE PROCEDURE trigger_set_timestamp();
					`,

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
						recorded_by VARCHAR(254) NOT NULL,
						created_at 	TIMESTAMP NOT NULL DEFAULT NOW(),
						updated_at 	TIMESTAMP NOT NULL DEFAULT NOW(),
						FOREIGN KEY(recorded_by) references users(username) ON DELETE CASCADE ON UPDATE CASCADE,
						FOREIGN KEY(owner) references owners(id) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY 	KEY(id)
					);
					
					CREATE TRIGGER set_timestamp
					BEFORE UPDATE ON properties
					FOR EACH ROW
					EXECUTE PROCEDURE trigger_set_timestamp();

					CREATE TRIGGER create_initial_invoice
					AFTER INSERT ON properties
					FOR EACH ROW
					EXECUTE PROCEDURE trigger_initial_invoice();
					`,

					`CREATE TABLE IF NOT EXISTS invoices (
						id SERIAL,
						amount 		NUMERIC (9, 2) NOT NULL DEFAULT (0),
						property 	TEXT,
						status 		VARCHAR(8) NOT NULL DEFAULT 'pending' CHECK(status in ('pending', 'payed')),
						created_at 	TIMESTAMP NOT NULL DEFAULT NOW(),
						updated_at 	TIMESTAMP NOT NULL DEFAULT NOW(),
						UNIQUE(id, amount),
						FOREIGN KEY(property) references properties(id) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY KEY(id)
					);

					CREATE TRIGGER set_timestamp
					BEFORE UPDATE ON invoices
					FOR EACH ROW
					EXECUTE PROCEDURE trigger_set_timestamp();
					`,

					`CREATE TABLE IF NOT EXISTS transactions (
						id 				UUID,
						madeby 			UUID,
						madefor			TEXT,
						invoice			SERIAL UNIQUE,
						amount    		NUMERIC (9, 2) NOT NULL DEFAULT (0),
						method  		VARCHAR(254),
						created_at 		TIMESTAMP NOT NULL DEFAULT NOW(),
						is_valid    	BOOLEAN DEFAULT false,
						FOREIGN KEY(madefor) references properties(id) ON DELETE CASCADE ON UPDATE CASCADE,
						FOREIGN KEY(madeby) references owners(id) ON DELETE CASCADE ON UPDATE CASCADE,
						FOREIGN KEY(invoice, amount) references invoices(id, amount) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY KEY	(id)
					);

					CREATE TRIGGER set_invoice_status
					AFTER INSERT ON transactions
					FOR EACH ROW
					EXECUTE PROCEDURE trigger_set_invoice_status();
					`,

					`CREATE TABLE IF NOT EXISTS messages (
						id 			UUID,
						title 		TEXT NOT NULL,
						body  		TEXT,
						hidden 		BOOLEAN DEFAULT false,
						creator		VARCHAR(15) NOT NULL,
						created_at 	TIMESTAMP NOT NULL DEFAULT NOW(),
						updated_at 	TIMESTAMP NOT NULL DEFAULT NOW(),
						PRIMARY KEY(id)
					);
					
					CREATE TRIGGER set_timestamp
					BEFORE UPDATE ON messages
					FOR EACH ROW
					EXECUTE PROCEDURE trigger_set_timestamp();
					`,

					`
					CREATE VIEW sectors_payment_view AS
					SELECT 
						SUM(invoices.amount) AS payed
					FROM invoices
						WHERE invoices.status='payed'
					`,
				},

				Down: []string{
					"DROP TABLE transactions",
					"DROP TABLE properties",
					"DROP TABLE owners",
					"DROP TABLE messages",
					"DROP TABLE admins",
					"DROP TABLE agents",
					"DROP TABLE managers",
					"DROP TABLE developers",
					"DROP TABLE users",
					"DROP TABLE accounts",
					"DROP TABLE sectors",
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
