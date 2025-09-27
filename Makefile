BINARY_NAME=bot
CURRENT_DIR=$(shell pwd)
TAG=$(shell git name-rev --tags --name-only $(shell git rev-parse HEAD))
DOCKER_REGISTRY=mxssl
export GO111MODULE=on

.PHONY: all build clean lint critic test

all: build

build:
	go build -v -o ${BINARY_NAME}

clean:
	rm -f ${BINARY_NAME}

lint:
	golangci-lint run -v

test:
	go test -v ./...

init:
	go mod init

tidy:
	go mod tidy

update-deps:
	go get -u ./...

github-release-dry:
	@echo "TAG: ${TAG}"
	goreleaser release --rm-dist --snapshot --skip-publish

github-release:
	@echo "TAG: ${TAG}"
	goreleaser release --rm-dist

docker-release:
	@echo "Registry: ${DOCKER_REGISTRY}"
	@echo "TAG: ${TAG}"
	docker build --tag ${DOCKER_REGISTRY}/tg-captcha-bot:${TAG} --tag ${DOCKER_REGISTRY}/tg-captcha-bot:latest .
	docker push ${DOCKER_REGISTRY}/tg-captcha-bot:${TAG}
	docker push ${DOCKER_REGISTRY}/tg-captcha-bot:latest

