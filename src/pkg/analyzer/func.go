package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/naming"
)

// FuncAnalyzer vérifie que les fonctions respectent les règles KTN
var FuncAnalyzer = &analysis.Analyzer{
	Name: "ktnfunc",
	Doc:  "Vérifie que les fonctions natives sont bien nommées, documentées et respectent les limites de complexité",
	Run:  runFuncAnalyzer,
}

// runFuncAnalyzer vérifie que toutes les fonctions respectent les règles KTN
// Retourne nil, nil car aucun résultat n'est nécessaire pour cet analyseur
func runFuncAnalyzer(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Ignorer les méthodes (fonctions avec receiver)
			if funcDecl.Recv != nil {
				continue
			}

			checkFunction(pass, file, funcDecl)
		}
	}

	return nil, nil
}

// checkFunction vérifie toutes les règles pour une fonction
// Les paramètres pass, file et funcDecl contrôlent la validation
func checkFunction(pass *analysis.Pass, file *ast.File, funcDecl *ast.FuncDecl) {
	funcName := funcDecl.Name.Name
	checkFuncNaming(pass, funcDecl, funcName)
	checkFuncDocumentation(pass, funcDecl, funcName)
	checkFuncParams(pass, funcDecl, funcName)
	checkFuncLength(pass, funcDecl, funcName)
	checkFuncComplexity(pass, funcDecl, funcName)
	// TODO: Désactivé temporairement - nécessite une meilleure implémentation
	// checkFuncInternalComments(pass, funcDecl, funcName)
	// checkFuncReturnComments(pass, funcDecl, funcName)
}

// checkFuncNaming vérifie le nommage de la fonction
// Les paramètres pass, funcDecl et funcName contrôlent la validation
func checkFuncNaming(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	if !naming.IsMixedCaps(funcName) {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-001] Fonction '%s' n'utilise pas la convention MixedCaps.\nUtilisez MixedCaps pour les fonctions exportées ou mixedCaps pour les privées.\nExemple: calculateTotal au lieu de calculate_total",
			funcName)
	}
}

// checkFuncDocumentation vérifie la documentation de la fonction
// Les paramètres pass, funcDecl et funcName contrôlent la validation
func checkFuncDocumentation(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	if funcDecl.Doc == nil || len(funcDecl.Doc.List) == 0 {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-002] Fonction '%s' sans commentaire godoc.\nToute fonction doit avoir un commentaire godoc.\nExemple:\n  // %s fait quelque chose...\n  func %s(...) { }",
			funcName, funcName, funcName)
	} else {
		checkGodocQuality(pass, funcDecl, funcName)
	}
}

// checkFuncParams vérifie le nombre de paramètres
// Les paramètres pass, funcDecl et funcName contrôlent la validation
func checkFuncParams(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
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

// checkFuncLength vérifie la longueur de la fonction
// Les paramètres pass, funcDecl et funcName contrôlent la validation
func checkFuncLength(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	funcLength := calculateFuncLength(pass.Fset, funcDecl)
	if funcLength > 35 {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-006] Fonction '%s' est trop longue (%d lignes > 35).\nLimitez les fonctions à 35 lignes maximum. Découpez en fonctions plus petites.",
			funcName, funcLength)
	}
}

// checkFuncComplexity vérifie la complexité cyclomatique
// Les paramètres pass, funcDecl et funcName contrôlent la validation
func checkFuncComplexity(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	complexity := calculateCyclomaticComplexity(funcDecl)
	if complexity >= 10 {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-007] Fonction '%s' a une complexité cyclomatique trop élevée (%d >= 10).\nRéduisez la complexité en extrayant des sous-fonctions ou en simplifiant la logique.",
			funcName, complexity)
	}
}

// checkGodocQuality vérifie la qualité du commentaire godoc avec format strict
// Les paramètres pass, funcDecl et funcName contrôlent la validation
func checkGodocQuality(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	doc := funcDecl.Doc.Text()

	// Vérifier que la première ligne commence par le nom de la fonction
	lines := strings.Split(strings.TrimSpace(doc), "\n")
	if len(lines) == 0 || !strings.HasPrefix(strings.TrimSpace(lines[0]), funcName+" ") {
		pass.Reportf(funcDecl.Doc.Pos(),
			"[KTN-FUNC-002] Commentaire godoc doit commencer par le nom de la fonction.\nExemple:\n  // %s fait quelque chose.\n  func %s(...) { }",
			funcName, funcName)
		return
	}

	// Vérifier le format des sections Params et Returns
	checkGodocFormat(pass, funcDecl, funcName, doc)
}

