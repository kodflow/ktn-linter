// Analyzer 005 for the ktntest package.
package ktntest

import (
	"go/ast"
	"slices"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// MIN_TEST_CASES est le nombre minimum de cas de test pour table-driven
	// Note: 3 assertions = pattern répétitif qui devrait utiliser table-driven
	// 2 assertions = normal (ex: vérifier erreur + résultat)
	MIN_TEST_CASES int = 3
)

// Analyzer005 checks that tests use table-driven test pattern
var Analyzer005 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktntest005",
	Doc:      "KTN-TEST-005: Les tests avec plusieurs cas doivent utiliser table-driven tests",
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

		// Vérifier si le test a plusieurs assertions sans table-driven
		if hasMultipleAssertions(funcDecl) && !hasTableDrivenPattern(funcDecl) {
			// Pas de table-driven test
			pass.Reportf(
				funcDecl.Pos(),
				"KTN-TEST-005: le test '%s' devrait utiliser table-driven tests",
				funcDecl.Name.Name,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// hasMultipleAssertions vérifie si le test a plusieurs assertions.
//
// Params:
//   - funcDecl: déclaration de fonction de test
//
// Returns:
//   - bool: true si plusieurs assertions
func hasMultipleAssertions(funcDecl *ast.FuncDecl) bool {
	assertionCount := 0

	// Parcourir le corps de la fonction
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		// Vérification de la condition
		if !ok {
			// Continue traversal
			return true
		}

		// Utiliser la nouvelle fonction isAssertionCall
		if isAssertionCall(call) {
			assertionCount++
		}

		// Continue traversal
		return true
	})

	// Retour du résultat
	return assertionCount >= MIN_TEST_CASES
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

// isAssertionCall vérifie si c'est un appel d'assertion (testing, testify/assert, testify/require).
//
// Params:
//   - call: appel de fonction AST
//
// Returns:
//   - bool: true si c'est une assertion
func isAssertionCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	// Vérification de la condition
	if !ok {
		// Pas un SelectorExpr
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	// Vérification de la condition
	if !ok {
		// Pas un Ident
		return false
	}

	methodName := sel.Sel.Name

	// Case 1: testing package (t.Error, t.Fatal, etc.)
	if ident.Name == "t" && isTestingMethod(methodName) {
		// Méthode testing détectée
		return true
	}

	// Case 2: testify/assert package
	if ident.Name == "assert" && isAssertMethod(methodName) {
		// Méthode assert détectée
		return true
	}

	// Case 3: testify/require package
	if ident.Name == "require" && isRequireMethod(methodName) {
		// Méthode require détectée
		return true
	}

	// Pas une assertion
	return false
}

// isTestingMethod vérifie si c'est une méthode du package testing.
//
// Params:
//   - methodName: nom de la méthode
//
// Returns:
//   - bool: true si c'est une méthode testing
func isTestingMethod(methodName string) bool {
	methods := []string{
		"Error", "Errorf", "Fatal", "Fatalf", "Fail", "FailNow",
		"Log", "Logf", "Skip", "Skipf", "SkipNow",
	}
	// Vérifier si c'est une méthode testing
	return slices.Contains(methods, methodName)
}

// isAssertMethod vérifie si c'est une méthode de testify/assert.
//
// Params:
//   - methodName: nom de la méthode
//
// Returns:
//   - bool: true si c'est une méthode assert
func isAssertMethod(methodName string) bool {
	// Liste des méthodes assert les plus courantes
	methods := []string{
		"Equal", "NotEqual", "Nil", "NotNil", "True", "False",
		"Empty", "NotEmpty", "Len", "Contains", "NotContains",
		"Greater", "GreaterOrEqual", "Less", "LessOrEqual",
		"Same", "NotSame", "Implements", "IsType", "Panics",
		"NotPanics", "WithinDuration", "InDelta", "InEpsilon",
		"JSONEq", "YAMLEq", "Error", "NoError", "ErrorIs", "ErrorAs",
		"ErrorContains", "Regexp", "NotRegexp", "Zero", "NotZero",
	}
	// Vérifier si c'est une méthode assert
	return slices.Contains(methods, methodName)
}

// isRequireMethod vérifie si c'est une méthode de testify/require.
//
// Params:
//   - methodName: nom de la méthode
//
// Returns:
//   - bool: true si c'est une méthode require
func isRequireMethod(methodName string) bool {
	// require a les mêmes méthodes que assert
	return isAssertMethod(methodName)
}
