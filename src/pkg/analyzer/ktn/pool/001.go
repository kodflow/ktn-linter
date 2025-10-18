package ktn_pool

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Rule001 vérifie que pool.Get() est suivi d'un defer pool.Put().
// KTN-POOL-001: Sans defer Put(), l'objet ne retourne jamais au pool, causant une fuite de ressources.
var Rule001 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_POOL_001",
	Doc:  "Vérifier que pool.Get() est suivi de defer pool.Put()",
	Run:  runRule001,
}

// runRule001 exécute l'analyse pour la règle POOL-001.
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
				// Continue traversing AST nodes.
				return true
			}

			checkPoolGetWithoutDeferPut(pass, funcDecl)
			// Continue traversing AST nodes.
			return true
		})
	}
	// Analysis completed successfully.
	return nil, nil
}

// checkPoolGetWithoutDeferPut vérifie si Get() est utilisé sans defer Put().
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction à analyser
func checkPoolGetWithoutDeferPut(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Map pour tracker les variables issues de pool.Get()
	poolVars := make(map[string]ast.Expr) // varName -> pool expression
	deferredPuts := make(map[string]bool) // varName -> has defer Put

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			TrackPoolGetAssignment(stmt, poolVars, pass)
		case *ast.DeferStmt:
			TrackDeferPut(stmt, deferredPuts)
		}
		// Continue traversing AST nodes.
		return true
	})

	// Vérifier les variables sans defer Put()
	for varName, poolExpr := range poolVars {
		if !deferredPuts[varName] {
			reportMissingDeferPut(pass, poolExpr, varName)
		}
	}
}

// TrackPoolGetAssignment détecte les assignations depuis pool.Get().
//
// Params:
//   - stmt: l'assignation à analyser
//   - poolVars: map des variables pool trackées
//   - pass: la passe d'analyse
func TrackPoolGetAssignment(stmt *ast.AssignStmt, poolVars map[string]ast.Expr, pass *analysis.Pass) {
	if len(stmt.Lhs) != 1 || len(stmt.Rhs) != 1 {
		// Early return from function.
		return
	}

	// Extraire l'expression sous-jacente (peut être un type assertion)
	rhsExpr := UnwrapTypeAssertion(stmt.Rhs[0])

	// Vérifier si RHS est un appel à pool.Get()
	if !IsPoolGetCall(rhsExpr, pass) {
		// Early return from function.
		return
	}

	// Extraire le nom de la variable
	varName := ExtractVarName(stmt.Lhs[0])
	if varName != "" {
		poolVars[varName] = rhsExpr
	}
}

// UnwrapTypeAssertion extrait l'expression d'un type assertion.
//
// Params:
//   - expr: l'expression (peut être TypeAssertExpr)
//
// Returns:
//   - ast.Expr: l'expression sous-jacente
func UnwrapTypeAssertion(expr ast.Expr) ast.Expr {
	if typeAssert, ok := expr.(*ast.TypeAssertExpr); ok {
		// Early return from function.
		return typeAssert.X
	}
	// Early return from function.
	return expr
}

// TrackDeferPut détecte les defer avec pool.Put().
//
// Params:
//   - stmt: le defer statement
//   - deferredPuts: map des Put() différés
func TrackDeferPut(stmt *ast.DeferStmt, deferredPuts map[string]bool) {
	callExpr := stmt.Call
	if callExpr == nil {
		// Early return from function.
		return
	}

	// Vérifier si c'est pool.Put(var)
	if !IsPoolPutCall(callExpr) {
		// Early return from function.
		return
	}

	// Extraire le nom de la variable passée à Put()
	if len(callExpr.Args) > 0 {
		varName := ExtractVarName(callExpr.Args[0])
		if varName != "" {
			deferredPuts[varName] = true
		}
	}
}

// IsPoolGetCall vérifie si l'expression est un appel à pool.Get().
//
// Params:
//   - expr: l'expression à vérifier
//   - pass: la passe d'analyse
//
// Returns:
//   - bool: true si c'est pool.Get()
func IsPoolGetCall(expr ast.Expr, pass *analysis.Pass) bool {
	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		// Condition not met, return false.
		return false
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok || selExpr.Sel.Name != "Get" {
		// Condition not met, return false.
		return false
	}

	// Vérifier le type de l'objet appelant
	if pass.TypesInfo != nil {
		if t := pass.TypesInfo.TypeOf(selExpr.X); t != nil {
			// Vérifier si c'est sync.Pool ou *sync.Pool
			typeStr := t.String()
			if strings.Contains(typeStr, "sync.Pool") {
				// Continue traversing AST nodes.
				return true
			}
		}
	}

	// Fallback: vérifier le nom de la variable
	if ident, ok := selExpr.X.(*ast.Ident); ok {
		varName := strings.ToLower(ident.Name)
		// Early return from function.
		return strings.Contains(varName, "pool")
	}

	// Condition not met, return false.
	return false
}

// IsPoolPutCall vérifie si l'appel est pool.Put().
//
// Params:
//   - callExpr: l'appel à vérifier
//
// Returns:
//   - bool: true si c'est pool.Put()
func IsPoolPutCall(callExpr *ast.CallExpr) bool {
	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok || selExpr.Sel.Name != "Put" {
		// Condition not met, return false.
		return false
	}

	// Vérifier le nom de la variable
	if ident, ok := selExpr.X.(*ast.Ident); ok {
		varName := strings.ToLower(ident.Name)
		// Early return from function.
		return strings.Contains(varName, "pool")
	}

	// Condition not met, return false.
	return false
}

// ExtractVarName extrait le nom de variable d'une expression.
//
// Params:
//   - expr: l'expression
//
// Returns:
//   - string: le nom de la variable ou ""
func ExtractVarName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		// Early return from function.
		return e.Name
	case *ast.CallExpr:
		// Type assertion: buf := pool.Get().([]byte)
		return ""
	}
	// Early return from function.
	return ""
}

// reportMissingDeferPut rapporte une violation KTN-POOL-001.
//
// Params:
//   - pass: la passe d'analyse
//   - expr: l'expression pool.Get()
//   - varName: le nom de la variable
func reportMissingDeferPut(pass *analysis.Pass, expr ast.Expr, varName string) {
	pass.Reportf(expr.Pos(),
		"[KTN-POOL-001] Variable '%s' obtenue via pool.Get() sans defer pool.Put().\n"+
			"Cela cause une fuite de ressources car l'objet ne retourne jamais au pool.\n"+
			"Utilisez 'defer pool.Put(%s)' immédiatement après Get().\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - fuite mémoire\n"+
			"  buf := bufferPool.Get().([]byte)\n"+
			"  process(buf)\n"+
			"  // buf n'est jamais retourné au pool\n"+
			"\n"+
			"  // ✅ CORRECT\n"+
			"  buf := bufferPool.Get().([]byte)\n"+
			"  defer bufferPool.Put(buf)\n"+
			"  process(buf)",
		varName, varName)
}
