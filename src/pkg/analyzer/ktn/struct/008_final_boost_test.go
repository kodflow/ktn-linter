package ktn_struct_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"

	ktn_struct "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/struct"
)

// Tests spécifiques pour atteindre 100% sur checkStructFieldCount
func TestRule004_FieldsWithoutNames(t *testing.T) {
	// Test avec champs embedded (sans nom)
	src := `package test

type Base1 struct {
	Value string
}

type Base2 struct {
	Count int
}

// StructWithEmbedded has embedded fields.
type StructWithEmbedded struct {
	Base1
	Base2
	RegularField string
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

func TestRule004_TooManyFieldsWithEmbedded(t *testing.T) {
	// Test avec champs embedded qui font dépasser la limite
	src := `package test

type Base struct {
	Field string
}

// TooBigWithEmbedded has too many fields including embedded.
type TooBigWithEmbedded struct {
	Base
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
		t.Error("expected diagnostic to be reported for struct with too many fields")
	}
}

func TestRule004_MultipleEmbeddedFields(t *testing.T) {
	// Test avec plusieurs champs embedded
	src := `package test

type Base1 struct {
	Field1 string
}

type Base2 struct {
	Field2 int
}

type Base3 struct {
	Field3 bool
}

// MultipleEmbedded has multiple embedded fields.
type MultipleEmbedded struct {
	Base1
	Base2
	Base3
	OwnField string
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

func TestRule003_FieldsWithoutNamesEmbedded(t *testing.T) {
	// Test pour rule003 avec champs embedded
	src := `package test

// Base is a base struct.
type Base struct {
	// Value is documented.
	Value string
}

// StructWithEmbedded has embedded fields.
type StructWithEmbedded struct {
	Base
	// ExportedField is documented.
	ExportedField string
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

// Tests pour couvrir les boucles multiples files
func TestAllRules_MultipleFilesInPass(t *testing.T) {
	src1 := `package test

type ValidStruct1 struct {
	Field string
}
`
	src2 := `package test

type ValidStruct2 struct {
	Value int
}
`
	src3 := `package test

type ValidStruct3 struct {
	Data bool
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
	file3, err := parser.ParseFile(fset, "test3.go", src3, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file1, file2, file3},
		Report: func(d analysis.Diagnostic) {
			// Peut avoir des rapports
		},
	}

	// Test Rule001 avec plusieurs fichiers
	_, err = ktn_struct.Rule001.Run(pass)
	if err != nil {
		t.Errorf("Rule001 unexpected error: %v", err)
	}

	// Test Rule002 avec plusieurs fichiers
	_, err = ktn_struct.Rule002.Run(pass)
	if err != nil {
		t.Errorf("Rule002 unexpected error: %v", err)
	}

	// Test Rule003 avec plusieurs fichiers
	_, err = ktn_struct.Rule003.Run(pass)
	if err != nil {
		t.Errorf("Rule003 unexpected error: %v", err)
	}

	// Test Rule004 avec plusieurs fichiers
	_, err = ktn_struct.Rule004.Run(pass)
	if err != nil {
		t.Errorf("Rule004 unexpected error: %v", err)
	}
}
