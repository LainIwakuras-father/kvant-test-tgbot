# syntax=docker/dockerfile:1.8

# ==================== Dev Stage ====================
FROM golang:1.25-alpine AS dev

WORKDIR /app

# Установка dev-зависимостей
RUN apk add --no-cache \
    git \
    bash \
    curl \
    make

# Установка стабильной версии air, совместимой с Go 1.24
RUN go install github.com/air-verse/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

# Копирование файлов зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Создание точки для дебаггера
EXPOSE 40000 2345

# Запуск air для hot reload
CMD ["air", "-c", ".air.toml"]