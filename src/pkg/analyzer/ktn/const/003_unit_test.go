package ktn_const_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	ktn_const "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/const"
	"golang.org/x/tools/go/analysis"
)

// TestRule003_DirectCoverage teste directement Rule003 pour améliorer la couverture.
func TestRule003_DirectCoverage(t *testing.T) {
	src := `package test

// Groupe avec plusieurs constantes sans commentaires individuels
const (
	BadNoComment1 = 42
	BadNoComment2 = "test"
	BadNoComment3 = 3.14
)

// Groupe avec commentaires individuels inline
const (
	GoodInline1 = 1  // inline comment
	GoodInline2 = 2  // inline comment
)

// Groupe avec commentaires doc individuels
const (
	// GoodDoc1 has doc
	GoodDoc1 = 10
	// GoodDoc2 has doc
	GoodDoc2 = 20
)
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

	_, err = ktn_const.Rule003.Run(pass)
	if err != nil {
		t.Errorf("Rule003.Run() error = %v", err)
	}
}
