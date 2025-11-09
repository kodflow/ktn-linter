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

// Analyzer009 checks that public function tests are only in external test files
var Analyzer009 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktntest009",
	Doc:      "KTN-TEST-009: Les tests de fonctions publiques (exportées) doivent être dans _external_test.go uniquement (black-box testing)",
	Run:      runTest009,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest009 exécute l'analyse KTN-TEST-009.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest009(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Map des fonctions publiques du package
	publicFunctions := make(map[string]bool)

	// Premier passage : collecter toutes les fonctions publiques
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

		// Vérifier si c'est une fonction publique
		if funcDecl.Name != nil && len(funcDecl.Name.Name) > 0 {
			firstRune := rune(funcDecl.Name.Name[0])
			// Vérification fonction publique
			if unicode.IsUpper(firstRune) {
				publicFunctions[funcDecl.Name.Name] = true

				// Pour les méthodes, ajouter aussi le pattern Type_Method
				if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
					receiverType := extractReceiverType(funcDecl.Recv.List[0].Type)
					// Vérification receiver valide
					if receiverType != "" {
						compositeKey := receiverType + "_" + funcDecl.Name.Name
						publicFunctions[compositeKey] = true
					}
				}
			}
		}
	})

	// Deuxième passage : vérifier les tests dans les fichiers _internal_test.go
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		baseName := filepath.Base(filename)

		// Vérifier si c'est un fichier _internal_test.go
		if !strings.HasSuffix(baseName, "_internal_test.go") {
			// Pas un fichier internal test
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

		// Vérifier si c'est un test de fonction publique
		if publicFunctions[testedFuncName] {
			pass.Reportf(
				funcDecl.Pos(),
				"KTN-TEST-009: le test '%s' dans '%s' teste une fonction publique '%s'. Les tests de fonctions publiques doivent être dans '%s' (black-box testing avec package xxx_test)",
				funcDecl.Name.Name,
				baseName,
				testedFuncName,
				strings.Replace(baseName, "_internal_test.go", "_external_test.go", 1),
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// extractReceiverType extrait le nom du type du receiver.
//
// Params:
//   - expr: expression du receiver
//
// Returns:
//   - string: nom du type
func extractReceiverType(expr ast.Expr) string {
	// Gérer les pointeurs
	if starExpr, ok := expr.(*ast.StarExpr); ok {
		// Retour récursif
		return extractReceiverType(starExpr.X)
	}

	// Gérer les identifiants
	if ident, ok := expr.(*ast.Ident); ok {
		// Retour du nom
		return ident.Name
	}

	// Type non géré
	return ""
}
