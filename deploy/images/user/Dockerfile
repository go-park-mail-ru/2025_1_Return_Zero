# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./gen ./gen

COPY ./internal/pkg/helpers/customErrors ./internal/pkg/helpers/customErrors

COPY ./internal/pkg/helpers/logger ./internal/pkg/helpers/logger

COPY ./init ./init

COPY ./config ./config

COPY ./microservices/user ./microservices/user

COPY ./microservices/interceptors ./microservices/interceptors

COPY ./config.yaml ./config.yaml

COPY ./.env ./.env

RUN go build -o main microservices/user/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/config.yaml . 

COPY --from=builder /app/.env .

EXPOSE 5005

CMD ["./main"]