// Package ktngeneric implements KTN linter rules for generic functions.
package ktngeneric

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeGeneric001 est le code de la regle KTN-GENERIC-001.
	ruleCodeGeneric001 string = "KTN-GENERIC-001"
)

// Analyzer001 checks that generic functions using == or != have comparable constraint.
var Analyzer001 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktngeneric001",
	Doc:      "KTN-GENERIC-001: Generic functions using == or != must have comparable constraint",
	Run:      runGeneric001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runGeneric001 execute l'analyse KTN-GENERIC-001.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: resultat de l'analyse
//   - error: erreur eventuelle
func runGeneric001(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeGeneric001) {
		// Regle desactivee
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeGeneric001, filename) {
			// File is excluded
			return
		}
		funcDecl := n.(*ast.FuncDecl)

		// Analyser la fonction generique
		analyzeGenericFunc(pass, funcDecl)
	})

	// Retour de la fonction
	return nil, nil
}

// analyzeGenericFunc analyse une fonction pour verifier les contraintes generiques.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: declaration de fonction a analyser
func analyzeGenericFunc(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Verifier si la fonction a des parametres de type
	if funcDecl.Type.TypeParams == nil {
		// Pas une fonction generique
		return
	}

	// Collecter les parametres de type avec contrainte "any"
	anyTypeParams := collectAnyTypeParams(funcDecl.Type.TypeParams)

	// Aucun type parameter avec "any"
	if len(anyTypeParams) == 0 {
		// Pas de parametre "any"
		return
	}

	// Verifier si == ou != est utilise sur ces type parameters
	checkEqualityUsage(pass, funcDecl, anyTypeParams)
}

// collectAnyTypeParams collecte les noms des type parameters avec contrainte "any".
//
// Params:
//   - typeParams: liste des parametres de type
//
// Returns:
//   - map[string]bool: map des noms de type parameters avec contrainte "any"
func collectAnyTypeParams(typeParams *ast.FieldList) map[string]bool {
	// Retourner map vide si aucun type parameter
	if typeParams == nil || len(typeParams.List) == 0 {
		return map[string]bool{}
	}

	// Initialiser la map avec capacite estimee
	result := make(map[string]bool, len(typeParams.List))

	// Parcourir les type parameters
	for _, field := range typeParams.List {
		// Verifier si la contrainte est "any"
		if isAnyConstraint(field.Type) {
			// Ajouter les noms de type parameters
			for _, name := range field.Names {
				result[name.Name] = true
			}
		}
	}

	// Retour de la map
	return result
}

// isAnyConstraint verifie si une contrainte est "any".
//
// Params:
//   - expr: expression de contrainte
//
// Returns:
//   - bool: true si la contrainte est "any"
func isAnyConstraint(expr ast.Expr) bool {
	// Verifier si c'est un identifiant "any"
	ident, ok := expr.(*ast.Ident)
	// Verifier l'identifiant
	if !ok {
		// Pas un identifiant
		return false
	}

	// Verifier si c'est "any"
	return ident.Name == "any"
}

// checkEqualityUsage verifie si == ou != est utilise sur les type parameters.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: declaration de fonction
//   - anyTypeParams: map des type parameters avec contrainte "any"
func checkEqualityUsage(pass *analysis.Pass, funcDecl *ast.FuncDecl, anyTypeParams map[string]bool) {
	// Collecter les parametres qui utilisent les type params "any"
	paramNames := collectParamNamesWithAnyType(funcDecl, anyTypeParams)

	// Si aucun parametre n'utilise les types "any"
	if len(paramNames) == 0 {
		// Pas de parametres avec type "any"
		return
	}

	// Parcourir le corps de la fonction
	if funcDecl.Body == nil {
		// Pas de corps de fonction
		return
	}

	// Creer le contexte pour le checking
	ctx := &anyTypeContext{
		paramNames:    paramNames,
		anyTypeParams: anyTypeParams,
		reported:      make(map[string]bool, 1),
	}

	// Parcourir les statements pour trouver == ou !=
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		binaryExpr, ok := n.(*ast.BinaryExpr)
		// Verifier si c'est une expression binaire
		if !ok {
			// Continuer le parcours
			return true
		}

		// Verifier si c'est == ou !=
		if binaryExpr.Op != token.EQL && binaryExpr.Op != token.NEQ {
			// Pas une comparaison d'egalite
			return true
		}

		// Verifier si un des operandes utilise un type parameter "any"
		reportIfUsesAnyTypeParam(pass, funcDecl, binaryExpr, ctx)

		// Continuer le parcours
		return true
	})
}

