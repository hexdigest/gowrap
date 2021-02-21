export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)

./bin:
	mkdir ./bin

./bin/gowrap: ./bin
	go install ./cmd/gowrap

./bin/golangci-lint: ./bin
	go install github.com/golangci/golangci-lint/cmd/golangci-lint

lint: ./bin/golangci-lint
	./bin/golangci-lint run --enable=goimports --disable=unused --exclude=S1023,"Error return value" ./...

test:
	 go test -race ./...

generate: ./bin/gowrap
	go generate ./...

all: ./bin/gowrap generate lint test
