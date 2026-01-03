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
	cfg := config.Get()
	if !cfg.IsRuleEnabled(ruleCodeStruct007) {
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.TypeSpec)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		analyzeTypeSpec007(pass, cfg, n)
	})

	return nil, nil
}

// analyzeTypeSpec007 analyse un TypeSpec pour la règle 007.
//
// Params:
//   - pass: contexte d'analyse
//   - cfg: configuration
//   - n: noeud AST à analyser
func analyzeTypeSpec007(pass *analysis.Pass, cfg *config.Config, n ast.Node) {
	typeSpec := n.(*ast.TypeSpec)
	filename := pass.Fset.Position(n.Pos()).Filename

	if cfg.IsFileExcluded(ruleCodeStruct007, filename) {
		return
	}

	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return
	}

	if !shared.IsSerializableStruct(structType, typeSpec.Name.Name) {
		return
	}

	checkExportedFieldsWithoutTags(pass, structType)
}

// checkExportedFieldsWithoutTags vérifie les champs exportés sans tags.
//
// Params:
//   - pass: contexte d'analyse
//   - structType: type de la struct
func checkExportedFieldsWithoutTags(pass *analysis.Pass, structType *ast.StructType) {
	if structType.Fields == nil {
		return
	}

	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			if !isExportedField(name.Name) {
				continue
			}

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
	if len(name) == 0 {
		return false
	}

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
	if field.Tag == nil || field.Tag.Value == "" {
		return false
	}

	tagValue := field.Tag.Value

	return strings.Contains(tagValue, "json:") || strings.Contains(tagValue, "xml:")
}
