package ktn_control_flow_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	ktn_control_flow "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/control_flow"
	"golang.org/x/tools/go/analysis"
)

// TestRuleRange001_UnderscoredVariables teste le cas où la variable range est _
func TestRuleRange001_UnderscoredVariables(t *testing.T) {
	src := `package test

func test() {
	items := []int{1, 2, 3}
	for _ = range items {
		go func() {
			println("no capture")
		}()
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
		Report: func(diag analysis.Diagnostic) {
			reported++
		},
	}

	_, err = ktn_control_flow.RuleRange001.Run(pass)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	// _ ne devrait pas être signalé comme capturé
	if reported > 0 {
		t.Errorf("_ ne devrait PAS être signalé, obtenu %d erreurs", reported)
	}
}

// TestRuleRange001_NoBody teste le cas où range n'a pas de body
func TestRuleRange001_NoBody(t *testing.T) {
	// Ce cas ne compile pas en Go, mais testons la robustesse
	src := `package test

func test() {
	items := []int{1, 2, 3}
	for range items {
		// Pas de closure
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
		Report: func(diag analysis.Diagnostic) {
			reported++
		},
	}

	_, err = ktn_control_flow.RuleRange001.Run(pass)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	if reported > 0 {
		t.Errorf("Pas de closure = pas d'erreur, obtenu %d erreurs", reported)
	}
}

// TestRuleRange001_NonIdentKey teste le cas où la clé n'est pas un Ident
func TestRuleRange001_NonIdentKey(t *testing.T) {
	src := `package test

type S struct {
	i int
}

func test() {
	items := []int{1, 2, 3}
	s := S{}
	for s.i = range items {
		// s.i n'est pas un *ast.Ident
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
		Report: func(diag analysis.Diagnostic) {
			reported++
		},
	}

	_, err = ktn_control_flow.RuleRange001.Run(pass)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	// Ne devrait pas crash, juste ignorer
	if reported > 0 {
		t.Logf("reported %d errors (attendu 0)", reported)
	}
}

// TestRuleRange001_ClosureDoesNotUseVariable teste closure qui n'utilise pas la variable
func TestRuleRange001_ClosureDoesNotUseVariable(t *testing.T) {
	src := `package test

func test() {
	items := []int{1, 2, 3}
	for _, item := range items {
		go func() {
			println("no use of item")
		}()
		_ = item // Utilise item hors closure
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
		Report: func(diag analysis.Diagnostic) {
			reported++
		},
	}

	_, err = ktn_control_flow.RuleRange001.Run(pass)
	if err != nil {
		t.Errorf("Run() error = %v", err)
	}

	// Closure n'utilise pas item = pas d'erreur
	if reported > 0 {
		t.Errorf("Closure sans utilisation ne devrait PAS être signalée, obtenu %d erreurs", reported)
	}
}
