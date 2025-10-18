package ktn_struct_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

// TestRule001_EmptyFile tests Rule001 with an empty file.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule001_EmptyFile(t *testing.T) {
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

	_, err = ktn_struct.Rule001.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule001_NonTypeDecl tests Rule001 with non-type declarations.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule001_NonTypeDecl(t *testing.T) {
	src := `package test
var x = 42
const y = "hello"
func foo() {}
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

	_, err = ktn_struct.Rule001.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule001_InterfaceType tests Rule001 with interface types.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule001_InterfaceType(t *testing.T) {
	src := `package test
type MyInterface interface {
	Method()
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

	_, err = ktn_struct.Rule001.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule001_ValidStructName tests Rule001 with a valid struct name.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule001_ValidStructName(t *testing.T) {
	src := `package test
type ValidStruct struct {
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

	_, err = ktn_struct.Rule001.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule001_InvalidStructName tests Rule001 with an invalid struct name.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule001_InvalidStructName(t *testing.T) {
	src := `package test
type invalid_struct struct {
	Field string
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

	_, err = ktn_struct.Rule001.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reported {
		t.Error("expected diagnostic to be reported")
	}
}
