package gowrap

import (
	"flag"
	"io"
	"io/ioutil"
	"text/template"
)

// Command interface represents gowrap subcommand
type Command interface {
	//FlagSet returns command specific flag set. If command doesn't have any flags nil should be returned.
	FlagSet() *flag.FlagSet

	//Run runs command
	Run(args []string, stdout io.Writer) error

	//Description returns short description of a command that is shown in the help message
	ShortDescription() string

	UsageLine() string
	HelpMessage(w io.Writer) error
}

// BaseCommand implements Command interface
type BaseCommand struct {
	Flags *flag.FlagSet
	Short string
	Usage string
	Help  string
}

//ShortDescription implements Command
func (b BaseCommand) ShortDescription() string {
	return b.Short
}

//UsageLine implements Command
func (b BaseCommand) UsageLine() string {
	return b.Usage
}

//FlagSet implements Command
func (b BaseCommand) FlagSet() *flag.FlagSet {
	return b.Flags
}

//HelpMessage implements Command
func (b BaseCommand) HelpMessage(w io.Writer) error {
	_, err := w.Write([]byte(b.Help))
	return err
}

var commands = map[string]Command{}

// RegisterCommand adds command to the global Commands map
func RegisterCommand(name string, cmd Command) {
	commands[name] = cmd
	if fs := cmd.FlagSet(); fs != nil {
		fs.Init("", flag.ContinueOnError)
		fs.SetOutput(ioutil.Discard)
	}
}

// GetCommand returns command from the global Commands map
func GetCommand(name string) Command {
	return commands[name]
}

// Usage writes gowrap usage message to w
func Usage(w io.Writer) error {
	return usageTemplate.Execute(w, commands)
}

var usageTemplate = template.Must(template.New("usage").Parse(`GoWrap is a tool for generating decorators for the Go interfaces

Usage:

	gowrap command [arguments]

The commands are:
{{ range $name, $cmd := . }}
	{{ printf "%-10s" $name }}{{ $cmd.ShortDescription }}
{{ end }}
Use "gowrap help [command]" for more information about a command.
`))
