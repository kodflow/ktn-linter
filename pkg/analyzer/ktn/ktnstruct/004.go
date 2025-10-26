package ktnstruct

import (
	"go/ast"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// MinDocLines nombre minimum de lignes de documentation pour une struct
	MinDocLines = 2
)

// Analyzer004 vérifie que les structs exportées ont une documentation
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct004",
	Doc:      "KTN-STRUCT-004: Toute struct exportée doit avoir une documentation complète",
	Run:      runStruct004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runStruct004 exécute l'analyse KTN-STRUCT-004.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct004(pass *analysis.Pass) (any, error) {
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
			if !hasValidDocumentation(genDecl.Doc, typeSpec.Name.Name) {
				pass.Reportf(
					typeSpec.Pos(),
					"KTN-STRUCT-004: la struct exportée '%s' doit avoir une documentation complète (au moins 2 lignes décrivant son rôle)",
					typeSpec.Name.Name,
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
//
// Returns:
//   - bool: true si documentation valide
func hasValidDocumentation(doc *ast.CommentGroup, structName string) bool {
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

	// Il faut au moins MinDocLines lignes de documentation réelle
	return realLines >= MinDocLines
}
