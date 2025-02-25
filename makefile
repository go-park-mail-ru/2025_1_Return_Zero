server:
	go run main.go tracks.go albums.go artists.go helpers.go

build:
	go build -o bin/server main.go tracks.go albums.go artists.go helpers.go

swag:
	swag init -g main.go

.PHONY: server build swag
