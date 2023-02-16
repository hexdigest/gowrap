package generator

import (
	"go/ast"
	"strings"
)

const (
	genericsSeparator = ", "

	genericsSquareBracketStart = "["
	genericsSquareBracketEnd   = "]"
)

// TemplateInputGenerics subset of generics interface information used for template generation
type TemplateInputGenerics struct {
	// Types of the interface when using generics (e.g. [I, O any])
	Types string

	// Params of the interface when using generics (e.g. [I, O])
	Params string
}

type genericsParams []genericsParam

type genericsParam struct {
	Name   string
	Params genericsParams
}

func (g genericsParam) String() string {
	name := g.Name
	var subParamNames []string
	for _, subParam := range g.Params {
		subParamNames = append(subParamNames, subParam.String())
	}
	if len(g.Params) > 0 {
		name += genericsSquareBracketStart + strings.Join(subParamNames, genericsSeparator) + genericsSquareBracketEnd
	}
	return name
}

type genericsTypes []genericsType

type genericsType struct {
	Type  string
	Names []string
}

func genericsWithBracketsBuild(t string) string {
	if t != "" {
		t = genericsSquareBracketStart + t + genericsSquareBracketEnd
	}
	return t
}

func (g genericsTypes) buildVars() (string, string) {
	var types, typesSep string
	var params, paramsSep string

	for _, genType := range g {
		var paramsByType, paramsByTypeSep string

		for _, name := range genType.Names {
			paramsByType += paramsByTypeSep + name
			params += paramsSep + name
			paramsSep = genericsSeparator
			paramsByTypeSep = genericsSeparator
		}

		if paramsByType != "" {
			types += typesSep + paramsByType + " " + genType.Type
			typesSep = genericsSeparator
		}
	}

	return genericsWithBracketsBuild(types), genericsWithBracketsBuild(params)
}

func genericsTypesBuild(ts *ast.TypeSpec) (types genericsTypes) {
	if ts.TypeParams != nil {
		for _, param := range ts.TypeParams.List {
			if param != nil {
				if gpt, ok := param.Type.(*ast.Ident); ok {
					var paramNames []string
					for _, name := range param.Names {
						if name != nil {
							paramNames = append(paramNames, name.Name)
						}
					}
					types = append(types, genericsType{
						Type:  gpt.Name,
						Names: paramNames,
					})
				}
			}
		}
	}
	return
}

func genericsBuildParamString(typeStr string, genericsTypes genericsTypes, genericsParams genericsParams) string {
	c := 0
	for _, genType := range genericsTypes {
		for _, name := range genType.Names {
			if name == typeStr {
				if len(genericsParams) > c {
					return genericsParams[c].String()
				}
			}
			c++
		}
	}
	return typeStr
}
