package ktn_method_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	ktn_method "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/method"
	"golang.org/x/tools/go/analysis"
)

// TestRule001_DirectCoverage tests Rule001 directly for coverage improvement.
//
// Params:
//   - t: testing instance
func TestRule001_DirectCoverage(t *testing.T) {
	src := `package test

type SimpleType struct {
	value int
}

type PointerType struct {
	data string
}

// Méthode avec receiver par valeur qui modifie
func (s SimpleType) Mutate() {
	s.value = 10
}

// Méthode avec receiver pointeur qui modifie
func (p *PointerType) Update() {
	p.data = "updated"
}

// Méthode avec receiver par valeur qui ne modifie pas
func (s SimpleType) Read() int {
	return s.value
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(diag analysis.Diagnostic) {
			t.Logf("Diagnostic: %s at %s", diag.Message, fset.Position(diag.Pos))
		},
	}

	_, err = ktn_method.Rule001.Run(pass)
	if err != nil {
		t.Errorf("Rule001.Run() error = %v", err)
	}
}
