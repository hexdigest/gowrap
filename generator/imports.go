package generator

import (
	"go/ast"
	"go/token"
	"strconv"
)

// removeDuplicateImports removes identical bindings across all import blocks.
// Different aliases for the same path are distinct bindings and must be kept.
func removeDuplicateImports(file *ast.File) (bool, error) {
	type binding struct {
		name string
		path string
	}
	seen := make(map[binding]bool)
	changed := false
	decls := file.Decls[:0]
	file.Imports = nil
	for _, decl := range file.Decls {
		imports, ok := decl.(*ast.GenDecl)
		if !ok || imports.Tok != token.IMPORT {
			decls = append(decls, decl)
			continue
		}
		specs := imports.Specs[:0]
		for _, spec := range imports.Specs {
			imp := spec.(*ast.ImportSpec)
			path, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				return false, err
			}
			key := binding{path: path}
			if imp.Name != nil {
				key.name = imp.Name.Name
			}
			if seen[key] {
				changed = true
				continue
			}
			seen[key] = true
			specs = append(specs, spec)
			file.Imports = append(file.Imports, imp)
		}
		imports.Specs = specs
		if len(specs) > 0 {
			decls = append(decls, decl)
		}
	}
	file.Decls = decls
	return changed, nil
}
