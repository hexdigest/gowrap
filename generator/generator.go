package generator

import (
	"bytes"
	"path/filepath"
	"strings"

	"go/ast"
	"go/token"
	"io"
	"text/template"

	"github.com/pkg/errors"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"

	"github.com/hexdigest/gowrap/pkg"
	"github.com/hexdigest/gowrap/printer"
)

//Generator generates decorators for the interface types
type Generator struct {
	Options

	headerTemplate *template.Template
	bodyTemplate   *template.Template
	srcPackage     *packages.Package
	dstPackage     *packages.Package
	methods        methodsList
	interfaceType  string
}

// TemplateInputs information passed to template for generation
type TemplateInputs struct {
	// Interface information for template
	Interface TemplateInputInterface
	// Vars additional vars to pass to the template, see Options.Vars
	Vars map[string]interface{}
}

// TemplateInputInterface subset of interface information used for template generation
type TemplateInputInterface struct {
	Name string
	// Type of the interface, with package name qualifier (e.g. sort.Interface)
	Type string
	// Methods name keyed map of method information
	Methods map[string]Method
}

type methodsList map[string]Method

//Options of the NewGenerator constructor
type Options struct {
	//InterfaceName is a name of interface type
	InterfaceName string

	//Imports from the file with interface definition
	Imports []string

	//SourcePackage is an import path or a relative path of the package that contains the source interface
	SourcePackage string

	//SourcePackageAlias is an import selector defauls is source package name
	SourcePackageAlias string

	//OutputFile name which is used to detect destination package name and also to fix imports in the resulting source
	OutputFile string

	//HeaderTemplate is used to generate package clause and comment over the generated source
	HeaderTemplate string

	//BodyTemplate generates import section, decorator constructor and methods
	BodyTemplate string

	//Vars additional vars that are passed to the templates from the command line
	Vars map[string]interface{}

	//HeaderVars header specific variables
	HeaderVars map[string]interface{}

	//Funcs is a map of helper functions that can be used within a template
	Funcs template.FuncMap
}

var errEmptyInterface = errors.New("interface has no methods")
var errUnexportedMethod = errors.New("unexported method")

//NewGenerator returns Generator initialized with options
func NewGenerator(options Options) (*Generator, error) {
	if options.Funcs == nil {
		options.Funcs = make(template.FuncMap)
	}

	headerTemplate, err := template.New("header").Funcs(options.Funcs).Parse(options.HeaderTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse header template")
	}

	bodyTemplate, err := template.New("body").Funcs(options.Funcs).Parse(options.BodyTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse body template")
	}

	if options.Vars == nil {
		options.Vars = make(map[string]interface{})
	}

	fs := token.NewFileSet()

	srcPackage, err := pkg.Load(options.SourcePackage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load source package")
	}

	dstPackagePath := filepath.Dir(options.OutputFile)
	if !strings.HasPrefix(dstPackagePath, "/") && !strings.HasPrefix(dstPackagePath, "./") {
		dstPackagePath = "./" + dstPackagePath
	}

	dstPackage, err := loadDestinationPackage(dstPackagePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load destination package: %s", dstPackagePath)
	}

	srcPackageAST, err := pkg.AST(fs, srcPackage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse source package")
	}

	interfaceType := srcPackage.Name + "." + options.InterfaceName
	if srcPackage.PkgPath == dstPackage.PkgPath {
		interfaceType = options.InterfaceName
		srcPackageAST.Name = ""
	} else if options.SourcePackageAlias != "" {
		srcPackageAST.Name = options.SourcePackageAlias
	}

	methods, imports, err := findInterface(fs, srcPackageAST, options.InterfaceName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse interface declaration")
	}

	if len(methods) == 0 {
		return nil, errEmptyInterface
	}

	for _, m := range methods {
		if srcPackageAST.Name != "" && []rune(m.Name)[0] == []rune(strings.ToLower(m.Name))[0] {
			return nil, errors.Wrap(errUnexportedMethod, m.Name)
		}
	}

	options.Imports = makeImports(imports)

	return &Generator{
		Options:        options,
		headerTemplate: headerTemplate,
		bodyTemplate:   bodyTemplate,
		srcPackage:     srcPackage,
		dstPackage:     dstPackage,
		interfaceType:  interfaceType,
		methods:        methods,
	}, nil
}

func makeImports(imports []*ast.ImportSpec) []string {
	result := make([]string, len(imports))
	for _, i := range imports {
		var name string
		if i.Name != nil {
			name = i.Name.Name
		}
		result = append(result, name+" "+i.Path.Value)
	}

	return result
}

func loadDestinationPackage(path string) (*packages.Package, error) {
	dstPackage, err := pkg.Load(path)
	if err != nil {
		//using directory name as a package name
		dstPackage, err = makePackage(path)
	}

	return dstPackage, err
}

var errNoPackageName = errors.New("failed to determine the destination package name")

func makePackage(path string) (*packages.Package, error) {
	name := filepath.Base(path)
	if name == string(filepath.Separator) || name == "." {
		return nil, errNoPackageName
	}

	return &packages.Package{
		Name: name,
	}, nil
}

//Generate generates code using header and body templates
func (g Generator) Generate(w io.Writer) error {
	buf := bytes.NewBuffer([]byte{})

	err := g.headerTemplate.Execute(buf, map[string]interface{}{
		"SourcePackage": g.srcPackage,
		"Package":       g.dstPackage,
		"Vars":          g.Options.Vars,
		"Options":       g.Options,
	})
	if err != nil {
		return err
	}

	err = g.bodyTemplate.Execute(buf, TemplateInputs{
		Interface: TemplateInputInterface{
			Name:    g.Options.InterfaceName,
			Type:    g.interfaceType,
			Methods: g.methods,
		},
		Vars: g.Options.Vars,
	})
	if err != nil {
		return err
	}

	processedSource, err := imports.Process(g.Options.OutputFile, buf.Bytes(), nil)
	if err != nil {
		return errors.Wrapf(err, "failed to format generated code:\n%s", buf)
	}

	_, err = w.Write(processedSource)
	return err
}

