package ktn_func

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Rule005 vérifie le nombre de paramètres d'une fonction.
//
// KTN-FUNC-005: Une fonction ne doit pas avoir plus de 5 paramètres.
// Si nécessaire, utilisez une struct de configuration.
//
// Incorrect:
//   func Process(a, b, c, d, e, f int) { } // 6 paramètres
//
// Correct:
//   type ProcessConfig struct {
//       A, B, C, D, E, F int
//   }
//   func Process(cfg ProcessConfig) { }
var Rule005 = &analysis.Analyzer{
	Name: "KTN_FUNC_005",
	Doc:  "Vérifie que les fonctions n'ont pas plus de 5 paramètres",
	Run:  runRule005,
}

// runRule005 exécute la vérification KTN-FUNC-005.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule005(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			checkParamsCount(pass, funcDecl)
		}
	}

	return nil, nil
}

// checkParamsCount vérifie le nombre de paramètres.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
func checkParamsCount(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	funcName := funcDecl.Name.Name

	if funcDecl.Type.Params == nil {
		return
	}

	paramCount := countParams(funcDecl.Type.Params)
	if paramCount > 5 {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-005] Fonction '%s' a trop de paramètres (%d > 5).\nLimitez à 5 paramètres maximum. Si nécessaire, utilisez une struct de configuration.\nExemple:\n  type %sConfig struct { ... }\n  func %s(cfg %sConfig) { }",
			funcName, paramCount, funcName, funcName, funcName)
	}
}
