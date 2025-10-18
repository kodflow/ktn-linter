package ktn_alloc

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer/utils"
	"golang.org/x/tools/go/analysis"
)

// Rule003 vérifie l'utilisation de new() avec des structs.
// KTN-ALLOC-003: Préférer &struct{} à new(struct) pour un code idiomatique Go.
var Rule003 = &analysis.Analyzer{
	Name: "KTN_ALLOC_003",
	Doc:  "Préférer composite literals &T{} à new(T) pour les structs",
	Run:  runRule003,
}

// runRule003 exécute l'analyse pour la règle ALLOC-003.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule003(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			ident, ok := callExpr.Fun.(*ast.Ident)
			if !ok || ident.Name != "new" {
				return true
			}

			// Vérifier que new() a exactement 1 argument
			if len(callExpr.Args) != 1 {
				return true
			}

			arg := callExpr.Args[0]

			// KTN-ALLOC-003 : Préférer &struct{} à new(struct)
			if utils.IsStructType(arg) {
				reportStructTypeViolation(pass, callExpr, arg)
			}

			return true
		})
	}
	return nil, nil
}

// reportStructTypeViolation rapporte une violation KTN-ALLOC-003.
//
// Params:
//   - pass: la passe d'analyse
//   - callExpr: l'appel à new()
//   - arg: l'argument de type struct
func reportStructTypeViolation(pass *analysis.Pass, callExpr *ast.CallExpr, arg ast.Expr) {
	typeName := utils.GetTypeName(arg)
	pass.Reportf(callExpr.Pos(),
		"[KTN-ALLOC-003] Utilisez le composite literal &%s{} au lieu de new(%s).\n"+
			"En Go idiomatique, on préfère les composite literals pour les structs.\n"+
			"Exemple:\n"+
			"  // ❌ NON-IDIOMATIQUE\n"+
			"  p := new(%s)\n"+
			"\n"+
			"  // ✅ IDIOMATIQUE GO\n"+
			"  p := &%s{}",
		typeName, typeName, typeName, typeName)
}
