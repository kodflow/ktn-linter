package analyzer

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/naming"
)

// FuncAnalyzer vérifie que les fonctions respectent les règles KTN strictes
// Règles implémentées :
// - KTN-FUNC-001 : Nom en MixedCaps/mixedCaps (pas de snake_case)
// - KTN-FUNC-002 : Fonction exportée doit avoir un commentaire godoc
// - KTN-FUNC-003 : Commentaire godoc doit documenter les paramètres
// - KTN-FUNC-004 : Commentaire godoc doit documenter les valeurs de retour
// - KTN-FUNC-005 : Maximum 5 paramètres
// - KTN-FUNC-006 : Maximum 35 lignes de code
// - KTN-FUNC-007 : Complexité cyclomatique < 10
// - KTN-FUNC-008 : Pas de préfixe "Get" inutile
// - KTN-FUNC-009 : Initialismes corrects (HTTP, URL, ID, etc.)
// - KTN-FUNC-010 : Context en premier paramètre
var FuncAnalyzer = &analysis.Analyzer{
	Name: "ktnfunc",
	Doc:  "Vérifie que les fonctions natives respectent les bonnes pratiques Go strictes",
	Run:  runFuncAnalyzer,
}

func runFuncAnalyzer(pass *analysis.Pass) (interface{}, error) {
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

			// Analyser la fonction native
			checkFunction(pass, file, funcDecl)
		}
	}

	return nil, nil
}

// checkFunction vérifie toutes les règles pour une fonction
func checkFunction(pass *analysis.Pass, file *ast.File, funcDecl *ast.FuncDecl) {
	funcName := funcDecl.Name.Name
	isExported := unicode.IsUpper(rune(funcName[0]))

	// KTN-FUNC-001 : Vérifier MixedCaps/mixedCaps (pas de snake_case)
	if !naming.IsMixedCaps(funcName) {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-001] Fonction '%s' n'utilise pas la convention MixedCaps.\nUtilisez MixedCaps pour les fonctions exportées ou mixedCaps pour les privées.\nExemple: ParseHTTPRequest au lieu de parse_http_request",
			funcName)
	}

	// KTN-FUNC-008 : Vérifier préfixe "Get" inutile
	if naming.HasGetterPrefix(funcName) {
		suggestedName := strings.TrimPrefix(funcName, "Get")
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-008] Fonction '%s' a un préfixe 'Get' inutile.\nEn Go, les getters n'utilisent pas le préfixe 'Get'.\nExemple: Utilisez '%s()' au lieu de '%s()'",
			funcName, suggestedName, funcName)
	}

	// KTN-FUNC-009 : Vérifier initialismes incorrects
	suggestions := naming.FixInitialisms(funcName)
	if len(suggestions) > 0 {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-009] Fonction '%s' contient des initialismes incorrects.\nLes initialismes Go doivent être en majuscules : HTTP, URL, ID, API, JSON, XML, etc.\nSuggestion: %s",
			funcName, suggestions[0])
	}

	// KTN-FUNC-002 : Vérifier commentaire godoc pour fonctions exportées
	if isExported {
		if funcDecl.Doc == nil || len(funcDecl.Doc.List) == 0 {
			pass.Reportf(funcDecl.Name.Pos(),
				"[KTN-FUNC-002] Fonction exportée '%s' sans commentaire godoc.\nToute fonction exportée doit avoir un commentaire godoc.\nExemple:\n  // %s fait quelque chose...\n  func %s(...) { }",
				funcName, funcName, funcName)
		} else {
			// Vérifier la qualité du commentaire godoc
			checkGodocQuality(pass, funcDecl, funcName)
		}
	}

	// KTN-FUNC-005 : Vérifier le nombre de paramètres
	if funcDecl.Type.Params != nil {
		paramCount := countParams(funcDecl.Type.Params)
		if paramCount > 5 {
			pass.Reportf(funcDecl.Name.Pos(),
				"[KTN-FUNC-005] Fonction '%s' a trop de paramètres (%d > 5).\nLimitez à 5 paramètres maximum. Si nécessaire, utilisez une struct de configuration.\nExemple:\n  type %sConfig struct { ... }\n  func %s(cfg %sConfig) { }",
				funcName, paramCount, funcName, funcName, funcName)
		}
	}

	// KTN-FUNC-010 : Vérifier que Context est en premier paramètre
	checkContextPosition(pass, funcDecl, funcName)

	// KTN-FUNC-006 : Vérifier la longueur de la fonction
	funcLength := calculateFuncLength(pass.Fset, funcDecl)
	if funcLength > 35 {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-006] Fonction '%s' est trop longue (%d lignes > 35).\nLimitez les fonctions à 35 lignes maximum. Découpez en fonctions plus petites.",
			funcName, funcLength)
	}

	// KTN-FUNC-007 : Vérifier la complexité cyclomatique
	complexity := calculateCyclomaticComplexity(funcDecl)
	if complexity >= 10 {
		pass.Reportf(funcDecl.Name.Pos(),
			"[KTN-FUNC-007] Fonction '%s' a une complexité cyclomatique trop élevée (%d >= 10).\nRéduisez la complexité en extrayant des sous-fonctions ou en simplifiant la logique.",
			funcName, complexity)
	}
}

