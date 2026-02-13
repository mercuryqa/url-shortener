# Создание бинарника
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY ../go.mod go.sum ./
RUN go mod download

COPY .. .

RUN go build -o url-service ./cmd/main.go

# Создание контейнера
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/url-service .
COPY --from=builder /app/internal/repo/migrations ./internal/repo/migrations

COPY ../.env .

EXPOSE 8082
CMD ["./url-service"]
