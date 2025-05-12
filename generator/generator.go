package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"

	"github.com/hexdigest/gowrap/pkg"
	"github.com/hexdigest/gowrap/printer"
)

// Generator generates decorators for the interface types
type Generator struct {
	Options

	headerTemplate *template.Template
	bodyTemplate   *template.Template
	srcPackage     *packages.Package
	dstPackage     *packages.Package
	methods        methodsList
	interfaceType  string
	genericTypes   string
	genericParams  string
	localPrefix    string
}

// TemplateInputs information passed to template for generation
type TemplateInputs struct {
	// Interface information for template
	Interface TemplateInputInterface
	// Vars additional vars to pass to the template, see Options.Vars
	Vars    map[string]interface{}
	Imports []string
}

// Import generates an import statement using a list of imports from the source file
// along with the ones from the template itself
func (t TemplateInputs) Import(imports ...string) string {
	allImports := make(map[string]struct{}, len(imports)+len(t.Imports))

	for _, i := range t.Imports {
		allImports[strings.TrimSpace(i)] = struct{}{}
	}

	for _, i := range imports {
		if len(i) == 0 {
			continue
		}

		i = strings.TrimSpace(i)

		if i[len(i)-1] != '"' {
			i += `"`
		}

		if i[0] != '"' {
			i = `"` + i
		}

		allImports[i] = struct{}{}
	}

	out := make([]string, 0, len(allImports))

	for i := range allImports {
		out = append(out, i)
	}

	sort.Strings(out)

	return "import (\n" + strings.Join(out, "\n") + ")\n"
}

// TemplateInputInterface subset of interface information used for template generation
type TemplateInputInterface struct {
	Name string
	// Type of the interface, with package name qualifier (e.g. sort.Interface)
	Type string
	// Generics of the interface when using generics
	Generics TemplateInputGenerics
	// Methods name keyed map of method information
	Methods map[string]Method
}

// Options of the NewGenerator constructor
type Options struct {
	//InterfaceName is a name of interface type
	InterfaceName string

	//Imports from the file with interface definition
	Imports []string

	//SourcePackage is an import path or a relative path of the package that contains the source interface
	SourcePackage string

	//SourcePackageInstance is the already loaded package (optional)
	SourcePackageInstance *packages.Package

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

	//LocalPrefix is a comma-separated string of import path prefixes, which, if set, instructs Process to sort the import
	//paths with the given prefixes into another group after 3rd-party packages.
	LocalPrefix string

	//IgnoreUnexported skip generation of unexported methods instead of return an error
	IgnoreUnexported bool
}

type methodsList map[string]Method

type processInput struct {
	fileSet        *token.FileSet
	currentPackage *packages.Package
	astPackage     *ast.Package
	targetName     string
	genericParams  genericParams
}

type targetProcessInput struct {
	processInput
	types        []*ast.TypeSpec
	typesPrefix  string
	imports      []*ast.ImportSpec
	genericTypes genericTypes
}

type processOutput struct {
	genericTypes genericTypes
	methods      methodsList
	imports      []*ast.ImportSpec
}

var errEmptyInterface = errors.New("interface has no methods")
var errUnexportedMethod = errors.New("unexported method")

// NewGenerator returns Generator initialized with options
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

	var srcPackage *packages.Package
	// Use the preloaded package if available, only load if not
	if options.SourcePackageInstance != nil {
		srcPackage = options.SourcePackageInstance
	} else {
		srcPackage, err = pkg.Load(options.SourcePackage)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load source package")
		}
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

	var interfaceType string
	if srcPackage.PkgPath == dstPackage.PkgPath {
		interfaceType = options.InterfaceName
		srcPackageAST.Name = ""
	} else {
		if options.SourcePackageAlias != "" {
			srcPackageAST.Name = options.SourcePackageAlias
		} else {
			srcPackageAST.Name = "_source" + cases.Title(language.Und, cases.NoLower).String(srcPackageAST.Name)
		}

		interfaceType = srcPackageAST.Name + "." + options.InterfaceName
		options.Imports = append(options.Imports, srcPackageAST.Name+` "`+srcPackage.PkgPath+`"`)
	}

	output, err := findTarget(processInput{
		fileSet:        fs,
		currentPackage: srcPackage,
		astPackage:     srcPackageAST,
		targetName:     options.InterfaceName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse interface declaration")
	}

	if len(output.methods) == 0 {
		return nil, errEmptyInterface
	}

	for n, m := range output.methods {
		unexported := srcPackageAST.Name != "" && []rune(m.Name)[0] == []rune(strings.ToLower(m.Name))[0]
		if !options.IgnoreUnexported && unexported {
			return nil, errors.Wrap(errUnexportedMethod, m.Name)
		}

		if unexported {
			delete(output.methods, n)
		}
	}

	options.Imports = append(options.Imports, makeImports(output.imports)...)

	genericTypes, genericParams := output.genericTypes.buildVars()

	return &Generator{
		Options:        options,
		headerTemplate: headerTemplate,
		bodyTemplate:   bodyTemplate,
		srcPackage:     srcPackage,
		dstPackage:     dstPackage,
		interfaceType:  interfaceType,
		genericTypes:   genericTypes,
		genericParams:  genericParams,
		methods:        output.methods,
		localPrefix:    options.LocalPrefix,
	}, nil
}

