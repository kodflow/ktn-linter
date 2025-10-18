package ktn_const_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	ktn_const "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/const"
	"golang.org/x/tools/go/analysis"
)

// TestRule004_DirectCoverage teste directement Rule004 pour améliorer la couverture.
func TestRule004_DirectCoverage(t *testing.T) {
	src := `package test

// Groupe sans iota, sans type explicite
const (
	BadNoType = 42
	BadNoType2 = "test"
)

// Groupe avec iota et différentes expressions
const (
	WithBinary = 1 << (10 * iota)
	WithUnary = -iota
	WithParen = (iota)
	WithCall = byte(iota)
	WithCallNoArgs = len("")  // CallExpr sans iota
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

	_, err = ktn_const.Rule004.Run(pass)
	if err != nil {
		t.Errorf("Rule004.Run() error = %v", err)
	}
}
