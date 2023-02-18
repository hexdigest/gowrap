package generator

import (
	"go/ast"
	"strings"
)

const (
	genericSeparator = ", "

	genericSquareBracketStart = "["
	genericSquareBracketEnd   = "]"
)

// TemplateInputGenerics subset of generics interface information used for template generation
type TemplateInputGenerics struct {
	// Types of the interface when using generics (e.g. [I, O any])
	Types string

	// Params of the interface when using generics (e.g. [I, O])
	Params string
}

type genericParams []genericParam

type genericParam struct {
	Name   string
	Params genericParams
}

func (g genericParam) String() string {
	name := g.Name
	var subParamNames []string
	for _, subParam := range g.Params {
		subParamNames = append(subParamNames, subParam.String())
	}
	if len(g.Params) > 0 {
		name += genericSquareBracketStart + strings.Join(subParamNames, genericSeparator) + genericSquareBracketEnd
	}
	return name
}

type genericTypes []genericType

type genericType struct {
	Type  string
	Names []string
}

func buildGenericsWithBrackets(t string) string {
	if t != "" {
		t = genericSquareBracketStart + t + genericSquareBracketEnd
	}
	return t
}

func (g genericTypes) buildVars() (string, string) {
	var types, typesSep string
	var params, paramsSep string

	for _, genType := range g {
		var paramsByType, paramsByTypeSep string

		for _, name := range genType.Names {
			paramsByType += paramsByTypeSep + name
			params += paramsSep + name
			paramsSep = genericSeparator
			paramsByTypeSep = genericSeparator
		}

		if paramsByType != "" {
			types += typesSep + paramsByType + " " + genType.Type
			typesSep = genericSeparator
		}
	}

	return buildGenericsWithBrackets(types), buildGenericsWithBrackets(params)
}

func buildGenericTypesFromSpec(ts *ast.TypeSpec) (types genericTypes) {
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
					types = append(types, genericType{
						Type:  gpt.Name,
						Names: paramNames,
					})
				}
			}
		}
	}
	return
}

func buildGenericParamsString(typeStr string, genericTypes genericTypes, genericParams genericParams) string {
	i := 0
	for _, genType := range genericTypes {
		for _, name := range genType.Names {
			if name == typeStr {
				if len(genericParams) > i {
					return genericParams[i].String()
				}
			}
			i++
		}
	}
	return typeStr
}
