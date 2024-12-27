# syntax=docker/dockerfile:1
FROM golang:1.20-alpine

WORKDIR /app

# Копируем go.mod/go.sum и ставим зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем всё остальное
COPY . .

# Собираем бинарник
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
