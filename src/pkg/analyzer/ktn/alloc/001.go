package ktn_alloc

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer/utils"
	"golang.org/x/tools/go/analysis"
)

// Rule001 vérifie que new() n'est pas utilisé avec des types référence (slice/map/chan).
// KTN-ALLOC-001: new() retourne un pointeur vers nil pour les types référence, causant des panics.
var Rule001 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_ALLOC_001",
	Doc:  "Interdiction de new() avec slice/map/chan (utiliser make() à la place)",
	Run:  runRule001,
}

// runRule001 exécute l'analyse pour la règle ALLOC-001.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				// Continue traversing AST nodes.
				return true
			}

			ident, ok := callExpr.Fun.(*ast.Ident)
			if !ok || ident.Name != "new" {
				// Continue traversing AST nodes.
				return true
			}

			// Vérifier que new() a exactement 1 argument
			if len(callExpr.Args) != 1 {
				// Continue traversing AST nodes.
				return true
			}

			arg := callExpr.Args[0]

			// KTN-ALLOC-001 : Interdire new() avec slice/map/chan
			if utils.IsReferenceType(arg) {
				reportReferenceTypeViolation(pass, callExpr, arg)
			}

			// Continue traversing AST nodes.
			return true
		})
	}
	// Analysis completed successfully.
	return nil, nil
}

// reportReferenceTypeViolation rapporte une violation KTN-ALLOC-001.
//
// Params:
//   - pass: la passe d'analyse
//   - callExpr: l'appel à new()
//   - arg: l'argument de type référence
func reportReferenceTypeViolation(pass *analysis.Pass, callExpr *ast.CallExpr, arg ast.Expr) {
	typeName := utils.GetTypeName(arg)
	pass.Reportf(callExpr.Pos(),
		"[KTN-ALLOC-001] Utilisation de new() avec un type référence (%s) interdite.\n"+
			"new() retourne un pointeur vers nil pour les types référence, ce qui cause des panics.\n"+
			"Utilisez make() à la place.\n"+
			"Exemple:\n"+
			"  // ❌ INTERDIT\n"+
			"  m := new(%s)  // *%s avec nil map/slice/chan\n"+
			"  (*m)[\"key\"] = value  // 💥 PANIC\n"+
			"\n"+
			"  // ✅ CORRECT\n"+
			"  m := make(%s)\n"+
			"  m[\"key\"] = value  // ✅ Fonctionne",
		typeName, typeName, typeName, typeName)
}
