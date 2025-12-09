// Analyzer 018 for the ktnvar package.
package ktnvar

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer018 checks that variables don't use snake_case naming
var Analyzer018 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar018",
	Doc:      "KTN-VAR-018: Les variables ne doivent pas utiliser snake_case (underscores)",
	Run:      runVar018,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar018 exécute l'analyse KTN-VAR-018.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar018(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

		// Only check var declarations
		if genDecl.Tok != token.VAR {
			// Continue traversing AST nodes
			return
		}

		// Iterate over specifications
		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)
			// Check each variable name
			checkVar018Names(pass, valueSpec)
		}
	})

	// Return analysis result
	return nil, nil
}

// checkVar018Names vérifie les noms de variables pour snake_case.
//
// Params:
//   - pass: contexte d'analyse
//   - valueSpec: spécification de variable
func checkVar018Names(pass *analysis.Pass, valueSpec *ast.ValueSpec) {
	// Iterate over variable names
	for _, name := range valueSpec.Names {
		varName := name.Name

		// Skip blank identifier
		if varName == "_" {
			continue
		}

		// Check for snake_case (lowercase with underscores)
		if isSnakeCase(varName) {
			pass.Reportf(
				name.Pos(),
				"KTN-VAR-018: la variable '%s' utilise snake_case. Utilisez camelCase (ex: '%s')",
				varName,
				snakeToCamel(varName),
			)
		}
	}
}

// isSnakeCase vérifie si un nom utilise snake_case (minuscules avec underscores).
// Exclut SCREAMING_SNAKE_CASE qui est déjà géré par VAR-001.
//
// Params:
//   - name: nom à vérifier
//
// Returns:
//   - bool: true si le nom est en snake_case
func isSnakeCase(name string) bool {
	// Must contain underscore
	if !strings.Contains(name, "_") {
		return false
	}

	// Check if it's NOT all uppercase (SCREAMING_SNAKE_CASE is handled by VAR-001)
	for _, ch := range name {
		// If we find a lowercase letter, it's snake_case
		if ch >= 'a' && ch <= 'z' {
			return true
		}
	}

	// All uppercase = SCREAMING_SNAKE_CASE, not snake_case
	return false
}

// snakeToCamel convertit un nom snake_case en camelCase.
//
// Params:
//   - name: nom en snake_case
//
// Returns:
//   - string: nom en camelCase
func snakeToCamel(name string) string {
	parts := strings.Split(name, "_")
	// Handle empty or single part
	if len(parts) <= 1 {
		return name
	}

	// Build camelCase name
	result := strings.ToLower(parts[0])
	// Iterate over remaining parts
	for i := 1; i < len(parts); i++ {
		part := parts[i]
		// Skip empty parts
		if len(part) == 0 {
			continue
		}
		// Capitalize first letter
		result += strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}

	// Return camelCase name
	return result
}
