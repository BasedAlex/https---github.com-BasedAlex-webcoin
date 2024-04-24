SHELL := /bin/bash

WEBCOIN_BINARY=webcoinApp

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images..
	docker-compose up -d 
	@echo Docker images started

## down: stop docker compose
down:
	@echo Stopping docker compose 
	docker-compose down
	@echo Done

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_webcoin
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when requireed) and starting docker images..."
	docker-compose up --build -d 
	@echo "Docker images built and started!"

## build_webcoin: builds the webcoin binary as a linux executable
build_webcoin:
	@echo Building webcoin binary...
	set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${WEBCOIN_BINARY} ./cmd/webcoin
	@echo Done

## runs webcoin linux executable 
run:
	./${WEBCOIN_BINARY}

## installs tools for linting
tools:
	go install github.com/daixiang0/gci@latest
	go install mvdan.cc/gofumpt@latest

## lint: runs golangci-lint on the app
lint:
	go mod tidy
	gofumpt -w .
	gci write . --skip-generated -s standard -s default
	golangci-lint run ./...

## runs lint and tool install
run_lint: tools lint

## builds and runs the app
run_build: build_webcoin run