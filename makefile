run:
	go run cmd/main.go

swag:
	swag init -g main.go

clean:
	$(RM) -rf *.out *.html
test:
	go test ./... -coverprofile=cover.out -coverpkg= . ./models
	go tool cover -html=cover.out -o cover.html

.PHONY: server build swag clean run test