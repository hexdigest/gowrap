package loader

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_gitRootPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "last element should be gowrap",
			want: "gowrap",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gitRootPath()
			require.NoError(t, err)

			assert.Equal(t, tt.want, path.Base(got))
		})
	}
}

func mockGitRootPath(mockPath string) func() (string, error) {
	return func() (string, error) {
		return mockPath, nil
	}
}
