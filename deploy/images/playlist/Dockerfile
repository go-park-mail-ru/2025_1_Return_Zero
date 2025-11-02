# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN go build -o main microservices/playlist/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

## config and environment provided at runtime via mounts/env_file

EXPOSE 5006

CMD ["./main"]
