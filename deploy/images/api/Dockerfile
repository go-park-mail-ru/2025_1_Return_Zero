# Build stage
FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=1

RUN apk add --no-cache gcc musl-dev libwebp-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main cmd/main.go

FROM alpine:latest

RUN apk add --no-cache libwebp

WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/config.yaml . 

COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]