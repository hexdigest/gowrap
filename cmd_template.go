package gowrap

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type writeFileFunc func(filename string, data []byte, perm os.FileMode) error

//remoteTemplateLoader returns loads template by an URL or a reference to the github repo
type remoteTemplateLoader interface {
	List() ([]string, error)
	Load(path string) (tmpl []byte, url string, err error)
}

//NewTemplateCommand creates TemplateCommand
func NewTemplateCommand(loader remoteTemplateLoader) *TemplateCommand {
	return &TemplateCommand{
		BaseCommand: BaseCommand{
			Short: "manage decorators templates",
			Usage: "subcommand [options]",
			Help: `
Subcommands are:

  list - list all template in gowrap repository

  copy - copy remote template to a local file, i.e.

    gowrap template copy fallback templates/fallback
`,
		},
		loader: loader,
	}
}

//TemplateCommand implements Command interface
type TemplateCommand struct {
	BaseCommand
	loader remoteTemplateLoader
}

var errExpectedSubcommand = CommandLineError("expected subcommand")
var errUnknownSubcommand = CommandLineError("unknown subcommand")

// Run implements Command interface
func (gc *TemplateCommand) Run(args []string, w io.Writer) error {
	if len(args) == 0 {
		return errExpectedSubcommand
	}

	subcommand := args[0]
	switch subcommand {
	case "list":
		return gc.list(w)
	case "copy":
		return gc.fetch(w, ioutil.WriteFile, args[1:])
	}
	return errUnknownSubcommand
}

var errNoTemplatesFound = errors.New("no remote templates found")

func (gc *TemplateCommand) list(w io.Writer) error {
	templates, err := gc.loader.List()
	if err != nil {
		return err
	}

	if len(templates) == 0 {
		return errNoTemplatesFound
	}

	fmt.Fprintln(w, "List of available remote templates:")
	for _, t := range templates {
		fmt.Fprintf(w, "  %s\n", t)
	}

	return nil
}

func (gc *TemplateCommand) fetch(w io.Writer, wf writeFileFunc, args []string) error {
	if len(args) < 2 {
		return CommandLineError("expected template and a local file name")
	}

	template, dstFileName := args[0], args[1]

	body, url, err := gc.loader.Load(template)
	if err != nil {
		return err
	}

	if err := wf(dstFileName, body, 0777); err != nil {
		return err
	}

	fmt.Fprintf(w, "successfully copied from %s\n", url)
	return nil
}
