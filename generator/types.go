package generator

import (
	"fmt"
	"go/ast"
	"strings"
)

type typePrinter interface {
	PrintType(ast.Node) (string, error)
}

// Method represents a method's signature
type Method struct {
	Name    string
	Params  ParamsSlice
	Results ParamsSlice

	ReturnsError   bool
	AcceptsContext bool
}

// Param represents fuction argument or result
type Param struct {
	Name     string
	Type     string
	Variadic bool
}

// ParamsSlice slice of parameters
type ParamsSlice []Param

// String implements fmt.Stringer
func (ps ParamsSlice) String() string {
	ss := []string{}
	for _, p := range ps {
		ss = append(ss, p.Name+" "+p.Type)
	}

	return strings.Join(ss, ", ")
}

// Pass returns comma separated params names to
// be passed to a function call with respect to
// variadic functions
func (ps ParamsSlice) Pass() string {
	params := []string{}
	for _, p := range ps {
		params = append(params, p.Pass())
	}

	return strings.Join(params, ", ")
}

// Pass returns a name of the parameter
// If parameter is variadic it returns a name followed by a ...
func (p Param) Pass() string {
	if p.Variadic {
		return p.Name + "..."
	}
	return p.Name
}

// NewMethod returns pointer to Signature struct or error
func NewMethod(name string, f *ast.FuncType, printer typePrinter) (*Method, error) {
	m := Method{Name: name}
	usedNames := map[string]bool{}

	//Always name the last return parameter as an "err" if it's of type "error"
	if f.Results != nil {
		ident, ok := f.Results.List[len(f.Results.List)-1].Type.(*ast.Ident)
		m.ReturnsError = ok && ident.Name == "error"
		usedNames["err"] = true
	}

	if len(f.Params.List) > 0 {
		if se, ok := f.Params.List[0].Type.(*ast.SelectorExpr); ok {
			m.AcceptsContext = ok && se.Sel.Name == "Context"
			usedNames["ctx"] = true
		}
	}

	var err error

	m.Params, err = makeParams(f.Params, usedNames, printer)
	if err != nil {
		return nil, err
	}
	m.Results, err = makeParams(f.Results, usedNames, printer)
	if err != nil {
		return nil, err
	}

	if m.ReturnsError {
		m.Results[len(m.Results)-1].Name = "err"
	}

	if m.AcceptsContext {
		m.Params[0].Name = "ctx"
	}

	return &m, nil
}

// NewParam returns Param struct
func NewParam(name string, typ ast.Expr, usedNames map[string]bool, printer typePrinter) (*Param, error) {
	if name == "" || usedNames[name] {
		name = genName(typePrefix(typ), 1, usedNames)
	}

	usedNames[name] = true

	typeStr, err := printer.PrintType(typ)
	if err != nil {
		return nil, err
	}

	_, variadic := typ.(*ast.Ellipsis)
	return &Param{
		Name:     name,
		Variadic: variadic,
		Type:     typeStr,
	}, nil
}

func makeParams(params *ast.FieldList, usedNames map[string]bool, printer typePrinter) (ParamsSlice, error) {
	if params == nil {
		return nil, nil
	}

	result := []Param{}
	for _, p := range params.List {
		//for anonymous parameters we generate params and results names
		//based on their type
		if p.Names == nil {
			param, err := NewParam("", p.Type, usedNames, printer)
			if err != nil {
				return nil, err
			}
			result = append(result, *param)
		} else {
			for _, ident := range p.Names {
				param, err := NewParam(ident.Name, p.Type, usedNames, printer)
				if err != nil {
					return nil, err
				}
				result = append(result, *param)
			}
		}
	}

	return result, nil
}

func typePrefix(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.SelectorExpr:
		return typePrefix(t.Sel)
	case *ast.StarExpr:
		return typePrefix(t.X) + "p" //*string -> sp (string pointer)
	case *ast.SliceExpr:
		return typePrefix(t.X) + "s" //[]string -> ss (string slice)
	case *ast.ArrayType:
		return typePrefix(t.Elt) + "a" //[2]string -> sa (string array)
	case *ast.MapType:
		return "m"
	case *ast.ChanType:
		return "ch"
	case *ast.StructType:
		return "st"
	case *ast.FuncType:
		return "f"
	case *ast.Ident:
		return strings.ToLower(t.Name[0:1])
	}

	return "p"
}

