package ktn_ops

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var RuleChan001 = &analysis.Analyzer{
	Name: "KTN_CHAN_002",
	Doc:  "Détecte close() appelé par le receiver",
	Run:  runRuleChan001,
}

func runRuleChan001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Vérifier si c'est un appel à close()
			ident, ok := callExpr.Fun.(*ast.Ident)
			if !ok || ident.Name != "close" || len(callExpr.Args) == 0 {
				return true
			}

			// Trouver la fonction englobante
			funcDecl := findEnclosingFunc(file, callExpr)
			if funcDecl == nil || funcDecl.Body == nil {
				return true
			}

			// Vérifier si la fonction reçoit du channel
			if functionReceivesFromChannel(funcDecl, callExpr.Args[0]) {
				pass.Reportf(callExpr.Pos(),
					"[KTN-CHAN-001] close() appelé par le receiver du channel.\n"+
						"C'est une mauvaise pratique: seul le sender devrait fermer un channel.\n"+
						"Fermer côté receiver peut causer des panics si le sender écrit encore.\n"+
						"Pattern: sender ferme, receiver détecte la fermeture.\n"+
						"Exemple:\n"+
						"  // ❌ MAUVAIS - receiver ferme\n"+
						"  func receive(ch chan int) {\n"+
						"      for v := range ch { process(v) }\n"+
						"      close(ch)  // 💥 Risque de panic\n"+
						"  }\n"+
						"\n"+
						"  // ✅ CORRECT - sender ferme\n"+
						"  func send(ch chan int) {\n"+
						"      for i := 0; i < 10; i++ { ch <- i }\n"+
						"      close(ch)  // ✅ Sender ferme\n"+
						"  }\n"+
						"  func receive(ch chan int) {\n"+
						"      for v := range ch { process(v) }  // ✅ Détecte fermeture\n"+
						"  }")
			}
			return true
		})
	}
	return nil, nil
}

func findEnclosingFunc(file *ast.File, target ast.Node) *ast.FuncDecl {
	var enclosing *ast.FuncDecl
	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if funcDecl.Body != nil && containsNodeInBlock(funcDecl.Body, target) {
				enclosing = funcDecl
				return false
			}
		}
		return true
	})
	return enclosing
}

func containsNodeInBlock(block *ast.BlockStmt, target ast.Node) bool {
	if block == nil {
		return false
	}
	found := false
	ast.Inspect(block, func(n ast.Node) bool {
		if n == target {
			found = true
			return false
		}
		return true
	})
	return found
}

func functionReceivesFromChannel(funcDecl *ast.FuncDecl, chanExpr ast.Expr) bool {
	chanName := extractChannelName(chanExpr)
	if chanName == "" {
		return false
	}

	receives := false
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		// Chercher <-chan
		if unaryExpr, ok := n.(*ast.UnaryExpr); ok {
			if unaryExpr.Op.String() == "<-" {
				if extractChannelName(unaryExpr.X) == chanName {
					receives = true
					return false
				}
			}
		}
		// Chercher range sur channel
		if rangeStmt, ok := n.(*ast.RangeStmt); ok {
			if extractChannelName(rangeStmt.X) == chanName {
				receives = true
				return false
			}
		}
		return true
	})
	return receives
}

func extractChannelName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		if x, ok := e.X.(*ast.Ident); ok {
			return x.Name + "." + e.Sel.Name
		}
	}
	return ""
}
