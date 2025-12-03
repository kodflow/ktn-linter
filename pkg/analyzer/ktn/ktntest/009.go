// Analyzer 009 for the ktntest package.
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
	// INITIAL_PUBLIC_FUNCS_CAP initial capacity for public funcs map
	INITIAL_PUBLIC_FUNCS_CAP int = 32
)

// Analyzer009 checks that public function tests are in external test files
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

	// Collecter toutes les fonctions publiques
	publicFunctions := collectPublicFunctions(pass, insp)

	// Vérifier les tests dans les fichiers _internal_test.go
	checkInternalTestsForPublicFunctions(pass, insp, publicFunctions)

	// Retour de la fonction
	return nil, nil
}

// collectPublicFunctions collecte les fonctions publiques du package.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//
// Returns:
//   - map[string]bool: map des noms de fonctions publiques
func collectPublicFunctions(pass *analysis.Pass, insp *inspector.Inspector) map[string]bool {
	publicFunctions := make(map[string]bool, INITIAL_PUBLIC_FUNCS_CAP)
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	// Parcourir toutes les déclarations de fonctions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			return
		}

		// Ajouter la fonction si elle est publique
		addPublicFunction(funcDecl, publicFunctions)
	})

	// Retour de la map
	return publicFunctions
}

// addPublicFunction ajoute une fonction publique à la map.
//
// Params:
//   - funcDecl: déclaration de fonction
//   - publicFunctions: map des fonctions publiques
func addPublicFunction(funcDecl *ast.FuncDecl, publicFunctions map[string]bool) {
	// Vérifier le nom de la fonction
	if funcDecl.Name == nil || len(funcDecl.Name.Name) == 0 {
		return
	}

	firstRune := rune(funcDecl.Name.Name[0])
	// Vérification fonction publique
	if !unicode.IsUpper(firstRune) {
		return
	}

	// Ajouter le nom de la fonction
	publicFunctions[funcDecl.Name.Name] = true

	// Pour les méthodes, ajouter aussi le pattern Type_Method
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		receiverType := extractReceiverType(funcDecl.Recv.List[0].Type)
		// Vérification receiver valide
		if receiverType != "" {
			publicFunctions[receiverType+"_"+funcDecl.Name.Name] = true
		}
	}
}

// checkInternalTestsForPublicFunctions vérifie les tests de fonctions publiques.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - publicFunctions: map des fonctions publiques
func checkInternalTestsForPublicFunctions(pass *analysis.Pass, insp *inspector.Inspector, publicFunctions map[string]bool) {
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}

	// Parcourir les fonctions de test
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		baseName := filepath.Base(filename)

		// Vérifier si c'est un test dans un fichier internal
		if !strings.HasSuffix(baseName, "_internal_test.go") || !shared.IsUnitTestFunction(funcDecl) {
			return
		}

		// Vérifier si c'est un test de fonction publique
		checkAndReportPublicFunctionTest(pass, funcDecl, baseName, publicFunctions)
	})
}

// checkAndReportPublicFunctionTest vérifie et reporte un test de fonction publique mal placé.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: déclaration de fonction de test
//   - baseName: nom de base du fichier
//   - publicFunctions: map des fonctions publiques
func checkAndReportPublicFunctionTest(pass *analysis.Pass, funcDecl *ast.FuncDecl, baseName string, publicFunctions map[string]bool) {
	// Extraire le nom de la fonction testée
	testedFuncName := strings.TrimPrefix(funcDecl.Name.Name, "Test")
	// Vérification nom valide
	if testedFuncName == "" {
		return
	}

	// Vérifier si c'est un test de fonction publique
	if publicFunctions[testedFuncName] {
		pass.Reportf(
			funcDecl.Pos(),
			"KTN-TEST-009: le test '%s' dans '%s' teste une fonction publique '%s'. Les tests de fonctions publiques doivent être dans '%s' (black-box testing avec package xxx_test)",
			funcDecl.Name.Name, baseName, testedFuncName,
			strings.Replace(baseName, "_internal_test.go", "_external_test.go", 1),
		)
	}
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
