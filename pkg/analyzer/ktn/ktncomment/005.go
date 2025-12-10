// Analyzer 005 for the ktncomment package.
package ktncomment

import (
	"go/ast"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeComment005 is the rule code for this analyzer
	ruleCodeComment005 string = "KTN-COMMENT-005"
	// defaultMinStructDocLines minimum lines for struct doc
	defaultMinStructDocLines int = 2
)

// Analyzer005 vérifie que les structs exportées ont une documentation
var Analyzer005 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktncomment005",
	Doc:      "KTN-COMMENT-005: Toute struct exportée doit avoir une documentation complète",
	Run:      runComment005,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runComment005 exécute l'analyse KTN-COMMENT-005.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runComment005(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeComment005) {
		// Règle désactivée
		return nil, nil
	}

	// Récupérer le seuil configuré
	minDocLines := cfg.GetThreshold(ruleCodeComment005, defaultMinStructDocLines)

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filtrer les nœuds de type GenDecl
	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Cast vers GenDecl
		genDecl := n.(*ast.GenDecl)

		// Ignorer les fichiers de test
		filename := pass.Fset.Position(genDecl.Pos()).Filename

		// Vérifier si le fichier est exclu
		if cfg.IsFileExcluded(ruleCodeComment005, filename) {
			// Fichier exclu
			return
		}

		// Vérification si fichier test
		if shared.IsTestFile(filename) {
			// Continue avec le nœud suivant
			return
		}

		// Vérifier chaque spec dans la déclaration
		for _, spec := range genDecl.Specs {
			// Vérifier si c'est une TypeSpec
			typeSpec, ok := spec.(*ast.TypeSpec)
			// Si ce n'est pas une TypeSpec, continuer
			if !ok {
				// Continue avec le spec suivant
				continue
			}

			// Vérifier si c'est une struct
			_, isStruct := typeSpec.Type.(*ast.StructType)
			// Si ce n'est pas une struct, continuer
			if !isStruct {
				// Continue avec le spec suivant
				continue
			}

			// Vérifier si c'est exporté
			if !ast.IsExported(typeSpec.Name.Name) {
				// Struct privée, pas besoin de doc
				continue
			}

			// Vérifier la documentation
			if !hasValidDocumentation(genDecl.Doc, typeSpec.Name.Name, minDocLines) {
				pass.Reportf(
					typeSpec.Pos(),
					"KTN-COMMENT-005: la struct exportée '%s' doit avoir une documentation complète (au moins %d lignes décrivant son rôle)",
					typeSpec.Name.Name,
					minDocLines,
				)
			}
		}
	})

	// Retour de la fonction
	return nil, nil
}

// hasValidDocumentation vérifie si la documentation est valide.
//
// Params:
//   - doc: groupe de commentaires
//   - structName: nom de la struct
//   - minLines: nombre minimum de lignes requis
//
// Returns:
//   - bool: true si documentation valide
func hasValidDocumentation(doc *ast.CommentGroup, structName string, minLines int) bool {
	// Si pas de documentation
	if doc == nil || len(doc.List) == 0 {
		// Documentation manquante
		return false
	}

	// Compter les lignes de documentation réelle
	var realLines int
	// Parcourir les commentaires
	for _, comment := range doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		text = strings.TrimSpace(strings.TrimPrefix(text, "/*"))
		text = strings.TrimSpace(strings.TrimSuffix(text, "*/"))

		// Ignorer les lignes vides
		if text == "" {
			// Continue avec le commentaire suivant
			continue
		}

		// Vérifier que la première ligne commence par le nom de la struct
		if realLines == 0 {
			// La première ligne doit commencer par le nom
			if !strings.HasPrefix(text, structName) {
				// Format invalide
				return false
			}
		}

		// Compter cette ligne
		realLines++
	}

	// Il faut au moins minLines lignes de documentation réelle
	return realLines >= minLines
}
