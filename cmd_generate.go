package gowrap

import (
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/hexdigest/gowrap/generator"
	"github.com/hexdigest/gowrap/pkg"
	"github.com/pkg/errors"
)

//GenerateCommand implements Command interface
type GenerateCommand struct {
	BaseCommand

	interfaceName string
	template      string
	outputFile    string
	sourceDir     string
	sourcePkg     string
	noGenerate    bool

	loader   templateLoader
	filepath fs
}

//NewGenerateCommand creates GenerateCommand
func NewGenerateCommand(l remoteTemplateLoader) *GenerateCommand {
	gc := &GenerateCommand{
		loader: loader{fileReader: ioutil.ReadFile, remoteLoader: l},
		filepath: fs{
			Rel:       filepath.Rel,
			Abs:       filepath.Abs,
			Dir:       filepath.Dir,
			WriteFile: ioutil.WriteFile,
		},
	}

	//this flagset loads flags values to the command fields
	fs := &flag.FlagSet{}
	fs.BoolVar(&gc.noGenerate, "g", false, "don't put //go:generate instruction to the generated code")
	fs.StringVar(&gc.interfaceName, "i", "", `interface name, i.e. "Reader"`)
	fs.StringVar(&gc.sourcePkg, "p", "", `source package import path, i.e. "io" or "github.com/hexdigest/gowrap"`)
	fs.StringVar(&gc.sourceDir, "d", "", "source package dir where to look for the interface declaration,\ndefault is a current directory")
	fs.StringVar(&gc.outputFile, "o", "", "output file name")
	fs.StringVar(&gc.template, "t", "", "template to use, it can be an HTTPS URL, local file or a\nreference to a template in gowrap repository,\n"+
		"run `gowrap template list` for details")

	gc.BaseCommand = BaseCommand{
		Short: "generate decorators",
		Usage: "-p package -i interfaceName -t template -o output_file.go",
		Flags: fs,
	}

	return gc
}

// Run implements Command interface
func (gc *GenerateCommand) Run(args []string, stdout io.Writer) error {
	if err := gc.FlagSet().Parse(args); err != nil {
		return CommandLineError(err.Error())
	}

	if err := gc.checkFlags(); err != nil {
		return err
	}

	generatorOptions, err := gc.getOptions()
	if err != nil {
		return err
	}

	gen, err := generator.NewGenerator(*generatorOptions)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer([]byte{})

	if err := gen.Generate(buf); err != nil {
		return err
	}

	return gc.filepath.WriteFile(gc.outputFile, buf.Bytes(), 0777)
}

var (
	errNoOutputFile                  = CommandLineError("output file is not specified")
	errNoInterfaceName               = CommandLineError("interface name is not specified")
	errNoTemplate                    = CommandLineError("no template specified")
	errEitherDirOrImportPathRequired = CommandLineError("either -d or -p option is expected")
)

func (gc *GenerateCommand) checkFlags() error {
	if gc.outputFile == "" {
		return errNoOutputFile
	}

	if gc.interfaceName == "" {
		return errNoInterfaceName
	}

	if gc.template == "" {
		return errNoTemplate
	}

	if gc.sourcePkg != "" && gc.sourceDir != "" {
		return errEitherDirOrImportPathRequired
	}
	return nil
}

func (gc *GenerateCommand) getOptions() (*generator.Options, error) {
	options := generator.Options{
		InterfaceName:  gc.interfaceName,
		OutputFile:     gc.outputFile,
		Funcs:          helperFuncs,
		HeaderTemplate: headerTemplate,
		Vars: map[string]interface{}{
			"OutputFile":        filepath.Base(gc.outputFile),
			"DisableGoGenerate": gc.noGenerate,
		},
	}

	outputFileDir, err := gc.filepath.Abs(gc.filepath.Dir(gc.outputFile))
	if err != nil {
		return nil, err
	}

	if gc.sourcePkg != "" {
		options.SourcePackageDir, err = pkg.Path(gc.sourcePkg)
		if err != nil {
			return nil, err
		}

		options.Vars["SourcePkg"] = gc.sourcePkg
	} else {
		if gc.sourceDir == "" {
			gc.sourceDir = "."
		}

		var sourceAbsPath string
		sourceAbsPath, err = gc.filepath.Abs(gc.sourceDir)
		if err != nil {
			return nil, err
		}

		options.Vars["SourceDir"], err = gc.filepath.Rel(outputFileDir, sourceAbsPath)
		if err != nil {
			return nil, err
		}

		options.SourcePackageDir = gc.sourceDir
	}

	options.BodyTemplate, options.Vars["Template"], err = gc.loadTemplate(outputFileDir)

	return &options, err
}

type readerFunc func(path string) ([]byte, error)

type loader struct {
	fileReader   readerFunc
	remoteLoader templateLoader
}

func (gc *GenerateCommand) loadTemplate(outputFileDir string) (contents, url string, err error) {
	body, url, err := gc.loader.Load(gc.template)
	if err != nil {
		return "", "", errors.Wrap(err, "failed to load template")
	}

	if !strings.HasPrefix(url, "https://") {
		templatePath, err := gc.filepath.Abs(url)
		if err != nil {
			return "", "", err
		}

		url, err = gc.filepath.Rel(outputFileDir, templatePath)
		if err != nil {
			return "", "", err
		}
	}

	return string(body), url, nil
}

// Load implements templateLoader
func (l loader) Load(template string) (tmpl []byte, url string, err error) {
	tmpl, err = l.fileReader(template)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}

		return l.remoteLoader.Load(template)
	}

	return tmpl, template, err
}

type templateLoader interface {
	Load(path string) (tmpl []byte, url string, err error)
}

type fs struct {
	Rel       func(string, string) (string, error)
	Abs       func(string) (string, error)
	Dir       func(string) string
	WriteFile func(string, []byte, os.FileMode) error
}

var helperFuncs = template.FuncMap{
	"up":   strings.ToUpper,
	"down": strings.ToLower,
}

const headerTemplate = `package {{.Package.Name}}

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using {{.Vars.Template}} template

{{if (not .Vars.DisableGoGenerate)}}
{{if .Vars.SourceDir}}
//{{"go:generate"}} gowrap gen -d {{.Vars.SourceDir}} -i {{.Options.InterfaceName}} -t {{.Vars.Template}} -o {{.Vars.OutputFile}}
{{else}}
//{{"go:generate"}} gowrap gen -p {{.Vars.SourcePkg}} -i {{.Options.InterfaceName}} -t {{.Vars.Template}} -o {{.Vars.OutputFile}}
{{end}}
{{end}}

`
