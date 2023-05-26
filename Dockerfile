FROM golang:1.20.4-alpine3.17 as builder

WORKDIR /go/src/github.com/momai/tg-captcha-bot
COPY . .

# Install external dependcies
RUN apk add --no-cache \
  ca-certificates \
  curl \
  git

# Compile binary
RUN CGO_ENABLED=0 \
  go build -v -o bot

# Copy compiled binary to clear Alpine Linux image
FROM alpine:3.18.0
WORKDIR /
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/momai/tg-captcha-bot .
RUN chmod +x bot
CMD ["./bot"]
