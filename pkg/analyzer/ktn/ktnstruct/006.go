// Analyzer 006 for the ktnstruct package.
package ktnstruct

import (
	"go/ast"
	"unicode"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeStruct006 code de la règle KTN-STRUCT-006
	ruleCodeStruct006 string = "KTN-STRUCT-006"
)

// Analyzer006 checks that DTO private fields don't have serialization tags
var Analyzer006 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct006",
	Doc:      "KTN-STRUCT-006: Les champs privés d'un DTO ne doivent pas avoir de tags de sérialisation",
	Run:      runStruct006,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runStruct006 exécute l'analyse KTN-STRUCT-006.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct006(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeStruct006) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		typeSpec := n.(*ast.TypeSpec)

		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeStruct006, filename) {
			// Fichier exclu
			return
		}

		// Vérifier si c'est une struct
		structType, ok := typeSpec.Type.(*ast.StructType)
		// Si pas une struct, ignorer
		if !ok {
			// Retour anticipé
			return
		}

		// Vérifier si c'est un DTO
		if !shared.IsSerializableStruct(structType, typeSpec.Name.Name) {
			// Pas un DTO, ignorer
			// Retour anticipé
			return
		}

		// Vérifier les champs privés avec tags
		checkPrivateFieldsWithTags(pass, structType, typeSpec.Name.Name)
	})

	// Retour de la fonction
	return nil, nil
}

// checkPrivateFieldsWithTags vérifie les champs privés avec tags.
//
// Params:
//   - pass: contexte d'analyse
//   - structType: type de la struct
//   - structName: nom de la struct
func checkPrivateFieldsWithTags(pass *analysis.Pass, structType *ast.StructType, structName string) {
	// Vérifier si la struct a des champs
	if structType.Fields == nil {
		// Retour anticipé
		return
	}

	// Parcourir les champs
	for _, field := range structType.Fields.List {
		// Vérifier si le champ a un tag
		if field.Tag == nil || field.Tag.Value == "" {
			continue
		}

		// Vérifier si c'est un champ privé
		for _, name := range field.Names {
			// Un champ est privé si son nom commence par une minuscule
			if isPrivateField(name.Name) {
				pass.Reportf(
					field.Pos(),
					"KTN-STRUCT-006: le champ privé '%s' du DTO '%s' ne devrait pas avoir de tag (tag ignoré lors de la sérialisation)",
					name.Name,
					structName,
				)
			}
		}
	}
}

// isPrivateField vérifie si un champ est privé.
//
// Params:
//   - name: nom du champ
//
// Returns:
//   - bool: true si le champ est privé
func isPrivateField(name string) bool {
	// Un champ est privé si son nom commence par une minuscule
	if len(name) == 0 {
		// Retour false pour nom vide
		return false
	}
	// Vérifier le premier caractère
	// Retour du résultat de la vérification
	return unicode.IsLower(rune(name[0]))
}
