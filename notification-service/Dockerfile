FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o notification-service ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/notification-service .
COPY --from=builder /app/public ./public

EXPOSE 8081

CMD ["./notification-service"]