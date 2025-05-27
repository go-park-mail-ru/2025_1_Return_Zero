# Build stage
FROM golang:1.24-alpine AS builder

ARG ENV_FILE

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./gen ./gen

COPY ./internal/pkg/helpers/customErrors ./internal/pkg/helpers/customErrors

COPY ./internal/pkg/helpers/logger ./internal/pkg/helpers/logger

COPY ./init ./init

COPY ./config ./config

COPY ./microservices/artist ./microservices/artist

COPY ./microservices/interceptors ./microservices/interceptors

COPY ./microservices/metrics ./microservices/metrics

COPY ./config.yaml ./config.yaml

RUN echo "$ENV_FILE" > .env

RUN go build -o main microservices/artist/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/config.yaml . 

COPY --from=builder /app/.env .

EXPOSE 5001

CMD ["./main"]