files = $(wildcard *.go)

server:
	go run $(files)

build:
	go build -o bin/server $(files)

swag:
	swag init -g main.go

.PHONY: server build swag
