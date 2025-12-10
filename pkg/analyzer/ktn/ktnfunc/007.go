// Analyzer 007 for the ktnfunc package.
package ktnfunc

import (
	"go/ast"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeFunc007 is the rule code for this analyzer
	ruleCodeFunc007 string = "KTN-FUNC-007"
	// lazyFieldsCap est la capacité initiale pour la map des champs lazy load
	lazyFieldsCap int = 4
)

// Analyzer007 checks that getter functions don't have side effects
var Analyzer007 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc007",
	Doc:      "KTN-FUNC-007: Les getters (Get*/Is*/Has*) ne doivent pas avoir de side effects (assignations, appels de fonctions modifiant l'état)",
	Run:      runFunc007,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc007 exécute l'analyse KTN-FUNC-007.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runFunc007(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeFunc007) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Vérifier si le fichier est exclu
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		if cfg.IsFileExcluded(ruleCodeFunc007, filename) {
			// Fichier exclu
			return
		}

		// Skip test functions
		if shared.IsTestFunction(funcDecl) {
			// Retour pour ignorer les fonctions de test
			return
		}
		funcName := funcDecl.Name.Name
		// Skip if not a getter (Get*, Is*, Has*)
		if !isGetter(funcName) {
			// Retour si pas un getter
			return
		}
		// Skip if no body (external functions)
		if funcDecl.Body == nil {
			// Retour si pas de corps
			return
		}
		// Check for side effects in getter
		checkGetterSideEffects(pass, funcDecl, funcName)
	})

	// Retour de la fonction
	return nil, nil
}

// checkGetterSideEffects vérifie les side effects dans un getter.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: déclaration de fonction
//   - funcName: nom de la fonction
func checkGetterSideEffects(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	lazyLoadFields := collectLazyLoadFields(funcDecl.Body)
	ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
		// Sélection selon la valeur
		switch stmt := node.(type) {
		// Traitement des assignations
		case *ast.AssignStmt:
			reportAssignSideEffect(pass, stmt, funcName, lazyLoadFields)
		// Traitement des incréments/décréments
		case *ast.IncDecStmt:
			reportIncDecSideEffect(pass, stmt, funcName)
		}
		// Continuer l'inspection
		return true
	})
}

// reportAssignSideEffect reporte les side effects d'assignation.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: assignation
//   - funcName: nom de la fonction
//   - lazyLoadFields: champs de lazy loading
func reportAssignSideEffect(pass *analysis.Pass, stmt *ast.AssignStmt, funcName string, lazyLoadFields map[string]bool) {
	// Vérifier chaque côté gauche
	for _, lhs := range stmt.Lhs {
		// Vérification si side effect
		if hasSideEffect(lhs) {
			// Skip if it's a lazy load assignment
			if isLazyLoadAssignment(lhs, lazyLoadFields) {
				// Continuer si lazy load valide
				continue
			}
			pass.Reportf(
				stmt.Pos(),
				"KTN-FUNC-007: le getter '%s' ne doit pas modifier l'état (assignation détectée)",
				funcName,
			)
		}
	}
}

// reportIncDecSideEffect reporte les side effects d'incrément/décrément.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: statement incr/decr
//   - funcName: nom de la fonction
func reportIncDecSideEffect(pass *analysis.Pass, stmt *ast.IncDecStmt, funcName string) {
	// Vérifier si side effect sur champ
	if hasSideEffect(stmt.X) {
		pass.Reportf(
			stmt.Pos(),
			"KTN-FUNC-007: le getter '%s' ne doit pas modifier l'état (incrémentation/décrémentation détectée)",
			funcName,
		)
	}
}

// isGetter checks if a function name suggests it's a getter
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si fonction getter
func isGetter(name string) bool {
	// Retour de la fonction
	return strings.HasPrefix(name, "Get") ||
		strings.HasPrefix(name, "Is") ||
		strings.HasPrefix(name, "Has")
}

