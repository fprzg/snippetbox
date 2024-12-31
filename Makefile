include .envrc

# =========================================================================== #
#
# Variables
#
# =========================================================================== #




# =========================================================================== #
#
# Rules 
#
# =========================================================================== #

## help: shows this message.
.PHONY: help
help: 
	@echo 'Usage'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


## confirm: asks the user for confirmation.
.PHONY: confirm
confirm: 
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]



# =========================================================================== #
#
# App
#
# =========================================================================== #

## app/install: Install the app (setup the database and certificates).
.PHONY: app/install
app/install: db/setup app/create-cert
	@cp .envrc.example .envrc && echo ".envrc file created."
	@echo "\nInstallation finished."

## app/create-cert: Create the PEM certificates. For development only.
.PHONY: app/create-cert
app/create-cert:
	@go run ${GOROOT}/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost && mv *.pem tls/

## app/run: Run the app (debug mode by default).
.PHONY: app/run
app/run: 
	echo ${SNIPPETBOX_ADDR}
	go run ./cmd/web -addr=${SNIPPETBOX_ADDR} -debug
	#go run ./cmd/web -addr=${SNIPPETBOX_ADDR} >>/tmp/snippetbox_info.log 2>> /tmp/snippetbox_error.log


## app/clean: Clean the app.
.PHONY: app/clean
app/clean: db/clean
	@echo "\nEnvironment cleaned"



# =========================================================================== #
#
# Database
#
# =========================================================================== #

define run_sql
	@mysql $(if $(2),-h localhost -u root $(2), -u root) < $(1) && \
		echo "$(notdir $(1)) execution succeded." || { \
		echo "$(notdir $(1)) execution failed."; exit 1; }
endef

# $(1) : Migration file
# $(2) : Datbase. If not provided runs the migration as root user.
# @sudo mysql $(if $(2), -h localhost -u root $(2), -u root) < $(1) &&
define run_migration
	sudo mysql -u root < $(1) && \
	echo "$(notdir $(1)) migration successful." || {\
	echo "$(notdir $(1)) migration failed."; exit 1; }
endef

## db/setup: Setup the database. Will abort with error if the database already exists.
.PHONY: db/setup
db/setup: 
	@for migration_script in $$(find migrations/ | sort | tail -n +2); do \
		$(call run_migration,$${migration_script}) \
	done
	@#$(call run_migration,migrations/0001-create-database.sql,)

## db/mysql: Open MySQL.
.PHONY: db/mysql
db/mysql: 
	@sudo mysql -u root snippetbox

## db/dummy-data: Insert dummy data (development).
.PHONY: db/dummy-data
db/dummy-data: 
	@sudo mysql -h localhost -u root snippetbox < scripts/db/insert-dummy-data.sql

## db/clean: Clean the database. Will exit with error if the database doesn't exist.
.PHONY: db/clean
db/clean: 
	@sudo mysql -u root < scripts/db/drop-database.sql && echo "Database and user erased"



# =========================================================================== #
#
# Production
#
# =========================================================================== #

## prod/deploy: Push to production.
.PHONY: prod/deploy
prod/deploy: confirm
	#whatever



# =========================================================================== #
#
# Test
#
# =========================================================================== #

test_coverage_profile="/tmp/snippetbox-profile.out"

## test/coverage: Shows the coverage of tests.
.PHONY: test/coverage
test/coverage: 
	go test -covermode=count -coverprofile=${test_coverage_profile} "./..."
	go tool cover -html=${test_coverage_profile}