SHELL := /bin/bash

start-deps:
	docker-compose -f test/docker-compose.yml up -d

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
	docker-compose -f examples/docker-compose.yml up -d
	@echo "Building corectl"
	go build
	./corectl build --config examples/corectl.yml

.PHONY: test