var errInterfaceNotFound = errors.New("interface type declaration not found")

// findInterface looks for the interface declaration in the given directory
// and returns a list of the interface's methods and a list of imports from the file
// where interface type declaration was found
func findInterface(fs *token.FileSet, p *ast.Package, interfaceName string) (methods methodsList, imports []*ast.ImportSpec, err error) {
	var found bool
	var types []*ast.TypeSpec
	var it *ast.InterfaceType

	//looking for the source interface declaration in all files in the dir
	//while doing this we also store all found type declarations to check if some of the
	//interface methods use unexported types
	for _, f := range p.Files {
		for _, ts := range typeSpecs(f) {
			types = append(types, ts)

			if i, ok := ts.Type.(*ast.InterfaceType); ok {
				if ts.Name.Name == interfaceName && !found {
					imports = f.Imports
					it = i
					found = true
				}
			}
		}
	}

	if !found {
		return nil, nil, errors.Wrap(errInterfaceNotFound, interfaceName)
	}

	methods, err = processInterface(fs, it, types, p.Name, imports)
	if err != nil {
		return nil, nil, err
	}

	return methods, imports, err
}

func typeSpecs(f *ast.File) []*ast.TypeSpec {
	result := []*ast.TypeSpec{}

	for _, decl := range f.Decls {
		if gd, ok := decl.(*ast.GenDecl); ok && gd.Tok == token.TYPE {
			for _, spec := range gd.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					result = append(result, ts)
				}
			}
		}
	}

	return result
}

func processInterface(fs *token.FileSet, it *ast.InterfaceType, types []*ast.TypeSpec, typesPrefix string, imports []*ast.ImportSpec) (methods methodsList, err error) {
	if it.Methods == nil {
		return nil, nil
	}

	methods = make(methodsList, len(it.Methods.List))

	for _, field := range it.Methods.List {
		var embeddedMethods methodsList
		var err error

		switch v := field.Type.(type) {
		case *ast.FuncType:
			var method *Method
			method, err = NewMethod(field.Names[0].Name, v, printer.New(fs, types, typesPrefix))
			if err == nil {
				methods[field.Names[0].Name] = *method
				continue
			}
		case *ast.SelectorExpr:
			embeddedMethods, err = processSelector(fs, v, imports)
		case *ast.Ident:
			embeddedMethods, err = processIdent(fs, v, types, typesPrefix, imports)
		}

		if err != nil {
			return nil, err
		}

		methods, err = mergeMethods(methods, embeddedMethods)
		if err != nil {
			return nil, err
		}
	}

	return methods, nil
}

func processSelector(fs *token.FileSet, se *ast.SelectorExpr, imports []*ast.ImportSpec) (methodsList, error) {
	interfaceName := se.Sel.Name
	packageSelector := se.X.(*ast.Ident).Name

	importPath, err := importPathByPrefix(imports, packageSelector)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to load %s.%s", packageSelector, interfaceName)
	}

	p, err := pkg.Load(importPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load imported package")
	}

	astPkg, err := pkg.AST(fs, p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to import package")
	}

	methods, _, err := findInterface(fs, astPkg, interfaceName)

	return methods, err
}

var errDuplicateMethod = errors.New("embedded interface has same method")

//mergeMethods merges two methods list, if there is a duplicate method name
//errDuplicateMethod is returned
func mergeMethods(ml1, ml2 methodsList) (methodsList, error) {
	if ml1 == nil || ml2 == nil {
		return ml1, nil
	}

	result := make(methodsList, len(ml1)+len(ml2))
	for k, v := range ml1 {
		result[k] = v
	}

	for name, signature := range ml2 {
		if _, ok := ml1[name]; ok {
			return nil, errors.Wrap(errDuplicateMethod, name)
		}

		result[name] = signature
	}

	return result, nil
}

var errEmbeddedInterfaceNotFound = errors.New("embedded interface not found")
var errNotAnInterface = errors.New("embedded type is not an interface")

func processIdent(fs *token.FileSet, i *ast.Ident, types []*ast.TypeSpec, typesPrefix string, imports []*ast.ImportSpec) (methodsList, error) {
	var embeddedInterface *ast.InterfaceType
	for _, t := range types {
		if t.Name.Name == i.Name {
			var ok bool
			embeddedInterface, ok = t.Type.(*ast.InterfaceType)
			if !ok {
				return nil, errors.Wrap(errNotAnInterface, t.Name.Name)
			}
			break
		}
	}

	if embeddedInterface == nil {
		return nil, errors.Wrap(errEmbeddedInterfaceNotFound, i.Name)
	}

	return processInterface(fs, embeddedInterface, types, typesPrefix, imports)
}

var errUnknownSelector = errors.New("unknown selector")

func importPathByPrefix(imports []*ast.ImportSpec, prefix string) (string, error) {
	for _, i := range imports {
		if i.Name != nil && i.Name.Name == prefix {
			return unquote(i.Path.Value), nil
		}
	}

	for _, i := range imports {
		p, err := pkg.Load(unquote(i.Path.Value))
		if err == nil && p.Name == prefix {
			return p.PkgPath, nil
		}
	}

	return "", errUnknownSelector
}

func unquote(s string) string {
	if s[0] == '"' {
		s = s[1:]
	}

	if s[len(s)-1] == '"' {
		s = s[0 : len(s)-1]
	}

	return s
}
