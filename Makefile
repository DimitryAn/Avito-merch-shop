include .env
export

dbLogin:
	psql $(DATABASE_URL)
	
dbCreateMigrations:
	migrate create -ext sql -dir migrations -seq schema

dbMigrations:
	migrate -database $(DATABASE_URL) -path migrations up

dbMigrAndLogin:
	migrate -database $(DATABASE_URL) -path migrations up
	psql $(DATABASE_URL)