package pkg

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestFromImport(t *testing.T) {
	fs := token.NewFileSet()

	t.Run("build import error", func(t *testing.T) {
		res, err := FromImport(fs, "fmt/unexisting-import-path")
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("filepath.Abs error", func(t *testing.T) {
		absError := errors.New("abs error")
		abs = func(string) (string, error) {
			return "", absError
		}
		res, err := FromImport(fs, "io")
		assert.Error(t, err)
		assert.Equal(t, absError, err)
		assert.Nil(t, res)
	})

	t.Run("success", func(t *testing.T) {
		abs = filepath.Abs
		res, err := FromImport(fs, "io")
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}

func TestName(t *testing.T) {
	t.Run("build import error", func(t *testing.T) {
		res, err := Name("fmt/unexisting-import-path")
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success", func(t *testing.T) {
		res, err := Name("io")
		assert.NoError(t, err)
		assert.Equal(t, "io", res)
	})
}

func TestPath(t *testing.T) {
	t.Run("build import error", func(t *testing.T) {
		res, err := Path("fmt/unexisting-import-path")
		assert.Error(t, err)
		assert.Empty(t, res)
	})

	t.Run("success", func(t *testing.T) {
		res, err := Path("io")
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
	})
}

func TestFromDir(t *testing.T) {
	fs := token.NewFileSet()

	t.Run("build import dir error", func(t *testing.T) {
		res, err := FromDir(fs, "unexisting-dir", nil)
		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("parseDir error", func(t *testing.T) {
		parseError := errors.New("parse dir error")
		parseDir = func(fset *token.FileSet, path string, filter func(os.FileInfo) bool, mode parser.Mode) (map[string]*ast.Package, error) {
			return nil, parseError
		}

		res, err := FromDir(fs, ".", nil)
		assert.Error(t, err)
		assert.Equal(t, parseError, err)
		assert.Nil(t, res)
	})

	t.Run("package not found", func(t *testing.T) {
		parseDir = func(fset *token.FileSet, path string, filter func(os.FileInfo) bool, mode parser.Mode) (map[string]*ast.Package, error) {
			return make(map[string]*ast.Package), nil
		}

		res, err := FromDir(fs, ".", nil)
		assert.Error(t, err)
		assert.Equal(t, errNotFound, err)
		assert.Nil(t, res)
	})
}

func TestNoTests(t *testing.T) {
	tests := []struct {
		name string
		fi   func(t minimock.Tester) os.FileInfo

		want1 bool
	}{
		{
			name: "no a Go file",
			fi: func(t minimock.Tester) os.FileInfo {
				return NewFileInfoMock(t).NameMock.Return("not.a.go.file")
			},
			want1: false,
		},
		{
			name: "test file",
			fi: func(t minimock.Tester) os.FileInfo {
				return NewFileInfoMock(t).NameMock.Return("go_test.go")
			},
			want1: false,
		},
		{
			name: "regular go file",
			fi: func(t minimock.Tester) os.FileInfo {
				return NewFileInfoMock(t).NameMock.Return("go.go")
			},
			want1: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			got1 := NoTests(tt.fi(mc))

			assert.Equal(t, tt.want1, got1, "NoTests returned unexpected result")
		})
	}
}
