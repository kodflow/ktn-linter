// Analyzer 013 for the ktnvar package.
package ktnvar

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
)

const (
	// ruleCodeVar013 is the rule code for this analyzer
	ruleCodeVar013 string = "KTN-VAR-013"
)

// Analyzer013 checks that variables are grouped together in a single var block
var Analyzer013 *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktnvar013",
	Doc:  "KTN-VAR-013: Vérifie que les variables de package sont groupées dans un seul bloc var ()",
	Run:  runVar013,
}

// runVar013 exécute l'analyse KTN-VAR-013.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar013(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar013) {
		// Règle désactivée
		return nil, nil
	}

	// Analyze each file independently
	for _, file := range pass.Files {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar013, pass.Fset.Position(file.Pos()).Filename) {
			// Fichier exclu
			continue
		}

		varGroups := collectVarGroups(file)
		checkVarGrouping(pass, varGroups)
	}

	// Retour de la fonction
	return nil, nil
}

// collectVarGroups collecte les déclarations var du fichier.
//
// Params:
//   - file: fichier à analyser
//
// Returns:
//   - []shared.DeclGroup: liste des groupes de variables
func collectVarGroups(file *ast.File) []shared.DeclGroup {
	var varGroups []shared.DeclGroup

	// Collect var declarations
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		// Vérification de la condition
		if !ok {
			continue
		}

		// Only collect var declarations
		if genDecl.Tok == token.VAR {
			varGroups = append(varGroups, shared.DeclGroup{
				Decl: genDecl,
				Pos:  genDecl.Pos(),
			})
		}
	}

	// Retour de la fonction
	return varGroups
}

// checkVarGrouping vérifie le groupement des variables.
//
// Params:
//   - pass: contexte d'analyse
//   - varGroups: groupes de variables à vérifier
func checkVarGrouping(pass *analysis.Pass, varGroups []shared.DeclGroup) {
	// If 0 or 1 var group, they're properly grouped
	if len(varGroups) <= 1 {
		// Retour de la fonction
		return
	}

	// Report all var groups except the first as scattered
	for i := 1; i < len(varGroups); i++ {
		pass.Reportf(
			varGroups[i].Pos,
			"KTN-VAR-013: les variables doivent être groupées ensemble dans un seul bloc var ()",
		)
	}
}
