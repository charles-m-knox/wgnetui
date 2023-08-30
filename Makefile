.PHONY: test coverage build lint

test:
	go test ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

build:
	go build

build-prod:
	go build -ldflags="-w -s -buildid=" -trimpath
	upx --best -o ./wgnetui-tmp wgnetui
	mv wgnetui-tmp wgnetui

lint:
	golangci-lint run

run:
	FYNE_SCALE=0.7 ./wgnetui
