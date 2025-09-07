include .envrc

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## run/api: run the cmd/api application
.PHONY: run/api
run/api:
	@echo "go run ./cmd/api -db-dsn=*** -smtp-host=*** -smtp-sender=*** -smtp-username=*** -smtp-password=***"
	@go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN} -smtp-host=${GREENLIGHT_SMTP_HOST} -smtp-sender=${GREENLIGHT_SMTP_SENDER} -smtp-username=${GREENLIGHT_SMTP_USERNAME} -smtp-password=${GREENLIGHT_SMTP_PASSWORD} 


## run/api/help: show help info for cmd/api application
.PHONY: run/api/help
run/api/help:
	go run ./cmd/api -help

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	psql ${GREENLIGHT_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## db/migrations/up count=$1: apply all up database migrations, or specify count
.PHONY: db/migrations/up
db/migrations/up: confirm
	@echo 'Running up migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} up ${count}

## db/migrations/down count=$1: apply all down database migrations, or specify count
.PHONY: db/migrations/down
db/migrations/down: confirm
	@echo 'Running down migrations...'
	migrate -path ./migrations -database ${GREENLIGHT_DB_DSN} down ${count}

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #

## tidy: tidy module dependencies and format all .go files
.PHONY: tidy
tidy:
	@echo 'Tidying module dependencies...'
	go mod tidy
	@echo 'Verifying and vendoring module dependencies...'
	go mod verify
	go mod vendor
	@echo 'Formatting .go files...'
	go fmt ./...

## audit: run quality control checks
.PHONY: audit
audit:
	@echo 'Checking module dependencies...'
	go mod tidy -diff
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	go tool staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

# ==================================================================================== #
# BUILD
# ==================================================================================== #

## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api

# ==================================================================================== #
# PRODUCTION
# ==================================================================================== #

## production/setup: setup production server
.PHONY: production/setup
production/setup:
	 rsync -rP --delete ./remote/setup root@${PRODUCTION_HOST_IP}:/root
	 ssh -t root@${PRODUCTION_HOST_IP} "bash /root/setup/01.sh"

## production/connect: connect to the production server
.PHONY: production/connect
production/connect:
	ssh greenlight@${PRODUCTION_HOST_IP}

## production/deploy/api: deploy the api to production
.PHONY: production/deploy/api
production/deploy/api:
	rsync -P ./bin/linux_amd64/api greenlight@${PRODUCTION_HOST_IP}:~
	rsync -rP --delete ./migrations greenlight@${PRODUCTION_HOST_IP}:~
	rsync -P ./remote/production/api.service greenlight@${PRODUCTION_HOST_IP}:~
	rsync -P ./remote/production/Caddyfile greenlight@${PRODUCTION_HOST_IP}:~
	ssh -t greenlight@${PRODUCTION_HOST_IP} '\
		migrate -path ~/migrations -database $$GREENLIGHT_DB_DSN up \
		&& sudo mv ~/api.service /etc/systemd/system/ \
		&& sudo systemctl enable api \
		&& sudo systemctl restart api \
		&& sudo mv ~/Caddyfile /etc/caddy/ \
		&& sudo systemctl reload caddy \
	'
