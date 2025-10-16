package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// GoroutineAnalyzer vérifie l'utilisation sécurisée des goroutines.
	GoroutineAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktngoroutine",
		Doc:  "Vérifie que les goroutines sont lancées de manière contrôlée avec synchronisation",
		Run:  runGoroutineAnalyzer,
	}
)

// runGoroutineAnalyzer exécute l'analyseur goroutine.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runGoroutineAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		checkGoroutineUsage(pass, file)
	}
	// Retourne nil car l'analyse est terminée
	return nil, nil
}

// checkGoroutineUsage vérifie l'utilisation des goroutines.
//
// Params:
//   - pass: la passe d'analyse
//   - file: le fichier à analyser
func checkGoroutineUsage(pass *analysis.Pass, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			// Retourne true pour continuer l'inspection
			return true
		}

		checkGoroutinesInFunction(pass, funcDecl)
		// Retourne true pour continuer l'inspection
		return true
	})
}

// checkGoroutinesInFunction vérifie les goroutines dans une fonction.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
func checkGoroutinesInFunction(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Tracker les WaitGroups et contexts dans la fonction
	hasWaitGroup := detectWaitGroup(funcDecl.Body)

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		goStmt, ok := n.(*ast.GoStmt)
		if !ok {
			// Retourne true pour continuer l'inspection
			return true
		}

		// KTN-GOROUTINE-001: Vérifier si go est dans une boucle SANS synchronisation
		// Si WaitGroup présent, c'est probablement un worker pool (pattern correct)
		if isInsideLoop(funcDecl.Body, goStmt) && !hasWaitGroup {
			reportGoroutineInLoop(pass, goStmt)
		}

		// KTN-GOROUTINE-002: Vérifier si go a un mécanisme de synchronisation
		if !hasWaitGroup && !hasContextParam(funcDecl) && !hasChannelSync(goStmt) {
			reportGoroutineWithoutSync(pass, goStmt)
		}

		// Retourne true pour continuer l'inspection
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
				// Retourne false pour arrêter l'inspection
				return false
			}
		case *ast.RangeStmt:
			if containsGoStmt(loop.Body, goStmt) {
				insideLoop = true
				// Retourne false pour arrêter l'inspection
				return false
			}
		}
		// Retourne true pour continuer l'inspection
		return true
	})

	// Retourne le résultat
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
			// Retourne false pour arrêter l'inspection
			return false
		}
		// Retourne true pour continuer l'inspection
		return true
	})

	// Retourne le résultat
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
				// Retourne false pour arrêter l'inspection
				return false
			}
		case *ast.CallExpr:
			if isWaitGroupMethod(decl.Fun) {
				hasWaitGroup = true
				// Retourne false pour arrêter l'inspection
				return false
			}
		}
		// Retourne true pour continuer l'inspection
		return true
	})

	// Retourne le résultat
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
		// Retourne false car pas de type
		return false
	}

	sel, ok := typeExpr.(*ast.SelectorExpr)
	if !ok {
		// Retourne false car pas un sélecteur
		return false
	}

	x, ok := sel.X.(*ast.Ident)
	// Retourne true si sync.WaitGroup
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
		// Retourne false car pas un sélecteur
		return false
	}

	// Retourne true si Add ou Wait
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
		// Retourne false car pas de paramètres
		return false
	}

	for _, param := range funcDecl.Type.Params.List {
		if sel, ok := param.Type.(*ast.SelectorExpr); ok {
			if x, ok := sel.X.(*ast.Ident); ok && x.Name == "context" && sel.Sel.Name == "Context" {
				// Retourne true car context trouvé
				return true
			}
		}
	}

	// Retourne false car context non trouvé
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
		// Retourne false car pas de corps
		return false
	}

	// Retourne le résultat de l'inspection des opérations channel
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
		// Retourne nil car ce n'est pas une func lit
		return nil
	}

	// Retourne le corps de la fonction
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
			// Retourne false pour arrêter l'inspection
			return false
		case *ast.UnaryExpr:
			if stmt.Op.String() == "<-" {
				hasChannel = true
				// Retourne false pour arrêter l'inspection
				return false
			}
		case *ast.SelectStmt:
			hasChannel = true
			// Retourne false pour arrêter l'inspection
			return false
		}
		// Retourne true pour continuer l'inspection
		return true
	})

	// Retourne le résultat
	return hasChannel
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
