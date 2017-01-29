PACKAGES := $$(go list ./... | grep -v /vendor/ | grep -v /cmd/)

test: setupdb
	@echo "Running tests..."
	DB_USER=postgres DB_NAME=bw_accounts_test go test $(PACKAGES)


setupdb:
	@echo "Setting up test database..."
	psql -U postgres -c "DROP DATABASE IF EXISTS bw_accounts_test;"
	psql -U postgres -c "CREATE DATABASE bw_accounts_test;"
	psql -U postgres bw_accounts_test < db/schema.sql
