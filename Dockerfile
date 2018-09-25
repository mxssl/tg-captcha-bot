FROM golang:alpine as builder

WORKDIR /go/src/github.com/mxssl/tg-captcha-bot
COPY . .

# install dep package manager
RUN apk add --no-cache ca-certificates curl git
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o bot

# Copy compiled binary to clear Alpine Linux image
FROM alpine:latest
WORKDIR /
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/mxssl/tg-captcha-bot .
RUN chmod +x bot
CMD ["./bot"]
