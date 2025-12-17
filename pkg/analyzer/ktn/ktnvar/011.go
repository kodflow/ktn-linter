// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar011 is the rule code for this analyzer
	ruleCodeVar011 string = "KTN-VAR-011"
)

var (
	// Analyzer011 détecte le shadowing de variables avec := au lieu de =.
	//
	// Le shadowing se produit quand on redéclare une variable avec := alors
	// qu'elle existe déjà dans un scope parent, créant une nouvelle variable
	// locale qui masque l'originale.
	//
	// Exceptions (patterns idiomatiques Go):
	// - "err" : réutilisation courante dans le chaînage d'erreurs
	// - "ok" : pattern idiomatique (_, ok := m[k] ou v, ok := x.(T))
	// - "ctx" : context souvent redéfini dans les sous-scopes
	Analyzer011 *analysis.Analyzer = &analysis.Analyzer{
		Name:     "ktnvar011",
		Doc:      "KTN-VAR-011: Vérifie le shadowing de variables avec := au lieu de =",
		Run:      runVar011,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	// allowedShadowing contains variable names that are allowed to shadow.
	// These are idiomatic Go patterns where shadowing is expected.
	allowedShadowing map[string]bool = map[string]bool{
		"err": true, // Error chaining pattern
		"ok":  true, // Map/type assertion pattern
		"ctx": true, // Context redefinition in sub-scopes
	}
)

// runVar011 exécute l'analyse de détection du shadowing.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: erreur éventuelle
func runVar011(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar011) {
		// Règle désactivée
		return nil, nil
	}

	// Récupération de l'inspecteur AST
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
	}

	// Parcours des assignations
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar011, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

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
			msg, _ := messages.Get(ruleCodeVar011)
			pass.Reportf(
				assign.Pos(),
				"%s: %s",
				ruleCodeVar011,
				msg.Format(config.Get().Verbose, ident.Name),
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
//   - bool: true si shadowing détecté (et non exempté)
func isShadowing(pass *analysis.Pass, ident *ast.Ident) bool {
	// Ignorer les blank identifiers
	if ident.Name == "_" {
		// Traitement
		return false
	}

	// Vérifier si c'est une variable exemptée (patterns idiomatiques)
	if allowedShadowing[ident.Name] {
		// Shadowing autorisé pour cette variable
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
