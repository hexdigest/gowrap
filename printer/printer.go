package printer

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"github.com/pkg/errors"
)

// Printer prints AST nodes
type Printer struct {
	fs          *token.FileSet
	types       []*ast.TypeSpec
	typesPrefix string
	buf         *bytes.Buffer
}

// New returns Printer
// If srcPackageName is not empty Printer assumes that printed node goes to the package
// different from the srcPackageName and all references to the types listed in typeSpecs
// should be prepended with the source package name selector
func New(fs *token.FileSet, typeSpecs []*ast.TypeSpec, typesPrefix string) *Printer {
	return &Printer{
		fs:          fs,
		types:       typeSpecs,
		typesPrefix: typesPrefix,
		buf:         bytes.NewBuffer([]byte{}),
	}
}

// Print prints AST node as is
func (p *Printer) Print(node ast.Node) (string, error) {
	if node == nil {
		return "", nil
	}
	defer p.buf.Reset()

	err := printer.Fprint(p.buf, p.fs, node)

	return p.buf.String(), err
}

// PrintType prints node that represents a type, i.e. type of a function argument or a or struct field
// If node references one of the types listed in p.types and srcPackageName is not empty
// it adds source package selector to the type identifier
// PrintType returns error if the srcPackageName is not empty and the printed type is unexported.
func (p *Printer) PrintType(node ast.Node) (string, error) {
	defer p.buf.Reset()

	switch t := node.(type) {
	case *ast.FuncType:
		return p.printFunc(t)
	case *ast.StarExpr:
		return p.printPointer(t)
	case *ast.Ellipsis:
		return p.printVariadicParam(t)
	case *ast.ChanType:
		return p.printChan(t)
	case *ast.ArrayType:
		return p.printArray(t)
	case *ast.MapType:
		return p.printMap(t)
	case *ast.StructType:
		return p.printStruct(t)
	case *ast.Ident:
		return p.printIdent(t)
	}

	err := printer.Fprint(p.buf, p.fs, node)
	return p.buf.String(), err
}

func (p *Printer) printArray(a *ast.ArrayType) (string, error) {
	sliceType, err := p.PrintType(a.Elt)
	if err != nil {
		return "", err
	}

	l, err := p.Print(a.Len)
	if err != nil {
		return "", err
	}

	return "[" + l + "]" + sliceType, nil
}

var chanTypes = map[ast.ChanDir]string{
	ast.SEND & ast.RECV: "chan ",
	ast.SEND:            "chan<- ",
	ast.RECV:            "<-chan ",
}

func (p *Printer) printChan(c *ast.ChanType) (string, error) {
	valueType, err := p.PrintType(c.Value)
	if err != nil {
		return "", err
	}

	return chanTypes[c.Dir] + valueType, nil
}

func (p *Printer) printFunc(ft *ast.FuncType) (string, error) {
	params, err := p.fieldList(ft.Params)
	if err != nil {
		return "", err
	}
	results, err := p.fieldList(ft.Results)
	if err != nil {
		return "", err
	}

	return "func(" + strings.Join(params, ", ") + ") (" + strings.Join(results, ", ") + ")", nil
}

func (p *Printer) printMap(mt *ast.MapType) (string, error) {
	keyType, err := p.PrintType(mt.Key)
	if err != nil {
		return "", err
	}
	valueType, err := p.PrintType(mt.Value)
	if err != nil {
		return "", err
	}
	return "map[" + keyType + "]" + valueType, nil
}

var errUnexportedType = errors.New("unexported type")

func (p *Printer) printIdent(i *ast.Ident) (string, error) {
	for _, ts := range p.types {

		if i.Name == ts.Name.Name {
			if len(p.typesPrefix) > 0 {
				//destination file is in another package
				//if the found type matches one of the types declared in a source package and
				//this type is unexported
				if []rune(ts.Name.Name)[0] == []rune(strings.ToLower(ts.Name.Name))[0] {
					return "", errors.Wrap(errUnexportedType, ts.Name.Name)
				}
				return p.typesPrefix + "." + i.Name, nil
			}

			return i.Name, nil
		}
	}

	err := printer.Fprint(p.buf, p.fs, i)
	return p.buf.String(), err
}

func (p *Printer) printPointer(pt *ast.StarExpr) (string, error) {
	pointerTo, err := p.PrintType(pt.X)
	if err != nil {
		return "", err
	}

	return "*" + pointerTo, nil
}

func (p *Printer) printStruct(s *ast.StructType) (string, error) {
	fields, err := p.fieldList(s.Fields)
	if err != nil {
		return "", err
	}

	return "struct{\n" + strings.Join(fields, "\n") + "\n}", nil
}

func (p *Printer) printVariadicParam(e *ast.Ellipsis) (string, error) {
	sliceType, err := p.PrintType(e.Elt)
	if err != nil {
		return "", err
	}

	return "..." + sliceType, nil
}

func (p *Printer) fieldList(fl *ast.FieldList) ([]string, error) {
	if fl == nil {
		return nil, nil
	}

	params := []string{}

	for _, param := range fl.List {
		names := []string{}
		for _, name := range param.Names {
			names = append(names, name.Name)
		}

		paramType, err := p.PrintType(param.Type)
		if err != nil {
			return nil, err
		}

		params = append(params, strings.Join(names, ", ")+" "+paramType)
	}

	return params, nil
}