// hasSideEffect checks if an expression modifies external state.
//
// Params:
//   - expr: expression à analyser
//
// Returns:
//   - bool: true si effet de bord détecté
func hasSideEffect(expr ast.Expr) bool {
	// Sélection selon la valeur
	switch e := expr.(type) {
	// Traitement
	case *ast.SelectorExpr:
		// Modifying a field is a side effect
		return true
	// Traitement
	case *ast.IndexExpr:
		// Modifying an index (array/map/slice element) could be a side effect
		// Check if the base is a selector
		if _, ok := e.X.(*ast.SelectorExpr); ok {
			// Retour de la fonction
			return true
		}
	}
	// Retour de la fonction
	return false
}

// collectLazyLoadFields collects field names that are checked for nil in if statements.
// These fields are candidates for lazy loading patterns.
//
// Params:
//   - body: function body to analyze
//
// Returns:
//   - map[string]bool: map of field names that could be lazy loaded
func collectLazyLoadFields(body *ast.BlockStmt) map[string]bool {
	lazyFields := make(map[string]bool, lazyFieldsCap)
	// Inspect the body for nil checks
	ast.Inspect(body, func(node ast.Node) bool {
		ifStmt, ok := node.(*ast.IfStmt)
		// Si pas un if statement, continuer
		if !ok {
			// Retour true pour continuer l'inspection
			return true
		}
		// Check if condition is "field == nil"
		binary, ok := ifStmt.Cond.(*ast.BinaryExpr)
		// Si pas une expression binaire, continuer
		if !ok {
			// Retour true pour continuer l'inspection
			return true
		}
		// Vérifier si c'est une comparaison avec nil
		if isNilComparison(binary) {
			fieldName := extractFieldName(binary)
			// Ajouter à la map si trouvé
			if fieldName != "" {
				lazyFields[fieldName] = true
			}
		}
		// Continuer l'inspection
		return true
	})
	// Retour de la map des champs
	return lazyFields
}

// isNilComparison checks if a binary expression is a nil comparison.
//
// Params:
//   - binary: expression binaire à vérifier
//
// Returns:
//   - bool: true si c'est une comparaison avec nil
func isNilComparison(binary *ast.BinaryExpr) bool {
	// Vérifier si l'un des côtés est nil
	if ident, ok := binary.Y.(*ast.Ident); ok && ident.Name == "nil" {
		// Retour true si Y est nil
		return true
	}
	// Vérifier l'autre côté aussi
	if ident, ok := binary.X.(*ast.Ident); ok && ident.Name == "nil" {
		// Retour true si X est nil
		return true
	}
	// Retour false si pas de nil
	return false
}

// extractFieldName extracts the field name from a nil comparison.
//
// Params:
//   - binary: expression de comparaison avec nil
//
// Returns:
//   - string: nom du champ ou chaîne vide
func extractFieldName(binary *ast.BinaryExpr) string {
	// Vérifier le côté gauche
	if sel, ok := binary.X.(*ast.SelectorExpr); ok {
		// Retour du nom depuis X
		return sel.Sel.Name
	}
	// Vérifier le côté droit
	if sel, ok := binary.Y.(*ast.SelectorExpr); ok {
		// Retour du nom depuis Y
		return sel.Sel.Name
	}
	// Retour chaîne vide
	return ""
}

// isLazyLoadAssignment checks if an assignment is a lazy load pattern.
//
// Params:
//   - lhs: left-hand side of assignment
//   - lazyFields: map of fields that could be lazy loaded
//
// Returns:
//   - bool: true si c'est un lazy load
func isLazyLoadAssignment(lhs ast.Expr, lazyFields map[string]bool) bool {
	// Vérifier si c'est un selector (accès à un champ)
	sel, ok := lhs.(*ast.SelectorExpr)
	// Si pas un selector
	if !ok {
		// Retour false
		return false
	}
	// Retour true si le champ est dans la liste des lazy load
	return lazyFields[sel.Sel.Name]
}
