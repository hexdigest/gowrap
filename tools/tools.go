// +build tools

package main

// https://github.com/go-modules-by-example/index/blob/master/010_tools/README.md
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
)
