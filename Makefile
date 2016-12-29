HEZILA_GO_EXECUTABLE ?= go
DIST_DIRS := find * -type d -exec
GIT_COMMIT=`git rev-parse --short HEAD`
GIT_BRANCH=`git rev-parse --abbrev-ref HEAD`
# NOTE: the `git tag` command is filtered through `grep .` so it returns non-zero when empty
GIT_TAG=`git tag --list "v*" --sort "v:refname" --points-at HEAD 2>/dev/null | tail -n 1 | grep . || echo "none"`
VERSION := $$(git describe --abbrev=0 --tags)
GOVERSION := 1.7

# Build Flags
BUILD_NUMBER ?= $(BUILD_NUMBER:)
BUILD_DATE = $(shell date -u)
BUILD_HASH = $(shell git rev-parse HEAD)
# If we don't set the build number it defaults to dev
ifeq ($(BUILD_NUMBER),)
	BUILD_NUMBER := dev
endif

# Golang Flags
GOPATH ?= $(GOPATH:):./vendor
GOFLAGS ?= $(GOFLAGS:)
GO=go

# Output paths
DIST_ROOT=dist
DIST_PATH=$(DIST_ROOT)/hezila

# Tests
TESTS=.

build:
	${HEZILA_GO_EXECUTABLE} build -o glide -ldflags "-X main.version=${VERSION}" hezila.go


setup: setup-ci
	go get -u github.com/githubnemo/CompileDaemon
	go get -u github.com/jstemmer/go-junit-report
	go get -u golang.org/x/tools/cmd/goimports

setup-ci:
	go get -u github.com/Masterminds/glide
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install
	glide install --strip-vendor


fmt:
	gofmt -w=true -s $$(find . -type f -name '*.go')
	goimports -w=true -d $$(find . -type f -name '*.go')

test:
	go test $$(glide nv)

test-race:
	go test -race $$(glide nv)

all: dist

dist: | check-style test

check-style:
	@echo Running GOFMT
	$(eval GOFMT_OUTPUT := $(shell gofmt -d -s cache/ math/ utils/ 2>&1))
	@echo "$(GOFMT_OUTPUT)"
	@if [ ! "$(GOFMT_OUTPUT)" ]; then \
		echo "gofmt sucess"; \
	else \
		echo "gofmt failure"; \
		exit 1; \
	fi

version:
	@echo $(REPO_VERSION)

clean:
	rm -f build/bin/*

.PHONY: all clean dev release test
