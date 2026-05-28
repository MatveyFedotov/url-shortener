FROM golang:1.26.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Статическая сборка для linux/amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/app

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app /root/app

# Добавляем права на выполнение (на всякий случай)
RUN chmod +x /root/app

CMD ["/root/app"]