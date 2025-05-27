# syntax=docker/dockerfile:1.7-labs
FROM golang:1.24-alpine AS builder

ARG ENV_FILE

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY --exclude=microservices --exclude=deploy . .

RUN echo "$ENV_FILE" > .env

RUN go build -o main cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/config.yaml . 

COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]