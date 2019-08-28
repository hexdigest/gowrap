package gowrap

import (
	"flag"
	"io"
	"testing"
	"time"

	minimock "github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestBaseCommand_ShortDescription(t *testing.T) {
	base := BaseCommand{Short: "short"}
	assert.Equal(t, "short", base.ShortDescription())
}

func TestBaseCommand_UsageLine(t *testing.T) {
	base := BaseCommand{Usage: "usage line"}
	assert.Equal(t, "usage line", base.UsageLine())
}

func TestBaseCommand_FlagSet(t *testing.T) {
	base := BaseCommand{}
	assert.Nil(t, base.FlagSet())
}

func TestBaseCommand_HelpMessage(t *testing.T) {
	type args struct {
		w io.Writer
	}
	tests := []struct {
		name string
		init func(t minimock.Tester) BaseCommand

		args func(t minimock.Tester) args

		wantErr bool
	}{
		{
			name: "success",
			init: func(t minimock.Tester) BaseCommand {
				return BaseCommand{Help: "help"}
			},
			args: func(t minimock.Tester) args {
				return args{w: NewWriterMock(t).WriteMock.Expect([]byte("help")).Return(0, nil)}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			tArgs := tt.args(mc)
			receiver := tt.init(mc)

			err := receiver.HelpMessage(tArgs.w)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestRegisterCommand(t *testing.T) {
	cmd := &GenerateCommand{BaseCommand: BaseCommand{Flags: flag.NewFlagSet("flagset", flag.ContinueOnError)}}
	RegisterCommand("TestRegisterCommand", cmd)
	assert.NotNil(t, commands["TestRegisterCommand"])
}

func TestGetCommand(t *testing.T) {
	commands["TestGetCommand"] = &GenerateCommand{}
	assert.NotNil(t, GetCommand("TestGetCommand"))
}

func TestUsage(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	w := NewWriterMock(mc).WriteMock.Return(0, nil)
	assert.NoError(t, Usage(w))
}
