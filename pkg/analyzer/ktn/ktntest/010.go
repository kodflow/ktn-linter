// Analyzer 010 for the ktntest package.
package ktntest

import (
	"go/ast"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// INITIAL_PRIVATE_FUNCTIONS_CAP initial cap for private funcs map
	INITIAL_PRIVATE_FUNCTIONS_CAP int = 32
)

// Analyzer010 checks private function tests are in internal test files
var Analyzer010 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktntest010",
	Doc:      "KTN-TEST-010: Les tests de fonctions privées (non-exportées) doivent être dans _internal_test.go uniquement (white-box testing)",
	Run:      runTest010,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest010 exécute l'analyse KTN-TEST-010.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest010(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecter toutes les fonctions privées
	privateFunctions := collectPrivateFunctions(pass, insp)

	// Vérifier les tests dans les fichiers _external_test.go
	checkExternalTestsForPrivateFunctions(pass, insp, privateFunctions)

	// Retour de la fonction
	return nil, nil
}

// collectPrivateFunctions collecte les fonctions privées du package.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//
// Returns:
//   - map[string]bool: map des noms de fonctions privées
func collectPrivateFunctions(pass *analysis.Pass, insp *inspector.Inspector) map[string]bool {
	privateFunctions := make(map[string]bool, INITIAL_PRIVATE_FUNCTIONS_CAP)
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	// Parcourir toutes les déclarations de fonctions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			return
		}

		// Ajouter la fonction si elle est privée
		if funcDecl.Name != nil && len(funcDecl.Name.Name) > 0 {
			firstRune := rune(funcDecl.Name.Name[0])
			// Vérification fonction privée
			if unicode.IsLower(firstRune) {
				privateFunctions[funcDecl.Name.Name] = true
			}
		}
	})

	// Retour de la map
	return privateFunctions
}

// checkExternalTestsForPrivateFunctions vérifie les tests de fonctions privées.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - privateFunctions: map des fonctions privées
func checkExternalTestsForPrivateFunctions(pass *analysis.Pass, insp *inspector.Inspector, privateFunctions map[string]bool) {
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	// Parcourir les fonctions de test
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		baseName := filepath.Base(filename)

		// Vérifier si c'est un test dans un fichier external
		if !strings.HasSuffix(baseName, "_external_test.go") || !shared.IsUnitTestFunction(funcDecl) {
			return
		}

		// Vérifier si c'est un test de fonction privée
		checkAndReportPrivateFunctionTest(pass, funcDecl, baseName, privateFunctions)
	})
}

// checkAndReportPrivateFunctionTest vérifie et reporte un test de fonction privée mal placé.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: déclaration de fonction de test
//   - baseName: nom de base du fichier
//   - privateFunctions: map des fonctions privées
func checkAndReportPrivateFunctionTest(pass *analysis.Pass, funcDecl *ast.FuncDecl, baseName string, privateFunctions map[string]bool) {
	// Extraire le nom de la fonction testée
	testedFuncName := strings.TrimPrefix(funcDecl.Name.Name, "Test")
	// Vérification nom valide
	if testedFuncName == "" {
		return
	}

	// Pour les tests de méthodes privées, extraire le nom de la fonction
	testedFuncName = extractPrivateFunctionName(testedFuncName)

	// Vérifier si c'est un test de fonction privée
	if isPrivateFunctionTested(testedFuncName, privateFunctions) {
		pass.Reportf(
			funcDecl.Pos(),
			"KTN-TEST-010: le test '%s' dans '%s' teste une fonction privée '%s'. Les tests de fonctions privées doivent être dans '%s' (white-box testing avec package xxx)",
			funcDecl.Name.Name, baseName, testedFuncName,
			strings.Replace(baseName, "_external_test.go", "_internal_test.go", 1),
		)
	}
}

// extractPrivateFunctionName extrait le nom de fonction privée du nom de test.
//
// Params:
//   - testedFuncName: nom extrait du test
//
// Returns:
//   - string: nom de la fonction privée
func extractPrivateFunctionName(testedFuncName string) string {
	// Pattern: TestType_privateMethod -> privateMethod
	parts := strings.Split(testedFuncName, "_")
	// Vérification pattern méthode
	if len(parts) > 1 {
		// Prendre le dernier élément
		return parts[len(parts)-1]
	}
	// Retour du nom original
	return testedFuncName
}

// isPrivateFunctionTested vérifie si une fonction privée est testée.
//
// Params:
//   - testedFuncName: nom de la fonction testée
//   - privateFunctions: map des fonctions privées
//
// Returns:
//   - bool: true si c'est une fonction privée testée
func isPrivateFunctionTested(testedFuncName string, privateFunctions map[string]bool) bool {
	// Vérifier le nom
	if len(testedFuncName) == 0 {
		// Nom vide, pas de fonction testée
		return false
	}
	// Vérifier si commence par minuscule et existe
	firstRune := rune(testedFuncName[0])
	// Retour de la vérification
	return unicode.IsLower(firstRune) && privateFunctions[testedFuncName]
}
