package generator

import (
	"go/ast"
	"reflect"
	"testing"
)

func Test_genericParam_String(t *testing.T) {
	type fields struct {
		Name   string
		Params genericParams
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "struct with generic params",
			fields: fields{
				Name: "somepkg.SomeGenericStruct",
				Params: genericParams{
					{
						Name: "int",
					},
					{
						Name: "string",
					},
				},
			},
			want: "somepkg.SomeGenericStruct[int, string]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := genericParam{
				Name:   tt.fields.Name,
				Params: tt.fields.Params,
			}
			if got := g.String(); got != tt.want {
				t.Errorf("genericParam.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_genericTypes_buildVars(t *testing.T) {
	tests := []struct {
		name  string
		g     genericTypes
		want  string
		want1 string
	}{
		{
			name: "[T any]",
			g: genericTypes{
				{
					Names: []string{"T"},
					Type:  "any",
				},
			},
			want:  "[T any]",
			want1: "[T]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.g.buildVars()
			if got != tt.want {
				t.Errorf("genericTypes.buildVars() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("genericTypes.buildVars() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_buildGenericTypesFromSpec(t *testing.T) {
	type args struct {
		ts          *ast.TypeSpec
		allTypes    []*ast.TypeSpec
		typesPrefix string
	}
	tests := []struct {
		name      string
		args      args
		wantTypes genericTypes
	}{
		{
			name: "build generic types any from spec",
			args: args{
				ts: &ast.TypeSpec{
					TypeParams: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.Ident{
									Name: "any",
								},
								Names: []*ast.Ident{
									{
										Name: "I",
									},
									{
										Name: "O",
									},
								},
							},
						},
					},
				},
			},
			wantTypes: genericTypes{
				{
					Type:  "any",
					Names: []string{"I", "O"},
				},
			},
		},
		{
			name: "build generic types foo.Bar from spec",
			args: args{
				ts: &ast.TypeSpec{
					TypeParams: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.SelectorExpr{
									X: &ast.Ident{
										Name: "foo",
									},
									Sel: &ast.Ident{
										Name: "Bar",
									},
								},
								Names: []*ast.Ident{
									{
										Name: "I",
									},
									{
										Name: "O",
									},
								},
							},
						},
					},
				},
			},
			wantTypes: genericTypes{
				{
					Type:  "foo.Bar",
					Names: []string{"I", "O"},
				},
			},
		},
		{
			name: "build generic types Bar from spec without types prefix",
			args: args{
				ts: &ast.TypeSpec{
					TypeParams: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.Ident{
									Name: "Bar",
								},
								Names: []*ast.Ident{
									{
										Name: "I",
									},
									{
										Name: "O",
									},
								},
							},
						},
					},
				},
				allTypes: []*ast.TypeSpec{
					{
						Name: &ast.Ident{
							Name: "Bar",
						},
					},
				},
			},
			wantTypes: genericTypes{
				{
					Type:  "Bar",
					Names: []string{"I", "O"},
				},
			},
		},
		{
			name: "build generic types Bar from spec with types prefix",
			args: args{
				ts: &ast.TypeSpec{
					TypeParams: &ast.FieldList{
						List: []*ast.Field{
							{
								Type: &ast.Ident{
									Name: "Bar",
								},
								Names: []*ast.Ident{
									{
										Name: "I",
									},
									{
										Name: "O",
									},
								},
							},
						},
					},
				},
				allTypes: []*ast.TypeSpec{
					{
						Name: &ast.Ident{
							Name: "Bar",
						},
					},
				},
				typesPrefix: "prefix",
			},
			wantTypes: genericTypes{
				{
					Type:  "prefix.Bar",
					Names: []string{"I", "O"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTypes := buildGenericTypesFromSpec(tt.args.ts, tt.args.allTypes, tt.args.typesPrefix); !reflect.DeepEqual(gotTypes, tt.wantTypes) {
				t.Errorf("buildGenericTypesFromSpec() = %v, want %v", gotTypes, tt.wantTypes)
			}
		})
	}
}

func Test_buildGenericParamsString(t *testing.T) {
	genTypes := genericTypes{{Names: []string{"A", "B"}, Type: "any"}}
	genParams := genericParams{{Name: "string"}, {Name: "int"}}

	type args struct {
		typeStr       string
		genericTypes  genericTypes
		genericParams genericParams
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "replace A by string",
			args: args{
				typeStr:       "A",
				genericTypes:  genTypes,
				genericParams: genParams,
			},
			want: "string",
		},
		{
			name: "replace B by int",
			args: args{
				typeStr:       "B",
				genericTypes:  genTypes,
				genericParams: genParams,
			},
			want: "int",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildGenericParamsString(tt.args.typeStr, tt.args.genericTypes, tt.args.genericParams); got != tt.want {
				t.Errorf("buildGenericParamsString() = %v, want %v", got, tt.want)
			}
		})
	}
}
