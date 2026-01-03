// Package ktngeneric implements KTN linter rules for generic functions.
package ktngeneric

import (
	"go/ast"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeGeneric002 est le code de la regle KTN-GENERIC-002.
	ruleCodeGeneric002 string = "KTN-GENERIC-002"
)

// Analyzer002 checks that generics are not used unnecessarily on interface types.
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktngeneric002",
	Doc:      "KTN-GENERIC-002: Detect unnecessary generics on interface types",
	Run:      runGeneric002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runGeneric002 execute l'analyse KTN-GENERIC-002.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: resultat de l'analyse
//   - error: erreur eventuelle
func runGeneric002(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeGeneric002) {
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
		if cfg.IsFileExcluded(ruleCodeGeneric002, filename) {
			// File is excluded
			return
		}
		funcDecl := n.(*ast.FuncDecl)

		// Analyser la fonction generique
		analyzeUnnecessaryGeneric(pass, funcDecl)
	})

	// Retour de la fonction
	return nil, nil
}

// analyzeUnnecessaryGeneric analyse une fonction pour detecter les generiques inutiles.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: declaration de fonction a analyser
func analyzeUnnecessaryGeneric(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Verifier si la fonction a des parametres de type
	if funcDecl.Type.TypeParams == nil {
		// Pas une fonction generique
		return
	}

	// Analyser chaque parametre de type
	for _, field := range funcDecl.Type.TypeParams.List {
		// Verifier si la contrainte est une interface simple
		if !isSingleInterfaceConstraint(pass, field.Type) {
			// Pas une contrainte interface simple
			continue
		}

		// Verifier chaque nom de type parameter dans le champ
		for _, name := range field.Names {
			// Verifier si le type parameter est utilise uniquement pour les signatures
			if isTypeParamOnlyForSignature(funcDecl, name.Name) {
				// Reporter l'erreur
				reportUnnecessaryGeneric(pass, funcDecl, name.Name, field.Type)
			}
		}
	}
}

// isSingleInterfaceConstraint verifie si la contrainte est une interface simple.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression de contrainte
//
// Returns:
//   - bool: true si c'est une interface simple (pas any/comparable)
func isSingleInterfaceConstraint(pass *analysis.Pass, expr ast.Expr) bool {
	// Exclure les contraintes generiques "any" et "comparable"
	if isGenericBuiltinConstraint(expr) {
		// Pas une interface specifique
		return false
	}

	// Guard contre nil (pour tests unitaires)
	if pass == nil || pass.TypesInfo == nil {
		// Pas d'info de type disponible
		return false
	}

	// Obtenir le type de la contrainte
	typeInfo := pass.TypesInfo.TypeOf(expr)
	// Verifier si le type est valide
	if typeInfo == nil {
		// Type non resolu
		return false
	}

	// Verifier si c'est une interface
	return isInterfaceType(typeInfo)
}

// isGenericBuiltinConstraint verifie si la contrainte est any ou comparable.
//
// Params:
//   - expr: expression de contrainte
//
// Returns:
//   - bool: true si c'est any ou comparable
func isGenericBuiltinConstraint(expr ast.Expr) bool {
	// Verifier si c'est un identifiant simple
	ident, ok := expr.(*ast.Ident)
	// Verifier l'identifiant
	if !ok {
		// Pas un identifiant simple
		return false
	}
	// Verifier si c'est "any" ou "comparable"
	return ident.Name == "any" || ident.Name == "comparable"
}

// isInterfaceType verifie si un type est une interface.
//
// Params:
//   - t: type a verifier
//
// Returns:
//   - bool: true si c'est une interface
func isInterfaceType(t types.Type) bool {
	// Check nil type
	if t == nil {
		// Return false for nil type
		return false
	}
	// Dereference le type nomme si necessaire
	underlying := t.Underlying()
	// Verifier si c'est une interface
	_, ok := underlying.(*types.Interface)
	// Retourner le resultat
	return ok
}

// isTypeParamOnlyForSignature verifie si le type parameter est utilise pour preservation.
//
// Params:
//   - funcDecl: declaration de fonction
//   - typeParamName: nom du type parameter
//
// Returns:
//   - bool: true si utilise uniquement dans la signature
func isTypeParamOnlyForSignature(funcDecl *ast.FuncDecl, typeParamName string) bool {
	// Verifier si le type est retourne
	usedInReturn := isTypeInReturnType(funcDecl, typeParamName)
	// Verifier si le type est utilise dans le corps
	usedInBody := isTypeUsedInBody(funcDecl, typeParamName)

	// Si utilise dans return OU dans le corps, c'est justifie
	if usedInReturn || usedInBody {
		// Usage justifie pour preservation de type
		return false
	}

	// Si non utilise ni dans return ni dans body, c'est inutile
	return true
}

// isTypeInReturnType verifie si le type parameter est dans le type de retour.
//
// Params:
//   - funcDecl: declaration de fonction
//   - typeParamName: nom du type parameter
//
// Returns:
//   - bool: true si le type est dans le retour
func isTypeInReturnType(funcDecl *ast.FuncDecl, typeParamName string) bool {
	// Verifier les resultats de la fonction
	if funcDecl.Type.Results == nil {
		// Pas de retour
		return false
	}

	// Parcourir les types de retour
	for _, field := range funcDecl.Type.Results.List {
		// Verifier si le type parameter est utilise
		if containsTypeParam(field.Type, typeParamName) {
			// Type parameter dans le retour
			return true
		}
	}

	// Type parameter non trouve dans le retour
	return false
}

