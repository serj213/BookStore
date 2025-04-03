# Используем мультиархитектурный образ с поддержкой C
FROM --platform=linux/amd64 golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем с включенным CGO
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/book-service ./cmd/main.go

# Финальный образ с зависимостями
FROM debian:stable-slim

# Устанавливаем runtime-зависимости
RUN apt-get update && apt-get install -y \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/book-service .
COPY --from=builder /app/config/ ./config
COPY --from=builder /app/migrations/ ./migrations
COPY --from=builder /app/docs/ ./docs
COPY --from=builder /app/cmd/wait-for-it.sh .

RUN chmod +x wait-for-it.sh
RUN chmod +x /app/book-service

EXPOSE 44044
EXPOSE 8082
CMD ["/app/book-service"] 