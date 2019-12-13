SHELL := /bin/bash

build:
	go build -ldflags "-X main.version=dev-$(shell git rev-parse HEAD | cut -c1-12)"

docs:
	@go build -ldflags "-X main.version=$(shell git describe --tag --abbrev=0)" -o corectl_temp main.go
	@./corectl_temp generate-docs
	@./corectl_temp generate-spec
	@rm ./corectl_temp

lint:
	go fmt ./...

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
