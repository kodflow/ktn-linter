package ktn_control_flow_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	ktn_control_flow "github.com/kodflow/ktn-linter/src/pkg/analyzer/ktn/control_flow"
)

// TestFindCopiedVars_MismatchedLengths teste le cas où len(Lhs) != len(Rhs)
func TestFindCopiedVars_MismatchedLengths(t *testing.T) {
	src := `package test
func test() {
	items := []int{1, 2, 3}
	for _, item := range items {
		item := item
		a, b := 1, 2, 3  // len(Lhs) = 2, len(Rhs) = 3
		_ = a
		_ = b
		go func() { println(item) }()
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	// Trouver le RangeStmt
	var rangeStmt *ast.RangeStmt
	var funcLit *ast.FuncLit
	ast.Inspect(file, func(n ast.Node) bool {
		if r, ok := n.(*ast.RangeStmt); ok {
			rangeStmt = r
		}
		if f, ok := n.(*ast.FuncLit); ok {
			funcLit = f
		}
		return true
	})

	if rangeStmt == nil || funcLit == nil {
		t.Fatal("RangeStmt ou FuncLit non trouvé")
	}

	rangeVars := []string{"item"}
	copiedVars := ktn_control_flow.FindCopiedVarsExported(rangeStmt.Body, funcLit, rangeVars)

	// Vérifier que "item" est bien détecté comme copié malgré l'assignation mal formée
	if !copiedVars["item"] {
		t.Errorf("item devrait être détecté comme copié")
	}
}

// TestFindCopiedVars_MultipleVars teste la copie de plusieurs variables
func TestFindCopiedVars_MultipleVars(t *testing.T) {
	src := `package test
func test() {
	items := []int{1, 2, 3}
	for i, item := range items {
		i := i
		item := item
		go func() { println(i, item) }()
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	var rangeStmt *ast.RangeStmt
	var funcLit *ast.FuncLit
	ast.Inspect(file, func(n ast.Node) bool {
		if r, ok := n.(*ast.RangeStmt); ok {
			rangeStmt = r
		}
		if f, ok := n.(*ast.FuncLit); ok {
			funcLit = f
		}
		return true
	})

	if rangeStmt == nil || funcLit == nil {
		t.Fatal("RangeStmt ou FuncLit non trouvé")
	}

	rangeVars := []string{"i", "item"}
	copiedVars := ktn_control_flow.FindCopiedVarsExported(rangeStmt.Body, funcLit, rangeVars)

	if !copiedVars["i"] {
		t.Errorf("i devrait être détecté comme copié")
	}
	if !copiedVars["item"] {
		t.Errorf("item devrait être détecté comme copié")
	}
}

// TestFindCopiedVars_Assignment teste le cas avec = au lieu de :=
func TestFindCopiedVars_Assignment(t *testing.T) {
	src := `package test
func test() {
	items := []int{1, 2, 3}
	var copied int
	for _, item := range items {
		copied = item  // = au lieu de :=, ne devrait PAS être détecté comme copie
		go func() { println(copied) }()
	}
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	var rangeStmt *ast.RangeStmt
	var funcLit *ast.FuncLit
	ast.Inspect(file, func(n ast.Node) bool {
		if r, ok := n.(*ast.RangeStmt); ok {
			rangeStmt = r
		}
		if f, ok := n.(*ast.FuncLit); ok {
			funcLit = f
		}
		return true
	})

	if rangeStmt == nil || funcLit == nil {
		t.Fatal("RangeStmt ou FuncLit non trouvé")
	}

	rangeVars := []string{"item"}
	copiedVars := ktn_control_flow.FindCopiedVarsExported(rangeStmt.Body, funcLit, rangeVars)

	// = n'est pas une déclaration, donc ne devrait pas être détecté
	if copiedVars["item"] {
		t.Errorf("item avec = ne devrait PAS être détecté comme copié")
	}
}

// TestFindCopiedVars_NonIdent teste le cas où lhs/rhs ne sont pas des Ident
func TestFindCopiedVars_NonIdent(t *testing.T) {
	src := `package test
type S struct{ v int }
func test() {
	items := []int{1, 2, 3}
	s := S{}
	for _, item := range items {
		s.v := item  // lhs n'est pas un Ident
		go func() { println(s.v) }()
	}
}
`
	// Ce code ne compile pas mais testons la robustesse
	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", src, parser.AllErrors)
	if file == nil {
		t.Skip("Code invalide")
	}

	var rangeStmt *ast.RangeStmt
	var funcLit *ast.FuncLit
	ast.Inspect(file, func(n ast.Node) bool {
		if r, ok := n.(*ast.RangeStmt); ok {
			rangeStmt = r
		}
		if f, ok := n.(*ast.FuncLit); ok {
			funcLit = f
		}
		return true
	})

	if rangeStmt == nil || funcLit == nil {
		t.Skip("RangeStmt ou FuncLit non trouvé")
	}

	rangeVars := []string{"item"}
	copiedVars := ktn_control_flow.FindCopiedVarsExported(rangeStmt.Body, funcLit, rangeVars)

	// s.v n'est pas un Ident, donc ne devrait pas crash
	if copiedVars["item"] {
		t.Logf("item détecté (inattendu mais OK)")
	}
}