// collectParamNamesWithAnyType collecte les noms de parametres utilisant un type "any".
//
// Params:
//   - funcDecl: declaration de fonction
//   - anyTypeParams: map des type parameters avec contrainte "any"
//
// Returns:
//   - map[string]string: map des noms de parametres vers leur type parameter
func collectParamNamesWithAnyType(funcDecl *ast.FuncDecl, anyTypeParams map[string]bool) map[string]string {
	// Initialiser avec capacite estimee
	result := make(map[string]string, len(anyTypeParams))

	// Verifier la liste des parametres
	if funcDecl.Type.Params == nil {
		// Pas de parametres
		return result
	}

	// Parcourir les parametres de la fonction
	for _, field := range funcDecl.Type.Params.List {
		typeName := extractTypeName(field.Type)
		// Verifier si le type utilise un type parameter "any"
		if anyTypeParams[typeName] {
			// Ajouter les noms de parametres
			for _, name := range field.Names {
				result[name.Name] = typeName
			}
		}
	}

	// Retour de la map
	return result
}

// extractTypeName extrait le nom du type d'une expression.
//
// Params:
//   - expr: expression de type
//
// Returns:
//   - string: nom du type ou chaine vide
func extractTypeName(expr ast.Expr) string {
	// Determiner le type de l'expression
	switch t := expr.(type) {
	// Type simple
	case *ast.Ident:
		// Retourner le nom du type
		return t.Name
	// Type tableau - extraire le type element
	case *ast.ArrayType:
		// Recursion pour extraire le type de base
		return extractTypeName(t.Elt)
	// Type non supporte
	default:
		// Retourner chaine vide
		return ""
	}
}

// anyTypeContext holds context for any type parameter checking.
type anyTypeContext struct {
	// paramNames maps parameter names to their type parameter
	paramNames map[string]string
	// anyTypeParams maps type parameters with "any" constraint
	anyTypeParams map[string]bool
	// reported tracks already reported errors
	reported map[string]bool
}

// reportIfUsesAnyTypeParam reporte une erreur si un operande utilise un type "any".
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: declaration de fonction
//   - binaryExpr: expression binaire
//   - ctx: context for any type checking
func reportIfUsesAnyTypeParam(
	pass *analysis.Pass,
	funcDecl *ast.FuncDecl,
	binaryExpr *ast.BinaryExpr,
	ctx *anyTypeContext,
) {
	// Verifier l'operande gauche
	leftUses := checkOperandUsesAnyType(binaryExpr.X, ctx.paramNames, ctx.anyTypeParams)
	// Verifier l'operande droit
	rightUses := checkOperandUsesAnyType(binaryExpr.Y, ctx.paramNames, ctx.anyTypeParams)

	// Si aucun operande n'utilise un type "any"
	if !leftUses && !rightUses {
		// Pas d'utilisation de type "any"
		return
	}

	// Creer la cle de deduplication
	funcName := funcDecl.Name.Name
	// Verifier si deja reporte
	if ctx.reported[funcName] {
		// Deja reporte
		return
	}

	// Marquer comme reporte
	ctx.reported[funcName] = true

	// Reporter l'erreur
	cfg := config.Get()
	msg, _ := messages.Get(ruleCodeGeneric001)
	pass.Reportf(
		funcDecl.Pos(),
		"%s: %s",
		ruleCodeGeneric001,
		msg.Format(cfg.Verbose, funcName),
	)
}

// checkOperandUsesAnyType verifie si un operande utilise un type parameter "any".
//
// Params:
//   - expr: expression operande
//   - paramNames: map des noms de parametres vers leur type parameter
//   - anyTypeParams: map des type parameters avec contrainte "any"
//
// Returns:
//   - bool: true si l'operande utilise un type "any"
func checkOperandUsesAnyType(expr ast.Expr, paramNames map[string]string, anyTypeParams map[string]bool) bool {
	// Analyser le type d'expression
	switch typedExpr := expr.(type) {
	// Verifier si c'est un parametre avec type "any"
	case *ast.Ident:
		_, usesAny := paramNames[typedExpr.Name]
		// Retour du resultat
		return usesAny
	// Expression d'indexation (ex: s[i])
	case *ast.IndexExpr:
		// Recursion sur l'expression indexee
		return checkOperandUsesAnyType(typedExpr.X, paramNames, anyTypeParams)
	// Type non supporte
	default:
		// Pas d'utilisation de type any
		return false
	}
}
