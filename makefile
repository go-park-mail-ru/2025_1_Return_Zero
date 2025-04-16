ifneq (,$(wildcard ./.env))
	include .env
	export
endif

run:
	go run cmd/main.go

swag:
	swag init -g cmd/main.go

test:
	./scripts/test_all.sh

migrate_up:
	tern migrate -c db/migrations/tern.conf --migrations db/migrations

migrate_down:
	tern migrate -c db/migrations/tern.conf --migrations db/migrations -d 0

populate:
	make migrate_down
	make migrate_up
	go run db/populate/main.go -file db/populate/data.sql

docker-up:
	cd deploy/ && make deploy

docker-remove:
	-docker stop $$(docker ps -q)             
	-docker rm -f $$(docker ps -aq)           
	-docker rmi -f $$(docker images -q)
	-docker image prune -f

clean:
	$(RM) -rf *.out *.html

build:
	go build -o bin/server cmd/main.go

.PHONY: server build swag clean run test