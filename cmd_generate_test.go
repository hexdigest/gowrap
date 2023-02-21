package gowrap

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/hexdigest/gowrap/generator"
)

func TestNewGenerateCommand(t *testing.T) {
	assert.NotNil(t, NewGenerateCommand(nil))
}

func Test_loader_Load(t *testing.T) {
	var unexpectedErr = errors.New("unexpected error")

	tests := []struct {
		name    string
		init    func(t minimock.Tester) loader
		inspect func(r loader, t *testing.T) //inspects loader after execution of Load

		template   string
		want1      []byte
		want2      string
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name: "file read error",
			init: func(t minimock.Tester) loader {
				return loader{
					fileReader: func(string) ([]byte, error) {
						return nil, unexpectedErr
					},
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, unexpectedErr, err)
			},
		},
		{
			name: "file read success",
			init: func(t minimock.Tester) loader {
				return loader{
					fileReader: func(string) ([]byte, error) {
						return []byte("success"), nil
					},
				}
			},
			template: "template.file",
			wantErr:  false,
			want1:    []byte("success"),
			want2:    "template.file",
		},
		{
			name: "remote template",
			init: func(t minimock.Tester) loader {
				return loader{
					fileReader: func(string) ([]byte, error) {
						return nil, os.ErrNotExist
					},
					remoteLoader: newRemoteTemplateLoaderMock(t).LoadMock.Expect("remote").Return([]byte("remote contents"), "remote-url", nil),
				}
			},
			template: "remote",
			wantErr:  false,
			want1:    []byte("remote contents"),
			want2:    "remote-url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			receiver := tt.init(mc)

			got1, got2, err := receiver.Load(tt.template)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "loader.Load returned unexpected result")

			assert.Equal(t, tt.want2, got2, "loader.Load returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestGenerateCommand_getOptions(t *testing.T) {
	var absError = errors.New("abs error")

	tests := []struct {
		name string
		init func(t minimock.Tester) *GenerateCommand

		want1      *generator.Options
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{

		{
			name: "unexisting output path",
			init: func(t minimock.Tester) *GenerateCommand {
				cmd := NewGenerateCommand(nil)
				cmd.filepath.Abs = func(string) (string, error) { return "", absError }
				return cmd
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, absError, errors.Cause(err))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			receiver := tt.init(mc)

			got1, err := receiver.getOptions()

			assert.Equal(t, tt.want1, got1, "GenerateCommand.parseArgs returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestGenerateCommand_Run(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *GenerateCommand
		inspect func(r *GenerateCommand, t *testing.T) //inspects *GenerateCommand after execution of Run

		args []string

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name: "parse args error",
			init: func(t minimock.Tester) *GenerateCommand {
				cmd := NewGenerateCommand(nil)
				cmd.BaseCommand.FlagSet().SetOutput(io.Discard)
				return cmd
			},
			args:    []string{"-pp"},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, "flag provided but not defined: -pp", err.Error())
			},
		},
		{
			name: "check flags error",
			init: func(t minimock.Tester) *GenerateCommand {
				return NewGenerateCommand(nil)
			},
			args:    []string{},
			wantErr: true,
		},
		{
			name: "get options error",
			init: func(t minimock.Tester) *GenerateCommand {
				return NewGenerateCommand(nil)
			},
			args:    []string{"-p", "unexistingpkg", "-i", "interface", "-o", "unexisting_dir/file.go", "-t", "unexisting_template"},
			wantErr: true,
		},
		{
			name: "package parse ok but template failed",
			init: func(t minimock.Tester) *GenerateCommand {
				loader := newRemoteTemplateLoaderMock(t).LoadMock.Expect("not_exists").Return(nil, "", errors.New("template load error"))
				return NewGenerateCommand(loader)
			},
			args:    []string{"-p", "io", "-i", "Writer", "-o", "file.go", "-t", "not_exists"},
			wantErr: true,
		},
		{
			name: "failed to create generator",
			args: []string{"-o", "pkg/out.file", "-i", "interface", "-t", "template/template"},
			init: func(t minimock.Tester) *GenerateCommand {
				loader := newRemoteTemplateLoaderMock(t).LoadMock.Return([]byte("{{."), "local/file", nil)
				return NewGenerateCommand(loader)
			},
			wantErr: true,
		},
		{
			name: "code generation error",
			args: []string{"-o", "out.file", "-i", "Command", "-t", "template/template"},
			init: func(t minimock.Tester) *GenerateCommand {
				loader := newRemoteTemplateLoaderMock(t).LoadMock.Return([]byte("body"), "local/file", nil)
				return NewGenerateCommand(loader)
			},
			wantErr: true,
		},
		{
			name: "success",
			args: []string{"-o", "out.file", "-i", "Command", "-t", "template/template"},
			init: func(t minimock.Tester) *GenerateCommand {
				cmd := NewGenerateCommand(nil)
				cmd.loader = newRemoteTemplateLoaderMock(t).LoadMock.Return([]byte("//comment"), "local/file", nil)
				cmd.filepath.WriteFile = func(string, []byte, os.FileMode) error { return nil }
				return cmd
			},
			wantErr: false,
		},
		{
			name: "success with local prefixes",
			args: []string{"-o", "out.file", "-i", "Command", "-t", "template/template", "-l", "foobar.com/pkg"},
			init: func(mt minimock.Tester) *GenerateCommand {
				cmd := NewGenerateCommand(nil)
				cmd.loader = newRemoteTemplateLoaderMock(mt).LoadMock.Return([]byte(`import (
	_ "foobar.com/pkg"
	_ "github.com/pkg/errors"
	_ "fmt"
)`), "local/file", nil)

				cmd.filepath.WriteFile = func(filename string, data []byte, perm os.FileMode) error {
					cmd.outputFile = filepath.Join(t.TempDir(), cmd.outputFile)
					return os.WriteFile(cmd.outputFile, data, perm)
				}
				return cmd
			},
			inspect: func(cmd *GenerateCommand, t *testing.T) {
				assert.EqualValues(t, "foobar.com/pkg", cmd.localPrefix)

				data, err := os.ReadFile(cmd.outputFile)
				assert.NoError(t, err)

				assert.Contains(t, string(data), `-l "foobar.com/pkg"`)
				assert.Contains(t, string(data), `import (
	_ "fmt"

	_ "github.com/pkg/errors"

	_ "foobar.com/pkg"
)`)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Finish()

			receiver := tt.init(mc)

			err := receiver.Run(tt.args, nil)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			if tt.wantErr {
				t.Logf("!!!\n\n%T: %v\n\n!!!", err, err)
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func Test_varsToArgs(t *testing.T) {
	tests := []struct {
		name  string
		v     vars
		want1 string
	}{
		{
			name:  "no vars",
			v:     nil,
			want1: "",
		},
		{
			name:  "two vars",
			v:     vars{varFlag{name: "key", value: "value"}, varFlag{name: "booleanKey", value: true}},
			want1: " -v key=value -v booleanKey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			got1 := varsToArgs(tt.v)

			assert.Equal(t, tt.want1, got1, "varsToArgs returned unexpected result")
		})
	}
}

func TestVars_toMap(t *testing.T) {
	tests := []struct {
		name  string
		vars  vars
		want1 map[string]interface{}
	}{
		{
			name: "success",
			vars: vars{{name: "key", value: "value"}, {name: "boolFlag", value: true}},
			want1: map[string]interface{}{
				"key":      "value",
				"boolFlag": true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got1 := tt.vars.toMap()

			assert.Equal(t, tt.want1, got1, "vars.toMap returned unexpected result")
		})
	}
}

func TestVars_Set(t *testing.T) {
	tests := []struct {
		name    string
		inspect func(r vars, t *testing.T) //inspects vars after execution of Set
		s       string
	}{
		{
			name: "bool var",
			s:    "boolVar",
			inspect: func(v vars, t *testing.T) {
				assert.Equal(t, vars{varFlag{name: "boolVar", value: true}}, v)
			},
		},

		{
			name: "string var",
			s:    "key=value",
			inspect: func(v vars, t *testing.T) {
				assert.Equal(t, vars{varFlag{name: "key", value: "value"}}, v)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := vars{}
			err := v.Set(tt.s)
			assert.NoError(t, err)

			tt.inspect(v, t)
		})
	}
}

func TestHelper_UpFirst(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "fist is lower-cased",
			in:   "typeName",
			out:  "TypeName",
		},
		{
			name: "single letter",
			in:   "v",
			out:  "V",
		},
		{
			name: "multi-bytes chars",
			in:   "йоу",
			out:  "Йоу",
		},
		{
			name: "empty string",
			in:   "",
			out:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := upFirst(tt.in)
			assert.Equal(t, tt.out, gotOut)
		})
	}
}

func TestHelper_DownFirst(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "fist is upper-cased",
			in:   "TypeName",
			out:  "typeName",
		},
		{
			name: "single letter",
			in:   "V",
			out:  "v",
		},
		{
			name: "multi-bytes chars",
			in:   "Йоу",
			out:  "йоу",
		},
		{
			name: "empty string",
			in:   "",
			out:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := downFirst(tt.in)
			assert.Equal(t, tt.out, gotOut)
		})
	}
}

func Test_toSnakeCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", ""},
		{"already_snake", "already_snake"},
		{"A", "a"},
		{"AA", "aa"},
		{"AaAa", "aa_aa"},
		{"HTTPRequest", "http_request"},
		{"BatteryLifeValue", "battery_life_value"},
		{"Id0Value", "id0_value"},
		{"ID0Value", "id0_value"},
	}

	for _, test := range tests {
		assert.Equal(t, test.want, toSnakeCase(test.input))
	}
}
