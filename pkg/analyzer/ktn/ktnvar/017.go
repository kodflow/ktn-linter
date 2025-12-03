// Analyzer 017 for the ktnvar package.
package ktnvar

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer017 détecte le shadowing de variables avec := au lieu de =.
//
// Le shadowing se produit quand on redéclare une variable avec := alors
// qu'elle existe déjà dans un scope parent, créant une nouvelle variable
// locale qui masque l'originale.
var Analyzer017 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar017",
	Doc:      "KTN-VAR-017: Vérifie le shadowing de variables avec := au lieu de =",
	Run:      runVar017,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar017 exécute l'analyse de détection du shadowing.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: erreur éventuelle
func runVar017(pass *analysis.Pass) (any, error) {
	// Récupération de l'inspecteur AST
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
	}

	// Parcours des assignations
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Vérification de l'assignation courte
		checkShortVarDecl(pass, n)
	})

	// Traitement
	return nil, nil
}

// checkShortVarDecl vérifie si une assignation courte fait du shadowing.
//
// Params:
//   - pass: contexte d'analyse
//   - n: nœud AST à analyser
func checkShortVarDecl(pass *analysis.Pass, n ast.Node) {
	// Cast en assignation
	assign, ok := n.(*ast.AssignStmt)
	// Vérification de la condition
	if !ok || assign.Tok != token.DEFINE {
		return // Pas une assignation courte (:=)
	}

	// Pour chaque variable assignée
	for _, lhs := range assign.Lhs {
		// Récupération de l'identifiant
		ident := extractIdent(lhs)
		// Vérification de la condition
		if ident == nil {
			continue
		}

		// Vérification du shadowing
		if isShadowing(pass, ident) {
			// Rapport d'erreur
			pass.Reportf(
				assign.Pos(),
				"KTN-VAR-017: shadowing de la variable '%s' avec ':=' au lieu de '='",
				ident.Name,
			)
		}
	}
}

// extractIdent extrait l'identifiant d'une expression.
//
// Params:
//   - expr: expression à analyser
//
// Returns:
//   - *ast.Ident: identifiant extrait ou nil
func extractIdent(expr ast.Expr) *ast.Ident {
	// Si c'est directement un identifiant
	if ident, ok := expr.(*ast.Ident); ok {
		// Traitement
		return ident
	}
	// Traitement
	return nil
}

// isShadowing vérifie si un identifiant fait du shadowing.
//
// Params:
//   - pass: contexte d'analyse
//   - ident: identifiant à vérifier
//
// Returns:
//   - bool: true si shadowing détecté
func isShadowing(pass *analysis.Pass, ident *ast.Ident) bool {
	// Ignorer les blank identifiers
	if ident.Name == "_" {
		// Traitement
		return false
	}

	// Récupération de l'objet défini (nouvelle définition avec :=)
	obj := pass.TypesInfo.Defs[ident]
	// Vérification de la condition
	if obj == nil {
		return false // Pas une nouvelle définition
	}

	// Récupération du scope de l'objet
	scope := obj.Parent()
	// Vérification de la condition
	if scope == nil {
		// Traitement
		return false
	}

	// Vérification dans le scope parent
	return lookupInParentScope(scope.Parent(), ident.Name)
}

// lookupInParentScope cherche une variable dans le scope parent.
//
// Params:
//   - scope: scope parent à vérifier
//   - name: nom de la variable
//
// Returns:
//   - bool: true si la variable existe dans le scope parent
func lookupInParentScope(scope *types.Scope, name string) bool {
	// Vérification de nil
	if scope == nil {
		// Traitement
		return false
	}

	// Recherche de la variable dans le scope courant
	obj := scope.Lookup(name)
	// Vérification de la condition
	if obj != nil {
		// Traitement
		return true
	}

	// Recherche récursive dans le scope parent
	return lookupInParentScope(scope.Parent(), name)
}
