BINARY_NAME=verity
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

.PHONY: all build clean test install build-all

all: build

build:
	go build ${LDFLAGS} -o ${BINARY_NAME} .

clean:
	rm -f ${BINARY_NAME}
	rm -rf dist/

test:
	go test -v ./...

install: build
	install -m 755 ${BINARY_NAME} /usr/local/bin/

# Build for all platforms
build-all: clean
	mkdir -p dist
	
	# macOS Intel
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-amd64 .
	
	# macOS Apple Silicon
	GOOS=darwin GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-darwin-arm64 .
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-amd64 .
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-linux-arm64 .
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o dist/${BINARY_NAME}-windows-amd64.exe .
	
	@echo "Built binaries:"
	@ls -lh dist/

run:
	go run . ${ARGS}
