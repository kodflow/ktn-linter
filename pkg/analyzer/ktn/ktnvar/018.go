package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// INITIAL_CONVERSIONS_CAP est la capacité initiale pour la map de conversions
	INITIAL_CONVERSIONS_CAP int = 10
	// MAX_ALLOWED_CONVERSIONS est le nombre maximum de conversions tolérées
	MAX_ALLOWED_CONVERSIONS int = 2
)

// Analyzer018 détecte les conversions string() répétées.
//
// Les conversions string([]byte) répétées créent des allocations inutiles.
// Il vaut mieux convertir une seule fois et réutiliser la variable.
var Analyzer018 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar018",
	Doc:      "KTN-VAR-018: Vérifie les conversions string() répétées",
	Run:      runVar018,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar018 exécute l'analyse de détection des conversions répétées.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: erreur éventuelle
func runVar018(pass *analysis.Pass) (interface{}, error) {
	// Récupération de l'inspecteur AST
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Parcours des fonctions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Vérification des conversions dans la fonction
		checkFuncForRepeatedConversions(pass, n)
	})

	// Traitement
	return nil, nil
}

// checkFuncForRepeatedConversions vérifie une fonction entière.
//
// Params:
//   - pass: contexte d'analyse
//   - n: nœud AST à analyser
func checkFuncForRepeatedConversions(pass *analysis.Pass, n ast.Node) {
	// Cast en fonction
	funcDecl, ok := n.(*ast.FuncDecl)
	// Vérification de la condition
	if !ok || funcDecl.Body == nil {
		// Traitement
		return
	}

	// Analyse des boucles
	checkLoopsForStringConversion(pass, funcDecl.Body)

	// Analyse des conversions multiples
	checkMultipleConversions(pass, funcDecl.Body)
}

// checkLoopsForStringConversion détecte string() dans les boucles.
//
// Params:
//   - pass: contexte d'analyse
//   - body: corps de la fonction
func checkLoopsForStringConversion(pass *analysis.Pass, body *ast.BlockStmt) {
	// Parcours de tous les statements
	ast.Inspect(body, func(n ast.Node) bool {
		// Vérification des boucles
		loop := extractLoop(n)
		// Vérification de la condition
		if loop == nil {
			// Traitement
			return true
		}

		// Vérification des conversions dans la boucle
		if hasStringConversion(loop) {
			// Rapport d'erreur
			pass.Reportf(
				loop.Pos(),
				"KTN-VAR-018: conversion string() répétée dans la boucle, préallouer hors de la boucle",
			)
		}

		// Traitement
		return true
	})
}

// extractLoop extrait le corps d'une boucle.
//
// Params:
//   - n: nœud AST
//
// Returns:
//   - ast.Node: corps de la boucle ou nil
func extractLoop(n ast.Node) ast.Node {
	// Switch sur le type de nœud
	switch loop := n.(type) {
	// Traitement
	case *ast.ForStmt:
		// Traitement
		return loop.Body
	// Traitement
	case *ast.RangeStmt:
		// Traitement
		return loop.Body
	// Traitement
	default:
		// Traitement
		return nil
	}
}

// hasStringConversion vérifie si un nœud contient string().
//
// Params:
//   - n: nœud AST
//
// Returns:
//   - bool: true si conversion trouvée
func hasStringConversion(n ast.Node) bool {
	// Recherche de conversions
	found := false

	ast.Inspect(n, func(node ast.Node) bool {
		// Vérification des appels de fonction
		if isStringConversion(node) {
			found = true
			// Traitement
			return false
		}
		// Traitement
		return true
	})

	// Traitement
	return found
}

// isStringConversion vérifie si un nœud est une conversion string().
//
// Params:
//   - n: nœud AST
//
// Returns:
//   - bool: true si c'est une conversion string()
func isStringConversion(n ast.Node) bool {
	// Cast en appel
	call, ok := n.(*ast.CallExpr)
	// Vérification de la condition
	if !ok {
		// Traitement
		return false
	}

	// Vérification du type de fonction
	ident, ok := call.Fun.(*ast.Ident)
	// Vérification de la condition
	if !ok {
		// Traitement
		return false
	}

	// Vérification que c'est "string"
	if ident.Name != "string" {
		// Traitement
		return false
	}

	// Vérification qu'il y a exactement 1 argument
	if len(call.Args) != 1 {
		// Traitement
		return false
	}

	// Vérification que l'argument est un []byte
	return true
}

// checkMultipleConversions détecte les conversions multiples.
//
// Params:
//   - pass: contexte d'analyse
//   - body: corps de la fonction
func checkMultipleConversions(pass *analysis.Pass, body *ast.BlockStmt) {
	// Map pour compter les conversions par variable
	conversions := make(map[string]int, INITIAL_CONVERSIONS_CAP)
	var firstPos map[string]ast.Node = make(map[string]ast.Node, INITIAL_CONVERSIONS_CAP)

	// Parcours pour compter
	ast.Inspect(body, func(n ast.Node) bool {
		// Skip les boucles (déjà traité)
		if extractLoop(n) != nil {
			// Traitement
			return false
		}

		// Vérification des conversions
		if call, ok := n.(*ast.CallExpr); ok && isStringConversion(call) {
			// Extraction de la variable convertie
			varName := utils.ExtractVarName(call.Args[0])
			// Vérification de la condition
			if varName != "" {
				conversions[varName]++
				// Vérification de la condition
				if _, exists := firstPos[varName]; !exists {
					firstPos[varName] = call
				}
			}
		}

		// Traitement
		return true
	})

	// Rapport des conversions multiples
	for varName, count := range conversions {
		// Vérification de la condition
		if count > MAX_ALLOWED_CONVERSIONS {
			pos := firstPos[varName]
			pass.Reportf(
				pos.Pos(),
				"KTN-VAR-018: conversion string() de '%s' répétée %d fois, stocker dans une variable",
				varName,
				count,
			)
		}
	}
}
