package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/naming"
)

// Analyzers
var (
	// FuncAnalyzer vérifie que les fonctions respectent les règles KTN
	FuncAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktnfunc",
		Doc:  "Vérifie que les fonctions natives sont bien nommées, documentées et respectent les limites de complexité",
		Run:  runFuncAnalyzer,
	}
)

// runFuncAnalyzer vérifie que toutes les fonctions respectent les règles KTN.
//
// Params:
//   - pass: la passe d'analyse contenant les fichiers à vérifier
//
// Returns:
//   - any: toujours nil car aucun résultat n'est nécessaire
//   - error: toujours nil, les erreurs sont rapportées via pass.Reportf
func runFuncAnalyzer(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Analyser toutes les fonctions ET méthodes
			checkFunction(pass, file, funcDecl)
		}
	}

	// Retourne nil car l'analyseur rapporte via pass.Reportf
	return nil, nil
}

// checkFunction vérifie toutes les règles pour une fonction.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - file: le fichier AST contenant la fonction
//   - funcDecl: la déclaration de fonction à vérifier
func checkFunction(pass *analysis.Pass, file *ast.File, funcDecl *ast.FuncDecl) {
	funcName := funcDecl.Name.Name
	isTestFile := isTestFile(pass)

	checkFuncNaming(pass, funcDecl, funcName)
	checkFuncDocumentation(pass, funcDecl, funcName)
	checkFuncParams(pass, funcDecl, funcName)
	checkFuncLength(pass, funcDecl, funcName, isTestFile)
	checkFuncComplexity(pass, funcDecl, funcName, isTestFile)
	checkNestingDepth(pass, funcDecl, funcName)
	checkReturnComments(pass, file, funcDecl, funcName)
}

// isTestFile vérifie si le fichier analysé est un fichier de test.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - bool: true si c'est un fichier de test (_test.go)
func isTestFile(pass *analysis.Pass) bool {
	for _, f := range pass.Files {
		pos := pass.Fset.Position(f.Pos())
		if strings.HasSuffix(pos.Filename, "_test.go") {
			// Retourne true car un fichier de test a été trouvé
			return true
		}
	}
	// Retourne false car aucun fichier de test n'a été trouvé
	return false
}

// checkFuncNaming vérifie le nommage de la fonction.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - funcDecl: la déclaration de fonction
//   - funcName: le nom de la fonction
func checkFuncNaming(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	if !naming.IsMixedCaps(funcName) {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-001] Fonction '%s' n'utilise pas la convention MixedCaps.\nUtilisez MixedCaps pour les fonctions exportées ou mixedCaps pour les privées.\nExemple: calculateTotal au lieu de calculate_total",
			funcName)
	}
}

// checkFuncDocumentation vérifie la documentation de la fonction.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - funcDecl: la déclaration de fonction
//   - funcName: le nom de la fonction
func checkFuncDocumentation(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	if funcDecl.Doc == nil || len(funcDecl.Doc.List) == 0 {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-002] Fonction '%s' sans commentaire godoc.\nToute fonction doit avoir un commentaire godoc.\nExemple:\n  // %s fait quelque chose...\n  func %s(...) { }",
			funcName, funcName, funcName)
	} else {
		checkGodocQuality(pass, funcDecl, funcName)
	}
}

// checkFuncParams vérifie le nombre de paramètres.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - funcDecl: la déclaration de fonction
//   - funcName: le nom de la fonction
func checkFuncParams(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	if funcDecl.Type.Params == nil {
		// Retourne car la fonction n'a pas de paramètres
		return
	}

	paramCount := countParams(funcDecl.Type.Params)
	if paramCount > 5 {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-005] Fonction '%s' a trop de paramètres (%d > 5).\nLimitez à 5 paramètres maximum. Si nécessaire, utilisez une struct de configuration.\nExemple:\n  type %sConfig struct { ... }\n  func %s(cfg %sConfig) { }",
			funcName, paramCount, funcName, funcName, funcName)
	}
}

// checkFuncLength vérifie la longueur de la fonction.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - funcDecl: la déclaration de fonction
//   - funcName: le nom de la fonction
//   - isTestFile: true si c'est un fichier de test
func checkFuncLength(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string, isTestFile bool) {
	funcLength := calculateFuncLength(pass.Fset, funcDecl)
	maxLength := 35
	if isTestFile {
		maxLength = 100 // Limite plus souple pour les tests
	}

	if funcLength > maxLength {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-006] Fonction '%s' est trop longue (%d lignes > %d).\nLimitez les fonctions à %d lignes maximum. Découpez en fonctions plus petites.",
			funcName, funcLength, maxLength, maxLength)
	}
}

