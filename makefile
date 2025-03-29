run:
	go run cmd/main.go

swag:
	swag init -g cmd/main.go

migrate_up:
	tern migrate -c db/migrations/tern.conf --migrations db/migrations

migrate_down:
	tern migrate -c db/migrations/tern.conf --migrations db/migrations -d -1

populate:
	make migrate_down
	make migrate_up
	go run db/populate/main.go -file db/populate/data.sql

docker_up:
	docker compose up

docker_remove:
	docker rm $(docker ps -a -q) && docker volume prune -f

clean:
	$(RM) -rf *.out *.html

build:
	go build -o bin/server cmd/main.go

test:
	go test ./... -coverprofile=cover.out -coverpkg= . ./models
	go tool cover -html=cover.out -o cover.html

.PHONY: server build swag clean run test