package ktn_goroutine

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Rule001 vérifie que les goroutines lancées dans des boucles ont une synchronisation appropriée.
// KTN-GOROUTINE-001: Goroutines illimitées dans une boucle causent des fuites mémoire.
var Rule001 = &analysis.Analyzer{
	Name: "KTN_GOROUTINE_001",
	Doc:  "Détecter les goroutines lancées dans une boucle sans limitation",
	Run:  runRule001,
}

// runRule001 exécute l'analyse pour la règle GOROUTINE-001.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok || funcDecl.Body == nil {
				return true
			}

			checkGoroutinesInFunction(pass, funcDecl)
			return true
		})
	}
	return nil, nil
}

// checkGoroutinesInFunction vérifie les goroutines dans une fonction.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
func checkGoroutinesInFunction(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Tracker les WaitGroups dans la fonction
	hasWaitGroup := detectWaitGroup(funcDecl.Body)

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		goStmt, ok := n.(*ast.GoStmt)
		if !ok {
			return true
		}

		// KTN-GOROUTINE-001: Vérifier si go est dans une boucle SANS synchronisation
		// Si WaitGroup présent, c'est probablement un worker pool (pattern correct)
		if isInsideLoop(funcDecl.Body, goStmt) && !hasWaitGroup {
			reportGoroutineInLoop(pass, goStmt)
		}

		return true
	})
}

// isInsideLoop vérifie si un GoStmt est à l'intérieur d'une boucle.
//
// Params:
//   - body: le corps de la fonction
//   - goStmt: le go statement
//
// Returns:
//   - bool: true si inside loop
func isInsideLoop(body *ast.BlockStmt, goStmt *ast.GoStmt) bool {
	insideLoop := false

	ast.Inspect(body, func(n ast.Node) bool {
		switch loop := n.(type) {
		case *ast.ForStmt:
			if containsGoStmt(loop.Body, goStmt) {
				insideLoop = true
				return false
			}
		case *ast.RangeStmt:
			if containsGoStmt(loop.Body, goStmt) {
				insideLoop = true
				return false
			}
		}
		return true
	})

	return insideLoop
}

// containsGoStmt vérifie si un bloc contient un GoStmt spécifique.
//
// Params:
//   - block: le bloc à vérifier
//   - target: le GoStmt recherché
//
// Returns:
//   - bool: true si trouvé
func containsGoStmt(block *ast.BlockStmt, target *ast.GoStmt) bool {
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

// detectWaitGroup détecte la présence d'un sync.WaitGroup dans le corps.
//
// Params:
//   - body: le corps de la fonction
//
// Returns:
//   - bool: true si WaitGroup détecté
func detectWaitGroup(body *ast.BlockStmt) bool {
	hasWaitGroup := false

	ast.Inspect(body, func(n ast.Node) bool {
		switch decl := n.(type) {
		case *ast.ValueSpec:
			if isWaitGroupType(decl.Type) {
				hasWaitGroup = true
				return false
			}
		case *ast.CallExpr:
			if isWaitGroupMethod(decl.Fun) {
				hasWaitGroup = true
				return false
			}
		}
		return true
	})

	return hasWaitGroup
}

// isWaitGroupType vérifie si un type est sync.WaitGroup.
//
// Params:
//   - typeExpr: l'expression de type
//
// Returns:
//   - bool: true si sync.WaitGroup
func isWaitGroupType(typeExpr ast.Expr) bool {
	if typeExpr == nil {
		return false
	}

	sel, ok := typeExpr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	x, ok := sel.X.(*ast.Ident)
	return ok && x.Name == "sync" && sel.Sel.Name == "WaitGroup"
}

// isWaitGroupMethod vérifie si un appel est Add/Wait.
//
// Params:
//   - fun: l'expression de fonction
//
// Returns:
//   - bool: true si Add ou Wait
func isWaitGroupMethod(fun ast.Expr) bool {
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	return sel.Sel.Name == "Add" || sel.Sel.Name == "Wait"
}

// reportGoroutineInLoop rapporte une violation KTN-GOROUTINE-001.
//
// Params:
//   - pass: la passe d'analyse
//   - goStmt: le go statement
func reportGoroutineInLoop(pass *analysis.Pass, goStmt *ast.GoStmt) {
	pass.Reportf(goStmt.Pos(),
		"[KTN-GOROUTINE-001] Goroutine lancée dans une boucle sans limitation.\n"+
			"Créer des goroutines illimitées cause des fuites mémoire et ralentissements.\n"+
			"Utilisez un worker pool avec buffered channel pour limiter la concurrence.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - goroutine leak\n"+
			"  for _, req := range requests {\n"+
			"      go handleRequest(req)  // Crée N goroutines\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - worker pool\n"+
			"  jobs := make(chan Request, 100)\n"+
			"  for i := 0; i < 10; i++ {\n"+
			"      go worker(jobs)  // Seulement 10 workers\n"+
			"  }\n"+
			"  for _, req := range requests {\n"+
			"      jobs <- req\n"+
			"  }")
}