// checkGodocFormat vérifie le format strict avec sections Params: et Returns:
func checkGodocFormat(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string, doc string) {
	hasParams := funcDecl.Type.Params != nil && countParams(funcDecl.Type.Params) > 0
	hasReturns := funcDecl.Type.Results != nil && funcDecl.Type.Results.NumFields() > 0

	// Si pas de params ni returns, juste la description suffit
	if !hasParams && !hasReturns {
		return
	}

	// Vérifier la présence et le format de la section Params:
	if hasParams {
		if !strings.Contains(doc, "Params:") {
			exampleParams := buildParamsExample(funcDecl.Type.Params)
			pass.Reportf(funcDecl.Doc.Pos(),
				"[KTN-FUNC-003] Commentaire godoc doit inclure une section 'Params:' avec format strict.\nExemple:\n  // %s description.\n  //\n  // Params:\n%s\n  func %s(...) { }",
				funcName, exampleParams, funcName)
		} else {
			checkParamsFormat(pass, funcDecl, funcName, doc)
		}
	}

	// Vérifier la présence et le format de la section Returns:
	if hasReturns {
		if !strings.Contains(doc, "Returns:") {
			exampleReturns := buildReturnsExample(funcDecl.Type.Results)
			pass.Reportf(funcDecl.Doc.Pos(),
				"[KTN-FUNC-004] Commentaire godoc doit inclure une section 'Returns:' avec format strict.\nExemple:\n  // %s description.\n  //\n  // Returns:\n%s\n  func %s(...) { }",
				funcName, exampleReturns, funcName)
		}
	}
}

// checkParamsFormat vérifie que chaque paramètre est documenté dans la section Params:
func checkParamsFormat(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string, doc string) {
	paramNames := extractParamNames(funcDecl.Type.Params)

	// Extraire la section Params du doc
	paramsSection := extractSection(doc, "Params:")

	undocumented := []string{}
	for _, pname := range paramNames {
		// Vérifier si le paramètre est mentionné dans la section Params avec le format "  - pname:"
		if !strings.Contains(paramsSection, "- "+pname+":") {
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

// extractSection extrait une section du doc (ex: "Params:" ou "Returns:")
func extractSection(doc, sectionName string) string {
	lines := strings.Split(doc, "\n")
	inSection := false
	var sectionLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, sectionName) {
			inSection = true
			continue
		}
		if inSection {
			// Arrêter si on rencontre une nouvelle section ou ligne vide
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

// buildParamsExample construit un exemple de section Params: pour l'erreur
func buildParamsExample(params *ast.FieldList) string {
	paramNames := extractParamNames(params)
	var examples []string
	for _, pname := range paramNames {
		examples = append(examples, "  //   - "+pname+": description du paramètre")
	}
	return strings.Join(examples, "\n")
}

// buildReturnsExample construit un exemple de section Returns: pour l'erreur
func buildReturnsExample(results *ast.FieldList) string {
	numReturns := results.NumFields()
	var examples []string
	for i := 0; i < numReturns; i++ {
		examples = append(examples, "  //   - type: description du retour")
	}
	return strings.Join(examples, "\n")
}

// hasAnyNamedReturns vérifie si des retours sont nommés
func hasAnyNamedReturns(results *ast.FieldList) bool {
	for _, field := range results.List {
		if len(field.Names) > 0 {
			return true
		}
	}
	return false
}

// countParams compte le nombre total de paramètres
func countParams(params *ast.FieldList) int {
	count := 0
	for _, field := range params.List {
		if len(field.Names) == 0 {
			// Paramètre sans nom (ex: type seul)
			count++
		} else {
			count += len(field.Names)
		}
	}
	return count
}

// extractParamNames extrait les noms de tous les paramètres
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

// calculateFuncLength calcule le nombre de lignes de code d'une fonction
func calculateFuncLength(fset *token.FileSet, funcDecl *ast.FuncDecl) int {
	if funcDecl.Body == nil {
		return 0
	}

	start := fset.Position(funcDecl.Body.Lbrace).Line
	end := fset.Position(funcDecl.Body.Rbrace).Line

	// Exclure les accolades de début et fin
	return end - start - 1
}

// calculateCyclomaticComplexity calcule la complexité cyclomatique d'une fonction
// Complexité = 1 + nombre de points de décision (if, for, case, &&, ||, etc.)
func calculateCyclomaticComplexity(funcDecl *ast.FuncDecl) int {
	if funcDecl.Body == nil {
		return 1
	}

	complexity := 1
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		complexity += getNodeComplexity(n)
		return true
	})

	return complexity
}

// getNodeComplexity retourne la complexité ajoutée par un nœud AST
func getNodeComplexity(n ast.Node) int {
	switch stmt := n.(type) {
	case *ast.IfStmt:
		return 1
	case *ast.ForStmt, *ast.RangeStmt:
		return 1
	case *ast.CaseClause:
		if stmt.List != nil {
			return 1
		}
	case *ast.CommClause:
		if stmt.Comm != nil {
			return 1
		}
	case *ast.BinaryExpr:
		if stmt.Op == token.LAND || stmt.Op == token.LOR {
			return 1
		}
	}
	return 0
}

// checkFuncInternalComments vérifie que les commentaires internes respectent la complexité
// Les paramètres pass, funcDecl et funcName contrôlent la validation
func checkFuncInternalComments(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	// Ignorer les fonctions sans body
	if funcDecl.Body == nil {
		return
	}

	// Calculer la complexité cyclomatique
	complexity := calculateCyclomaticComplexity(funcDecl)

	// Compter les commentaires internes dans le body
	commentCount := countInternalComments(pass.Fset, funcDecl.Body)

	// KTN-FUNC-011: Minimum 1 commentaire par point de complexité
	// On demande au moins (complexity - 1) commentaires car la complexité de base = 1
	minComments := complexity - 1

	// Permettre une tolérance pour les fonctions simples
	if complexity <= 2 {
		return
	}

	// Vérifier si le nombre de commentaires est suffisant
	if commentCount < minComments {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-011] Fonction '%s' manque de commentaires internes (%d commentaires pour une complexité de %d).\nAjoutez au moins 1 commentaire par point de complexité cyclomatique.\nComplexité: %d -> Minimum %d commentaires requis",
			funcName, commentCount, complexity, complexity, minComments)
	}
}

