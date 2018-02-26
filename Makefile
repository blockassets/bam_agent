BINARY=bam_agent
BINARY_LINUX=bam_agent-linux-arm

.DEFAULT_GOAL:=all

DATE=$(shell date -u '+%Y-%m-%d %H:%M:%S')
COMMIT=$(shell git log --format=%h -1)

build: VERSION=main.version=$(COMMIT) $(DATE)
build: COMPILE_FLAGS=-o $(BINARY) -ldflags="-X '$(VERSION)'"
build:
	go build $(COMPILE_FLAGS)

arm-build: GOOS=linux
arm-build: GOARCH=arm
arm-build: GOARM=7
arm-build: VERSION=main.version=$(TRAVIS_BUILD_NUMBER) $(COMMIT) $(DATE) $(GOOS) $(GOARCH)
arm-build: COMPILE_FLAGS=-o $(BINARY_LINUX) -ldflags="-s -w -X '$(VERSION)'" # -s -w makes binary size smaller
arm-build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build $(COMPILE_FLAGS)

# | ensures execution order
arm: | clean arm-build rice gzip

gzip:
	gzip -9 $(BINARY_LINUX)

test:
	@go test ./...

dep:
	@dep ensure
	@go get github.com/GeertJohan/go.rice/rice

rice:
	rice append --exec $(BINARY_LINUX)

clean:
	@rm -f $(BINARY) $(BINARY_LINUX).gz

all: clean test build
