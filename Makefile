
all: test build

build:
	go build ./cmd/morgothd

test:
	go test ./...
