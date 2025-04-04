# Используем мультиархитектурный образ с поддержкой C
FROM --platform=linux/amd64 golang:1.22 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Устанавливаем зависимости для librdkafka
RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    pkg-config \
    && rm -rf /var/lib/apt/lists/*

COPY . .

# Собираем с включенным CGO
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/kafka-service ./cmd/main.go

# Финальный образ с зависимостями
FROM debian:stable-slim

# Устанавливаем runtime-зависимости
RUN apt-get update && apt-get install -y \
    librdkafka1 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/kafka-service .
COPY --from=builder /app/config/ ./config
COPY --from=builder /app/cmd/wait-for-it.sh .

RUN apt-get update && apt-get install -y librdkafka1 && rm -rf /var/lib/apt/lists/*

RUN chmod +x /app/kafka-service
RUN chmod +x wait-for-it.sh

CMD ["/app/kafka-service"] 