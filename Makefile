.DEFAULT_GOAL := serve

TARGET := rest-and-go
PWD := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
LINT_CLI := run -v

lint:
	docker run -t --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.55.2 golangci-lint ${LINT_CLI}

format:
	gofmt -l -s -w .

build:
	go build -o ${TARGET} cmd/main.go && mv ${TARGET} cmd/build/

serve: build
	./cmd/build/${TARGET}
