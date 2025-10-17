package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// ComparisonAnalyzer vérifie les comparaisons redondantes.
	ComparisonAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktncomparison",
		Doc:  "Vérifie les comparaisons redondantes (bool == true/false)",
		Run:  runComparisonAnalyzer,
	}
)

// runComparisonAnalyzer exécute l'analyseur de comparaisons.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runComparisonAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			binaryExpr, ok := n.(*ast.BinaryExpr)
			if !ok {
				// Retourne true pour continuer
				return true
			}

			// Vérifier si c'est une comparaison == ou !=
			if binaryExpr.Op == token.EQL || binaryExpr.Op == token.NEQ {
				checkBooleanComparison(pass, binaryExpr)
			}

			// Retourne true pour continuer l'inspection
			return true
		})
	}
	// Retourne nil car l'analyse est terminée
	return nil, nil
}

// checkBooleanComparison vérifie les comparaisons booléennes redondantes.
//
// Params:
//   - pass: la passe d'analyse
//   - binary: l'expression binaire
func checkBooleanComparison(pass *analysis.Pass, binary *ast.BinaryExpr) {
	// Vérifier si un des côtés est true ou false
	xBool := isBooleanLiteral(binary.X)
	yBool := isBooleanLiteral(binary.Y)

	if !xBool && !yBool {
		// Aucun des côtés n'est un booléen littéral
		// Retourne
		return
	}

	// Déterminer quel côté est le booléen et quel côté est l'expression
	var boolSide ast.Expr
	var exprSide ast.Expr

	if xBool {
		boolSide = binary.X
		exprSide = binary.Y
	} else {
		boolSide = binary.Y
		exprSide = binary.X
	}

	// Vérifier si l'autre côté pourrait être booléen
	if !couldBeBoolean(exprSide) {
		// L'expression n'est probablement pas booléenne
		// Retourne
		return
	}

	reportRedundantBoolComparison(pass, binary, boolSide, exprSide, binary.Op)
}

// isBooleanLiteral vérifie si une expression est true ou false.
//
// Params:
//   - expr: l'expression
//
// Returns:
//   - bool: true si c'est true ou false
func isBooleanLiteral(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "true" || ident.Name == "false"
}

// couldBeBoolean vérifie si une expression pourrait être booléenne.
//
// Params:
//   - expr: l'expression
//
// Returns:
//   - bool: true si potentiellement booléenne
func couldBeBoolean(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		// Variable ou fonction, pourrait être bool
		return true
	case *ast.CallExpr:
		// Appel de fonction, pourrait retourner bool
		return true
	case *ast.SelectorExpr:
		// Champ ou méthode, pourrait être bool
		return true
	case *ast.BinaryExpr:
		// Expression binaire, pourrait être bool (comparaison)
		return e.Op == token.EQL || e.Op == token.NEQ ||
			e.Op == token.LSS || e.Op == token.LEQ ||
			e.Op == token.GTR || e.Op == token.GEQ ||
			e.Op == token.LAND || e.Op == token.LOR
	case *ast.UnaryExpr:
		// Expression unaire, ! donne bool
		return e.Op == token.NOT
	default:
		return false
	}
}

// reportRedundantBoolComparison rapporte une violation KTN-COMP-001.
//
// Params:
//   - pass: la passe d'analyse
//   - binary: l'expression binaire
//   - boolSide: le côté booléen (true/false)
//   - exprSide: le côté expression
//   - op: l'opérateur
func reportRedundantBoolComparison(pass *analysis.Pass, binary *ast.BinaryExpr, boolSide, exprSide ast.Expr, op token.Token) {
	boolIdent, ok := boolSide.(*ast.Ident)
	if !ok {
		return
	}
	boolValue := boolIdent.Name

	// Construire la suggestion
	var suggestion string
	if op == token.EQL {
		if boolValue == "true" {
			suggestion = "utilisez directement l'expression"
		} else {
			suggestion = "utilisez !expression"
		}
	} else { // NEQ
		if boolValue == "true" {
			suggestion = "utilisez !expression"
		} else {
			suggestion = "utilisez directement l'expression"
		}
	}

	pass.Reportf(binary.Pos(),
		"[KTN-COMP-001] Comparaison booléenne redondante avec %s.\n"+
			"Comparer un booléen à true/false est inutile et nuit à la lisibilité.\n"+
			"Utilisez directement l'expression ou sa négation.\n"+
			"Suggestion: %s\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - comparaisons redondantes\n"+
			"  if isValid == true { }     // Lourd\n"+
			"  if isValid != false { }    // Confusant\n"+
			"  if isValid == false { }    // Verbeux\n"+
			"\n"+
			"  // ✅ CORRECT - direct et clair\n"+
			"  if isValid { }             // Simple\n"+
			"  if isValid { }             // Clair\n"+
			"  if !isValid { }            // Élégant",
		boolValue, suggestion)
}