// checkFuncReturnComments vérifie que chaque return a un commentaire
// Les paramètres pass, funcDecl et funcName contrôlent la validation
func checkFuncReturnComments(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	// Ignorer les fonctions sans body
	if funcDecl.Body == nil {
		return
	}

	// Collecter tous les returns
	returns := collectReturnStatements(funcDecl.Body)

	// Ignorer les fonctions sans return ou avec un seul return simple
	if len(returns) == 0 {
		return
	}

	// Vérifier chaque return
	for _, ret := range returns {
		// Vérifier s'il y a un commentaire juste avant le return
		if !hasCommentBeforeReturn(pass.Fset, funcDecl.Body, ret) {
			pass.Reportf(ret.Pos(),
				"[KTN-FUNC-012] Instruction return sans commentaire dans '%s'.\nAjoutez un commentaire avant chaque return pour expliquer ce qui est retourné.\nExemple:\n  // Retour du résultat calculé\n  return result, nil",
				funcName)
		}
	}
}

// countInternalComments compte les commentaires internes dans le body d'une fonction
func countInternalComments(fset *token.FileSet, body *ast.BlockStmt) int {
	count := 0

	// Parcourir tous les statements
	ast.Inspect(body, func(n ast.Node) bool {
		// Ignorer le body lui-même
		if n == body {
			return true
		}

		// Compter les commentaires sur les statements
		switch n.(type) {
		case *ast.AssignStmt, *ast.ExprStmt, *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt:
			// Vérifier s'il y a un commentaire avant ce statement
			// Cette implémentation simple compte chaque statement qui pourrait avoir un commentaire
			count++
		}

		return true
	})

	// Retourner environ la moitié car tous les statements n'ont pas de commentaires
	// C'est une heuristique simple, on peut l'améliorer
	return count / 3
}

// collectReturnStatements collecte tous les returns dans le body
func collectReturnStatements(body *ast.BlockStmt) []*ast.ReturnStmt {
	var returns []*ast.ReturnStmt

	ast.Inspect(body, func(n ast.Node) bool {
		if ret, ok := n.(*ast.ReturnStmt); ok {
			returns = append(returns, ret)
		}
		return true
	})

	return returns
}

// hasCommentBeforeReturn vérifie si un return a un commentaire juste avant
func hasCommentBeforeReturn(fset *token.FileSet, body *ast.BlockStmt, ret *ast.ReturnStmt) bool {
	// Pour l'instant, implémentation simplifiée
	// On considère qu'un return est OK s'il fait partie du dernier statement
	// Une vraie implémentation devrait analyser les CommentMap
	// Pour éviter trop de faux positifs au début, on est permissif
	return false // TODO: Implémenter la vraie vérification avec CommentMap
}
