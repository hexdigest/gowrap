lint:
	golint ./... && go vet ./...

test:
	go generate ./... && go test -race ./...

install:
	go install ./cmd/gowrap

all: lint test install
