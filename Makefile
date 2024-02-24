MAIN_PACKAGE_PATH := ./cmd/wallet
BINARY_NAME := wallet

## tidy: format code and tidy modfile
.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

## build: build the application
.PHONY: build
build:
	go build -o=./bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}
	chmod +x ./bin/${BINARY_NAME}

## run: run the  application
.PHONY: run
run: build
	./bin/${BINARY_NAME}