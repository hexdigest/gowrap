package gowrap

import (
	"bytes"
	"flag"
	"fmt"
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
	sourcePkg     string
	noGenerate    bool
	vars          vars

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
	fs.StringVar(&gc.interfaceName, "i", "", `the source interface name, i.e. "Reader"`)
	fs.StringVar(&gc.sourcePkg, "p", "", "the source package import path, i.e. \"io\", \"github.com/hexdigest/gowrap\" or\na relative import path like \"./generator\"")
	fs.StringVar(&gc.outputFile, "o", "", "the output file name")
	fs.StringVar(&gc.template, "t", "", "the template to use, it can be an HTTPS URL, local file or a\nreference to a template in gowrap repository,\n"+
		"run `gowrap template list` for details")
	fs.Var(&gc.vars, "v", "a key-value pair to parametrize the template,\narguments without an equal sign are treated as a bool values,\ni.e. -v foo=bar -v disableChecks")

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
	errNoOutputFile    = CommandLineError("output file is not specified")
	errNoInterfaceName = CommandLineError("interface name is not specified")
	errNoTemplate      = CommandLineError("no template specified")
	errNoSourcePackage = CommandLineError("no source package specified")
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

	return nil
}

func (gc *GenerateCommand) getOptions() (*generator.Options, error) {
	options := generator.Options{
		InterfaceName:  gc.interfaceName,
		OutputFile:     gc.outputFile,
		Funcs:          helperFuncs,
		HeaderTemplate: headerTemplate,
		HeaderVars: map[string]interface{}{
			"DisableGoGenerate": gc.noGenerate,
			"OutputFileName":    filepath.Base(gc.outputFile),
			"VarsArgs":          varsToArgs(gc.vars),
		},
		Vars: gc.vars.toMap(),
	}

	outputFileDir, err := gc.filepath.Abs(gc.filepath.Dir(gc.outputFile))
	if err != nil {
		return nil, err
	}

	if gc.sourcePkg == "" {
		gc.sourcePkg = "./"
	}

	sourcePackage, err := pkg.Load(gc.sourcePkg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load source package")
	}

	options.SourcePackage = sourcePackage.PkgPath
	options.BodyTemplate, options.HeaderVars["Template"], err = gc.loadTemplate(outputFileDir)

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

type varFlag struct {
	name  string
	value interface{}
}

//vars is a helper type that implements flag.Value to read multiple vars from the command line
type vars []varFlag

//String implements flag.Value
func (v vars) String() string {
	return fmt.Sprintf("%#v", v)
}

func (v *vars) Set(s string) error {
	chunks := strings.SplitN(s, "=", 2)
	switch len(chunks) {
	case 1:
		*v = append(*v, varFlag{name: chunks[0], value: true})
	case 2:
		*v = append(*v, varFlag{name: chunks[0], value: chunks[1]})
	}

	return nil
}

func (v vars) toMap() map[string]interface{} {
	m := make(map[string]interface{}, len(v))
	for _, vf := range v {
		m[vf.name] = vf.value
	}

	return m
}

func varsToArgs(v vars) string {
	if len(v) == 0 {
		return ""
	}

	var ss []string

	for _, vf := range v {
		switch typedValue := vf.value.(type) {
		case string:
			ss = append(ss, vf.name+"="+typedValue)
		case bool:
			ss = append(ss, vf.name)
		}
	}

	return " -v " + strings.Join(ss, " -v ")
}

var helperFuncs = template.FuncMap{
	"up":   strings.ToUpper,
	"down": strings.ToLower,
}

const headerTemplate = `package {{.Package.Name}}

// DO NOT EDIT!
// This code is generated with http://github.com/hexdigest/gowrap tool
// using {{.Options.HeaderVars.Template}} template

{{if (not .Options.HeaderVars.DisableGoGenerate)}}
//{{"go:generate"}} gowrap gen -p {{.SourcePackage.PkgPath}} -i {{.Options.InterfaceName}} -t {{.Options.HeaderVars.Template}} -o {{.Options.HeaderVars.OutputFileName}}{{.Options.HeaderVars.VarsArgs}}
{{end}}

`
