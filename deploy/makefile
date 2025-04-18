ifneq (,$(wildcard ../.env))
	include ../.env
	export
endif

# deploy
deploy:
	docker compose up -d --build

# docker down
undeploy:
	docker compose down --volumes --remove-orphans --rmi all

# docker remove
docker-remove:
	docker rm $(docker ps) && docker volume prune -f

clean:
	$(RM) -rf *.out *.html

build:
	go build -o bin/server cmd/main.go

test:
	go test ./... -coverprofile=cover.out -coverpkg= . ./models
	go tool cover -html=cover.out -o cover.html

.PHONY: server build swag clean run test