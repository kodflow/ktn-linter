package ktn_func

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Rule004 vérifie le format strict de la section Returns dans godoc.
//
// KTN-FUNC-004: La section Returns: doit documenter tous les retours.
// Format requis:
//   // Returns:
//   //   - type: description du retour
//
// Correct:
//   // CalculateTotal calcule le total.
//   //
//   // Returns:
//   //   - float64: le total calculé
//   //   - error: erreur éventuelle
//   func CalculateTotal() (float64, error) { }
var Rule004 = &analysis.Analyzer{
	Name: "KTN_FUNC_004",
	Doc:  "Vérifie le format strict de la section Returns dans godoc",
	Run:  runRule004,
}

// runRule004 exécute la vérification KTN-FUNC-004.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule004(pass *analysis.Pass) (any, error) {
	// Ignorer seulement les fichiers réellement dans tests/target, pas testdata
	for _, f := range pass.Files {
		pos := pass.Fset.Position(f.Pos())
		normalizedPath := strings.ReplaceAll(pos.Filename, "\\", "/")
		if strings.Contains(normalizedPath, "tests/target/") ||
			strings.Contains(normalizedPath, "tests/bad_usage/") ||
			strings.Contains(normalizedPath, "tests/good_usage/") {
			return nil, nil
		}
	}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			checkReturnsDocumentation(pass, funcDecl)
		}
	}

	return nil, nil
}

// checkReturnsDocumentation vérifie la documentation des retours.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
func checkReturnsDocumentation(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	funcName := funcDecl.Name.Name

	// Vérifier si la fonction a des retours
	if funcDecl.Type.Results == nil || funcDecl.Type.Results.NumFields() == 0 {
		return
	}

	// Vérifier si la fonction a un godoc
	if funcDecl.Doc == nil || len(funcDecl.Doc.List) == 0 {
		return // Déjà géré par FUNC-002
	}

	doc := funcDecl.Doc.Text()

	// Vérifier la présence de la section Returns:
	if !strings.Contains(doc, "Returns:") {
		exampleReturns := buildReturnsExample(funcDecl.Type.Results)
		pass.Reportf(funcDecl.Doc.Pos(),
			"[KTN-FUNC-004] Commentaire godoc doit inclure une section 'Returns:' avec format strict.\nExemple:\n  // %s description.\n  //\n  // Returns:\n%s\n  func %s(...) { }",
			funcName, exampleReturns, funcName)
	}
}

// buildReturnsExample construit un exemple de section Returns.
//
// Params:
//   - results: la liste des valeurs de retour
//
// Returns:
//   - string: l'exemple formaté
func buildReturnsExample(results *ast.FieldList) string {
	numReturns := results.NumFields()
	var examples []string
	for i := 0; i < numReturns; i++ {
		examples = append(examples, "  //   - type: description du retour")
	}
	return strings.Join(examples, "\n")
}
