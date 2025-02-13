package printer

import (
	"bytes"
	"go/ast"
	"go/token"
	"testing"
	"time"

	minimock "github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, New(nil, nil, ""))
}

func TestPrinter_Print(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T) //inspects *Printer after execution of Print

		node ast.Node

		want1 string
	}{
		{
			name: "nil node",
			init: func(t minimock.Tester) *Printer {
				return &Printer{}
			},
			node:  nil,
			want1: "",
		},
		{
			name: "success",
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:  token.NewFileSet(),
					buf: bytes.NewBuffer([]byte{}),
				}
			},
			node:  &ast.Ident{Name: "name"},
			want1: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.Print(tt.node)
			require.NoError(t, err)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.Print returned unexpected result")
		})
	}
}

func TestPrinter_fieldList(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T) //inspects *Printer after execution of fieldList

		fl func(t minimock.Tester) *ast.FieldList

		want1      []string
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation
	}{
		{
			name: "nil field list",
			fl:   func(t minimock.Tester) *ast.FieldList { return nil },
			init: func(t minimock.Tester) *Printer {
				return &Printer{}
			},
			want1:   nil,
			wantErr: false,
		},
		{
			name: "print error",
			fl: func(t minimock.Tester) *ast.FieldList {
				return &ast.FieldList{
					List: []*ast.Field{
						{Type: &ast.Ident{Name: "unexported"}},
					},
				}
			},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "unexported"}}},
					typesPrefix: "otherPackage",
				}
			},
			want1:   nil,
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexportedType, errors.Cause(err))
			},
		},
		{
			name: "success",
			fl: func(t minimock.Tester) *ast.FieldList {
				return &ast.FieldList{
					List: []*ast.Field{
						{Names: []*ast.Ident{{Name: "param"}}, Type: &ast.Ident{Name: "ExportedType"}},
					},
				}
			},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:    token.NewFileSet(),
					buf:   bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{{Name: &ast.Ident{Name: "ExportedType"}}},
				}
			},
			want1: []string{"param ExportedType"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.fieldList(tt.fl(mc))

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.fieldList returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestPrinter_printArray(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T)

		array *ast.ArrayType

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name:  "unexported type",
			array: &ast.ArrayType{Elt: &ast.Ident{Name: "unexported"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "unexported"}}},
					typesPrefix: "otherPackage",
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexportedType, errors.Cause(err))
			},
		},
		{
			name:  "success",
			array: &ast.ArrayType{Elt: &ast.Ident{Name: "Exported"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:    token.NewFileSet(),
					buf:   bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{{Name: &ast.Ident{Name: "Exported"}}},
				}
			},
			want1:   "[]Exported",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.printArray(tt.array)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.printArray returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestPrinter_printChan(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T)

		ch *ast.ChanType

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "unexported type",
			ch:   &ast.ChanType{Value: &ast.Ident{Name: "unexported"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "unexported"}}},
					typesPrefix: "otherPackage",
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexportedType, errors.Cause(err))
			},
		},
		{
			name: "bidirectional channel",
			ch:   &ast.ChanType{Value: &ast.Ident{Name: "Exported"}, Dir: ast.SEND | ast.RECV},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:    token.NewFileSet(),
					buf:   bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{{Name: &ast.Ident{Name: "Exported"}}},
				}
			},
			want1:   "chan Exported",
			wantErr: false,
		},
		{
			name: "recv",
			ch:   &ast.ChanType{Value: &ast.Ident{Name: "Recv"}, Dir: ast.RECV},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:    token.NewFileSet(),
					buf:   bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{{Name: &ast.Ident{Name: "Recv"}}},
				}
			},
			want1:   "<-chan Recv",
			wantErr: false,
		},
		{
			name: "send",
			ch:   &ast.ChanType{Value: &ast.Ident{Name: "Send"}, Dir: ast.SEND},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:    token.NewFileSet(),
					buf:   bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{{Name: &ast.Ident{Name: "Send"}}},
				}
			},
			want1:   "chan<- Send",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.printChan(tt.ch)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.printArray returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestPrinter_printFunc(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T)

		f *ast.FuncType

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "unexported param type",
			f:    &ast.FuncType{Params: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "unexported"}}}}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "unexported"}}},
					typesPrefix: "otherPackage",
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexportedType, errors.Cause(err))
			},
		},
		{
			name: "unexported result type",
			f:    &ast.FuncType{Results: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "unexported"}}}}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "unexported"}}},
					typesPrefix: "otherPackage",
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexportedType, errors.Cause(err))
			},
		},
		{
			name: "success",
			f:    &ast.FuncType{},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:    token.NewFileSet(),
					buf:   bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{{Name: &ast.Ident{Name: "Exported"}}},
				}
			},
			want1:   "func() ()",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.printFunc(tt.f)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.printArray returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestPrinter_printMap(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T)

		m *ast.MapType

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "unexported key type",
			m:    &ast.MapType{Key: &ast.Ident{Name: "unexported"}, Value: &ast.Ident{Name: "Exported"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "unexported"}}},
					typesPrefix: "otherPackage",
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexportedType, errors.Cause(err))
			},
		},
		{
			name: "unexported value type",
			m:    &ast.MapType{Key: &ast.Ident{Name: "Exported"}, Value: &ast.Ident{Name: "unexported"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "unexported"}}},
					typesPrefix: "otherPackage",
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexportedType, errors.Cause(err))
			},
		},
		{
			name: "success",
			m:    &ast.MapType{Key: &ast.Ident{Name: "Exported"}, Value: &ast.Ident{Name: "Exported"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:    token.NewFileSet(),
					buf:   bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{{Name: &ast.Ident{Name: "Exported"}}},
				}
			},
			want1:   "map[Exported]Exported",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.printMap(tt.m)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.printArray returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestPrinter_printPointer(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T)

		p *ast.StarExpr

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "pointer unexported type",
			p:    &ast.StarExpr{X: &ast.Ident{Name: "unexported"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "unexported"}}},
					typesPrefix: "otherPackage",
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexportedType, errors.Cause(err))
			},
		},
		{
			name: "success",
			p:    &ast.StarExpr{X: &ast.Ident{Name: "Exported"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:    token.NewFileSet(),
					buf:   bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{{Name: &ast.Ident{Name: "Exported"}}},
				}
			},
			want1:   "*Exported",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.printPointer(tt.p)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.printPointer returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestPrinter_printStruct(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T)

		s *ast.StructType

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "contains a fields with unexported type",
			s:    &ast.StructType{Fields: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "unexported"}}}}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "unexported"}}},
					typesPrefix: "otherPackage",
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexportedType, errors.Cause(err))
			},
		},
		{
			name: "success",
			s:    &ast.StructType{Fields: &ast.FieldList{List: []*ast.Field{{Type: &ast.Ident{Name: "Exported"}}}}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:  token.NewFileSet(),
					buf: bytes.NewBuffer([]byte{}),
				}
			},
			want1:   "struct{\n Exported\n}",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.printStruct(tt.s)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.printStruct returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestPrinter_printVariadicParam(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T)

		vp *ast.Ellipsis

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "unexported type",
			vp:   &ast.Ellipsis{Elt: &ast.Ident{Name: "unexported"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "unexported"}}},
					typesPrefix: "otherPackage",
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexportedType, errors.Cause(err))
			},
		},
		{
			name: "success",
			vp:   &ast.Ellipsis{Elt: &ast.Ident{Name: "Exported"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:    token.NewFileSet(),
					buf:   bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{{Name: &ast.Ident{Name: "Exported"}}},
				}
			},
			want1:   "...Exported",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.printVariadicParam(tt.vp)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.printArray returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestPrinter_printIdent(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T)

		ident *ast.Ident

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name:  "unexported type",
			ident: &ast.Ident{Name: "unexported"},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "unexported"}}},
					typesPrefix: "otherpackage",
				}
			},
			wantErr: true,
			inspectErr: func(err error, t *testing.T) {
				assert.Equal(t, errUnexportedType, errors.Cause(err))
			},
		},
		{
			name:  "success",
			ident: &ast.Ident{Name: "Exported"},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					typesPrefix: "prefix",
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "Exported"}}},
				}
			},
			want1:   "prefix.Exported",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.printIdent(tt.ident)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.printIdent returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestPrinter_printGeneric(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T)

		indexExpr *ast.IndexExpr

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "success",
			indexExpr: &ast.IndexExpr{
				X: &ast.Ident{
					Name: "Bar",
				},
				Index: &ast.Ident{
					Name: "Baz",
				},
			},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					typesPrefix: "prefix",
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "Bar"}}, {Name: &ast.Ident{Name: "Baz"}}},
				}
			},
			want1:   "prefix.Bar[prefix.Baz]",
			wantErr: false,
		},
		{
			name: "success, generic from other package",
			indexExpr: &ast.IndexExpr{
				X: &ast.Ident{
					Name: "Bar",
				},
				Index: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "otherpkg",
					},
					Sel: &ast.Ident{
						Name: "Baz",
					},
				},
			},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					typesPrefix: "prefix",
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types:       []*ast.TypeSpec{{Name: &ast.Ident{Name: "Bar"}}},
				}
			},
			want1:   "prefix.Bar[otherpkg.Baz]",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.printGeneric(tt.indexExpr)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.printGeneric returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestPrinter_printGenericList(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T)

		indexListExpr *ast.IndexListExpr

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "success",
			indexListExpr: &ast.IndexListExpr{
				X: &ast.Ident{
					Name: "Bar",
				},
				Indices: []ast.Expr{
					&ast.Ident{
						Name: "Baz",
					},
					&ast.Ident{
						Name: "Bak",
					},
				},
			},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					typesPrefix: "prefix",
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{
						{Name: &ast.Ident{Name: "Bar"}},
						{Name: &ast.Ident{Name: "Baz"}},
						{Name: &ast.Ident{Name: "Bak"}},
					},
				}
			},
			want1:   "prefix.Bar[prefix.Baz, prefix.Bak]",
			wantErr: false,
		},
		{
			name: "success, generic from other package",
			indexListExpr: &ast.IndexListExpr{
				X: &ast.Ident{
					Name: "Bar",
				},
				Indices: []ast.Expr{
					&ast.Ident{
						Name: "Baz",
					},
					&ast.SelectorExpr{
						X: &ast.Ident{
							Name: "otherpkg",
						},
						Sel: &ast.Ident{
							Name: "Bak",
						},
					},
				},
			},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					typesPrefix: "prefix",
					fs:          token.NewFileSet(),
					buf:         bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{
						{Name: &ast.Ident{Name: "Bar"}},
						{Name: &ast.Ident{Name: "Baz"}},
					},
				}
			},
			want1:   "prefix.Bar[prefix.Baz, otherpkg.Bak]",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.printGenericList(tt.indexListExpr)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.printGenericList returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}

