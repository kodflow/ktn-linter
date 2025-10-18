package ktn_func

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Rule006 vérifie la longueur des fonctions.
//
// KTN-FUNC-006: Une fonction ne doit pas dépasser 35 lignes (100 pour les tests).
// Les fonctions trop longues doivent être découpées en sous-fonctions.
//
// Incorrect: fonction de 50 lignes
// Correct: fonction de 30 lignes, ou plusieurs fonctions plus petites
var Rule006 = &analysis.Analyzer{
	Name: "KTN_FUNC_006",
	Doc:  "Vérifie que les fonctions ne dépassent pas 35 lignes (100 pour tests)",
	Run:  runRule006,
}

// runRule006 exécute la vérification KTN-FUNC-006.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule006(pass *analysis.Pass) (any, error) {
	isTestFile := isTestFile(pass)

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			checkFuncLength(pass, funcDecl, isTestFile)
		}
	}

	return nil, nil
}

// checkFuncLength vérifie la longueur de la fonction.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
//   - isTestFile: true si c'est un fichier de test
func checkFuncLength(pass *analysis.Pass, funcDecl *ast.FuncDecl, isTestFile bool) {
	funcName := funcDecl.Name.Name
	funcLength := calculateFuncLength(pass.Fset, funcDecl)
	maxLength := 35
	if isTestFile {
		maxLength = 100
	}

	if funcLength > maxLength {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-006] Fonction '%s' est trop longue (%d lignes > %d).\nLimitez les fonctions à %d lignes maximum. Découpez en fonctions plus petites.",
			funcName, funcLength, maxLength, maxLength)
	}
}

// calculateFuncLength calcule le nombre de lignes de code d'une fonction.
//
// Params:
//   - fset: le FileSet
//   - funcDecl: la déclaration de fonction
//
// Returns:
//   - int: le nombre de lignes
func calculateFuncLength(fset *token.FileSet, funcDecl *ast.FuncDecl) int {
	if funcDecl.Body == nil {
		return 0
	}

	start := fset.Position(funcDecl.Body.Lbrace).Line
	end := fset.Position(funcDecl.Body.Rbrace).Line

	return end - start - 1
}

// isTestFile vérifie si le fichier est un fichier de test.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - bool: true si c'est un fichier _test.go
func isTestFile(pass *analysis.Pass) bool {
	for _, f := range pass.Files {
		pos := pass.Fset.Position(f.Pos())
		if strings.HasSuffix(pos.Filename, "_test.go") {
			return true
		}
	}
	return false
}
