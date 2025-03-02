files = $(wildcard *.go)

server:
	go run $(files)

build:
	go build -o bin/server $(files)

swag:
	swag init -g main.go

clean:
	$(RM) -rf *.out *.html

run:
	go run .

test:
	go test -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

.PHONY: server build swag clean run test