// containsTypeParam verifie si une expression contient un type parameter.
//
// Params:
//   - expr: expression a analyser
//   - typeParamName: nom du type parameter
//
// Returns:
//   - bool: true si le type parameter est present
func containsTypeParam(expr ast.Expr, typeParamName string) bool {
	found := false
	// Parcourir l'expression
	ast.Inspect(expr, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		// Verifier si c'est l'identifiant recherche
		if ok && ident.Name == typeParamName {
			// Type parameter trouve
			found = true
			// Arreter le parcours
			return false
		}
		// Continuer le parcours
		return true
	})
	// Retourner le resultat
	return found
}

// isTypeUsedInBody verifie si le type parameter est utilise dans le corps.
//
// Params:
//   - funcDecl: declaration de fonction
//   - typeParamName: nom du type parameter
//
// Returns:
//   - bool: true si utilise dans le corps
func isTypeUsedInBody(funcDecl *ast.FuncDecl, typeParamName string) bool {
	// Verifier le corps de la fonction
	if funcDecl.Body == nil {
		// Pas de corps
		return false
	}

	found := false
	// Parcourir le corps
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		// Detecter les assertions de type avec le type parameter
		if checkTypeAssertionUsage(n, typeParamName) {
			// Usage trouve
			found = true
			// Arreter le parcours
			return false
		}
		// Detecter les conversions explicites
		if checkTypeConversionUsage(n, typeParamName) {
			// Usage trouve
			found = true
			// Arreter le parcours
			return false
		}
		// Continuer le parcours
		return true
	})
	// Retourner le resultat
	return found
}

// checkTypeAssertionUsage verifie si un noeud est une assertion de type.
//
// Params:
//   - n: noeud a verifier
//   - typeParamName: nom du type parameter
//
// Returns:
//   - bool: true si c'est une assertion avec le type parameter
func checkTypeAssertionUsage(n ast.Node, typeParamName string) bool {
	typeAssert, ok := n.(*ast.TypeAssertExpr)
	// Verifier si c'est une assertion de type
	if !ok {
		// Pas une assertion de type
		return false
	}
	// Verifier si le type parameter est utilise
	return containsTypeParam(typeAssert.Type, typeParamName)
}

// checkTypeConversionUsage verifie si un noeud est une conversion de type.
//
// Params:
//   - n: noeud a verifier
//   - typeParamName: nom du type parameter
//
// Returns:
//   - bool: true si c'est une conversion avec le type parameter
func checkTypeConversionUsage(n ast.Node, typeParamName string) bool {
	callExpr, ok := n.(*ast.CallExpr)
	// Verifier si c'est un appel
	if !ok {
		// Pas un appel
		return false
	}
	// Verifier si la fonction est un type parameter
	ident, ok := callExpr.Fun.(*ast.Ident)
	// Verifier l'identifiant
	if !ok {
		// Pas un identifiant simple
		return false
	}
	// Verifier si c'est le type parameter
	return ident.Name == typeParamName
}

// reportUnnecessaryGeneric reporte une erreur pour un generique inutile.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: declaration de fonction
//   - typeParamName: nom du type parameter
//   - constraintExpr: expression de contrainte
func reportUnnecessaryGeneric(
	pass *analysis.Pass,
	funcDecl *ast.FuncDecl,
	typeParamName string,
	constraintExpr ast.Expr,
) {
	// Guard contre nil (pour tests unitaires)
	if pass == nil {
		// Pas de contexte pour reporter
		return
	}

	// Obtenir le nom de la contrainte
	constraintName := extractConstraintName(constraintExpr)
	// Construire le message
	cfg := config.Get()
	msg, _ := messages.Get(ruleCodeGeneric002)
	// Reporter l'erreur
	pass.Reportf(
		funcDecl.Pos(),
		"%s: %s",
		ruleCodeGeneric002,
		msg.Format(cfg.Verbose, funcDecl.Name.Name, typeParamName, constraintName),
	)
}

// extractConstraintName extrait le nom de la contrainte d'une expression.
//
// Params:
//   - expr: expression de contrainte
//
// Returns:
//   - string: nom de la contrainte
func extractConstraintName(expr ast.Expr) string {
	// Analyser le type d'expression pour extraire le nom
	switch typedExpr := expr.(type) {
	// Identifiant simple
	case *ast.Ident:
		// Retourner le nom de l'identifiant
		return typedExpr.Name
	// Expression qualifiee (ex: io.Reader)
	case *ast.SelectorExpr:
		// Extraire le nom qualifie complet
		return extractSelectorName(typedExpr)
	// Type non supporte
	default:
		// Retourner nom generique
		return "interface"
	}
}

// extractSelectorName extrait le nom d'une expression selector.
//
// Params:
//   - sel: expression selector
//
// Returns:
//   - string: nom qualifie (ex: io.Reader)
func extractSelectorName(sel *ast.SelectorExpr) string {
	// Obtenir le nom du package
	pkgIdent, ok := sel.X.(*ast.Ident)
	// Verifier si c'est un identifiant
	if !ok {
		// Pas un identifiant simple
		return sel.Sel.Name
	}
	// Retourner le nom qualifie
	return pkgIdent.Name + "." + sel.Sel.Name
}
