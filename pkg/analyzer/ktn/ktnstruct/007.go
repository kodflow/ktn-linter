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
	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeStruct007) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.TypeSpec)(nil)}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		analyzeTypeSpec007(pass, cfg, n)
	})

	// Fin de l'analyse
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

	// Vérifier si le fichier est exclu
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

	// Vérifier si la struct est sérialisable (DTO)
	if !shared.IsSerializableStruct(structType, typeSpec.Name.Name) {
		// Pas un DTO
		return
	}

	checkExportedFieldsWithoutTags(pass, cfg, structType)
}

// checkExportedFieldsWithoutTags vérifie les champs exportés sans tags.
//
// Params:
//   - pass: contexte d'analyse
//   - cfg: configuration
//   - structType: type de la struct
func checkExportedFieldsWithoutTags(
	pass *analysis.Pass,
	cfg *config.Config,
	structType *ast.StructType,
) {
	// Vérifier si la struct a des champs
	if structType.Fields == nil {
		// Pas de champs
		return
	}

	tags := cfg.SerializationTags()

	// Parcourir les champs de la struct
	for _, field := range structType.Fields.List {
		// Parcourir les noms du champ
		for _, name := range field.Names {
			// Ignorer les champs non exportés
			if !isExportedField(name.Name) {
				continue
			}

			// Vérifier si le champ a un tag de sérialisation
			if !hasSerializationTag(field, tags) {
				msg, _ := messages.Get(ruleCodeStruct007)
				pass.Reportf(
					field.Pos(),
					"%s: %s",
					ruleCodeStruct007,
					msg.Format(cfg.Verbose, name.Name),
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
	// Vérifier si le nom est vide
	if len(name) == 0 {
		// Nom vide
		return false
	}

	// Champ exporté si commence par une majuscule
	return unicode.IsUpper(rune(name[0]))
}

// hasSerializationTag vérifie si un champ a un tag de sérialisation.
//
// Params:
//   - field: champ à vérifier
//   - tags: liste des tags de sérialisation reconnus
//
// Returns:
//   - bool: true si le champ a un tag de sérialisation reconnu
func hasSerializationTag(field *ast.Field, tags []string) bool {
	// Vérifier si le champ a un tag
	if field.Tag == nil || field.Tag.Value == "" {
		// Pas de tag
		return false
	}

	tagValue := field.Tag.Value

	// Parcourir les tags de sérialisation reconnus
	for _, tag := range tags {
		// Vérifier si le tag contient un tag de sérialisation
		if strings.Contains(tagValue, tag) {
			// Tag trouvé
			return true
		}
	}

	// Aucun tag de sérialisation trouvé
	return false
}
