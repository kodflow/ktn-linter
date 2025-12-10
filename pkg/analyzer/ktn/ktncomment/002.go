// Analyzer 002 for the ktncomment package.
package ktncomment

import (
	"go/ast"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
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
			// Reporter à la position de la déclaration package
			pass.Reportf(
				file.Name.Pos(),
				"KTN-COMMENT-002: le fichier doit avoir un commentaire descriptif avant 'package %s'. Ajouter une description légère du contenu/rôle du fichier",
				file.Name.Name,
			)
		}
	}

	// Retour de la fonction
	return nil, nil
}

// checkFileComment vérifie si le fichier a un commentaire avant la déclaration package.
//
// Params:
//   - file: fichier AST à analyser
//   - minLength: longueur minimale requise
//
// Returns:
//   - bool: true si le fichier a un commentaire valide
func checkFileComment(file *ast.File, minLength int) bool {
	// Vérifier s'il y a des commentaires dans le fichier
	if file.Doc == nil || len(file.Doc.List) == 0 {
		// Pas de commentaire de package
		return false
	}

	// Vérifier que le commentaire n'est pas vide et contient du texte utile
	for _, comment := range file.Doc.List {
		text := strings.TrimSpace(comment.Text)

		// Enlever les // ou /* */
		text = strings.TrimPrefix(text, "//")
		text = strings.TrimPrefix(text, "/*")
		text = strings.TrimSuffix(text, "*/")
		text = strings.TrimSpace(text)

		// Si le commentaire contient du texte (au moins minLength chars)
		if len(text) >= minLength {
			// Commentaire valide trouvé
			return true
		}
	}

	// Aucun commentaire valide
	return false
}
