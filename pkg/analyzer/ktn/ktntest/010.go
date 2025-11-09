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

// Analyzer010 checks that private function tests are only in internal test files
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

	// Map des fonctions privées du package
	privateFunctions := make(map[string]bool)

	// Premier passage : collecter toutes les fonctions privées
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Skip test files
			return
		}

		// Vérifier si c'est une fonction privée
		if funcDecl.Name != nil && len(funcDecl.Name.Name) > 0 {
			firstRune := rune(funcDecl.Name.Name[0])
			// Vérification fonction privée (commence par minuscule)
			if unicode.IsLower(firstRune) {
				privateFunctions[funcDecl.Name.Name] = true
			}
		}
	})

	// Deuxième passage : vérifier les tests dans les fichiers _external_test.go
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		baseName := filepath.Base(filename)

		// Vérifier si c'est un fichier _external_test.go
		if !strings.HasSuffix(baseName, "_external_test.go") {
			// Pas un fichier external test
			return
		}

		// Vérifier si c'est une fonction de test
		if !shared.IsUnitTestFunction(funcDecl) {
			// Pas une fonction de test
			return
		}

		// Extraire le nom de la fonction testée
		testedFuncName := strings.TrimPrefix(funcDecl.Name.Name, "Test")
		// Vérification nom valide
		if testedFuncName == "" {
			// Nom vide
			return
		}

		// Pour les tests de méthodes privées, extraire le nom de la fonction
		// Pattern: TestType_privateMethod -> privateMethod
		parts := strings.Split(testedFuncName, "_")
		// Vérification pattern méthode
		if len(parts) > 1 {
			// Prendre le dernier élément (nom de la méthode)
			testedFuncName = parts[len(parts)-1]
		}

		// Vérifier si c'est un test de fonction privée
		// (commence par minuscule)
		if len(testedFuncName) > 0 {
			firstRune := rune(testedFuncName[0])
			// Vérification fonction privée testée
			if unicode.IsLower(firstRune) && privateFunctions[testedFuncName] {
				pass.Reportf(
					funcDecl.Pos(),
					"KTN-TEST-010: le test '%s' dans '%s' teste une fonction privée '%s'. Les tests de fonctions privées doivent être dans '%s' (white-box testing avec package xxx)",
					funcDecl.Name.Name,
					baseName,
					testedFuncName,
					strings.Replace(baseName, "_external_test.go", "_internal_test.go", 1),
				)
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}
