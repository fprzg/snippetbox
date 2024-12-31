include .envrc



#################################################
#
# Rules 
#
#################################################

## help: shows this message.
.PHONY: help
help: 
	@echo 'Usage'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


## confirm: asks the user for confirmation.
.PHONY: confirm
confirm: 
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]


## app/install: Install the app (setup the database and certificates).
.PHONY: app/install
app/install: confirm


## app/run: Create the PEM certificates.
.PHONY app/cert
app/cert:

## app/run: Run the app (debug mode by default).
.PHONY: app/run
app/run: 
	go run ./cmd/web -addr=${SNIPPETBOX_ADDR} -debug
	#go run ./cmd/web -addr=$SNIPPETBOX_ADDR >>/tmp/snippetbox_info.log 2>> /tmp/snippetbox_error.log


## app/clean: Clean the app (database and certificates).
.PHONY: app/clean
app/clean: confirm

## db/setup: Setup the database. Will abort with error if the database already exists.
.PHONY: db/setup
db/setup: 


## db/mysql: Open MySQL.
.PHONY: db/mysql
db/mysql: 


## db/clean: Clean the database. Will exit with error if the database doesn't exist.
.PHONY: db/clean
db/clean: 


