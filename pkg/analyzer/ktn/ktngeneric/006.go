// Package ktngeneric implements KTN linter rules for generic functions.
package ktngeneric

import (
	"go/ast"
	"go/token"
	"maps"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeGeneric006 est le code de la regle KTN-GENERIC-006.
	ruleCodeGeneric006 string = "KTN-GENERIC-006"
)

// Analyzer006 checks that generic functions using ordered/arithmetic ops have proper constraint.
var Analyzer006 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktngeneric006",
	Doc:      "KTN-GENERIC-006: Generic functions using <, >, +, -, *, /, % must have cmp.Ordered constraint",
	Run:      runGeneric006,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runGeneric006 execute l'analyse KTN-GENERIC-006.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: resultat de l'analyse
//   - error: erreur eventuelle
func runGeneric006(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeGeneric006) {
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
		if cfg.IsFileExcluded(ruleCodeGeneric006, filename) {
			// File is excluded
			return
		}
		funcDecl := n.(*ast.FuncDecl)

		// Analyser la fonction generique pour les operateurs ordered
		analyzeGenericFuncOrdered(pass, funcDecl)
	})

	// Retour de la fonction
	return nil, nil
}

// analyzeGenericFuncOrdered analyse une fonction pour les operateurs ordered/arithmetiques.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: declaration de fonction a analyser
func analyzeGenericFuncOrdered(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
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

	// Verifier si des operateurs ordered sont utilises sur ces type parameters
	checkOrderedUsage(pass, funcDecl, anyTypeParams)
}

// isOrderedOp verifie si un operateur est un operateur ordered ou arithmetique.
//
// Params:
//   - op: operateur token
//
// Returns:
//   - bool: true si l'operateur est ordered ou arithmetique
func isOrderedOp(op token.Token) bool {
	// Verifier si l'operateur necessite cmp.Ordered
	switch op {
	// Operateur de comparaison ordered
	case token.LSS, token.LEQ, token.GTR, token.GEQ:
		// Necessite cmp.Ordered
		return true
	// Operateur arithmetique
	case token.ADD, token.SUB, token.MUL, token.QUO, token.REM:
		// Necessite cmp.Ordered
		return true
	// Autre operateur
	default:
		// Ne necessite pas cmp.Ordered
		return false
	}
}

// checkOrderedUsage verifie si des operateurs ordered sont utilises sur les type parameters.
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: declaration de fonction
//   - anyTypeParams: map des type parameters avec contrainte "any"
func checkOrderedUsage(pass *analysis.Pass, funcDecl *ast.FuncDecl, anyTypeParams map[string]bool) {
	// Collecter les parametres qui utilisent les type params "any"
	paramNames := collectParamNamesWithAnyType(funcDecl, anyTypeParams)

	// Collecter les variables locales qui ont un type "any"
	localVars := collectLocalVarsWithAnyType(funcDecl, anyTypeParams)

	// Merger les maps
	allNames := mergeStringMaps(paramNames, localVars)

	// Si aucun parametre n'utilise les types "any"
	if len(allNames) == 0 {
		// Pas de parametres avec type "any"
		return
	}

	// Parcourir le corps de la fonction
	if funcDecl.Body == nil {
		// Pas de corps de fonction
		return
	}

	// Creer le contexte pour le checking
	ctx := &orderedTypeContext{
		paramNames:    allNames,
		anyTypeParams: anyTypeParams,
		reported:      make(map[string]bool, 1),
	}

	// Parcourir les statements pour trouver les operateurs ordered
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		binaryExpr, ok := n.(*ast.BinaryExpr)
		// Verifier si c'est une expression binaire
		if !ok {
			// Continuer le parcours
			return true
		}

		// Verifier si c'est un operateur ordered ou arithmetique
		if !isOrderedOp(binaryExpr.Op) {
			// Pas un operateur ordered
			return true
		}

		// Verifier si un des operandes utilise un type parameter "any"
		reportIfUsesAnyTypeParamOrdered(pass, funcDecl, binaryExpr, ctx)

		// Continuer le parcours
		return true
	})
}

// collectLocalVarsWithAnyType collecte les variables locales avec type "any".
//
// Params:
//   - funcDecl: declaration de fonction
//   - anyTypeParams: map des type parameters avec contrainte "any"
//
// Returns:
//   - map[string]string: map des noms de variables vers leur type parameter
func collectLocalVarsWithAnyType(funcDecl *ast.FuncDecl, anyTypeParams map[string]bool) map[string]string {
	// Initialiser avec capacite estimee
	result := make(map[string]string, len(anyTypeParams))

	// Verifier le corps de la fonction
	if funcDecl.Body == nil {
		// Pas de corps
		return result
	}

	// Parcourir les statements
	for _, stmt := range funcDecl.Body.List {
		// Extraire les declarations de variables
		extractVarDeclsFromStmt(stmt, anyTypeParams, result)
	}

	// Retour de la map
	return result
}

