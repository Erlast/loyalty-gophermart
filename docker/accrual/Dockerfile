# Stage 1: Сборка приложения
FROM golang:1.21 AS builder

WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальные файлы проекта и собираем приложение
COPY . .

# Собираем приложение с отключением CGO и для Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /accrual ./cmd/accrual/main.go

# Stage 2: Создание минимального образа для запуска
FROM debian:buster-slim

WORKDIR /

# Копируем собранный бинарный файл из предыдущего стейджа
COPY --from=builder /accrual /accrual

EXPOSE 8090

# Указываем точку входа для приложения
ENTRYPOINT ["/accrual"]

