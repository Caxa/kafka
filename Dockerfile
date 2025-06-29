FROM golang:1.24-alpine

RUN apk add --no-cache git

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем зависимости и качаем модули
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Устанавливаем рабочую директорию с main.go
WORKDIR /app/cmd

# Собираем приложение
RUN go build -o /app/main .

# Запускаем
CMD ["/app/main"]
