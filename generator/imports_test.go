package generator

import (
	"bytes"
	"go/parser"
	"go/token"
	"testing"
	"text/template"
)

func TestGenerateDeduplicatesImportsAcrossTemplates(t *testing.T) {
	for _, header := range []string{
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
}
