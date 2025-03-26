run:
	go run cmd/main.go

swag:
	swag init -g cmd/main.go

clean:
	$(RM) -rf *.out *.html

build:
	go build -o bin/server cmd/main.go

test:
	go test ./... -coverprofile=cover.out -coverpkg= . ./models
	go tool cover -html=cover.out -o cover.html

.PHONY: server build swag clean run test