func makeImports(imports []*ast.ImportSpec) []string {
	result := make([]string, len(imports))
	for i, im := range imports {
		var name string
		if im.Name != nil {
			name = im.Name.Name
		}
		result[i] = name + " " + im.Path.Value
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

// Generate generates code using header and body templates
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
			Name: g.Options.InterfaceName,
			Generics: TemplateInputGenerics{
				Types:  g.genericTypes,
				Params: g.genericParams,
			},
			Type:    g.interfaceType,
			Methods: g.methods,
		},
		Imports: g.Options.Imports,
		Vars:    g.Options.Vars,
	})
	if err != nil {
		return err
	}

	imports.LocalPrefix = g.localPrefix
	processedSource, err := imports.Process(g.Options.OutputFile, buf.Bytes(), nil)
	if err != nil {
		return errors.Wrapf(err, "failed to format generated code:\n%s", buf)
	}

	_, err = w.Write(processedSource)
	return err
}

var errTargetNotFound = errors.New("target declaration not found")

func findTarget(input processInput) (output processOutput, err error) {
	ts, imports, types := iterateFiles(input.astPackage, input.targetName)
	if ts == nil {
		return processOutput{}, errors.Wrap(errTargetNotFound, input.targetName)
	}

	output.imports = imports
	output.genericTypes = buildGenericTypesFromSpec(ts, types, input.astPackage.Name)
	output.methods, err = findMethods(ts, targetProcessInput{
		processInput: input,
		types:        types,
		typesPrefix:  input.astPackage.Name,
		imports:      output.imports,
		genericTypes: output.genericTypes,
	})
	if err != nil {
		return processOutput{}, err
	}

	return
}

func findMethods(selectedType *ast.TypeSpec, input targetProcessInput) (methods methodsList, err error) {
	switch t := selectedType.Type.(type) {
	case *ast.InterfaceType:
		methods, err = processInterface(t, input)
		if err != nil {
			return methodsList{}, err
		}
	case *ast.SelectorExpr:
		ident, ok := t.X.(*ast.Ident)
		if !ok {
			return
		}
		srcPackagePath := findSourcePackage(ident, input.imports)

		return getMethods(t, srcPackagePath)
	case *ast.Ident:
		if t.Obj == nil {
			return
		}
		if ts, ok := t.Obj.Decl.(*ast.TypeSpec); ok {
			return findMethods(ts, input)
		}
	}

	return
}

