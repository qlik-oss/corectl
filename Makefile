test:
	go test ./...

build:
	go build -o ./cli main.go

.DEFAULT_GOAL := build