// extractVarDeclsFromStmt extrait les declarations de variables d'un statement.
//
// Params:
//   - stmt: statement a analyser
//   - anyTypeParams: map des type parameters avec contrainte "any"
//   - result: map resultat a remplir
func extractVarDeclsFromStmt(stmt ast.Stmt, anyTypeParams map[string]bool, result map[string]string) {
	// Determiner le type de statement
	switch s := stmt.(type) {
	// Declaration statement
	case *ast.DeclStmt:
		extractFromDeclStmt(s, anyTypeParams, result)
	// Range statement (for _, v := range ...)
	case *ast.RangeStmt:
		extractFromRangeStmt(s, anyTypeParams, result)
	}
}

// extractFromDeclStmt extrait les variables d'un DeclStmt.
//
// Params:
//   - s: DeclStmt a analyser
//   - anyTypeParams: map des type parameters
//   - result: map resultat
func extractFromDeclStmt(s *ast.DeclStmt, anyTypeParams map[string]bool, result map[string]string) {
	genDecl, ok := s.Decl.(*ast.GenDecl)
	// Verifier si c'est une declaration generale
	if !ok {
		// Pas une GenDecl
		return
	}

	// Parcourir les specs
	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		// Verifier si c'est une ValueSpec
		if !ok {
			// Pas une ValueSpec
			continue
		}

		// Sans type explicite (ex: var x = ...), on ne peut pas inferer
		if valueSpec.Type == nil {
			// Skip untyped var specs
			continue
		}

		typeName := extractTypeName(valueSpec.Type)
		// Verifier si le type est un type parameter "any"
		if anyTypeParams[typeName] {
			// Ajouter les noms
			for _, name := range valueSpec.Names {
				result[name.Name] = typeName
			}
		}
	}
}

// extractFromRangeStmt extrait les variables d'un RangeStmt.
// Note: Sans acces a pass.TypesInfo, on ne peut pas determiner le type element.
//
// Params:
//   - _s: RangeStmt a analyser (non utilise - necessite TypesInfo)
//   - _anyTypeParams: map des type parameters (non utilise - necessite TypesInfo)
//   - _result: map resultat (non utilise - necessite TypesInfo)
func extractFromRangeStmt(_s *ast.RangeStmt, _anyTypeParams map[string]bool, _result map[string]string) {
	// Note: Implementation complete necessite pass.TypesInfo
	// Pour l'instant, on ne peut pas determiner le type des variables de range
	// sans information de type du compilateur
}

// mergeStringMaps fusionne deux maps string->string.
//
// Params:
//   - m1: premiere map
//   - m2: deuxieme map
//
// Returns:
//   - map[string]string: map fusionnee
func mergeStringMaps(m1, m2 map[string]string) map[string]string {
	// Cloner m1 comme base
	result := maps.Clone(m1)

	// maps.Copy panic si result est nil
	if result == nil {
		// Creer une map avec capacite de m2
		result = make(map[string]string, len(m2))
	}

	// Ajouter m2 au resultat
	maps.Copy(result, m2)

	// Retour de la map
	return result
}

// orderedTypeContext holds context for ordered type parameter checking.
type orderedTypeContext struct {
	// paramNames maps parameter names to their type parameter
	paramNames map[string]string
	// anyTypeParams maps type parameters with "any" constraint
	anyTypeParams map[string]bool
	// reported tracks already reported errors
	reported map[string]bool
}

// reportIfUsesAnyTypeParamOrdered reporte une erreur si un operande utilise un type "any".
//
// Params:
//   - pass: contexte d'analyse
//   - funcDecl: declaration de fonction
//   - binaryExpr: expression binaire
//   - ctx: context for ordered type checking
func reportIfUsesAnyTypeParamOrdered(
	pass *analysis.Pass,
	funcDecl *ast.FuncDecl,
	binaryExpr *ast.BinaryExpr,
	ctx *orderedTypeContext,
) {
	// Guard contre nil (pour tests unitaires)
	if pass == nil {
		// Pas de contexte pour reporter
		return
	}

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
	msg, _ := messages.Get(ruleCodeGeneric006)
	pass.Reportf(
		funcDecl.Pos(),
		"%s: %s",
		ruleCodeGeneric006,
		msg.Format(cfg.Verbose, funcName),
	)
}
