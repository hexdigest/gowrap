export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)
export GOFLAGS := -mod=mod

GORELEASER_VERSION := "v1.25.0"

UNAME_S := $(shell uname -s)
UNAME_P := $(shell uname -p)
UNAME_M := $(shell uname -m)

GORELEASER_URL := "https://github.com/goreleaser/goreleaser/releases/download/$(GORELEASER_VERSION)/goreleaser_$(UNAME_S)_$(UNAME_P).tar.gz"

./bin:
	mkdir ./bin

./bin/gowrap: ./bin
	go install ./cmd/gowrap

./bin/golangci-lint: ./bin
	go install -modfile tools/go.mod github.com/golangci/golangci-lint/cmd/golangci-lint

./bin/goreleaser:
	curl -L $(GORELEASER_URL) -o ./bin/goreleaser.tar.gz
	@tar -xvf ./bin/goreleaser.tar.gz -C ./bin
	@touch ./bin/goreleaser

./bin/minimock: ./bin
	go install -modfile tools/go.mod github.com/gojuno/minimock/v3/cmd/minimock

.PHONY:
lint: ./bin/golangci-lint
	./bin/golangci-lint run --enable=goimports --disable=unused --exclude=S1023,"Error return value" ./...

.PHONY:
test:
	 go test -race -v ./...

.PHONY:
generate: ./bin/gowrap ./bin/minimock
	go generate ./...

.PHONY:
tidy:
	go mod tidy
	cd tools && go mod tidy
	cd templates_tests && go mod tidy

.PHONY:
all: ./bin/gowrap generate lint test

.PHONY:
release: ./bin/goreleaser
	GITHUB_TOKEN=`cat .gh_token` goreleaser release

.PHONY:
build: ./bin/goreleaser
	goreleaser build --snapshot --rm-dist

.PHONY:
clean:
	rm -rf ./bin
