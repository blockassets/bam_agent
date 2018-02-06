DATE=$(shell date -u '+%Y-%m-%d %H:%M:%S')
COMMIT=$(shell git log --format=%h -1)
VERSION=main.version=${TRAVIS_BUILD_NUMBER} ${COMMIT} ${DATE}
BINARY=bam_agent
COMPILE_FLAGS=-o ${BINARY} -ldflags="-X '${VERSION}'"

rice:
	@rice append --exec ${BINARY}

build:
	@go build ${COMPILE_FLAGS}

arm-compile:
	@GOOS=linux GOARCH=arm GOARM=7 go build ${COMPILE_FLAGS}

# | ensures execution order
arm: | arm-compile rice

test:
	@go test .

dep:
	@dep ensure
	@go get github.com/GeertJohan/go.rice/rice

clean:
	@rm -f bam_agent

all: clean test build

.PHONY: build