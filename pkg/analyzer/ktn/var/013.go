package ktnvar

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer013 checks for slice/map allocations inside loops
var Analyzer013 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar013",
	Doc:      "KTN-VAR-013: Évite les allocations de slices/maps dans les boucles chaudes",
	Run:      runVar013,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar013 exécute l'analyse KTN-VAR-013.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar013(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ForStmt)(nil),
		(*ast.RangeStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Récupération du corps de la boucle
		var body *ast.BlockStmt
		// Vérification du type de boucle
		switch loop := n.(type) {
		// Cas d'une boucle for classique
		case *ast.ForStmt:
			// Boucle for classique
			body = loop.Body
		// Cas d'une boucle range
		case *ast.RangeStmt:
			// Boucle range
			body = loop.Body
		// Cas par défaut
		default:
			// Type de boucle non supporté
			return
		}

		// Parcours des instructions du corps
		checkLoopBodyForAlloc(pass, body)
	})

	// Retour de la fonction
	return nil, nil
}

// checkLoopBodyForAlloc vérifie les allocations dans le corps d'une boucle.
//
// Params:
//   - pass: contexte d'analyse
//   - body: corps de la boucle à vérifier
func checkLoopBodyForAlloc(pass *analysis.Pass, body *ast.BlockStmt) {
	// Vérification du corps de la boucle
	if body == nil {
		// Corps de boucle vide
		return
	}

	// Parcours des instructions
	for _, stmt := range body.List {
		checkStmtForAlloc(pass, stmt)
	}
}

// checkStmtForAlloc vérifie une instruction pour détecter allocations.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: instruction à vérifier
func checkStmtForAlloc(pass *analysis.Pass, stmt ast.Stmt) {
	// Vérification du type d'instruction
	switch s := stmt.(type) {
	// Cas d'une affectation
	case *ast.AssignStmt:
		// Vérification des affectations
		checkAssignForAlloc(pass, s)
	// Cas d'une déclaration
	case *ast.DeclStmt:
		// Vérification des déclarations
		checkDeclForAlloc(pass, s)
	// Note: les boucles imbriquées sont déjà gérées par Preorder
	}
}

// checkAssignForAlloc vérifie une affectation pour détecter allocations.
//
// Params:
//   - pass: contexte d'analyse
//   - assign: affectation à vérifier
func checkAssignForAlloc(pass *analysis.Pass, assign *ast.AssignStmt) {
	// Parcours des valeurs affectées
	for _, rhs := range assign.Rhs {
		// Vérification si allocation de slice ou map
		if isSliceOrMapAlloc(rhs) {
			// Allocation de slice/map détectée
			pass.Reportf(
				rhs.Pos(),
				"KTN-VAR-013: évitez d'allouer des slices/maps dans une boucle",
			)
		}
	}
}

// checkDeclForAlloc vérifie une déclaration pour détecter allocations.
//
// Params:
//   - pass: contexte d'analyse
//   - decl: déclaration à vérifier
func checkDeclForAlloc(pass *analysis.Pass, decl *ast.DeclStmt) {
	genDecl, ok := decl.Decl.(*ast.GenDecl)
	// Vérification du type de déclaration
	if !ok {
		// Pas une déclaration générale
		return
	}

	var valueSpec *ast.ValueSpec
	// Parcours des spécifications
	for _, spec := range genDecl.Specs {
		valueSpec, ok = spec.(*ast.ValueSpec)
		// Vérification de la spécification de valeur
		if !ok {
			// Pas une spécification de valeur
			continue
		}

		// Parcours des valeurs
		for _, value := range valueSpec.Values {
			// Vérification si allocation de slice ou map
			if isSliceOrMapAlloc(value) {
				// Allocation de slice/map détectée
				pass.Reportf(
					value.Pos(),
					"KTN-VAR-013: évitez d'allouer des slices/maps dans une boucle",
				)
			}
		}
	}
}

// isSliceOrMapAlloc vérifie si une expression est une allocation de slice/map.
//
// Params:
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si allocation détectée
func isSliceOrMapAlloc(expr ast.Expr) bool {
	// Vérification du type d'expression
	switch e := expr.(type) {
	// Cas d'un littéral composite
	case *ast.CompositeLit:
		// Vérification du type composite
		if isSliceOrMapType(e.Type) {
			// Allocation de slice/map sous forme de littéral
			return true
		}
	// Cas d'un appel de fonction
	case *ast.CallExpr:
		// Vérification des appels make()
		if ident, ok := e.Fun.(*ast.Ident); ok && ident.Name == "make" {
			// Appel à make() détecté
			return true
		}
	}
	// Pas d'allocation détectée
	return false
}

// isSliceOrMapType vérifie si un type est un slice ou une map.
//
// Params:
//   - typ: type à vérifier
//
// Returns:
//   - bool: true si slice ou map
func isSliceOrMapType(typ ast.Expr) bool {
	// Vérification du type
	switch typ.(type) {
	// Cas d'un type array ou map
	case *ast.ArrayType, *ast.MapType:
		// Type slice ou map détecté
		return true
	}
	// Pas un type slice ou map
	return false
}
