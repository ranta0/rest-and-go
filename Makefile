.DEFAULT_GOAL := serve

TARGET := rest-and-go

build:
	go build -o ${TARGET} cmd/main.go && mv ${TARGET} cmd/build/

serve: build
	./cmd/build/${TARGET}
