package ktn_struct_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

// Tests pour couvrir les cas de boucles avec déclarations de type mixtes

// TestRule001_MultipleTypeSpecsInOneDecl tests Rule001 with multiple type specs in one declaration.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule001_MultipleTypeSpecsInOneDecl(t *testing.T) {
	src := `package test
type (
	ValidStruct struct {
		Field string
	}
	invalid_struct struct {
		Field string
	}
	MyInt int
	AnotherStruct struct {
		Value int
	}
)
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	reportCount := 0
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = ktn_struct.Rule001.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if reportCount != 1 {
		t.Errorf("expected 1 diagnostic, got %d", reportCount)
	}
}

// TestRule002_MultipleTypeSpecsInOneDecl tests Rule002 with multiple type specs in one declaration.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule002_MultipleTypeSpecsInOneDecl(t *testing.T) {
	src := `package test

// ValidStruct is documented.
type ValidStruct struct {
	Field string
}

type UndocumentedStruct struct {
	Field string
}

type MyInt int
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	reportCount := 0
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = ktn_struct.Rule002.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if reportCount != 1 {
		t.Errorf("expected 1 diagnostic, got %d", reportCount)
	}
}

// TestRule003_MultipleTypeSpecsInOneDecl tests Rule003 with multiple type specs in one declaration.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule003_MultipleTypeSpecsInOneDecl(t *testing.T) {
	src := `package test

// GoodStruct is documented.
type GoodStruct struct {
	// Field is documented.
	Field string
}

// BadStruct is documented but fields are not.
type BadStruct struct {
	ExportedField string
}

type MyInt int
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	reportCount := 0
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = ktn_struct.Rule003.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if reportCount != 1 {
		t.Errorf("expected 1 diagnostic, got %d", reportCount)
	}
}

// TestRule004_MultipleTypeSpecsInOneDecl tests Rule004 with multiple type specs in one declaration.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule004_MultipleTypeSpecsInOneDecl(t *testing.T) {
	src := `package test

// GoodStruct has few fields.
type GoodStruct struct {
	Field1 string
	Field2 int
}

// BadStruct has too many fields.
type BadStruct struct {
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

type MyInt int
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	reportCount := 0
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			reportCount++
		},
	}

	_, err = ktn_struct.Rule004.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if reportCount != 1 {
		t.Errorf("expected 1 diagnostic, got %d", reportCount)
	}
}

// TestRule004_StructWithMultipleFieldsPerDeclaration tests Rule004 with multiple fields per declaration.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule004_StructWithMultipleFieldsPerDeclaration(t *testing.T) {
	src := `package test
// MultiFields has multiple fields per line.
type MultiFields struct {
	Field1, Field2, Field3, Field4, Field5 string
	Field6, Field7, Field8, Field9, Field10 int
	Field11, Field12, Field13, Field14, Field15 bool
	Field16 float64
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
		t.Error("expected diagnostic to be reported for 16 fields")
	}
}
