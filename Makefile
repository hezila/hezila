REPO_VERSION := $$(git describe --abbrev=0 --tags)
BUILD_DATE := $$(date +%Y-%m-%d-%H:%M)
GOVERSION := 1.7


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
	# goimports -w=true -d $$(find . -type f -name '*.go')

test:
	go test $$(glide nv)

test-race:
	go test -race $$(glide nv)

version:
	@echo $(REPO_VERSION)

clean:
	rm -f build/bin/*
