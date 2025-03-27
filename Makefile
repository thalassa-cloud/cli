BINARY = tcloud
GOARCH = arm64

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

VERSION=0.0.0
IMAGE=thalassa-cloud/tcloud

ifneq (${BRANCH}, release)
	BRANCH := -${BRANCH}
else
	BRANCH :=
endif

PKG_LIST := $(shell go list ./...)
LDFLAGS = -ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

all: link clean linux darwin

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY}-linux-${GOARCH} . ;

darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY}-darwin-${GOARCH} . ;

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o bin/${BINARY}-windows-${GOARCH}.exe . ;

build:
	CGO_ENABLED=0 go build ${LDFLAGS} -o bin/${BINARY} . ;
	chmod +x bin/${BINARY};

test: ## Run unittests
	@go test -short ${PKG_LIST}

docs:
	go run tools/docs.go

clean:
	-rm -f bin/${BINARY}-* bin/${BINARY}

.PHONY: link linux darwin windows test fmt clean docs
