.PHONY=test coverage build lint

BUILDDIR=build
VER=0.0.1
BIN=$(BUILDDIR)/wgnetui-v$(VER)

build-dev:
	CGO_ENABLED=1 go build -v

mkbuilddir:
	mkdir -p $(BUILDDIR)

build-prod: mkbuilddir
	CGO_ENABLED=1 go build -v -o $(BIN) -ldflags="-w -s -buildid=" -trimpath

run-dev:
	FYNE_SCALE=0.7 ./wgnetui

run-prod:
	FYNE_SCALE=0.7 ./$(BIN)

lint:
	golangci-lint run ./...

test:
	go test ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

compress-prod: mkbuilddir
	rm -f $(BIN)-compressed
	upx --best -o ./$(BIN)-compressed $(BIN)

# upx does not support mac currently

# rm -f $(BIN)-darwin-arm64-compressed
# note for mac m1 - this seems to taint the binary, it doesn't work;
# you'll probably have to do without upx for now
# upx --best -o ./$(BIN)-darwin-arm64-compressed $(BIN)-darwin-arm64
build-mac-arm64: mkbuilddir
	CGO_ENABLED=1 GOARCH=arm64 GOOS=darwin go build -v -o $(BIN)-darwin-arm64 -ldflags="-w -s -buildid=" -trimpath
	rm -f $(BIN)-darwin-arm64.xz
	xz -9 -e -T 12 -vv $(BIN)-darwin-arm64

# rm -f $(BIN)-darwin-amd64-compressed
# upx --best -o ./$(BIN)-darwin-arm64-compressed $(BIN)-darwin-amd64
build-mac-amd64: mkbuilddir
	CGO_ENABLED=1 GOARCH=amd64 GOOS=darwin go build -v -o $(BIN)-darwin-amd64 -ldflags="-w -s -buildid=" -trimpath
	rm -f $(BIN)-darwin-amd64.xz
	xz -9 -e -T 12 -vv $(BIN)-darwin-amd64

build-win-amd64: mkbuilddir
	CGO_ENABLED=1 GOARCH=amd64 GOOS=windows go build -v -o $(BIN)-win-amd64-uncompressed -ldflags="-w -s -buildid=" -trimpath
	rm -f $(BIN)-win-amd64
	upx --best -o ./$(BIN)-win-amd64 $(BIN)-win-amd64-uncompressed

build-linux-arm64: mkbuilddir
	CGO_ENABLED=1 GOARCH=arm64 GOOS=linux go build -v -o $(BIN)-linux-arm64-uncompressed -ldflags="-w -s -buildid=" -trimpath
	rm -f $(BIN)-linux-arm64
	upx --best -o ./$(BIN)-linux-arm64 $(BIN)-linux-arm64-uncompressed

build-linux-amd64: mkbuilddir
	CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -v -o $(BIN)-linux-amd64-uncompressed -ldflags="-w -s -buildid=" -trimpath
	rm -f $(BIN)-linux-amd64
	upx --best -o ./$(BIN)-linux-amd64 $(BIN)-linux-amd64-uncompressed

build-all: mkbuilddir build-linux-amd64 build-win-amd64 build-mac-amd64 # build-mac-arm64 build-linux-arm64

delete-builds:
	rm $(BUILDDIR)/*
