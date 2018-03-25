all: test lint

lint:
	golint ./... && go vet ./...

test:
	go test -race ./...

install:
	go install ./cmd/gowrap

all: lint test install
