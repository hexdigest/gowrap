package pkg

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// Package contains package meta information
type Package struct {
	Name       string
	ImportPath string
	Path       string
}

type filterFunc func(fi os.FileInfo) bool

var abs = filepath.Abs
var parseDir = parser.ParseDir

// FromImport returns package with the given import path
func FromImport(fs *token.FileSet, importPath string) (*ast.Package, error) {
	p, err := build.Default.Import(importPath, "", build.FindOnly)
	if err != nil {
		return nil, err
	}

	dir, err := abs(filepath.Join(p.SrcRoot, p.ImportPath))
	if err != nil {
		return nil, err
	}

	return FromDir(fs, dir, NoTests)
}

// Name returns name of the package in given import path
func Name(importPath string) (string, error) {
	p, err := build.Default.Import(importPath, "./", build.ImportComment)
	if err != nil {
		return "", err
	}

	return p.Name, nil
}

//Path returns package absolute path by given import path
func Path(importPath string) (string, error) {
	p, err := build.Default.Import(importPath, "", build.ImportComment)
	if err != nil {
		return "", err
	}

	return abs(filepath.Join(p.SrcRoot, p.ImportPath))
}

var errNotFound = errors.New("package not found")

// FromDir returns package from the dir.
// If dir does not contains any files that pass filter lowercased dir name is used as a package name
// If filter is nil then NoTests filter is used
func FromDir(fs *token.FileSet, dir string, filter filterFunc) (*ast.Package, error) {
	p, err := build.Default.ImportDir(dir, build.ImportComment)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		filter = NoTests
	}

	pkgs, err := parseDir(fs, dir, filter, parser.DeclarationErrors | parser.ParseComments)
	if err != nil {
		return nil, err
	}

	astPkg, ok := pkgs[p.Name]
	if !ok {
		return nil, errNotFound
	}

	return astPkg, nil
}

// NoTests returns true for all files with .go suffix except _test.go
func NoTests(fi os.FileInfo) bool {
	return strings.HasSuffix(fi.Name(), ".go") && !strings.HasSuffix(fi.Name(), "_test.go")
}