// checkFuncComplexity vérifie la complexité cyclomatique.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - funcDecl: la déclaration de fonction
//   - funcName: le nom de la fonction
//   - isTestFile: true si c'est un fichier de test
func checkFuncComplexity(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string, isTestFile bool) {
	complexity := calculateCyclomaticComplexity(funcDecl)
	maxComplexity := 10
	if isTestFile {
		maxComplexity = 15 // Limite plus souple pour les tests
	}

	if complexity > maxComplexity {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-007] Fonction '%s' a une complexité cyclomatique trop élevée (%d > %d).\nRéduisez la complexité en extrayant des sous-fonctions ou en simplifiant la logique.",
			funcName, complexity, maxComplexity)
	}
}

// checkNestingDepth vérifie la profondeur d'imbrication des blocs.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - funcDecl: la déclaration de fonction
//   - funcName: le nom de la fonction
func checkNestingDepth(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	if funcDecl.Body == nil {
		// Retourne car la fonction n'a pas de corps
		return
	}

	maxDepth := calculateMaxNestingDepth(funcDecl.Body, 0)
	if maxDepth > 3 {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-010] Fonction '%s' a une profondeur d'imbrication trop élevée (%d > 3).\nLimitez l'imbrication à 3 niveaux maximum. Extrayez des sous-fonctions pour réduire la complexité.",
			funcName, maxDepth)
	}
}

// calculateMaxNestingDepth calcule la profondeur maximale d'imbrication dans un bloc.
//
// Params:
//   - node: le nœud AST à analyser
//   - currentDepth: la profondeur actuelle
//
// Returns:
//   - int: la profondeur maximale trouvée
func calculateMaxNestingDepth(node ast.Node, currentDepth int) int {
	maxDepth := currentDepth
	ast.Inspect(node, func(n ast.Node) bool {
		maxDepth = inspectNestingNode(n, maxDepth, currentDepth)
		// Retourne le résultat de shouldContinueInspection pour continuer ou non
		return shouldContinueInspection(n)
	})
	// Retourne la profondeur maximale trouvée
	return maxDepth
}

// inspectNestingNode met à jour la profondeur max pour un nœud d'imbrication.
//
// Params:
//   - n: le nœud à inspecter
//   - currentMax: profondeur maximale actuelle
//   - depth: profondeur courante
//
// Returns:
//   - int: nouvelle profondeur maximale
func inspectNestingNode(n ast.Node, currentMax, depth int) int {
	switch stmt := n.(type) {
	case *ast.IfStmt:
		currentMax = updateMaxDepth(currentMax, stmt.Body, depth)
		if stmt.Else != nil {
			currentMax = updateMaxDepth(currentMax, stmt.Else, depth)
		}
	case *ast.ForStmt:
		currentMax = updateMaxDepth(currentMax, stmt.Body, depth)
	case *ast.RangeStmt:
		currentMax = updateMaxDepth(currentMax, stmt.Body, depth)
	case *ast.SwitchStmt:
		currentMax = updateMaxDepth(currentMax, stmt.Body, depth)
	case *ast.SelectStmt:
		currentMax = updateMaxDepth(currentMax, stmt.Body, depth)
	}
	// Retourne la profondeur maximale mise à jour
	return currentMax
}

// shouldContinueInspection détermine si l'inspection doit continuer.
//
// Params:
//   - n: le nœud à vérifier
//
// Returns:
//   - bool: false pour les structures imbriquées, true sinon
func shouldContinueInspection(n ast.Node) bool {
	switch n.(type) {
	case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.SwitchStmt, *ast.SelectStmt:
		// Retourne false pour arrêter l'inspection des structures imbriquées
		return false
	}
	// Retourne true pour continuer l'inspection
	return true
}

// updateMaxDepth met à jour la profondeur maximale pour un nœud donné.
//
// Params:
//   - currentMax: la profondeur maximale actuelle
//   - node: le nœud à analyser
//   - depth: la profondeur actuelle
//
// Returns:
//   - int: la nouvelle profondeur maximale
func updateMaxDepth(currentMax int, node ast.Node, depth int) int {
	newDepth := calculateMaxNestingDepth(node, depth+1)
	if newDepth > currentMax {
		// Retourne la nouvelle profondeur car elle est supérieure
		return newDepth
	}
	// Retourne la profondeur actuelle car elle reste maximale
	return currentMax
}

