// Analyzer 002 for the ktnvar package.
package ktnvar

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer002 checks that package-level variables have explicit type AND value
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar002",
	Doc:      "KTN-VAR-002: Les variables de package doivent avoir le format 'var name type = value'",
	Run:      runVar002,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar002 exécute l'analyse KTN-VAR-002.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar002(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter for File nodes to access package-level declarations
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

		// Check package-level declarations only
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			// Skip if not a GenDecl
			if !ok {
				// Not a general declaration
				continue
			}

			// Only check var declarations
			if genDecl.Tok != token.VAR {
				// Continue traversing AST nodes.
				continue
			}

			// Itération sur les spécifications
			for _, spec := range genDecl.Specs {
				valueSpec := spec.(*ast.ValueSpec)
				// Vérifier si le type est explicite ou visible dans la valeur
				checkVarSpec(pass, valueSpec)
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}

// checkVarSpec vérifie une spécification de variable.
// Style obligatoire: var name type = value
//
// Params:
//   - pass: contexte d'analyse
//   - valueSpec: spécification de variable
func checkVarSpec(pass *analysis.Pass, valueSpec *ast.ValueSpec) {
	hasExplicitType := valueSpec.Type != nil
	hasValues := len(valueSpec.Values) > 0

	// Cas 1: Pas de type explicite = erreur
	if !hasExplicitType {
		// Parcourir les noms
		for _, name := range valueSpec.Names {
			// Ignorer les blank identifiers
			if name.Name == "_" {
				continue
			}

			pass.Reportf(
				name.Pos(),
				"KTN-VAR-002: la variable '%s' doit avoir un type explicite (format: var name type = value)",
				name.Name,
			)
		}
		return
	}

	// Cas 2: Pas de valeur = erreur
	if !hasValues {
		// Parcourir les noms
		for _, name := range valueSpec.Names {
			// Ignorer les blank identifiers
			if name.Name == "_" {
				continue
			}

			pass.Reportf(
				name.Pos(),
				"KTN-VAR-002: la variable '%s' doit être initialisée (format: var name type = value)",
				name.Name,
			)
		}
	}
	// Cas 3: Type explicite ET valeur = OK (style obligatoire respecté)
}

