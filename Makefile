.PHONY: test coverage build lint

test:
	go test ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

build:
	go build

lint:
	golangci-lint run
