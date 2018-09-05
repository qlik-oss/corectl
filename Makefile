test:
	go test ./...

build:
	go build -o corectl -v

.DEFAULT_GOAL := build