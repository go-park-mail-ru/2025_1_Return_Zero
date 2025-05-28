# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . . 

RUN go build -o main microservices/artist/cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/config.yaml . 

COPY --from=builder /app/.env .

EXPOSE 5001

CMD ["./main"]