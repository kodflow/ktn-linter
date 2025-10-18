package ktn_control_flow_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	ktn_control_flow "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/control_flow"
	"golang.org/x/tools/go/analysis"
)

// TestRuleIf001_NoElse teste le cas sans else (pas simplifiable)
func TestRuleIf001_NoElse(t *testing.T) {
	src := `package test

func test(x bool) bool {
	if x {
		return true
	}
	return false
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
		Report: func(_ analysis.Diagnostic) {
			reported = true
		},
	}

	_, err = ktn_control_flow.RuleIf001.Run(pass)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	if reported {
		t.Error("if sans else ne devrait PAS être signalé")
	}
}

// TestRuleIf001_ElseIfChain teste le cas avec else if
func TestRuleIf001_ElseIfChain(t *testing.T) {
	src := `package test

func test(x bool) bool {
	if x {
		return true
	} else if !x {
		return false
	}
	return false
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
		Report: func(_ analysis.Diagnostic) {
			reported = true
		},
	}

	_, err = ktn_control_flow.RuleIf001.Run(pass)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	// else if n'est PAS un *ast.BlockStmt donc ne devrait pas être détecté
	if reported {
		t.Error("else if ne devrait PAS être signalé")
	}
}

// TestRuleIf001_MultipleStatementsInElse teste else avec plusieurs statements
func TestRuleIf001_MultipleStatementsInElse(t *testing.T) {
	src := `package test

func test(x bool) bool {
	if x {
		return true
	} else {
		println("debug")
		return false
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
		Report: func(_ analysis.Diagnostic) {
			reported = true
		},
	}

	_, err = ktn_control_flow.RuleIf001.Run(pass)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	if reported {
		t.Error("else avec plusieurs statements ne devrait PAS être signalé")
	}
}
