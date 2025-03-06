

run:
	go run .

build:
	go build -o bin/server .

swag:
	swag init -g main.go

clean:
	$(RM) -rf *.out *.html
test:
	go test -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

.PHONY: server build swag clean run test