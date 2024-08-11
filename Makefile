MAIN_PACKAGE_PATH := ./cmd/wallet
BINARY_NAME := wallet

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## audit: run quality control checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go test -race -buildvcs -vet=off ./...

## db/migrations/new name=$1: create a new migration
.PHONY: db/docker
db/docker:
	docker run --name  ${BINARY_NAME}-pg  -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=${BINARY_NAME} -d -v "./resources/migrations":/docker-entrypoint-initdb.d postgres

## db/migrations/new name=$1: create a new migration
.PHONY: db/migrations/new
db/migrations/new:
	go run -tags 'pg' github.com/golang-migrate/migrate/v4/cmd/migrate@latest create -seq -ext=.sql -dir=./resources/migrations ${name}

## db/generate: run jet code generation
.PHONY: db/generate
db/generate:
	go run -tags 'jet' github.com/go-jet/jet/v2/cmd/jet@latest -dsn=postgresql://postgres:postgres@localhost:5432/${BINARY_NAME}?sslmode=disable -schema=public -path=./internal/repository/jet

## grpcui: start grpcui on port 50051
.PHONY: grpcui
grpcui:
	grpcui -plaintext :50051

## gogenerate:      run go codegen
.PHONY: gogenerate
gogenerate:: 
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/*.proto
	go generate ./...

## build: build the application
.PHONY: build
build:
	go build -o=./bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}
	chmod +x ./bin/${BINARY_NAME}

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=./bin/coverage.out ./...
	go tool cover -html=./bin/coverage.out

## run: run the  application
.PHONY: run
	go run ${MAIN_PACKAGE_PATH}
