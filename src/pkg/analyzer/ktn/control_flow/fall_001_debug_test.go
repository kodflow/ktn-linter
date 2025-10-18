package ktn_control_flow_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// TestDebug_IsInsideSwitchCase analyse comment isInsideSwitchCase fonctionne
func TestDebug_IsInsideSwitchCase(t *testing.T) {
	src := `package test

func test() {
	x := 1
	switch x {
	case 1:
		println("one")
		fallthrough
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

	// Compter les fallthrough statements
	fallthroughCount := 0
	var fallthroughStmt *ast.BranchStmt

	ast.Inspect(file, func(n ast.Node) bool {
		if branch, ok := n.(*ast.BranchStmt); ok && branch.Tok == token.FALLTHROUGH {
			fallthroughCount++
			fallthroughStmt = branch
			t.Logf("Trouvé fallthrough à la position %v", fset.Position(branch.Pos()))
		}
		return true
	})

	if fallthroughCount != 1 {
		t.Errorf("Attendu 1 fallthrough, obtenu %d", fallthroughCount)
	}

	// Vérifier si on peut le trouver dans le switch
	found := false
	ast.Inspect(file, func(n ast.Node) bool {
		if switchStmt, ok := n.(*ast.SwitchStmt); ok {
			t.Logf("Switch trouvé avec %d cases", len(switchStmt.Body.List))
			for i, stmt := range switchStmt.Body.List {
				if caseClause, ok := stmt.(*ast.CaseClause); ok {
					t.Logf("  Case %d a %d statements", i, len(caseClause.Body))
					for j, s := range caseClause.Body {
						t.Logf("    Statement %d: %T", j, s)
						if s == fallthroughStmt {
							found = true
							t.Logf("    -> C'EST LE FALLTHROUGH!")
						}
					}
				}
			}
		}
		return true
	})

	if !found {
		t.Error("fallthrough non trouvé dans le case body")
	}
}
