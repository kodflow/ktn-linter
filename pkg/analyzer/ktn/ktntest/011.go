// Analyzer 011 for the ktntest package.
package ktntest

import (
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCode est le code de la règle.
	ruleCodeTest011 string = "KTN-TEST-011"
)

// Analyzer011 checks package naming convention for internal/external test files
var Analyzer011 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktntest011",
	Doc:      "KTN-TEST-011: Les fichiers _internal_test.go doivent utiliser 'package xxx', les fichiers _external_test.go doivent utiliser 'package xxx_test'",
	Run:      runTest011,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest011 exécute l'analyse KTN-TEST-011.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest011(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeTest011) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Valider les fichiers de test
	validateTestFiles011(pass, cfg)

	// Définir le filtre de nœuds
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	// Parcourir les nœuds
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Rien à faire dans preorder
	})

	// Retour de la fonction
	return nil, nil
}

// validateTestFiles011 valide les conventions de package pour les fichiers de test.
//
// Params:
//   - pass: contexte d'analyse
//   - cfg: configuration
func validateTestFiles011(pass *analysis.Pass, cfg *config.Config) {
	// Itération sur les fichiers
	for _, file := range pass.Files {
		// Obtenir le chemin du fichier
		filename := pass.Fset.Position(file.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeTest011, filename) {
			// Fichier exclu
			continue
		}

		// Extraire le nom de base
		basename := filepath.Base(filename)

		// Vérification si test
		if !shared.IsTestFile(basename) {
			// Continuer si pas test
			continue
		}

		// Extraire le package actuel
		actualPkg := file.Name.Name

		// Valider le package selon le type de fichier
		validatePackageConvention011(pass, file, basename, actualPkg)
	}
}

// validatePackageConvention011 valide la convention de package selon le type de fichier test.
//
// Params:
//   - pass: contexte d'analyse
//   - file: nœud AST du fichier
//   - basename: nom de base du fichier
//   - actualPkg: nom du package actuel
func validatePackageConvention011(pass *analysis.Pass, file *ast.File, basename, actualPkg string) {
	// Vérification suffixe internal
	if strings.HasSuffix(basename, "_internal_test.go") {
		// Vérification package _test
		if basePkg, ok := strings.CutSuffix(actualPkg, "_test"); ok {
			// Signaler erreur package
			pass.Reportf(
				file.Name.Pos(),
				"KTN-TEST-011: le fichier '%s' doit utiliser 'package %s' (white-box testing) au lieu de 'package %s'. Les fichiers _internal_test.go testent les fonctions privées et doivent partager le même package",
				basename,
				basePkg,
				actualPkg,
			)
		}
	} else {
		// Vérification suffixe external
		if strings.HasSuffix(basename, "_external_test.go") {
			// Cas alternatif: external
			// Vérification package sans _test
			if !strings.HasSuffix(actualPkg, "_test") {
				// Extraire package attendu
				expectedPkg := extractExpectedPackageFromFilename(basename)
				// Signaler erreur package
				pass.Reportf(
					file.Name.Pos(),
					"KTN-TEST-011: le fichier '%s' doit utiliser 'package %s_test' (black-box testing) au lieu de 'package %s'. Les fichiers _external_test.go testent l'API publique et doivent utiliser un package externe",
					basename,
					expectedPkg,
					actualPkg,
				)
			}
		}
	}
}

// extractExpectedPackageFromFilename extrait le nom de package attendu depuis le nom de fichier.
//
// Params:
//   - filename: nom du fichier (ex: calculator_external_test.go)
//
// Returns:
//   - string: nom du package attendu (ex: calculator)
func extractExpectedPackageFromFilename(filename string) string {
	// Retirer _external_test.go
	baseName := strings.TrimSuffix(filename, "_external_test.go")
	// Retour du nom de base
	return baseName
}
