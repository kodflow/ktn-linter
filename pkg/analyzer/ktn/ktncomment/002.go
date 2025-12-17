// Package ktncomment provides analyzers for comment formatting rules.
package ktncomment

import (
	"go/ast"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
)

const (
	// ruleCodeComment002 is the rule code for this analyzer
	ruleCodeComment002 string = "KTN-COMMENT-002"
	// defaultMinPackageCommentLength minimum length for a valid package comment
	defaultMinPackageCommentLength int = 3
)

// Analyzer002 checks that each Go file has a package description comment.
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktncomment002",
	Doc:  "KTN-COMMENT-002: chaque fichier .go doit avoir un commentaire descriptif avant la déclaration package",
	Run:  runComment002,
}

// runComment002 exécute l'analyse KTN-COMMENT-002.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runComment002(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeComment002) {
		// Règle désactivée
		return nil, nil
	}

	// Récupérer le seuil configuré
	minLength := cfg.GetThreshold(ruleCodeComment002, defaultMinPackageCommentLength)

	// Parcourir tous les fichiers
	for _, file := range pass.Files {
		// Récupérer le nom du fichier
		filename := pass.Fset.Position(file.Pos()).Filename

		// Vérifier si le fichier est exclu
		if cfg.IsFileExcluded(ruleCodeComment002, filename) {
			// Fichier exclu
			continue
		}

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Skip les fichiers de test
			continue
		}

		// Vérifier si le fichier a un commentaire de package
		hasFileComment := checkFileComment(file, minLength)

		// Si pas de commentaire, reporter l'erreur
		if !hasFileComment {
			msg, _ := messages.Get(ruleCodeComment002)
			// Reporter à la position de la déclaration package
			pass.Reportf(
				file.Name.Pos(),
				"%s: %s",
				ruleCodeComment002,
				msg.Format(cfg.Verbose, file.Name.Name),
			)
		}
	}

	// Retour de la fonction
	return nil, nil
}

// checkFileComment vérifie si le fichier a un commentaire avant la déclaration package.
// Le commentaire doit respecter le format Go: "Package <name> ..."
//
// Params:
//   - file: fichier AST à analyser
//   - minLength: longueur minimale requise
//
// Returns:
//   - bool: true si le fichier a un commentaire valide au format "Package <name> ..."
func checkFileComment(file *ast.File, minLength int) bool {
	// Vérifier s'il y a des commentaires dans le fichier
	if file.Doc == nil || len(file.Doc.List) == 0 {
		// Pas de commentaire de package
		return false
	}

	// Obtenir le nom du package pour validation
	pkgName := file.Name.Name
	expectedPrefix := "Package " + pkgName

	// Vérifier que le commentaire respecte le format Go standard
	for _, comment := range file.Doc.List {
		text := strings.TrimSpace(comment.Text)

		// Enlever les // ou /* */
		text = strings.TrimPrefix(text, "//")
		text = strings.TrimPrefix(text, "/*")
		text = strings.TrimSuffix(text, "*/")
		text = strings.TrimSpace(text)

		// Si le commentaire contient du texte (au moins minLength chars)
		if len(text) >= minLength {
			// Vérifier le format "Package <name> ..." (Go standard)
			if strings.HasPrefix(text, expectedPrefix) {
				// Format Go standard respecté
				return true
			}
		}
	}

	// Aucun commentaire valide au format Go standard
	return false
}
