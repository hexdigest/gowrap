package generator

import (
	"bytes"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"
	"text/template"
)

func TestGenerateDeduplicatesImportsAcrossTemplates(t *testing.T) {
	for _, header := range []string{
		"package generated\nimport clock `time`\n",
		`package generated
import clock "time"
`,
		`package generated
import (
clock "time"
)
`,
	} {
		t.Run(header, func(t *testing.T) {
			g := Generator{
				Options:        Options{Imports: []string{`clock "time"`}},
				headerTemplate: template.Must(template.New("header").Parse(header)),
				bodyTemplate: template.Must(template.New("body").Parse(`
{{.Import}}
func Now() clock.Time { return clock.Now() }
`)),
			}
			var output bytes.Buffer
			if err := g.Generate(&output); err != nil {
				t.Fatal(err)
			}
			file, err := parser.ParseFile(token.NewFileSet(), "generated.go", output.Bytes(), 0)
			if err != nil {
				t.Fatal(err)
			}
			if len(file.Imports) != 1 {
				t.Fatalf("expected one import, got %d:\n%s", len(file.Imports), &output)
			}
			if file.Imports[0].Name.Name != "clock" {
				t.Fatalf("import alias was changed:\n%s", &output)
			}
			checkGeneratedCode(t, output.Bytes())
		})
	}
}

func TestGeneratePreservesDistinctImportAliases(t *testing.T) {
	g := Generator{
		Options: Options{Imports: []string{`second "time"`}},
		headerTemplate: template.Must(template.New("header").Parse(`package generated
import first "time"
`)),
		bodyTemplate: template.Must(template.New("body").Parse(`
{{.Import}}
func Convert(value first.Time) second.Time { return value }
`)),
	}
	var output bytes.Buffer
	if err := g.Generate(&output); err != nil {
		t.Fatal(err)
	}
	file, err := parser.ParseFile(token.NewFileSet(), "generated.go", output.Bytes(), 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(file.Imports) != 2 {
		t.Fatalf("expected distinct aliases, got %d:\n%s", len(file.Imports), &output)
	}
	checkGeneratedCode(t, output.Bytes())
}

func TestGenerateDeduplicatesImportsPreservingComments(t *testing.T) {
	g := Generator{
		headerTemplate: template.Must(template.New("header").Parse(`//go:build !ignore

// Code generated for testing. DO NOT EDIT.
package generated

// Header imports.
import clock "time" // Header clock.
`)),
		bodyTemplate: template.Must(template.New("body").Parse(`
// Body imports.
import (
clock "time" // Body clock.
)

// Now returns the current time.
func Now() clock.Time { return clock.Now() }
`)),
	}
	var output bytes.Buffer
	if err := g.Generate(&output); err != nil {
		t.Fatal(err)
	}
	for _, comment := range []string{
		"//go:build !ignore\n\n",
		"// Code generated for testing. DO NOT EDIT.",
		"// Header imports.",
		"// Header clock.",
		"// Body imports.",
		"// Body clock.",
		"// Now returns the current time.\nfunc Now",
	} {
		if !strings.Contains(output.String(), comment) {
			t.Errorf("comment missing or moved: %q:\n%s", comment, &output)
		}
	}
	file := checkGeneratedCode(t, output.Bytes())
	if len(file.Imports) != 1 {
		t.Fatalf("expected one import, got %d:\n%s", len(file.Imports), &output)
	}
}

func TestGenerateDeduplicatesImportsWithLocalPrefix(t *testing.T) {
	g := Generator{
		Options: Options{Imports: []string{
			`clock "time"`,
			`"github.com/pkg/errors"`,
			`"github.com/hexdigest/gowrap/pkg"`,
		}},
		localPrefix: "github.com/hexdigest/gowrap",
		headerTemplate: template.Must(template.New("header").Parse(`package generated
import (
clock "time"
"github.com/hexdigest/gowrap/pkg"
)
`)),
		bodyTemplate: template.Must(template.New("body").Parse(`
{{.Import}}
func Now(p pkg.Package) (clock.Time, error) { return clock.Now(), errors.New(p.Name) }
`)),
	}
	var output bytes.Buffer
	if err := g.Generate(&output); err != nil {
		t.Fatal(err)
	}
	const wantImports = `import (
	clock "time"

	"github.com/pkg/errors"

	"github.com/hexdigest/gowrap/pkg"
)`
	if !strings.Contains(output.String(), wantImports) {
		t.Fatalf("expected standard, third-party, and local import groups:\n%s", &output)
	}
}

func checkGeneratedCode(t *testing.T, source []byte) *ast.File {
	t.Helper()
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "generated.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	config := types.Config{Importer: importer.Default()}
	if _, err := config.Check("generated", fileSet, []*ast.File{file}, nil); err != nil {
		t.Fatalf("generated code does not type-check: %v:\n%s", err, source)
	}
	return file
}
