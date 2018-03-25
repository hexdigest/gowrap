package gowrap

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/hexdigest/gowrap/generator"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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
				cmd.BaseCommand.FlagSet().SetOutput(ioutil.Discard)
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
			args: []string{"-o", "pkg/out.file", "-i", "interface", "-t", "template/template", "-d", "pkg"},
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
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
