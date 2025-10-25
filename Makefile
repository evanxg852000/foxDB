APP_NAME=foxdb

.PHONY: build run test clean

build:
	go build -o ./cmd/$(APP_NAME) ./cmd/main.go

run: build
	./cmd/$(APP_NAME)

test:
	go test ./...

clean:
	rm -f ./cmd/$(APP_NAME)
