package ktn_struct_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

// Test pour couvrir les branches GenDecl non-TYPE (IMPORT, CONST, VAR)
func TestRule001_IgnoresNonTypeDeclInGenDecl(t *testing.T) {
	src := `package test

import (
	"fmt"
	"os"
)

const (
	MaxSize = 100
	MinSize = 10
)

var (
	GlobalVar = 42
	AnotherVar = "test"
)

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

func TestRule002_IgnoresNonTypeDeclInGenDecl(t *testing.T) {
	src := `package test

import "fmt"

const MaxSize = 100

var GlobalVar = 42

// ValidStruct is documented.
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

func TestRule003_IgnoresNonTypeDeclInGenDecl(t *testing.T) {
	src := `package test

import "fmt"

const MaxSize = 100

var GlobalVar = 42

// ValidStruct is documented.
type ValidStruct struct {
	// Field is documented.
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

	_, err = ktn_struct.Rule003.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestRule004_IgnoresNonTypeDeclInGenDecl(t *testing.T) {
	src := `package test

import "fmt"

const MaxSize = 100

var GlobalVar = 42

// ValidStruct is documented.
type ValidStruct struct {
	Field1 string
	Field2 int
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

// Test pour les TypeSpec qui ne sont pas des *ast.TypeSpec
func TestRule001_IgnoresNonStructTypes(t *testing.T) {
	src := `package test

type (
	ValidStruct struct {
		Field string
	}
	MyInt int
	MyString string
	MyFloat float64
)
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
			// Des diagnostics peuvent être émis pour les structs invalides
		},
	}

	_, err = ktn_struct.Rule001.Run(pass)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// Test pour couvrir checkStructFieldCount avec struct fields nil
func TestRule004_StructWithNilFieldsList(t *testing.T) {
	// Note: Ce cas est difficile à créer avec le parser Go standard
	// car une struct valide a toujours un Fields non-nil
	// Mais le test existe pour documenter le comportement
	src := `package test

// EmptyStruct is empty.
type EmptyStruct struct {}

// OnlyOnefield has one field.
type OnlyOneField struct {
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
