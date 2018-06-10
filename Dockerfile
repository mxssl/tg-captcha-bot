FROM golang:alpine as builder

# Устанавливаем зависимости
RUN apk add --no-cache ca-certificates curl git

# Получаем внешние пакеты и компилируем бинарник
WORKDIR /go/src/go-app
COPY . .
RUN go get -u gopkg.in/tucnak/telebot.v2
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o gobot

# Копируем бинарник в чистый образ Alpine Linux
FROM alpine:latest
WORKDIR /
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/go-app/gobot .
RUN chmod +x gobot
CMD ["./gobot"]
