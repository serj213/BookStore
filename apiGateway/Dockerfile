# Используем мультиархитектурный образ с поддержкой C
FROM --platform=linux/amd64 golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем с включенным CGO
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/api-gateway ./cmd/main.go

# Финальный образ с зависимостями
FROM debian:stable-slim

# Устанавливаем runtime-зависимости
RUN apt-get update && apt-get install -y \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/api-gateway .
COPY --from=builder /app/config/ ./config

RUN chmod +x /app/api-gateway

EXPOSE 8083
CMD ["/app/api-gateway"] 