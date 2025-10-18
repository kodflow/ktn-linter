package ktn_struct_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

// TestRule002_EmptyFile tests Rule002 with an empty file.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule002_EmptyFile(t *testing.T) {
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

	_, err = ktn_struct.Rule002.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule002_NonStructType tests Rule002 with non-struct types.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule002_NonStructType(t *testing.T) {
	src := `package test
type MyInt int
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

	_, err = ktn_struct.Rule002.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule002_DocumentedStruct tests Rule002 with a documented struct.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule002_DocumentedStruct(t *testing.T) {
	src := `package test
// ValidStruct represents a valid struct.
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

	_, err = ktn_struct.Rule002.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule002_UndocumentedStruct tests Rule002 with an undocumented struct.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule002_UndocumentedStruct(t *testing.T) {
	src := `package test
type UndocumentedStruct struct {
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

	_, err = ktn_struct.Rule002.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !reported {
		t.Error("expected diagnostic to be reported")
	}
}
