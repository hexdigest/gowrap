all: test lint

lint:
	golint ./... && go vet ./...

test:
	go generate ./... && go test -race ./...
