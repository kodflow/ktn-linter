package ktnpackage

import (
	"go/ast"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

// Analyzer001 vérifie que chaque fichier Go (non-test) a une description avant la déclaration package.
var Analyzer001 = &analysis.Analyzer{
	Name: "ktnpackage001",
	Doc:  "KTN-PACKAGE-001: chaque fichier .go doit avoir un commentaire descriptif avant la déclaration package",
	Run:  runPackage001,
}

// runPackage001 exécute l'analyse KTN-PACKAGE-001.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runPackage001(pass *analysis.Pass) (any, error) {
	// Parcourir tous les fichiers
	for _, file := range pass.Files {
		// Récupérer le nom du fichier
		filename := pass.Fset.Position(file.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			// Skip les fichiers de test
			continue
		}

		// Vérifier si le fichier a un commentaire de package
		hasFileComment := checkFileComment(file)

		// Si pas de commentaire, reporter l'erreur
		if !hasFileComment {
			// Reporter à la position de la déclaration package
			pass.Reportf(
				file.Name.Pos(),
				"KTN-PACKAGE-001: le fichier doit avoir un commentaire descriptif avant 'package %s'. Ajouter une description légère du contenu/rôle du fichier",
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
//
// Returns:
//   - bool: true si le fichier a un commentaire valide
func checkFileComment(file *ast.File) bool {
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

		// Si le commentaire contient du texte (au moins 3 caractères), c'est valide
		if len(text) >= 3 {
			// Commentaire valide trouvé
			return true
		}
	}

	// Aucun commentaire valide
	return false
}
