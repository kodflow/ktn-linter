// Analyzer 004 for the ktntest package.
package ktntest

import (
	"go/ast"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// MIN_ERROR_CHECKS est le nombre minimum de vérifications d'erreur
	MIN_ERROR_CHECKS int = 1
)

// Analyzer004 checks that tests cover error cases
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktntest004",
	Doc:      "KTN-TEST-004: Les tests doivent couvrir les cas d'erreur et exceptions",
	Run:      runTest004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest004 exécute l'analyse KTN-TEST-004.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest004(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Analyser les fonctions de test
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Vérification de la condition
		if !shared.IsTestFile(filename) {
			// Pas un fichier de test
			return
		}

		// Vérifier si c'est une fonction de test unitaire (Test*)
		if !shared.IsUnitTestFunction(funcDecl) {
			// Pas une fonction de test unitaire
			return
		}

		// Vérifier si le test couvre les cas d'erreur
		if !hasErrorCaseCoverage(funcDecl) {
			// Pas de couverture des cas d'erreur
			pass.Reportf(
				funcDecl.Pos(),
				"KTN-TEST-004: le test '%s' devrait couvrir les cas d'erreur/exceptions",
				funcDecl.Name.Name,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// hasErrorCaseCoverage vérifie si le test couvre les cas d'erreur.
//
// Params:
//   - funcDecl: déclaration de fonction de test
//
// Returns:
//   - bool: true si le test couvre les cas d'erreur
func hasErrorCaseCoverage(funcDecl *ast.FuncDecl) bool {
	errorCaseIndicators := 0

	// Parcourir le corps de la fonction
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		// Vérification de la condition
		switch node := n.(type) {
		// Cas d'un composite literal (tests table)
		case *ast.CompositeLit:
			// Vérifier si c'est un test table avec des cas d'erreur
			if hasErrorTestCases(node) {
				errorCaseIndicators++
			}
		// Cas d'un ident (variables)
		case *ast.Ident:
			// Vérifier les noms de variables indicateurs d'erreur
			if isErrorIndicatorName(node.Name) {
				errorCaseIndicators++
			}
		// Cas d'un string literal
		case *ast.BasicLit:
			// Vérifier les strings avec "error", "invalid", "fail"
			if node.Kind.String() == "STRING" {
				value := strings.ToLower(node.Value)
				// Vérification de la condition
				if strings.Contains(value, "error") ||
					strings.Contains(value, "invalid") ||
					strings.Contains(value, "fail") {
					errorCaseIndicators++
				}
			}
		}
		// Continue traversal
		return true
	})

	// Retour du résultat
	return errorCaseIndicators >= MIN_ERROR_CHECKS
}

// hasErrorTestCases vérifie si un composite literal a des cas d'erreur.
//
// Params:
//   - lit: composite literal
//
// Returns:
//   - bool: true si des cas d'erreur sont présents
func hasErrorTestCases(lit *ast.CompositeLit) bool {
	// Parcourir les éléments
	for _, elt := range lit.Elts {
		// Vérifier si c'est un KeyValueExpr
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			// Vérifier si la clé est "name" ou "want"
			if ident, identOk := kv.Key.(*ast.Ident); identOk {
				// Vérification de la condition
				if ident.Name == "name" {
					// Vérifier si la valeur contient "error", "invalid", "fail"
					if basic, basicOk := kv.Value.(*ast.BasicLit); basicOk {
						value := strings.ToLower(basic.Value)
						// Vérification de la condition
						if strings.Contains(value, "error") ||
							strings.Contains(value, "invalid") ||
							strings.Contains(value, "fail") {
							// Cas d'erreur trouvé
							return true
						}
					}
				}
			}
		}
	}

	// Pas de cas d'erreur
	return false
}

// isErrorIndicatorName vérifie si un nom indique un cas d'erreur.
//
// Params:
//   - name: nom de la variable
//
// Returns:
//   - bool: true si c'est un indicateur d'erreur
func isErrorIndicatorName(name string) bool {
	lowerName := strings.ToLower(name)
	errorIndicators := []string{
		"err",
		"error",
		"invalid",
		"fail",
		"bad",
		"wrong",
	}

	// Parcours des indicateurs
	for _, indicator := range errorIndicators {
		// Vérification de la condition
		if strings.Contains(lowerName, indicator) {
			// Indicateur trouvé
			return true
		}
	}

	// Pas d'indicateur
	return false
}