// checkGodocQuality vérifie la qualité du commentaire godoc
func checkGodocQuality(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	doc := funcDecl.Doc.Text()

	// KTN-FUNC-003 : Vérifier que les paramètres sont documentés (si > 2 paramètres)
	if funcDecl.Type.Params != nil {
		paramCount := countParams(funcDecl.Type.Params)
		if paramCount > 2 {
			// Extraire les noms de paramètres
			paramNames := extractParamNames(funcDecl.Type.Params)
			undocumented := []string{}
			for _, pname := range paramNames {
				if !strings.Contains(doc, pname) {
					undocumented = append(undocumented, pname)
				}
			}

			if len(undocumented) > 0 {
				pass.Reportf(funcDecl.Doc.Pos(),
					"[KTN-FUNC-003] Commentaire godoc de '%s' ne documente pas les paramètres: %s.\nDocumentez chaque paramètre pour clarifier leur rôle.",
					funcName, strings.Join(undocumented, ", "))
			}
		}
	}

	// KTN-FUNC-004 : Vérifier que les valeurs de retour sont documentées (si retours nommés ou > 1 retour)
	if funcDecl.Type.Results != nil && funcDecl.Type.Results.NumFields() > 0 {
		numResults := funcDecl.Type.Results.NumFields()
		hasNamedReturns := false
		for _, field := range funcDecl.Type.Results.List {
			if len(field.Names) > 0 {
				hasNamedReturns = true
				break
			}
		}

		// Vérifier si "error" ou "retourne" est mentionné
		hasReturnDoc := strings.Contains(strings.ToLower(doc), "retourne") ||
			strings.Contains(strings.ToLower(doc), "return") ||
			strings.Contains(strings.ToLower(doc), "erreur") ||
			strings.Contains(strings.ToLower(doc), "error")

		// Si plus de 1 retour OU retours nommés, documentation obligatoire
		if (numResults > 1 || hasNamedReturns) && !hasReturnDoc {
			pass.Reportf(funcDecl.Doc.Pos(),
				"[KTN-FUNC-004] Commentaire godoc de '%s' ne documente pas les valeurs de retour.\nDocumentez ce que retourne la fonction et dans quels cas.",
				funcName)
		}
	}
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

// checkContextPosition vérifie que context.Context est en premier paramètre
func checkContextPosition(pass *analysis.Pass, funcDecl *ast.FuncDecl, funcName string) {
	if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 {
		return
	}

	// Chercher si context.Context est présent
	contextIndex := -1
	for i, field := range funcDecl.Type.Params.List {
		if isContextType(field.Type) {
			contextIndex = i
			break
		}
	}

	// Si context.Context est présent mais pas en premier
	if contextIndex > 0 {
		pass.Reportf(funcDecl.Type.Params.List[contextIndex].Pos(),
			"[KTN-FUNC-010] Paramètre 'context.Context' de '%s' n'est pas en première position.\nLe paramètre context.Context doit toujours être le premier paramètre.\nExemple: func %s(ctx context.Context, ...) { }",
			funcName, funcName)
	}
}

// isContextType vérifie si un type est context.Context
func isContextType(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	return ident.Name == "context" && sel.Sel.Name == "Context"
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
		switch stmt := n.(type) {
		case *ast.IfStmt:
			complexity++
		case *ast.ForStmt, *ast.RangeStmt:
			complexity++
		case *ast.CaseClause:
			// Chaque case ajoute 1 (sauf default)
			if stmt.List != nil {
				complexity++
			}
		case *ast.CommClause:
			// Chaque case de select ajoute 1
			if stmt.Comm != nil {
				complexity++
			}
		case *ast.BinaryExpr:
			// && et || ajoutent de la complexité
			if stmt.Op == token.LAND || stmt.Op == token.LOR {
				complexity++
			}
		}
		return true
	})

	return complexity
}
