package ktn_goroutine

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Rule002 vérifie que les goroutines ont un mécanisme de synchronisation.
// KTN-GOROUTINE-002: Sans sync.WaitGroup, context.Context ou channel, la goroutine peut leak.
var Rule002 = &analysis.Analyzer{
	Name: "KTN_GOROUTINE_002",
	Doc:  "Détecter les goroutines sans mécanisme de synchronisation",
	Run:  runRule002,
}

// runRule002 exécute l'analyse pour la règle GOROUTINE-002.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule002(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok || funcDecl.Body == nil {
				return true
			}

			checkGoroutineSynchronization(pass, funcDecl)
			return true
		})
	}
	return nil, nil
}

// checkGoroutineSynchronization vérifie les mécanismes de synchronisation des goroutines.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
func checkGoroutineSynchronization(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Tracker les WaitGroups et contexts dans la fonction
	hasWaitGroup := detectWaitGroupInFunc(funcDecl.Body)
	hasContext := hasContextParam(funcDecl)

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		goStmt, ok := n.(*ast.GoStmt)
		if !ok {
			return true
		}

		// KTN-GOROUTINE-002: Vérifier si go a un mécanisme de synchronisation
		if !hasWaitGroup && !hasContext && !hasChannelSync(goStmt) {
			reportGoroutineWithoutSync(pass, goStmt)
		}

		return true
	})
}

// detectWaitGroupInFunc détecte la présence d'un sync.WaitGroup dans le corps.
//
// Params:
//   - body: le corps de la fonction
//
// Returns:
//   - bool: true si WaitGroup détecté
func detectWaitGroupInFunc(body *ast.BlockStmt) bool {
	hasWaitGroup := false

	ast.Inspect(body, func(n ast.Node) bool {
		switch decl := n.(type) {
		case *ast.ValueSpec:
			if isWaitGroupTypeCheck(decl.Type) {
				hasWaitGroup = true
				return false
			}
		case *ast.CallExpr:
			if isWaitGroupMethodCheck(decl.Fun) {
				hasWaitGroup = true
				return false
			}
		}
		return true
	})

	return hasWaitGroup
}

// isWaitGroupTypeCheck vérifie si un type est sync.WaitGroup.
//
// Params:
//   - typeExpr: l'expression de type
//
// Returns:
//   - bool: true si sync.WaitGroup
func isWaitGroupTypeCheck(typeExpr ast.Expr) bool {
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

// isWaitGroupMethodCheck vérifie si un appel est Add/Wait.
//
// Params:
//   - fun: l'expression de fonction
//
// Returns:
//   - bool: true si Add ou Wait
func isWaitGroupMethodCheck(fun ast.Expr) bool {
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	return sel.Sel.Name == "Add" || sel.Sel.Name == "Wait"
}

// hasContextParam vérifie si une fonction a un paramètre context.Context.
//
// Params:
//   - funcDecl: la déclaration de fonction
//
// Returns:
//   - bool: true si context trouvé
func hasContextParam(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Type == nil || funcDecl.Type.Params == nil {
		return false
	}

	for _, param := range funcDecl.Type.Params.List {
		if sel, ok := param.Type.(*ast.SelectorExpr); ok {
			if x, ok := sel.X.(*ast.Ident); ok && x.Name == "context" && sel.Sel.Name == "Context" {
				return true
			}
		}
	}

	return false
}

// hasChannelSync vérifie si un GoStmt utilise un channel pour synchronisation.
//
// Params:
//   - goStmt: le go statement
//
// Returns:
//   - bool: true si channel utilisé
func hasChannelSync(goStmt *ast.GoStmt) bool {
	funcBody := extractGoroutineFuncBody(goStmt)
	if funcBody == nil {
		return false
	}

	return inspectForChannelOps(funcBody)
}

// extractGoroutineFuncBody extrait le corps de la fonction d'un GoStmt.
//
// Params:
//   - goStmt: le go statement
//
// Returns:
//   - *ast.BlockStmt: le corps de la fonction ou nil
func extractGoroutineFuncBody(goStmt *ast.GoStmt) *ast.BlockStmt {
	funcLit, ok := goStmt.Call.Fun.(*ast.FuncLit)
	if !ok {
		return nil
	}

	return funcLit.Body
}

// inspectForChannelOps inspecte un corps de fonction pour détecter des opérations channel.
//
// Params:
//   - body: le corps de fonction
//
// Returns:
//   - bool: true si opération channel trouvée
func inspectForChannelOps(body *ast.BlockStmt) bool {
	hasChannel := false

	ast.Inspect(body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.SendStmt:
			hasChannel = true
			return false
		case *ast.UnaryExpr:
			if stmt.Op.String() == "<-" {
				hasChannel = true
				return false
			}
		case *ast.SelectStmt:
			hasChannel = true
			return false
		}
		return true
	})

	return hasChannel
}

// reportGoroutineWithoutSync rapporte une violation KTN-GOROUTINE-002.
//
// Params:
//   - pass: la passe d'analyse
//   - goStmt: le go statement
func reportGoroutineWithoutSync(pass *analysis.Pass, goStmt *ast.GoStmt) {
	pass.Reportf(goStmt.Pos(),
		"[KTN-GOROUTINE-002] Goroutine lancée sans mécanisme de synchronisation.\n"+
			"Sans sync.WaitGroup, context.Context, ou channel, la goroutine peut leak.\n"+
			"La fonction parent peut se terminer avant la goroutine, causant perte de données.\n"+
			"Utilisez l'un de ces mécanismes de synchronisation.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - pas de synchronisation\n"+
			"  go func() {\n"+
			"      process(data)\n"+
			"  }()\n"+
			"  return  // goroutine peut ne jamais finir\n"+
			"\n"+
			"  // ✅ CORRECT - avec WaitGroup\n"+
			"  var wg sync.WaitGroup\n"+
			"  wg.Add(1)\n"+
			"  go func() {\n"+
			"      defer wg.Done()\n"+
			"      process(data)\n"+
			"  }()\n"+
			"  wg.Wait()  // Attend la fin\n"+
			"\n"+
			"  // ✅ CORRECT - avec context\n"+
			"  func Process(ctx context.Context) {\n"+
			"      go func() {\n"+
			"          select {\n"+
			"          case <-ctx.Done():\n"+
			"              return\n"+
			"          case result <- process(data):\n"+
			"          }\n"+
			"      }()\n"+
			"  }")
}
