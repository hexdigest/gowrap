export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)
export GOFLAGS := -mod=mod

./bin:
	mkdir ./bin

./bin/gowrap: ./bin
	go install ./cmd/gowrap

./bin/golangci-lint: ./bin
	go install -modfile tools/go.mod github.com/golangci/golangci-lint/cmd/golangci-lint

./bin/goreleaser:
	go install -modfile tools/go.mod github.com/goreleaser/goreleaser

./bin/minimock: ./bin
	go install -modfile tools/go.mod github.com/gojuno/minimock/v3/cmd/minimock

.PHONY:
lint: ./bin/golangci-lint
	./bin/golangci-lint run --enable=goimports --disable=unused --exclude=S1023,"Error return value" ./...

.PHONY:
test:
	 go test -race ./...

.PHONY:
generate: ./bin/gowrap ./bin/minimock
	go generate ./...

.PHONY:
tidy:
	go mod tidy -compat=1.17
	cd tools && go mod tidy -compat=1.17

.PHONY:
all: ./bin/gowrap generate lint test

.PHONY:
release: ./bin/goreleaser
	goreleaser release

.PHONY:
build: ./bin/goreleaser
	goreleaser build --snapshot --rm-dist

.PHONY:
clean:
	rm -rf ./bin
