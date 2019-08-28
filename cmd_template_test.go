package gowrap

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	minimock "github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestNewTemplateCommand(t *testing.T) {
	cmd := NewTemplateCommand(nil)
	assert.NotNil(t, cmd)
}

func TestTemplateCommand_list(t *testing.T) {
	errUnexpected := errors.New("unexpected error")
	type args struct {
		w io.Writer
	}
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *TemplateCommand
		inspect func(r *TemplateCommand, t *testing.T) //inspects *TemplateCommand after execution of list

		args args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name: "loader error",
			init: func(t minimock.Tester) *TemplateCommand {
				return &TemplateCommand{
					loader: newRemoteTemplateLoaderMock(t).ListMock.Return(nil, errUnexpected),
				}
			},
			wantErr: true,
		},
		{
			name: "no remote templates",
			init: func(t minimock.Tester) *TemplateCommand {
				return &TemplateCommand{
					loader: newRemoteTemplateLoaderMock(t).ListMock.Return(nil, nil),
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errNoTemplatesFound, err)
			},
		},
		{
			name: "success",
			init: func(t minimock.Tester) *TemplateCommand {
				return &TemplateCommand{
					loader: newRemoteTemplateLoaderMock(t).ListMock.Return([]string{"template"}, nil),
				}
			},
			args: args{
				w: bytes.NewBuffer([]byte{}),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			err := receiver.list(tt.args.w)

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

func TestTemplateCommand_fetch(t *testing.T) {
	errUnexpected := errors.New("unexpected error")
	type args struct {
		w    io.Writer
		wf   writeFileFunc
		args []string
	}
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *TemplateCommand
		inspect func(r *TemplateCommand, t *testing.T) //inspects *TemplateCommand after execution of fetch
		args    func(t minimock.Tester) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name: "invalid number of arguments",
			args: func(minimock.Tester) args {
				return args{}
			},
			wantErr: true,
			init: func(t minimock.Tester) *TemplateCommand {
				return &TemplateCommand{}
			},
			inspectErr: func(err error, t *testing.T) {
				assert.Contains(t, err.Error(), "expected template and a local file name")
			},
		},
		{
			name: "loader error",
			args: func(t minimock.Tester) args {
				return args{
					w:    nil,
					args: []string{"tmpl", "out.file"},
				}
			},
			wantErr: true,
			init: func(t minimock.Tester) *TemplateCommand {
				return &TemplateCommand{
					loader: newRemoteTemplateLoaderMock(t).LoadMock.Expect("tmpl").Return([]byte(""), "template", errUnexpected),
				}
			},
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexpected, err)
			},
		},
		{
			name: "write file error",
			args: func(t minimock.Tester) args {
				return args{
					w: NewWriterMock(t).WriteMock.Return(0, nil),
					wf: func(filename string, data []byte, perm os.FileMode) error {
						return errUnexpected
					},
					args: []string{"tmpl", "out.file"},
				}
			},
			wantErr: true,
			init: func(t minimock.Tester) *TemplateCommand {
				return &TemplateCommand{
					loader: newRemoteTemplateLoaderMock(t).LoadMock.Expect("tmpl").Return(nil, "", nil),
				}
			},
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexpected, err)
			},
		},
		{
			name: "success",
			args: func(t minimock.Tester) args {
				return args{
					w: NewWriterMock(t).WriteMock.Return(0, nil),
					wf: func(filename string, data []byte, perm os.FileMode) error {
						return nil
					},
					args: []string{"tmpl", "out.file"},
				}
			},
			init: func(t minimock.Tester) *TemplateCommand {
				return &TemplateCommand{
					loader: newRemoteTemplateLoaderMock(t).LoadMock.Expect("tmpl").Return(nil, "", nil),
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			tArgs := tt.args(t)

			err := receiver.fetch(tArgs.w, tArgs.wf, tArgs.args)

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

func TestTemplateCommand_Run(t *testing.T) {
	errUnexpected := errors.New("unexpected error")
	type args struct {
		args []string
		w    io.Writer
	}
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *TemplateCommand
		inspect func(r *TemplateCommand, t *testing.T) //inspects *TemplateCommand after execution of Run

		args func(t minimock.Tester) args

		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name: "invalid arguments",
			init: func(t minimock.Tester) *TemplateCommand {
				return &TemplateCommand{}
			},
			args: func(t minimock.Tester) args {
				return args{w: nil, args: nil}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errExpectedSubcommand, err)
			},
		},
		{
			name: "unknown subcommand",
			init: func(t minimock.Tester) *TemplateCommand {
				return &TemplateCommand{}
			},
			args: func(t minimock.Tester) args {
				return args{w: nil, args: []string{"unknown"}}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnknownSubcommand, err)
			},
		},
		{
			name: "list subcommand",
			init: func(t minimock.Tester) *TemplateCommand {
				return &TemplateCommand{
					loader: newRemoteTemplateLoaderMock(t).ListMock.Return(nil, errUnexpected),
				}
			},
			args: func(t minimock.Tester) args {
				return args{w: nil, args: []string{"list"}}
			},
			wantErr: true,
		},
		{
			name: "copy subcommand",
			init: func(t minimock.Tester) *TemplateCommand {
				return &TemplateCommand{}
			},
			args: func(t minimock.Tester) args {
				return args{w: nil, args: []string{"copy"}}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			tArgs := tt.args(mc)
			receiver := tt.init(mc)

			err := receiver.Run(tArgs.args, tArgs.w)

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
