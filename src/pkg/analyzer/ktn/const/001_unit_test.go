package ktn_const_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	ktn_const "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/const"
	"golang.org/x/tools/go/analysis"
)

// TestRule001_DirectCoverage tests Rule001 directly for coverage improvement.
//
// Params:
//   - t: testing instance
func TestRule001_DirectCoverage(t *testing.T) {
	src := `package test

// Constante ungrouped
const UngroupedConst = 42

// Autre constante ungrouped
const AnotherUngrouped = "test"
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

	_, err = ktn_const.Rule001.Run(pass)
	if err != nil {
		t.Errorf("Rule001.Run() error = %v", err)
	}
}
