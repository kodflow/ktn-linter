// Analyzer 002 for the ktnreturn package.
package ktnreturn

import (
	"go/ast"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeReturn002 rule code for RETURN-002
	ruleCodeReturn002 string = "KTN-RETURN-002"
)

// Analyzer002 detects nil returns for slice and map types.
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnreturn002",
	Doc:      "KTN-RETURN-002: préférer slice/map vide à nil",
	Run:      runReturn002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runReturn002 analyzes return statements for nil slice/map returns.
// Params:
//   - pass: Analysis pass containing type information
//
// Returns:
//   - any: always nil
//   - error: analysis error if any
func runReturn002(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeReturn002) {
		// Règle désactivée
		return nil, nil
	}

	inspectResult := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspectResult.Preorder(nodeFilter, func(n ast.Node) {
		// Vérifier si le fichier est exclu
		filename := pass.Fset.Position(n.Pos()).Filename
		// Vérification si le fichier courant est exclu par la configuration de la règle
		if cfg.IsFileExcluded(ruleCodeReturn002, filename) {
			// Fichier exclu
			return
		}

		funcDecl := n.(*ast.FuncDecl)

		// Skip if function has no return type
		if funcDecl.Type == nil || funcDecl.Type.Results == nil {
			// Retour anticipé si pas de type de retour
			return
		}

		// Check all return types
		for _, result := range funcDecl.Type.Results.List {
			// Verification de la condition
			if isSliceOrMapType(pass, result.Type) {
				// Analyze function body for nil returns
				checkNilReturns(pass, funcDecl)
				break
			}
		}
	})

	// Retour sans erreur
	return nil, nil
}

// isSliceOrMapType checks if expression is slice or map type.
//
// Params:
//   - pass: Analysis pass
//   - expr: Expression to check
//
// Returns:
//   - bool: true if expression is slice or map type
func isSliceOrMapType(pass *analysis.Pass, expr ast.Expr) bool {
	typeInfo := pass.TypesInfo.TypeOf(expr)
	// Return false if type information is unavailable
	if typeInfo == nil {
		// Retour si information de type indisponible
		return false
	}

	// Check underlying type
	switch typeInfo.Underlying().(type) {
	// Verification de la condition
	case *types.Slice, *types.Map:
		// Retour si type slice ou map détecté
		return true
	}
	// Retour par défaut si type non-slice/map
	return false
}

// checkNilReturns analyzes function body for nil returns.
// Params:
//   - pass: Analysis pass
//   - funcDecl: Function declaration to analyze
func checkNilReturns(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Skip if function has no body
	if funcDecl.Body == nil {
		// Retour anticipé si pas de corps de fonction
		return
	}

	// Collect slice/map return types for better messages
	returnTypes := collectSliceMapReturnTypes(pass, funcDecl)

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		retStmt, ok := n.(*ast.ReturnStmt)
		// Continue traversal if not return statement
		if !ok {
			// Retour pour continuer la traversée
			return true
		}

		// Check each return value
		for i, result := range retStmt.Results {
			// Verification de la condition
			if isNilIdent(result) && i < len(returnTypes) {
				typeInfo := returnTypes[i]
				// Report avec message adapté au type
				if typeInfo != "" {
					pass.Reportf(
						retStmt.Pos(),
						"KTN-RETURN-002: préférer %s à nil",
						typeInfo,
					)
				}
			}
		}

		// Retour pour continuer la traversée
		return true
	})
}

// collectSliceMapReturnTypes collecte les types slice/map pour chaque position de retour.
//
// Params:
//   - pass: Analysis pass
//   - funcDecl: Function declaration
//
// Returns:
//   - []string: suggestion de syntaxe pour chaque position de retour
func collectSliceMapReturnTypes(pass *analysis.Pass, funcDecl *ast.FuncDecl) []string {
	// Vérification des résultats
	if funcDecl.Type.Results == nil {
		// Retour d'une slice vide si pas de résultats
		return []string{}
	}

	// Collecte des types avec expansion des noms multiples
	var result []string
	// Itération sur les résultats
	for _, field := range funcDecl.Type.Results.List {
		typeInfo := pass.TypesInfo.TypeOf(field.Type)
		// Vérification du type
		if typeInfo == nil {
			result = append(result, "")
			continue
		}

		// Génération de la suggestion
		suggestion := getSuggestionForType(typeInfo)

		// Si plusieurs noms dans le champ, répéter la suggestion
		count := len(field.Names)
		// Au moins une fois si pas de noms
		if count == 0 {
			count = 1
		}
		// Ajout des suggestions
		for range count {
			result = append(result, suggestion)
		}
	}

	// Retour de la liste des suggestions collectées
	return result
}

// getSuggestionForType retourne la suggestion de syntaxe pour un type.
//
// Params:
//   - t: type à analyser
//
// Returns:
//   - string: suggestion de syntaxe ou chaîne vide
func getSuggestionForType(t types.Type) string {
	// Vérification du type sous-jacent
	switch underlying := t.Underlying().(type) {
	// Cas slice
	case *types.Slice:
		elemType := underlying.Elem().String()
		// Retour de la suggestion pour slice
		return "[]" + elemType + "{}"
	// Cas map
	case *types.Map:
		keyType := underlying.Key().String()
		elemType := underlying.Elem().String()
		// Retour de la suggestion pour map
		return "map[" + keyType + "]" + elemType + "{}"
	}
	// Retour vide si type non slice/map
	return ""
}

// isNilIdent checks if expression is nil identifier.
//
// Params:
//   - expr: Expression to check
//
// Returns:
//   - bool: true if expression is nil identifier
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	// Return false if not identifier
	if !ok {
		// Retour si l'expression n'est pas un identifiant
		return false
	}
	// Verification de la condition
	// Retour du résultat de la comparaison
	return ident.Name == "nil"
}
