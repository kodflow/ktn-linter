// Analyzer 005 for the ktnstruct package.
package ktnstruct

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeStruct005 code de la règle KTN-STRUCT-005
	ruleCodeStruct005 string = "KTN-STRUCT-005"
)

// Analyzer005 vérifie que les champs exportés sont avant les champs privés
var Analyzer005 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnstruct005",
	Doc:      "KTN-STRUCT-005: Les champs exportés doivent être placés avant les champs privés dans une struct",
	Run:      runStruct005,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runStruct005 exécute l'analyse KTN-STRUCT-005.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct005(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeStruct005) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filtrer les nœuds de type TypeSpec
	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Cast vers TypeSpec
		typeSpec := n.(*ast.TypeSpec)

		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeStruct005, filename) {
			// Fichier exclu
			return
		}

		// Vérifier que c'est une struct
		structType, ok := typeSpec.Type.(*ast.StructType)
		// Si ce n'est pas une struct, continuer
		if !ok {
			// Continue avec le nœud suivant
			return
		}

		// Ignorer les fichiers de test
		filename = pass.Fset.Position(typeSpec.Pos()).Filename
		// Vérification si fichier test
		if shared.IsTestFile(filename) {
			// Continue avec le nœud suivant
			return
		}

		// Vérifier l'ordre des champs
		checkFieldOrder(pass, structType)
	})

	// Retour de la fonction
	return nil, nil
}

// checkFieldOrder vérifie que les champs exportés sont avant les privés.
//
// Params:
//   - pass: contexte d'analyse
//   - structType: type struct
func checkFieldOrder(pass *analysis.Pass, structType *ast.StructType) {
	// Si pas de champs
	if structType.Fields == nil || len(structType.Fields.List) == 0 {
		// Retour anticipé
		return
	}

	var fields []fieldInfo
	// Parcourir tous les champs
	for _, field := range structType.Fields.List {
		// Pour chaque nom de champ
		for _, name := range field.Names {
			fields = append(fields, fieldInfo{
				name:     name.Name,
				exported: ast.IsExported(name.Name),
				pos:      name.Pos(),
			})
		}
	}

	// Vérifier l'ordre
	foundPrivate := false
	msg, _ := messages.Get(ruleCodeStruct005)
	// Parcourir les champs
	for _, f := range fields {
		// Si on trouve un champ exporté après un privé
		if f.exported && foundPrivate {
			pass.Reportf(
				f.pos,
				"%s: %s",
				ruleCodeStruct005,
				msg.Format(config.Get().Verbose),
			)
		}

		// Marquer qu'on a trouvé un champ privé
		if !f.exported {
			foundPrivate = true
		}
	}
}

// fieldInfo stocke les informations d'un champ.
type fieldInfo struct {
	name     string
	exported bool
	pos      token.Pos
}
