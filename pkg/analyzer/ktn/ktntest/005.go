// Analyzer 005 for the ktntest package.
package ktntest

import (
	"go/ast"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer005 checks that tests use table-driven test pattern
var Analyzer005 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktntest005",
	Doc:      "KTN-TEST-005: TOUS les tests doivent utiliser le pattern table-driven",
	Run:      runTest005,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest005 exécute l'analyse KTN-TEST-005.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest005(pass *analysis.Pass) (any, error) {
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

		// Vérifier si le test utilise le pattern table-driven (obligatoire)
		if !hasTableDrivenPattern(funcDecl) {
			// Pas de table-driven test
			pass.Reportf(
				funcDecl.Pos(),
				"KTN-TEST-005: le test '%s' DOIT utiliser le pattern table-driven",
				funcDecl.Name.Name,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// isTestsVariableName vérifie si le nom correspond à une variable de tests.
//
// Params:
//   - name: nom de la variable
//
// Returns:
//   - bool: true si c'est un nom de variable de tests
func isTestsVariableName(name string) bool {
	lowerName := strings.ToLower(name)
	// Vérifier les noms courants
	return lowerName == "tests" || lowerName == "testcases" || lowerName == "cases"
}

// checkAssignStmt vérifie si l'AssignStmt contient une variable de tests.
//
// Params:
//   - node: nœud AssignStmt
//
// Returns:
//   - bool: true si contient variable de tests
func checkAssignStmt(node *ast.AssignStmt) bool {
	// Parcourir les variables à gauche
	for _, lhs := range node.Lhs {
		// Vérifier si c'est un identifiant
		if ident, ok := lhs.(*ast.Ident); ok {
			// Vérifier si c'est une variable de tests
			if isTestsVariableName(ident.Name) {
				// Variable de tests trouvée
				return true
			}
		}
	}
	// Pas de variable de tests
	return false
}

// checkRangeStmt vérifie si le RangeStmt itère sur une variable de tests.
//
// Params:
//   - node: nœud RangeStmt
//
// Returns:
//   - bool: true si itère sur variable de tests
func checkRangeStmt(node *ast.RangeStmt) bool {
	// Vérifier si on itère sur un identifiant
	if ident, ok := node.X.(*ast.Ident); ok {
		// Vérifier si c'est une variable de tests
		return isTestsVariableName(ident.Name)
	}
	// Pas une boucle sur variable de tests
	return false
}

// hasTableDrivenPattern vérifie si le test utilise table-driven pattern.
//
// Params:
//   - funcDecl: déclaration de fonction de test
//
// Returns:
//   - bool: true si table-driven pattern
func hasTableDrivenPattern(funcDecl *ast.FuncDecl) bool {
	hasTestsVar := false
	hasRangeLoop := false

	// Parcourir le corps de la fonction
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		// Vérification du type de nœud
		switch node := n.(type) {
		// Cas d'une déclaration de variable
		case *ast.AssignStmt:
			// Vérifier si contient variable de tests
			if checkAssignStmt(node) {
				hasTestsVar = true
			}
		// Cas d'une boucle range
		case *ast.RangeStmt:
			// Vérifier si itère sur variable de tests
			if checkRangeStmt(node) {
				hasRangeLoop = true
			}
		}
		// Continue traversal
		return true
	})

	// Retour du résultat
	return hasTestsVar && hasRangeLoop
}
