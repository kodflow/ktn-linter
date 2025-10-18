package ktn_control_flow_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	ktn_control_flow "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/control_flow"
	"golang.org/x/tools/go/analysis"
)

// TestRule Fall001_FallthroughInsideSwitch teste le cas où fallthrough EST dans un switch (pas d'erreur)
// TestRuleFall001_FallthroughInsideSwitch tests the functionality of the corresponding implementation.
func TestRuleFall001_FallthroughInsideSwitch(t *testing.T) {
	src := `package test

func test() {
	x := 1
	switch x {
	case 1:
		println("one")
		fallthrough  // OK - dans un switch
	case 2:
		println("two")
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	reported := false
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(diag analysis.Diagnostic) {
			reported = true
			t.Errorf("Ne devrait PAS rapporter d'erreur pour fallthrough dans switch: %s", diag.Message)
		},
	}

	_, err = ktn_control_flow.RuleFall001.Run(pass)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	if reported {
		t.Error("fallthrough dans switch ne devrait PAS être signalé")
	}
}

// TestRuleFall001_MultipleFallthrough teste plusieurs fallthrough dans un switch
func TestRuleFall001_MultipleFallthrough(t *testing.T) {
	src := `package test

func test() {
	x := 1
	switch x {
	case 1:
		fallthrough
	case 2:
		fallthrough
	case 3:
		println("done")
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	reported := 0
	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(_ analysis.Diagnostic) {
			reported++
		},
	}

	_, err = ktn_control_flow.RuleFall001.Run(pass)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	if reported > 0 {
		t.Errorf("Ne devrait PAS rapporter d'erreurs, obtenu %d erreurs", reported)
	}
}
