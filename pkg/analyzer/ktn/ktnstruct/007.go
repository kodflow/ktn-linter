// Package ktnstruct provides analyzers for struct-related lint rules.
package ktnstruct

import (
	"go/ast"
	"strings"
	"unicode"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeStruct007 code de la règle KTN-STRUCT-007
	ruleCodeStruct007 string = "KTN-STRUCT-007"
)

// Analyzer007 checks that DTO exported fields have serialization tags.
var Analyzer007 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct007",
	Doc:      "KTN-STRUCT-007: Les champs exportés d'un DTO doivent avoir des tags json/xml",
	Run:      runStruct007,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runStruct007 exécute l'analyse KTN-STRUCT-007.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct007(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeStruct007) {
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
		if cfg.IsFileExcluded(ruleCodeStruct007, filename) {
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
			return
		}

		// Vérifier les champs exportés sans tags
		checkExportedFieldsWithoutTags(pass, structType)
	})

	// Retour de la fonction
	return nil, nil
}

// checkExportedFieldsWithoutTags vérifie les champs exportés sans tags.
//
// Params:
//   - pass: contexte d'analyse
//   - structType: type de la struct
func checkExportedFieldsWithoutTags(pass *analysis.Pass, structType *ast.StructType) {
	// Vérifier si la struct a des champs
	if structType.Fields == nil {
		// Retour anticipé
		return
	}

	// Parcourir les champs
	for _, field := range structType.Fields.List {
		// Vérifier chaque nom de champ
		for _, name := range field.Names {
			// Un champ est exporté si son nom commence par une majuscule
			if !isExportedField(name.Name) {
				// Champ privé, ignorer
				continue
			}

			// Vérifier si le champ a un tag de sérialisation
			if !hasSerializationTag(field) {
				msg, _ := messages.Get(ruleCodeStruct007)
				pass.Reportf(
					field.Pos(),
					"%s: %s",
					ruleCodeStruct007,
					msg.Format(config.Get().Verbose, name.Name),
				)
			}
		}
	}
}

// isExportedField vérifie si un champ est exporté.
//
// Params:
//   - name: nom du champ
//
// Returns:
//   - bool: true si le champ est exporté
func isExportedField(name string) bool {
	// Un champ est exporté si son nom commence par une majuscule
	if len(name) == 0 {
		// Retour false pour nom vide
		return false
	}
	// Vérifier le premier caractère
	return unicode.IsUpper(rune(name[0]))
}

// hasSerializationTag vérifie si un champ a un tag de sérialisation.
//
// Params:
//   - field: champ à vérifier
//
// Returns:
//   - bool: true si le champ a un tag json ou xml
func hasSerializationTag(field *ast.Field) bool {
	// Pas de tag
	if field.Tag == nil || field.Tag.Value == "" {
		// Retour false
		return false
	}

	tagValue := field.Tag.Value
	// Vérifier si le tag contient json ou xml
	hasJSON := strings.Contains(tagValue, "json:")
	hasXML := strings.Contains(tagValue, "xml:")

	// Retourner true si json ou xml présent
	return hasJSON || hasXML
}
