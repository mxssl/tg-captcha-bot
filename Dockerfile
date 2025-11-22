FROM golang:1.25-alpine as builder

WORKDIR /go/src/github.com/mxssl/tg-captcha-bot
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
FROM alpine:3.22
WORKDIR /
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/mxssl/tg-captcha-bot .
RUN chmod +x bot
CMD ["./bot"]
