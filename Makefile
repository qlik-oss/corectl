SHELL := /bin/bash

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git rev-parse HEAD | cut -c1-12)

build:
	go build -ldflags "-X main.version=$(shell ./bump.sh)-dev -X main.branch=$(BRANCH) -X main.commit=$(COMMIT)" -o corectl main.go

install:
	go install -ldflags "-X main.version=$(shell ./bump.sh)-dev -X main.branch=$(BRANCH) -X main.commit=$(COMMIT)"

docs: build
	@./corectl generate-docs
	@./corectl generate-spec

lint:
	go fmt ./...
	@$(shell which golint || go get -u golang.org/x/lint/golint)
	golint -set_exit_status

start-deps:
	AcceptEUAL=yes docker-compose -f test/docker-compose.yml up -d

test: start-deps
	go test ./... -tags=integration -count=1 -race
	docker-compose -f test/docker-compose.yml down

c.out:
	./coverage.sh

coverage: c.out
	go tool cover -html=c.out -o coverage.html
	rm c.out

example:
	@echo "Starting engine in docker"
	AcceptEUAL=yes docker-compose -f examples/docker-compose.yml up -d
	@echo "Building corectl"
	go build
	./corectl build --config examples/corectl.yml

.PHONY: test docs
