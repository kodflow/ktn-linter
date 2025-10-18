package ktn_func

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Rule003 vérifie le format strict de la section Params dans godoc.
//
// KTN-FUNC-003: La section Params: doit documenter tous les paramètres.
// Format requis:
//   // Params:
//   //   - paramName: description du paramètre
//
// Correct:
//   // CalculateTotal calcule le total.
//   //
//   // Params:
//   //   - items: liste des éléments
//   //   - tax: taux de taxe
//   func CalculateTotal(items []int, tax float64) { }
var Rule003 = &analysis.Analyzer{
	Name: "KTN_FUNC_003",
	Doc:  "Vérifie le format strict de la section Params dans godoc",
	Run:  runRule003,
}

// runRule003 exécute la vérification KTN-FUNC-003.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule003(pass *analysis.Pass) (any, error) {
	if isTargetTestFile(pass) {
		return nil, nil
	}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			checkParamsDocumentation(pass, funcDecl)
		}
	}

	return nil, nil
}

// checkParamsDocumentation vérifie la documentation des paramètres.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
func checkParamsDocumentation(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	funcName := funcDecl.Name.Name

	// Vérifier si la fonction a des paramètres
	if funcDecl.Type.Params == nil || countParams(funcDecl.Type.Params) == 0 {
		return
	}

	// Vérifier si la fonction a un godoc
	if funcDecl.Doc == nil || len(funcDecl.Doc.List) == 0 {
		return // Déjà géré par FUNC-002
	}

	doc := funcDecl.Doc.Text()

	// Vérifier la présence de la section Params:
	if !strings.Contains(doc, "Params:") {
		exampleParams := buildParamsExample(funcDecl.Type.Params)
		pass.Reportf(funcDecl.Doc.Pos(),
			"[KTN-FUNC-003] Commentaire godoc doit inclure une section 'Params:' avec format strict.\nExemple:\n  // %s description.\n  //\n  // Params:\n%s\n  func %s(...) { }",
			funcName, exampleParams, funcName)
		return
	}

	// Vérifier que chaque paramètre est documenté
	checkEachParam(pass, funcDecl, doc)
}

// checkEachParam vérifie que chaque paramètre est documenté.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
//   - doc: le texte du commentaire godoc
func checkEachParam(pass *analysis.Pass, funcDecl *ast.FuncDecl, doc string) {
	paramNames := extractParamNames(funcDecl.Type.Params)
	paramsSection := extractSection(doc, "Params:")

	var undocumented []string
	for _, pname := range paramNames {
		found := false
		for _, line := range strings.Split(paramsSection, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "- "+pname+":") {
				found = true
				break
			}
		}
		if !found {
			undocumented = append(undocumented, pname)
		}
	}

	if len(undocumented) > 0 {
		exampleParams := buildParamsExample(funcDecl.Type.Params)
		pass.Reportf(funcDecl.Doc.Pos(),
			"[KTN-FUNC-003] Paramètres non documentés dans la section 'Params:': %s.\nFormat requis:\n  // Params:\n%s",
			strings.Join(undocumented, ", "), exampleParams)
	}
}

// extractSection extrait une section du doc.
//
// Params:
//   - doc: le texte du commentaire godoc
//   - sectionName: le nom de la section
//
// Returns:
//   - string: le contenu de la section
func extractSection(doc, sectionName string) string {
	lines := strings.Split(doc, "\n")
	inSection := false
	var sectionLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == sectionName {
			inSection = true
			continue
		}
		if inSection {
			if strings.HasSuffix(trimmed, ":") && len(trimmed) > 1 {
				break
			}
			if trimmed == "" {
				break
			}
			sectionLines = append(sectionLines, line)
		}
	}

	return strings.Join(sectionLines, "\n")
}

// buildParamsExample construit un exemple de section Params.
//
// Params:
//   - params: la liste des paramètres
//
// Returns:
//   - string: l'exemple formaté
func buildParamsExample(params *ast.FieldList) string {
	paramNames := extractParamNames(params)
	var examples []string
	for _, pname := range paramNames {
		examples = append(examples, "  //   - "+pname+": description du paramètre")
	}
	return strings.Join(examples, "\n")
}

// countParams compte le nombre total de paramètres.
//
// Params:
//   - params: la liste des paramètres
//
// Returns:
//   - int: le nombre total de paramètres
func countParams(params *ast.FieldList) int {
	count := 0
	for _, field := range params.List {
		if len(field.Names) == 0 {
			count++
		} else {
			count += len(field.Names)
		}
	}
	return count
}

// extractParamNames extrait les noms de tous les paramètres.
//
// Params:
//   - params: la liste des paramètres
//
// Returns:
//   - []string: la liste des noms de paramètres
func extractParamNames(params *ast.FieldList) []string {
	var names []string
	for _, field := range params.List {
		for _, name := range field.Names {
			if name.Name != "_" {
				names = append(names, name.Name)
			}
		}
	}
	return names
}
