package ktn_struct_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

// TestRule001_GroupedTypeDeclarations tests Rule001 with grouped type declarations.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule001_GroupedTypeDeclarations(t *testing.T) {
	src := `package test

// Multiple type declarations in one var
var (
	x = 1
)

const (
	y = 2
)

type (
	// ValidStruct1 is a valid struct.
	ValidStruct1 struct {
		Field string
	}

	invalid_struct_2 struct {
		Field string
	}

	MyInt int
	MyString string

	// ValidStruct3 is another valid struct.
	ValidStruct3 struct {
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

// TestRule002_GroupedTypeDeclarations tests Rule002 with grouped type declarations.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule002_GroupedTypeDeclarations(t *testing.T) {
	src := `package test

// DocumentedStruct is documented.
type DocumentedStruct struct {
	Field string
}

type UndocumentedStruct1 struct {
	Field string
}

type MyInt int

type UndocumentedStruct2 struct {
	Value int
}
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

	if reportCount != 2 {
		t.Errorf("expected 2 diagnostics, got %d", reportCount)
	}
}

// TestRule003_GroupedTypeDeclarations tests Rule003 with grouped type declarations.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule003_GroupedTypeDeclarations(t *testing.T) {
	src := `package test

// GoodStruct is documented.
type GoodStruct struct {
	// Field is documented.
	Field string
}

// BadStruct is documented.
type BadStruct1 struct {
	UndocumentedField string
}

type MyInt int

// BadStruct2 is documented.
type BadStruct2 struct {
	AnotherUndocumented int
}
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

	if reportCount != 2 {
		t.Errorf("expected 2 diagnostics, got %d", reportCount)
	}
}

// TestRule004_GroupedTypeDeclarations tests Rule004 with grouped type declarations.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule004_GroupedTypeDeclarations(t *testing.T) {
	src := `package test

// SmallStruct is small.
type SmallStruct struct {
	Field1 string
	Field2 int
}

// BigStruct has too many fields.
type BigStruct struct {
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

// TestRule001_MultipleFiles tests Rule001 with multiple files.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule001_MultipleFiles(t *testing.T) {
	src1 := `package test

type ValidStruct struct {
	Field string
}
`
	src2 := `package test

type invalid_struct struct {
	Field string
}
`

	fset := token.NewFileSet()
	file1, err := parser.ParseFile(fset, "test1.go", src1, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	file2, err := parser.ParseFile(fset, "test2.go", src2, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	reportCount := 0
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file1, file2},
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
