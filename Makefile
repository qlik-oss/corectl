test:
	go test ./...

build:
	go build -o cli -v

.DEFAULT_GOAL := build