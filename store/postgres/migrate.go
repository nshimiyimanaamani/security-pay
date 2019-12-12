package postgres

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

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
				},

				Down: []string{
					"DROP TABLE users",
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
						id			TEXT,
						owner		UUID,
						sector		VARCHAR(254) NOT NULL,
						cell		VARCHAR(254) NOT NULL,
						village		VARCHAR(254) NOT NULL,
						FOREIGN KEY(owner) references owners(id) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY 	KEY(id)
					)`,
					`CREATE TABLE IF NOT EXISTS transactions (
						id 				UUID,
						madeby 			UUID,
						madefor			TEXT,
						amount    		VARCHAR(254),
						method  		VARCHAR(254),
						date_modified 	TIMESTAMP,
						is_valid    	BOOLEAN DEFAULT false,
						FOREIGN KEY(madefor) references properties(id) ON DELETE CASCADE ON UPDATE CASCADE,
						FOREIGN KEY(madeby) references owners(id) ON DELETE CASCADE ON UPDATE CASCADE,
						PRIMARY KEY	(id)
					)`,
				},

				Down: []string{
					"DROP TABLE owners",
					"DROP TABLE properties",
					"DROP TABLE transactions",
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
			{
				Id: "paypack_5",

				Up: []string{
					`ALTER TABLE users ADD COLUMN cell VARCHAR(254);`,
				},

				Down: []string{
					`ALTER TABLE users DROP COLUMN  cell;`,
				},
			},
			{
				Id: "paypack_6",

				Up: []string{
					`ALTER TABLE owners 
						ADD UNIQUE (phone);
					`,
				},

				Down: []string{
					`ALTER TABLE owners 
						DROP COLUMN  password;
					`,
				},
			},
			{
				Id: "paypack_7",

				Up: []string{
					`CREATE TABLE IF NOT EXISTS messages (
						id 			UUID,
						title 		TEXT NOT NULL,
						body  		TEXT,
						hidden 		BOOLEAN DEFAULT false,
						creator	VARCHAR(15) NOT NULL,
						created_at 	TIMESTAMP,
						updated_at  TIMESTAMP,
						PRIMARY KEY(id)
					)`,
				},

				Down: []string{
					`DROP TABLE messages;
					`,
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
			{
				Id: "paypack_8",

				Up: []string{
					`ALTER TABLE properties
						ADD COLUMN occupied BOOLEAN DEFAULT FALSE,
						ADD COLUMN recorded_by UUID NOT NULL,
						ADD CONSTRAINT recorded_by FOREIGN KEY(recorded_by) references users(id);
					`,
				},

				Down: []string{
					`ALTER TABLE properties
						DROP COLUMN occupied;
						DROP CONSTRAINT recorded_by;
						DROP COLUMN recorded_by;
					`,
				},
			},
		},
	}
	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	return err
}
