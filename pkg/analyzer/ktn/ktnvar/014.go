// Analyzer 014 for the ktnvar package.
package ktnvar

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar014 is the rule code for this analyzer
	ruleCodeVar014 string = "KTN-VAR-014"
)

// Analyzer014 checks that package-level variables are declared after constants
var Analyzer014 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar014",
	Doc:      "KTN-VAR-014: Vérifie que les variables de package sont déclarées après les constantes",
	Run:      runVar014,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar014 exécute l'analyse KTN-VAR-014.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar014(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar014) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	// Process each file
	insp.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar014, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}
		varSeen := false

		// Check declarations in order
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			// Skip non-GenDecl (functions, etc.)
			if !ok {
				continue
			}

			// Track variable declarations
			if genDecl.Tok == token.VAR {
				varSeen = true
			}

			// Error: const after var
			if genDecl.Tok == token.CONST && varSeen {
				pass.Reportf(
					genDecl.Pos(),
					"KTN-VAR-014: les constantes doivent être déclarées avant les variables",
				)
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}