func getMethods(sel *ast.SelectorExpr, srcPackagePath string) (methods methodsList, err error) {
	srcPkg, err := pkg.Load(srcPackagePath)
	if err != nil {
		return nil, errors.Wrapf(err, "cant load %s package", srcPackagePath)
	}

	fs := token.NewFileSet()
	srcAst, err := pkg.AST(fs, srcPkg)
	if err != nil {
		return nil, errors.Wrapf(err, "cant ast %s package", srcPackagePath)
	}

	out, err := findTarget(processInput{
		fileSet:        fs,
		currentPackage: srcPkg,
		astPackage:     srcAst,
		targetName:     sel.Sel.Name,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find target in %s package", srcPackagePath)
	}

	return out.methods, nil
}

func findSourcePackage(ident *ast.Ident, imports []*ast.ImportSpec) string {
	for _, imp := range imports {
		cleanPath := strings.Trim(imp.Path.Value, "\"")
		if imp.Name != nil {
			if ident.Name == imp.Name.Name {
				return cleanPath
			}

			continue
		}

		slash := strings.LastIndex(cleanPath, "/")
		if ident.Name == cleanPath[slash+1:] {
			return cleanPath
		}
	}

	return ""
}

func iterateFiles(p *ast.Package, name string) (selectedType *ast.TypeSpec, imports []*ast.ImportSpec, types []*ast.TypeSpec) {
	for _, f := range p.Files {
		if f != nil {
			for _, ts := range typeSpecs(f) {
				types = append(types, ts)
				if ts.Name.Name == name {
					selectedType = ts
					imports = f.Imports
					continue
				}
			}
		}
	}
	return
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

func getEmbeddedMethods(t ast.Expr, pr typePrinter, input targetProcessInput, checkInterface bool) (param genericParam, methods methodsList, err error) {
	param.Name, err = pr.PrintType(t)
	if err != nil {
		return
	}

	switch v := t.(type) {
	case *ast.SelectorExpr:
		methods, err = processSelector(v, input)
		return

	case *ast.Ident:
		methods, err = processIdent(v, input, checkInterface)
		return
	}
	return
}

func processEmbedded(t ast.Expr, pr typePrinter, input targetProcessInput, checkInterface bool) (genericParam genericParam, embeddedMethods methodsList, err error) {
	var x ast.Expr
	var hasGenericsParams bool
	var genericParams genericParams

	switch v := t.(type) {
	case *ast.IndexExpr:
		x = v.X
		hasGenericsParams = true
		//	Don't check if embedded interface's generic params are also interfaces, e.g. given the interface:
		//		type SomeInterface {
		//	      EmbeddedGenericInterface[Bar]
		//		}
		//	we won't be checking if Bar is also an interface
		genericParam, _, err = processEmbedded(v.Index, pr, input, false)
		if err != nil {
			return
		}
		if genericParam.Name != "" {
			genericParams = append(genericParams, genericParam)
		}

	case *ast.IndexListExpr:
		x = v.X
		hasGenericsParams = true

		if v.Indices != nil {
			for _, index := range v.Indices {
				//	Don't check if embedded interface's generic params are also interfaces, e.g. given the interface:
				//		type SomeInterface {
				//	      EmbeddedGenericInterface[Bar]
				//		}
				//	we won't be checking if Bar is also an interface
				genericParam, _, err = processEmbedded(index, pr, input, false)
				if err != nil {
					return
				}
				if genericParam.Name != "" {
					genericParams = append(genericParams, genericParam)
				}
			}
		}
	default:
		x = v
	}

	input.genericParams = genericParams
	genericParam, embeddedMethods, err = getEmbeddedMethods(x, pr, input, checkInterface)
	if err != nil {
		return
	}

	if hasGenericsParams {
		genericParam.Params = genericParams
	}

	return
}

func processInterface(it *ast.InterfaceType, targetInput targetProcessInput) (methods methodsList, err error) {
	if it.Methods == nil {
		return nil, nil
	}

	methods = make(methodsList, len(it.Methods.List))

	pr := printer.New(targetInput.fileSet, targetInput.types, targetInput.typesPrefix)

	for _, field := range it.Methods.List {
		var embeddedMethods methodsList
		var err error

		switch v := field.Type.(type) {
		case *ast.FuncType:
			var method *Method

			method, err = NewMethod(field.Names[0].Name, field, pr, targetInput.genericTypes, targetInput.genericParams)
			if err == nil {
				methods[field.Names[0].Name] = *method
				continue
			}

		default:
			_, embeddedMethods, err = processEmbedded(v, pr, targetInput, true)
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

func processSelector(se *ast.SelectorExpr, input targetProcessInput) (methodsList, error) {
	selectedName := se.Sel.Name
	packageSelector := se.X.(*ast.Ident).Name

	importPath, err := findImportPathForName(packageSelector, input.imports, input.currentPackage)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to find package %s", packageSelector)
	}

	p, ok := input.currentPackage.Imports[importPath]
	if !ok {
		return nil, fmt.Errorf("unable to find package %s", packageSelector)
	}

	astPkg, err := pkg.AST(input.fileSet, p)
	if err != nil {
		return nil, errors.Wrap(err, "failed to import package")
	}

	output, err := findTarget(processInput{
		fileSet:        input.fileSet,
		currentPackage: p,
		astPackage:     astPkg,
		targetName:     selectedName,
		genericParams:  input.genericParams,
	})

	return output.methods, err
}

// mergeMethods merges two methods list. Retains overlapping methods from the
// parent list
func mergeMethods(methods, embeddedMethods methodsList) (methodsList, error) {
	if methods == nil || embeddedMethods == nil {
		return methods, nil
	}

	result := make(methodsList, len(methods)+len(embeddedMethods))
	for name, signature := range embeddedMethods {
		result[name] = signature
	}

	for name, signature := range methods {
		result[name] = signature
	}

	return result, nil
}

var errNotAnInterface = errors.New("embedded type is not an interface")

func processIdent(i *ast.Ident, input targetProcessInput, checkInterface bool) (methodsList, error) {
	var embeddedInterface *ast.InterfaceType
	var genericsTypes genericTypes
	for _, t := range input.types {
		if t.Name.Name == i.Name {
			var ok bool
			embeddedInterface, ok = t.Type.(*ast.InterfaceType)
			if ok {
				genericsTypes = buildGenericTypesFromSpec(t, input.types, input.typesPrefix)
				break
			}

			if !checkInterface {
				break
			}

			return nil, errors.Wrap(errNotAnInterface, t.Name.Name)
		}
	}

	if embeddedInterface == nil {
		return nil, nil
	}

	input.genericTypes = genericsTypes
	return processInterface(embeddedInterface, input)
}

var errUnknownSelector = errors.New("unknown selector")

func findImportPathForName(name string, imports []*ast.ImportSpec, currentPackage *packages.Package) (string, error) {
	for _, i := range imports {
		if i.Name != nil && i.Name.Name == name {
			return unquote(i.Path.Value), nil
		}
	}

	for path, pkg := range currentPackage.Imports {
		if pkg.Name == name {
			return path, nil
		}
	}

	return "", errors.Wrapf(errUnknownSelector, name)
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
