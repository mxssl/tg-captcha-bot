BINARY_NAME=bot

.PHONY: all build clean lint test dep build-linux build-darwin

all: build

build:
	go build -o ${BINARY_NAME} -v

clean:
	rm -f ${BINARY_NAME}

lint:
	golangci-lint run -v
	
test:
	go test -v ./...

dep:
	dep ensure
