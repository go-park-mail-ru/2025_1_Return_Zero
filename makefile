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
	tern migrate -c db/migrations/tern.conf --migrations db/migrations -d 1

populate:
	make migrate_up

docker-up:
	cd deploy/ && make deploy

docker-remove:
	-docker stop $$(docker ps -q)             
	-docker rm -f $$(docker ps -aq)           
	-docker rmi -f $$(docker images -q)
	-docker image prune -f

protogen-all:
	protoc -I proto proto/**/*.proto --go_out=gen --go-grpc_out=gen

mock-all:
	mockgen -source=gen/auth/auth_grpc.pb.go -destination=mocks/mock_auth_client.go -package=mocks AuthServiceClient
	mockgen -source=gen/user/user_grpc.pb.go -destination=mocks/mock_user_client.go -package=mocks UserServiceClient
	mockgen -source=gen/album/album_grpc.pb.go -destination=mocks/mock_album_client.go -package=mocks AlbumServiceClient
	mockgen -source=gen/playlist/playlist_grpc.pb.go -destination=mocks/mock_playlist_client.go -package=mocks PlaylistServiceClient
	mockgen -source=gen/artist/artist_grpc.pb.go -destination=mocks/mock_artist_client.go -package=mocks ArtistServiceClient
	mockgen -source=gen/track/track_grpc.pb.go -destination=mocks/mock_track_client.go -package=mocks TrackServiceClient

clean:
	$(RM) -rf *.out *.html

build:
	go build -o bin/server cmd/main.go

easyjson:
	easyjson -all internal/pkg/model/delivery

.PHONY: server build swag clean run test mock-all easyjson