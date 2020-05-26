package postgres

import (
	"database/sql"

	migrate "github.com/rubenv/sql-migrate"
)

func migrateDB(db *sql.DB) error {
	migrations := &migrate.MemoryMigrationSource{
		Migrations: []*migrate.Migration{
			{
				Id: "001_initial",

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

					`
					CREATE OR REPLACE FUNCTION refresh_payment_metrics()
					RETURNS TRIGGER AS $$
					BEGIN
						refresh materialized view concurrently sector_payment_metrics;
						refresh materialized view concurrently cell_payment_metrics;
						refresh materialized view concurrently village_payment_metrics;
						RETURN NULL;
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
						password 	VARCHAR(60)	 NOT NULL,
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

					CREATE TRIGGER refresh_payment_view
					AFTER INSERT OR UPDATE OR DELETE ON invoices
					FOR EACH ROW
					EXECUTE PROCEDURE refresh_payment_metrics();
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
						FOREIGN KEY(creator) REFERENCES owners(phone) ON UPDATE CASCADE ON DELETE CASCADE,
						PRIMARY KEY(id)
					);
					
					CREATE TRIGGER set_timestamp
					BEFORE UPDATE ON messages
					FOR EACH ROW
					EXECUTE PROCEDURE trigger_set_timestamp();
					`,

					`
					create view payment_metrics as
						select 
							property,
							properties.sector,
							properties.cell,
							properties.village,
							date_trunc('month', invoices.created_at) as period,
							count(*) filter (where status='pending') as pending,
							count(*) filter (where status='payed') as payed,
							coalesce(sum(amount) filter(where status='pending') ,0) as pending_amount,
							coalesce(sum(amount) filter(where status='payed') ,0) as payed_amount
						from invoices
							join properties on invoices.property=properties.id
						group by 
							property,
							period,
							properties.sector, 
							properties.cell, 
							properties.village
						order by property; 
					`,

					`
					create materialized view sector_payment_metrics as
						select 
							sector,
							period,
							sum(pending) as pending_count, 
							sum(payed) as payed_count,
							coalesce(sum(pending_amount),0) as pending_amount,
							coalesce(sum(payed_amount),0) as payed_amount
						from payment_metrics group by sector, period;

					create unique index on  sector_payment_metrics(sector, period);
					`,

					`
					create materialized view cell_payment_metrics as
						select 
							cell,
							sector, 
							period,
							sum(pending) as pending_count, 
							sum(payed) as payed_count,
							coalesce(sum(pending_amount),0) as pending_amount,
							coalesce(sum(payed_amount),0) as payed_amount
						from payment_metrics group by cell, sector, period;
						
					create unique index on  cell_payment_metrics(cell, period);
					`,

					`
					create materialized view village_payment_metrics as
						select 
							village,
							cell,
							period,
							sum(pending) as pending_count, 
							sum(payed) as payed_count,
							coalesce(sum(pending_amount),0) as pending_amount,
							coalesce(sum(payed_amount),0) as payed_amount
						from payment_metrics group by village, cell, period;
					
					create unique index on  village_payment_metrics(village, period);
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
			{
				Id: "002_update_invoice",
				Up: []string{
					`
					CREATE  OR REPLACE FUNCTION correct_invoice_amount() RETURNS void AS
					$$
					DECLARE
						rec record;
					BEGIN
    					FOR rec IN
							select 
								invoices.id as invoice_id, 
								due, amount 
							from 
								properties 
							inner join invoices on properties.id=invoices.property
						LOOP
        					IF rec.due != rec.amount THEN
		    					UPDATE invoices SET amount=rec.due WHERE id=rec.invoice_id;
        					END IF;
						END LOOP;
					END;
					$$ LANGUAGE plpgsql;
					`,

					`SELECT correct_invoice_amount();`,

					`ALTER TABLE properties 
						ADD UNIQUE(id, due)`,

					`ALTER TABLE invoices
						DROP CONSTRAINT invoices_property_fkey,
						ADD FOREIGN KEY(property, amount) REFERENCES properties(id, due) ON UPDATE CASCADE ON DELETE CASCADE;
					`,
				},
			},
			{
				Id: "003_fix_metrics_indexes",
				Up: []string{
					`CREATE unique index on cell_payment_metrics(cell, sector, period);`,
					`CREATE unique index on village_payment_metrics(village, cell, period);`,

					`DROP index cell_payment_metrics_cell_period_idx;`,
					`DROP index village_payment_metrics_village_period_idx;`,
				},
			},

			{
				Id: "004_create_unique_index_on_invoice_created_at",
				Up: []string{
					`
					CREATE OR REPLACE FUNCTION start_of_month(p_date TIMESTAMP)
  					returns DATE
					AS
					$$
  				 		SELECT DATE_TRUNC('month', p_date)::date;
					$$
					language SQL
					immutable;
					`,

					`
					CREATE unique index single_invoice_per_property_per_month 
						ON invoices(property, start_of_month(created_at))
					`,
				},
			},
			{
				Id: "005_update_trigger_initial_invoice()",
				Up: []string{
					`
					CREATE OR REPLACE FUNCTION trigger_initial_invoice()
					RETURNS TRIGGER AS $$
					BEGIN
						INSERT INTO invoices (
							amount, property, created_at, updated_at
						) VALUES (
							NEW.due, NEW.id, NEW.created_at, New.updated_at
						);
						RETURN NEW;
					END;
					$$ LANGUAGE plpgsql;`,
				},
			},
			{
				Id: "006_create_one_month_old_properties_view",

				Up: []string{
					`
					CREATE materialized view one_month_old_properties_view AS
						select 
							id, due, created_at
						from 
							properties
						where 
							created_at < date_trunc('month', now())::date
						order 
							by id asc
					`,

					`CREATE unique index ON one_month_old_properties_view(id);`,

					// Make sure to refreshe the view for existing databases
					`refresh materialized view concurrently one_month_old_properties_view;`,

					`
					CREATE OR REPLACE FUNCTION refresh_one_month_old_properties_view()
					RETURNS TRIGGER AS $$
					BEGIN
						refresh materialized view concurrently one_month_old_properties_view;
						RETURN NULL;
					END;
					$$ LANGUAGE plpgsql;`,

					`
					CREATE TRIGGER trigger_refresh_one_month_old_properties_view
					AFTER INSERT OR UPDATE OR DELETE ON properties
					FOR EACH ROW
					EXECUTE PROCEDURE refresh_one_month_old_properties_view();
					`,
				},
			},
			{
				Id: "007_create_audit_func",
				Up: []string{
					`
					CREATE  OR REPLACE FUNCTION audit_func(x INT, y INT) RETURNS INT AS
					$$
					DECLARE
						count INT := 0;
						rec RECORD;
					BEGIN
						FOR rec IN
							select 
								id, due 
							from 
								one_month_old_properties_view
							order by id offset x limit y
						LOOP
							INSERT INTO invoices (property, amount) VALUES (rec.id, rec.due);
							count:= count + 1;
						END LOOP;

						RETURN count;
					END;
					$$ LANGUAGE plpgsql;
					`,
				},
			},
			{
				Id: "008_create_earliest_pending_invoices_view",
				Up: []string{
					`CREATE MATERIALIZED VIEW earliest_pending_invoices_view AS
						SELECT 
							invoices.*
						FROM
							(
								SELECT
									property, MIN(created_at) as created_at 
								FROM 
									invoices WHERE status='pending' GROUP BY property
							) AS pending
						INNER JOIN
							invoices
						ON 
							invoices.property=pending.property AND invoices.created_at=pending.created_at 
						ORDER BY 
							invoices.id
					`,

					`CREATE unique index ON earliest_pending_invoices_view(property);`,

					`REFRESH MATERIALIZED VIEW concurrently earliest_pending_invoices_view;`,

					`
					CREATE OR REPLACE FUNCTION refresh_earliest_pending_invoices_view()
					RETURNS TRIGGER AS $$
					BEGIN
						refresh materialized view concurrently earliest_pending_invoices_view;
						RETURN NULL;
					END;
					$$ LANGUAGE plpgsql;`,

					`
					CREATE TRIGGER trigger_refresh_earliest_pending_invoices_view
					AFTER INSERT OR UPDATE OR DELETE ON invoices
					FOR EACH ROW
					EXECUTE PROCEDURE refresh_earliest_pending_invoices_view();
					`,
				},
			},
			{
				Id: "009_drop_accounts_sector_fk_contraint",
				Up: []string{
					`ALTER TABLE accounts DROP CONSTRAINT accounts_id_fkey;`,
				},
			},
		},
	}
	_, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	return err
}
