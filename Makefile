# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./
BINARY_NAME := stickeroutline

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

## build: build the application
.PHONY: build
build:
	go build -o=./bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## run: run the  application
.PHONY: run
run: build
	./bin/${BINARY_NAME}
