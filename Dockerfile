FROM golang:1.24-alpine

RUN apk add --no-cache git bash postgresql-client netcat-openbsd

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd
RUN go build -o /app/main .

# Явно копируем скрипт в корень контейнера
COPY wait-for-postgres.sh /app/wait-for-postgres.sh
RUN chmod +x /app/wait-for-postgres.sh

WORKDIR /app
CMD ["/bin/sh", "-c", "./wait-for-postgres.sh postgres:5432 -- /app/main"]
