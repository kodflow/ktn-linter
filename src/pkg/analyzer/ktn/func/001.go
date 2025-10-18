package ktn_func

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer/utils"
)

// Rule001 vérifie le nommage MixedCaps des fonctions.
//
// KTN-FUNC-001: Les fonctions doivent utiliser la convention MixedCaps.
// - Fonctions exportées : MixedCaps (ex: CalculateTotal)
// - Fonctions privées : mixedCaps (ex: calculateTotal)
//
// Incorrect: calculate_total, Calculate_Total
// Correct: calculateTotal, CalculateTotal
var Rule001 = &analysis.Analyzer{
	Name: "KTN_FUNC_001",
	Doc:  "Vérifie que les fonctions utilisent la convention MixedCaps",
	Run:  runRule001,
}

// runRule001 exécute la vérification KTN-FUNC-001.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			checkFuncNaming(pass, funcDecl)
		}
	}

	return nil, nil
}

// checkFuncNaming vérifie le nommage de la fonction.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - funcDecl: la déclaration de fonction
func checkFuncNaming(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	funcName := funcDecl.Name.Name

	if !utils.IsMixedCaps(funcName) {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-001] Fonction '%s' n'utilise pas la convention MixedCaps.\nUtilisez MixedCaps pour les fonctions exportées ou mixedCaps pour les privées.\nExemple: calculateTotal au lieu de calculate_total",
			funcName)
	}
}
