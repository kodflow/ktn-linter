package ktn_struct_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

// TestRule004_EmptyFile tests Rule004 with an empty file.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule004_EmptyFile(t *testing.T) {
	src := `package test`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			t.Errorf("unexpected diagnostic: %v", d)
		},
	}

	_, err = ktn_struct.Rule004.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule004_EmptyStruct tests Rule004 with an empty struct.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule004_EmptyStruct(t *testing.T) {
	src := `package test
// EmptyStruct has no fields.
type EmptyStruct struct {
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			t.Errorf("unexpected diagnostic: %v", d)
		},
	}

	_, err = ktn_struct.Rule004.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule004_StructWithNilFields tests Rule004 with minimal fields.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule004_StructWithNilFields(t *testing.T) {
	src := `package test
// MinimalStruct is a minimal struct.
type MinimalStruct struct {
	Field string
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			t.Errorf("unexpected diagnostic: %v", d)
		},
	}

	_, err = ktn_struct.Rule004.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule004_StructWithEmbeddedField tests Rule004 with embedded fields.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule004_StructWithEmbeddedField(t *testing.T) {
	src := `package test
// EmbeddedStruct has an embedded field.
type Base struct {
	Field1 string
}

type EmbeddedStruct struct {
	Base
	Field2 string
	Field3 int
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			t.Errorf("unexpected diagnostic: %v", d)
		},
	}

	_, err = ktn_struct.Rule004.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule004_StructWithExactly15Fields tests Rule004 with exactly 15 fields.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule004_StructWithExactly15Fields(t *testing.T) {
	src := `package test
// MaxFieldsStruct has exactly 15 fields.
type MaxFieldsStruct struct {
	Field1  string
	Field2  int
	Field3  bool
	Field4  float64
	Field5  string
	Field6  int
	Field7  bool
	Field8  float64
	Field9  string
	Field10 int
	Field11 bool
	Field12 float64
	Field13 string
	Field14 int
	Field15 bool
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			t.Errorf("unexpected diagnostic: %v", d)
		},
	}

	_, err = ktn_struct.Rule004.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule004_StructWithTooManyFields tests Rule004 with too many fields.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule004_StructWithTooManyFields(t *testing.T) {
	src := `package test
// TooManyFieldsStruct has 16 fields.
type TooManyFieldsStruct struct {
	Field1  string
	Field2  int
	Field3  bool
	Field4  float64
	Field5  string
	Field6  int
	Field7  bool
	Field8  float64
	Field9  string
	Field10 int
	Field11 bool
	Field12 float64
	Field13 string
	Field14 int
	Field15 bool
	Field16 string
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	reported := false
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			reported = true
		},
	}

	_, err = ktn_struct.Rule004.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reported {
		t.Error("expected diagnostic to be reported")
	}
}

// TestRule004_NonStructType tests Rule004 with non-struct types.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule004_NonStructType(t *testing.T) {
	src := `package test
type MyInt int
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			t.Errorf("unexpected diagnostic: %v", d)
		},
	}

	_, err = ktn_struct.Rule004.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule004_MultipleFieldsPerLine tests Rule004 with multiple fields per line.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule004_MultipleFieldsPerLine(t *testing.T) {
	src := `package test
// MultiFields has multiple fields declared on the same line.
type MultiFields struct {
	Field1, Field2, Field3 string
	Field4, Field5 int
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			t.Errorf("unexpected diagnostic: %v", d)
		},
	}

	_, err = ktn_struct.Rule004.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
