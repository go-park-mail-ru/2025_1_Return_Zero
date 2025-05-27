FROM golang:1.24-alpine AS builder

ARG ENV_FILE_B64

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./gen ./gen

COPY ./internal/pkg/helpers/customErrors ./internal/pkg/helpers/customErrors

COPY ./internal/pkg/helpers/logger ./internal/pkg/helpers/logger

COPY ./init ./init

COPY ./config ./config

COPY ./microservices/auth ./microservices/auth

COPY ./microservices/interceptors ./microservices/interceptors

COPY ./microservices/metrics ./microservices/metrics

COPY ./config.yaml ./config.yaml

RUN echo "$ENV_FILE_B64" | base64 -d > .env

RUN go build -o main microservices/auth/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/config.yaml . 

COPY --from=builder /app/.env .

EXPOSE 5004

CMD ["./main"]