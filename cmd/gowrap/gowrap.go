package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/hexdigest/gowrap"
	"github.com/hexdigest/gowrap/loader"
)

func init() {
	ldr := loader.New(nil)

	gowrap.RegisterCommand("gen", gowrap.NewGenerateCommand(ldr))
	gowrap.RegisterCommand("template", gowrap.NewTemplateCommand(ldr))
}

func main() {
	if len(os.Args) < 2 {
		if err := gowrap.Usage(os.Stderr); err != nil {
			die(1, err.Error())
		}
		os.Exit(2)
	}

	flag.CommandLine.Usage = func() {
		die(2, "Run 'gowrap help' for usage.")
	}

	flag.Parse()
	args := flag.Args()

	if args[0] == "help" {
		if err := help(args[1:], os.Stdout); err != nil {
			die(2, err.Error())
		}
		return
	}

	command := gowrap.GetCommand(args[0])
	if command == nil {
		die(2, "gowrap: unknown subcommand %q\nRun 'gowrap help' for usage.", args[0])
	}

	if err := command.Run(args[1:], os.Stdout); err != nil {
		if _, ok := err.(gowrap.CommandLineError); ok {
			die(2, "%s\nRun 'gowrap help %s' for usage.\n", err.Error(), args[0])
		}
		die(1, err.Error())
	}
}

func die(exitCode int, format string, args ...interface{}) {
	if _, err := fmt.Fprintf(os.Stderr, format+"\n", args...); err != nil {
		os.Exit(1)
	}
	os.Exit(exitCode)
}

func help(args []string, w io.Writer) error {
	if len(args) > 1 {
		return errors.New("usage: gowrap help command\n\nToo many arguments given")
	}

	if len(args) == 0 {
		return gowrap.Usage(w)
	}

	command := gowrap.GetCommand(args[0])
	if command == nil {
		return fmt.Errorf(fmt.Sprintf("gounit: unknown subcommand %q\nRun 'gounit help' for usage", args[0]))
	}

	if _, err := fmt.Fprintf(w, "Usage: gowrap %s %s\n", args[0], command.UsageLine()); err != nil {
		return err
	}

	if fs := command.FlagSet(); fs != nil {
		fs.SetOutput(w)
		fs.PrintDefaults()
	}

	return command.HelpMessage(w)
}