func TestPrinter_PrintType(t *testing.T) {
	tests := []struct {
		name    string
		init    func(t minimock.Tester) *Printer
		inspect func(r *Printer, t *testing.T)

		node ast.Node

		want1      string
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "func type",
			node: &ast.FuncType{},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:  token.NewFileSet(),
					buf: bytes.NewBuffer([]byte{}),
				}
			},
			want1: "func() ()",
		},
		{
			name: "pointer type",
			node: &ast.StarExpr{X: &ast.Ident{Name: "string"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:  token.NewFileSet(),
					buf: bytes.NewBuffer([]byte{}),
				}
			},
			want1: "*string",
		},
		{
			name: "variadic type",
			node: &ast.Ellipsis{Elt: &ast.Ident{Name: "string"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:  token.NewFileSet(),
					buf: bytes.NewBuffer([]byte{}),
				}
			},
			want1: "...string",
		},
		{
			name: "array type",
			node: &ast.ArrayType{Elt: &ast.Ident{Name: "string"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:  token.NewFileSet(),
					buf: bytes.NewBuffer([]byte{}),
				}
			},
			want1: "[]string",
		},
		{
			name: "map type",
			node: &ast.MapType{Key: &ast.Ident{Name: "string"}, Value: &ast.Ident{Name: "int"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:  token.NewFileSet(),
					buf: bytes.NewBuffer([]byte{}),
				}
			},
			want1: "map[string]int",
		},
		{
			name: "chan type",
			node: &ast.ChanType{Value: &ast.Ident{Name: "string"}, Dir: ast.SEND | ast.RECV},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:  token.NewFileSet(),
					buf: bytes.NewBuffer([]byte{}),
				}
			},
			want1: "chan string",
		},
		{
			name: "struct type",
			node: &ast.StructType{},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:  token.NewFileSet(),
					buf: bytes.NewBuffer([]byte{}),
				}
			},
			want1: "struct{\n\n}",
		},
		{
			name: "selector expression (imported type)",
			node: &ast.SelectorExpr{X: &ast.Ident{Name: "package"}, Sel: &ast.Ident{Name: "Identifier"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:  token.NewFileSet(),
					buf: bytes.NewBuffer([]byte{}),
				}
			},
			want1: "package.Identifier",
		},
		{
			name: "generic type",
			node: &ast.IndexExpr{X: &ast.Ident{Name: "Bar"}, Index: &ast.Ident{Name: "string"}},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:    token.NewFileSet(),
					buf:   bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{{Name: &ast.Ident{Name: "Bar"}}},
				}
			},
			want1: "Bar[string]",
		},
		{
			name: "generic list type",
			node: &ast.IndexListExpr{
				X: &ast.Ident{Name: "Bar"},
				Indices: []ast.Expr{
					&ast.Ident{Name: "string"},
					&ast.Ident{Name: "int"},
				},
			},
			init: func(t minimock.Tester) *Printer {
				return &Printer{
					fs:    token.NewFileSet(),
					buf:   bytes.NewBuffer([]byte{}),
					types: []*ast.TypeSpec{{Name: &ast.Ident{Name: "Bar"}}},
				}
			},
			want1: "Bar[string, int]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			defer mc.Wait(time.Second)

			receiver := tt.init(mc)

			got1, err := receiver.PrintType(tt.node)

			if tt.inspect != nil {
				tt.inspect(receiver, t)
			}

			assert.Equal(t, tt.want1, got1, "Printer.PrintType returned unexpected result")

			if tt.wantErr {
				if assert.Error(t, err) && tt.inspectErr != nil {
					tt.inspectErr(err, t)
				}
			} else {
				assert.NoError(t, err)
			}

		})
	}
}
