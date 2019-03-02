BINARY_NAME=bot
CURRENT_DIR=$(shell pwd)
export GO111MODULE=on

.PHONY: all build clean lint critic test dep

all: dep build

build:
	go build -o ${BINARY_NAME} -v

clean:
	rm -f ${BINARY_NAME}

lint:
	golangci-lint run

critic:
	gocritic check-project ${CURRENT_DIR}

test:
	go test -v ./...

init:
	go mod init

tidy:
	go mod tidy

release-dry:
	goreleaser release --rm-dist --snapshot --skip-publish

release:
	goreleaser release --rm-dist

release-docker:
	docker build --tag mxssl/tg-captcha-bot:$(shell git tag --list | tail -n 1) .
	docker push mxssl/tg-captcha-bot:$(shell git tag --list | tail -n 1)