// checkGodocQuality vérifie la qualité du commentaire godoc avec format strict.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - funcDecl: la déclaration de fonction
//   - funcName: le nom de la fonction
func checkGodocQuality(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	doc := funcDecl.Doc.Text()

	// Vérifier que la première ligne commence par le nom de la fonction
	lines := strings.Split(strings.TrimSpace(doc), "\n")
	if len(lines) == 0 || !strings.HasPrefix(strings.TrimSpace(lines[0]), funcName+" ") {
		pass.Reportf(funcDecl.Doc.Pos(),
			"[KTN-FUNC-002] Commentaire godoc doit commencer par le nom de la fonction.\nExemple:\n  // %s fait quelque chose.\n  func %s(...) { }",
			funcName, funcName)
		// Retourne car le format de base du godoc est invalide
		return
	}

	// Vérifier le format des sections Params et Returns
	checkGodocFormat(pass, funcDecl, funcName, doc)
}

// checkGodocFormat vérifie le format strict avec sections Params: et Returns:.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - funcDecl: la déclaration de fonction
//   - funcName: le nom de la fonction
//   - doc: le texte du commentaire godoc
func checkGodocFormat(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string, doc string) {
	hasParams := funcDecl.Type.Params != nil && countParams(funcDecl.Type.Params) > 0
	hasReturns := funcDecl.Type.Results != nil && funcDecl.Type.Results.NumFields() > 0

	// Si pas de params ni returns, juste la description suffit
	if !hasParams && !hasReturns {
		// Retourne car aucune section Params/Returns n'est requise
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

// checkParamsFormat vérifie que chaque paramètre est documenté dans la section Params:.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - funcDecl: la déclaration de fonction
//   - funcName: le nom de la fonction
//   - doc: le texte du commentaire godoc
func checkParamsFormat(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string, doc string) {
	paramNames := extractParamNames(funcDecl.Type.Params)

	// Extraire la section Params du doc
	paramsSection := extractSection(doc, "Params:")

	undocumented := []string{}
	for _, pname := range paramNames {
		// Vérifier si le paramètre est mentionné dans la section Params
		// On accepte les espaces optionnels devant le tiret: "  - pname:" ou "- pname:"
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

// extractSection extrait une section du doc (ex: "Params:" ou "Returns:").
//
// Params:
//   - doc: le texte du commentaire godoc complet
//   - sectionName: le nom de la section à extraire (ex: "Params:" ou "Returns:")
//
// Returns:
//   - string: le contenu de la section extraite
func extractSection(doc, sectionName string) string {
	lines := strings.Split(doc, "\n")
	inSection := false
	var sectionLines []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Vérifier si c'est exactement la ligne de section (pas juste qu'elle contient le mot)
		if trimmed == sectionName {
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

	// Retourne le contenu de la section extraite sous forme de chaîne
	return strings.Join(sectionLines, "\n")
}

// buildParamsExample construit un exemple de section Params: pour l'erreur.
//
// Params:
//   - params: la liste des paramètres de la fonction
//
// Returns:
//   - string: l'exemple formaté de section Params: avec tous les paramètres
func buildParamsExample(params *ast.FieldList) string {
	paramNames := extractParamNames(params)
	var examples []string
	for _, pname := range paramNames {
		examples = append(examples, "  //   - "+pname+": description du paramètre")
	}
	// Retourne l'exemple formaté de section Params
	return strings.Join(examples, "\n")
}

// buildReturnsExample construit un exemple de section Returns: pour l'erreur.
//
// Params:
//   - results: la liste des valeurs de retour de la fonction
//
// Returns:
//   - string: l'exemple formaté de section Returns: avec tous les retours
func buildReturnsExample(results *ast.FieldList) string {
	numReturns := results.NumFields()
	var examples []string
	for i := 0; i < numReturns; i++ {
		examples = append(examples, "  //   - type: description du retour")
	}
	// Retourne l'exemple formaté de section Returns
	return strings.Join(examples, "\n")
}

// countParams compte le nombre total de paramètres.
//
// Params:
//   - params: la liste des paramètres à compter
//
// Returns:
//   - int: le nombre total de paramètres (gère les déclarations groupées)
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
	// Retourne le nombre total de paramètres comptés
	return count
}

// extractParamNames extrait les noms de tous les paramètres.
//
// Params:
//   - params: la liste des paramètres dont extraire les noms
//
// Returns:
//   - []string: la liste des noms de paramètres (ignore les underscores)
func extractParamNames(params *ast.FieldList) []string {
	var names []string
	for _, field := range params.List {
		for _, name := range field.Names {
			if name.Name != "_" {
				names = append(names, name.Name)
			}
		}
	}
	// Retourne la liste des noms de paramètres extraits
	return names
}

// calculateFuncLength calcule le nombre de lignes de code d'une fonction.
//
// Params:
//   - fset: le FileSet pour obtenir les positions dans le code source
//   - funcDecl: la déclaration de fonction à mesurer
//
// Returns:
//   - int: le nombre de lignes de code (excluant les accolades de début/fin)
func calculateFuncLength(fset *token.FileSet, funcDecl *ast.FuncDecl) int {
	if funcDecl.Body == nil {
		// Retourne 0 car la fonction n'a pas de corps
		return 0
	}

	start := fset.Position(funcDecl.Body.Lbrace).Line
	end := fset.Position(funcDecl.Body.Rbrace).Line

	// Retourne le nombre de lignes excluant les accolades de début et fin
	return end - start - 1
}

// calculateCyclomaticComplexity calcule la complexité cyclomatique d'une fonction.
// Complexité = 1 + nombre de points de décision (if, for, case, &&, ||, etc.)
//
// Params:
//   - funcDecl: la déclaration de fonction à analyser
//
// Returns:
//   - int: la complexité cyclomatique (minimum 1 pour une fonction sans branchement)
func calculateCyclomaticComplexity(funcDecl *ast.FuncDecl) int {
	if funcDecl.Body == nil {
		// Retourne 1 car une fonction vide a une complexité minimale de 1
		return 1
	}

	complexity := 1
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		complexity += getNodeComplexity(n)
		// Retourne true pour continuer à inspecter tous les nœuds
		return true
	})

	// Retourne la complexité cyclomatique totale calculée
	return complexity
}

// getNodeComplexity retourne la complexité ajoutée par un nœud AST.
//
// Params:
//   - n: le nœud AST à évaluer
//
// Returns:
//   - int: la complexité ajoutée (1 pour if/for/case/&&/||, 0 sinon)
func getNodeComplexity(n ast.Node) int {
	switch stmt := n.(type) {
	case *ast.IfStmt:
		// Retourne 1 car if ajoute un point de décision
		return 1
	case *ast.ForStmt, *ast.RangeStmt:
		// Retourne 1 car les boucles ajoutent un point de décision
		return 1
	case *ast.CaseClause:
		if stmt.List != nil {
			// Retourne 1 car case ajoute un point de décision
			return 1
		}
	case *ast.CommClause:
		if stmt.Comm != nil {
			// Retourne 1 car select case ajoute un point de décision
			return 1
		}
	case *ast.BinaryExpr:
		if stmt.Op == token.LAND || stmt.Op == token.LOR {
			// Retourne 1 car && et || ajoutent un point de décision
			return 1
		}
	}
	// Retourne 0 car ce nœud n'ajoute pas de complexité
	return 0
}

// checkReturnComments vérifie que tous les return statements ont des commentaires.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - file: le fichier AST contenant la fonction
//   - funcDecl: la déclaration de fonction
//   - funcName: le nom de la fonction
func checkReturnComments(pass *analysis.Pass, file *ast.File, funcDecl *ast.FuncDecl, funcName string) {
	if funcDecl.Body == nil {
		// Retourne car la fonction n'a pas de corps
		return
	}

	// Parcourir tous les statements et vérifier les return
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		returnStmt, ok := n.(*ast.ReturnStmt)
		if !ok {
			// Retourne true pour continuer l'inspection
			return true
		}

		// Vérifier si le return a un commentaire juste au-dessus
		if !hasCommentAbove(file, pass.Fset, returnStmt) {
			pass.Reportf(returnStmt.Pos(),
				"[KTN-FUNC-008] Return statement sans commentaire explicatif.\n"+
					"Tout return doit avoir un commentaire juste au-dessus expliquant ce qui est retourné.\n"+
					"Exemple:\n"+
					"  // Erreur de traitement\n"+
					"  return err\n"+
					"\n"+
					"  // Succès\n"+
					"  return nil")
		}
		// Retourne true pour continuer l'inspection des autres returns
		return true
	})
}

// hasCommentAbove vérifie si un return statement a un commentaire juste au-dessus.
//
// Params:
//   - file: le fichier AST contenant le return
//   - fset: le FileSet pour obtenir les positions
//   - returnStmt: le statement return à vérifier
//
// Returns:
//   - bool: true si un commentaire existe juste au-dessus du return
func hasCommentAbove(file *ast.File, fset *token.FileSet, returnStmt *ast.ReturnStmt) bool {
	returnLine := fset.Position(returnStmt.Pos()).Line

	// Parcourir tous les groupes de commentaires du fichier
	for _, commentGroup := range file.Comments {
		if commentGroup == nil || len(commentGroup.List) == 0 {
			continue
		}

		// Vérifier la dernière ligne du groupe de commentaires
		lastComment := commentGroup.List[len(commentGroup.List)-1]
		commentEndLine := fset.Position(lastComment.End()).Line

		// Le commentaire doit être juste au-dessus du return (ligne précédente)
		if commentEndLine == returnLine-1 {
			// Retourne true car un commentaire existe juste au-dessus
			return true
		}
	}

	// Retourne false car aucun commentaire n'a été trouvé juste au-dessus
	return false
}
