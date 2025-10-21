package ktnfunc

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer003 checks for magic numbers (hardcoded numeric literals)
var Analyzer003 = &analysis.Analyzer{
	Name:     "ktnfunc003",
	Doc:      "KTN-FUNC-003: Les nombres littéraux doivent être des constantes nommées (pas de magic numbers)",
	Run:      runFunc003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc003 exécute l'analyse KTN-FUNC-003.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runFunc003(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Map des valeurs autorisées (non magic numbers)
	allowedValues := getAllowedValues()

	// Collecter les littéraux autorisés (const, array sizes)
	allowedLiterals := collectAllowedLiterals(inspect)

	// Vérifier les magic numbers
	checkMagicNumbers(inspect, pass, allowedValues, allowedLiterals)

	// Retour succès
	return nil, nil
}

// getAllowedValues retourne les valeurs numériques autorisées (non magic).
//
// Returns:
//   - map[string]bool: map des valeurs autorisées
func getAllowedValues() map[string]bool {
	// Retour de la map des valeurs autorisées
	return map[string]bool{
		"0":  true,
		"1":  true,
		"-1": true,
	}
}

// collectAllowedLiterals collecte les littéraux dans const declarations.
//
// Params:
//   - inspect: inspecteur AST
//
// Returns:
//   - map[ast.Node]bool: map des littéraux autorisés
func collectAllowedLiterals(inspect *inspector.Inspector) map[ast.Node]bool {
	allowedLiterals := make(map[ast.Node]bool)

	// Filter pour GenDecl seulement
	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

		// Si c'est une déclaration const
		if genDecl.Tok == token.CONST {
			ast.Inspect(genDecl, func(inner ast.Node) bool {
				// Si c'est un littéral
				if lit, ok := inner.(*ast.BasicLit); ok {
					allowedLiterals[lit] = true
				}
				// Continuer l'inspection
				return true
			})
		}
	})

	// Retour de la map
	return allowedLiterals
}

// checkMagicNumbers vérifie et rapporte les magic numbers.
//
// Params:
//   - inspect: inspecteur AST
//   - pass: contexte d'analyse
//   - allowedValues: valeurs autorisées
//   - allowedLiterals: littéraux autorisés
func checkMagicNumbers(inspect *inspector.Inspector, pass *analysis.Pass, allowedValues map[string]bool, allowedLiterals map[ast.Node]bool) {
	// Filter pour les littéraux
	nodeFilter := []ast.Node{
		(*ast.BasicLit)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		lit := n.(*ast.BasicLit)

		// Vérifier si c'est un nombre (INT ou FLOAT)
		if lit.Kind != token.INT && lit.Kind != token.FLOAT {
			// Pas un nombre, ignorer
			return
		}

		// Vérifier si c'est une valeur autorisée
		if allowedValues[lit.Value] {
			// Valeur autorisée, ignorer
			return
		}

		// Vérifier si c'est dans les littéraux autorisés
		if allowedLiterals[lit] {
			// Littéral autorisé, ignorer
			return
		}

		// Reporter l'erreur
		pass.Reportf(
			lit.Pos(),
			"KTN-FUNC-003: le nombre '%s' devrait être une constante nommée (magic number)",
			lit.Value,
		)
	})
}
