include .envrc

.PHONY: help
## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'


# Migration Related
MIGRATE=migrate -database $(DATABASE_URL) -path=./migrations

.PHONY: migrate/version/get
## migrate/version/get: get current migration version
migrate/version/get:
	$(MIGRATE) version


.PHONY: migrate/version/goto
## migrate/version/goto: go to the specified version for ${ver}
migrate/version/goto:
	$(MIGRATE) goto $(ver)


.PHONY: migrate/version/force
## migrate/version/force: forcefully go to specific version for ${ver}
migrate/version/force:
	$(MIGRATE) force $(ver)


.PHONY: migrate/create
## migrate/create: create a new migration
migrate/create:
	migrate create -seq -ext=.sql -dir=./migrations $(name)

.PHONY: migrate/up
## migrate/up: migrate to the next version
migrate/up:
	$(MIGRATE) up


.PHONY: migrate/down
## migrate/down: migrate to the previous version
migrate/down:
	echo "y" | $(MIGRATE) down;



.PHONY: migrate/reset
## migrate/reset: reset all table
migrate/reset:
	echo "y" | $(MIGRATE) down;
	$(MIGRATE) up;


# Golang related
.PHONY: dev
## dev: start go-air, enable hot reload
dev:
	air

## tidy: format all .go files and tidy module dependencies
.PHONY: tidy
tidy:
	@echo 'Formatting .go files...'
	go fmt ./...
	@echo 'Tidying module dependencies...'
	go mod tidy

## audit: run quality control checks
.PHONY: audit
audit:
	@echo 'Checking module dependencies'
	go mod tidy -diff
	go mod verify
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...


## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo 'Building cmd/api...'
	go build -ldflags='-s' -o=./bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/api ./cmd/api


## build/docker: build docker image of the application
.PHONY: build/docker
build/docker:
	@echo 'Building docker image...'
	docker build -t golang-app .
