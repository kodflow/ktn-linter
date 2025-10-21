package ktnvar

import (
	"go/ast"
	"go/token"
	"regexp"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	// Analyzer003 checks that package variables use camelCase naming (not SCREAMING_SNAKE_CASE)
	Analyzer003 *analysis.Analyzer = &analysis.Analyzer{
		Name:     "ktnvar003",
		Doc:      "KTN-VAR-003: Vérifie que les variables de package utilisent camelCase (pas SCREAMING_SNAKE_CASE)",
		Run:      runVar003,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	// screamingSnakeCasePattern matches SCREAMING_SNAKE_CASE names
	// Must contain at least one uppercase letter followed by an underscore or vice versa
	screamingSnakeCasePattern *regexp.Regexp = regexp.MustCompile(`^[A-Z][A-Z0-9_]*[A-Z0-9]$|^[A-Z][A-Z0-9_]*_[A-Z0-9_]*$`)
)

// runVar003 exécute l'analyse KTN-VAR-003.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar003(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

		// Only check var declarations
		if genDecl.Tok != token.VAR {
			// Continue traversing AST nodes
			return
		}

		// Itération sur les spécifications
		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)

			// Itération sur les noms de variables
			for _, name := range valueSpec.Names {
				varName := name.Name

				// Skip blank identifiers
				if varName == "_" {
					continue
				}

				// Check if the variable name uses SCREAMING_SNAKE_CASE (which is wrong for vars)
				if isScreamingSnakeCase(varName) {
					pass.Reportf(
						name.Pos(),
						"KTN-VAR-003: la variable '%s' ne doit pas utiliser SCREAMING_SNAKE_CASE (réservé aux constantes). Utilisez camelCase pour les variables privées ou PascalCase pour les variables exportées",
						varName,
					)
				}
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}

// isScreamingSnakeCase vérifie si un nom utilise SCREAMING_SNAKE_CASE.
//
// Params:
//   - name: nom de la variable à vérifier
//
// Returns:
//   - bool: true si le nom est en SCREAMING_SNAKE_CASE
func isScreamingSnakeCase(name string) bool {
	// Must match the pattern: uppercase letters, digits, and underscores only
	// AND must have at least one underscore OR be all uppercase with length > 1
	if !screamingSnakeCasePattern.MatchString(name) {
		// Not SCREAMING_SNAKE_CASE pattern
		return false
	}

	// Check if name contains underscore (main indicator of SCREAMING_SNAKE_CASE)
	hasUnderscore := false
	// Itération sur les caractères du nom
	for _, ch := range name {
		// Vérification de la présence d'underscore
		if ch == '_' {
			hasUnderscore = true
			break
		}
	}

	// Retour de la fonction
	return hasUnderscore
}
