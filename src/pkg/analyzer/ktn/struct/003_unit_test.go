package ktn_struct_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

// TestRule003_EmptyFile tests Rule003 with an empty file.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule003_EmptyFile(t *testing.T) {
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

	_, err = ktn_struct.Rule003.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule003_EmptyStruct tests Rule003 with an empty struct.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule003_EmptyStruct(t *testing.T) {
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

	_, err = ktn_struct.Rule003.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule003_PrivateFields tests Rule003 with private fields only.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule003_PrivateFields(t *testing.T) {
	src := `package test
// PrivateFields has only private fields.
type PrivateFields struct {
	field1 string
	field2 int
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

	_, err = ktn_struct.Rule003.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule003_DocumentedExportedFields tests Rule003 with documented exported fields.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule003_DocumentedExportedFields(t *testing.T) {
	src := `package test
// ValidStruct has documented exported fields.
type ValidStruct struct {
	// Field1 is a string field.
	Field1 string
	// Field2 is an int field.
	Field2 int
	private string
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

	_, err = ktn_struct.Rule003.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestRule003_UndocumentedExportedFields tests Rule003 with undocumented exported fields.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule003_UndocumentedExportedFields(t *testing.T) {
	src := `package test
// BadStruct has undocumented exported fields.
type BadStruct struct {
	Field1 string
	Field2 int
	private string
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

// TestRule003_NonStructType tests Rule003 with non-struct types.
// nolint:KTN-FUNC-001 // Test naming convention
func TestRule003_NonStructType(t *testing.T) {
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

	_, err = ktn_struct.Rule003.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
