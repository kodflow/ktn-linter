package ktn_func

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Rule002 vérifie la documentation godoc des fonctions.
//
// KTN-FUNC-002: Toute fonction doit avoir un commentaire godoc.
// Le commentaire doit commencer par le nom de la fonction.
//
// Incorrect: fonction sans commentaire ou commentaire mal formaté
// Correct:
//
//	// CalculateTotal calcule le total des éléments.
//	func CalculateTotal() { }
var Rule002 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_FUNC_002",
	Doc:  "Vérifie que les fonctions ont un commentaire godoc",
	Run:  runRule002,
}

// runRule002 exécute la vérification KTN-FUNC-002.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule002(pass *analysis.Pass) (any, error) {
	// Ignorer les fichiers de test dans tests/target/**
	if isTargetTestFile(pass) {
		return nil, nil
	}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			checkFuncDocumentation(pass, funcDecl)
		}
	}

	return nil, nil
}

// checkFuncDocumentation vérifie la documentation de la fonction.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - funcDecl: la déclaration de fonction
func checkFuncDocumentation(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	funcName := funcDecl.Name.Name

	if funcDecl.Doc == nil || len(funcDecl.Doc.List) == 0 {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-002] Fonction '%s' sans commentaire godoc.\nToute fonction doit avoir un commentaire godoc.\nExemple:\n  // %s fait quelque chose...\n  func %s(...) { }",
			funcName, funcName, funcName)
		return
	}

	// Vérifier que le commentaire commence par le nom de la fonction
	doc := funcDecl.Doc.Text()
	lines := strings.Split(strings.TrimSpace(doc), "\n")
	if len(lines) == 0 || !strings.HasPrefix(strings.TrimSpace(lines[0]), funcName+" ") {
		pass.Reportf(funcDecl.Doc.Pos(),
			"[KTN-FUNC-002] Commentaire godoc doit commencer par le nom de la fonction.\nExemple:\n  // %s fait quelque chose.\n  func %s(...) { }",
			funcName, funcName)
	}
}

// isTargetTestFile vérifie si le fichier analysé est dans tests/.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - bool: true si c'est un fichier dans tests/
func isTargetTestFile(pass *analysis.Pass) bool {
	for _, f := range pass.Files {
		pos := pass.Fset.Position(f.Pos())
		normalizedPath := strings.ReplaceAll(pos.Filename, "\\", "/")
		if strings.Contains(normalizedPath, "tests/target/") ||
			strings.Contains(normalizedPath, "tests/bad_usage/") ||
			strings.Contains(normalizedPath, "tests/good_usage/") {
			return true
		}
	}
	return false
}