func genName(prefix string, n int, usedNames map[string]bool) string {
	name := fmt.Sprintf("%s%d", prefix, n)
	if usedNames[name] {
		return genName(prefix, n+1, usedNames)
	}

	return name
}

// Call returns a string with the method call
func (m Method) Call() string {
	params := []string{}
	for _, p := range m.Params {
		params = append(params, p.Pass())
	}

	return m.Name + "(" + strings.Join(params, ", ") + ")"
}

// Pass returns a return statement followed by the method call
// If method does not have any results it returns a method call followed by return statement
func (m Method) Pass(prefix string) string {
	if len(m.Results) > 0 {
		return "return " + prefix + m.Call()
	}

	return prefix + m.Call() + "\nreturn"
}

// ParamsNames returns a list of method params names
func (m Method) ParamsNames() string {
	ss := []string{}
	for _, p := range m.Params {
		ss = append(ss, p.Name)
	}
	return strings.Join(ss, ", ")
}

// ResultsNames returns a list of method results names
func (m Method) ResultsNames() string {
	ss := []string{}
	for _, r := range m.Results {
		ss = append(ss, r.Name)
	}
	return strings.Join(ss, ", ")
}

// ParamsStruct returns a struct type with fields corresponding
// to the method params
func (m Method) ParamsStruct() string {
	ss := []string{}
	for _, p := range m.Params {
		if p.Variadic {
			ss = append(ss, p.Name+" "+strings.Replace(p.Type, "...", "[]", 1))
		} else {
			ss = append(ss, p.Name+" "+p.Type)
		}
	}
	return "struct{\n" + strings.Join(ss, "\n ") + "}"
}

// ResultsStruct returns a struct type with fields corresponding
// to the method results
func (m Method) ResultsStruct() string {
	ss := []string{}
	for _, r := range m.Results {
		ss = append(ss, r.Name+" "+r.Type)
	}
	return "struct{\n" + strings.Join(ss, "\n ") + "}"
}

// ParamsMap returns a string representation of the map[string]interface{}
// filled with method's params
func (m Method) ParamsMap() string {
	ss := []string{}
	for _, p := range m.Params {
		ss = append(ss, `"`+p.Name+`": `+p.Name)
	}
	return "map[string]interface{}{\n" + strings.Join(ss, ",\n ") + "}"
}

// ResultsMap returns a string representation of the map[string]interface{}
// filled with method's results
func (m Method) ResultsMap() string {
	ss := []string{}
	for _, r := range m.Results {
		ss = append(ss, `"`+r.Name+`": `+r.Name)
	}
	return "map[string]interface{}{\n" + strings.Join(ss, ",\n ") + "}"
}

// HasParams returns true if method has params
func (m Method) HasParams() bool {
	return len(m.Params) > 0
}

// HasResults returns true if method has results
func (m Method) HasResults() bool {
	return len(m.Results) > 0
}

// ReturnStruct returns return statement with the return params
// taken from the structName
func (m Method) ReturnStruct(structName string) string {
	if len(m.Results) == 0 {
		return "return"
	}

	ss := []string{}
	for _, r := range m.Results {
		ss = append(ss, structName+"."+r.Name)
	}
	return "return " + strings.Join(ss, ", ")
}

// Signature returns comma separated method's params followed by the comma separated
// method's results
func (m Method) Signature() string {
	params := []string{}
	for _, p := range m.Params {
		params = append(params, p.Name+" "+p.Type)
	}

	results := []string{}
	for _, r := range m.Results {
		results = append(results, r.Name+" "+r.Type)
	}

	return "(" + strings.Join(params, ", ") + ") (" + strings.Join(results, ", ") + ")"
}

// Declaration returns a method name followed by it's signature
func (m Method) Declaration() string {
	return m.Name + m.Signature()
}
