FROM golang:1.12.0-alpine as builder

ENV GO111MODULE=on

WORKDIR /go/src/github.com/mxssl/tg-captcha-bot
COPY . .

# Install external dependcies
RUN apk add --no-cache \
  ca-certificates \
  curl \
  git

# Compile binary
RUN CGO_ENABLED=0 \
  GOOS=`go env GOHOSTOS` \
  GOARCH=`go env GOHOSTARCH` \
  go build -o bot

# Copy compiled binary to clear Alpine Linux image
FROM alpine:3.9
WORKDIR /
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/mxssl/tg-captcha-bot .
RUN chmod +x bot
CMD ["./bot"]
