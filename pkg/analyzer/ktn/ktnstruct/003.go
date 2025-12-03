// Analyzer 003 for the ktnstruct package.
package ktnstruct

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer003 vérifie que les champs exportés sont avant les champs privés
var Analyzer003 = &analysis.Analyzer{
	Name:     "ktnstruct003",
	Doc:      "KTN-STRUCT-003: Les champs exportés doivent être placés avant les champs privés dans une struct",
	Run:      runStruct003,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runStruct003 exécute l'analyse KTN-STRUCT-003.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runStruct003(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filtrer les nœuds de type TypeSpec
	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Cast vers TypeSpec
		typeSpec := n.(*ast.TypeSpec)

		// Vérifier que c'est une struct
		structType, ok := typeSpec.Type.(*ast.StructType)
		// Si ce n'est pas une struct, continuer
		if !ok {
			// Continue avec le nœud suivant
			return
		}

		// Ignorer les fichiers de test
		filename := pass.Fset.Position(typeSpec.Pos()).Filename
		// Vérification si fichier test
		if shared.IsTestFile(filename) {
			// Continue avec le nœud suivant
			return
		}

		// Vérifier l'ordre des champs
		checkFieldOrder(pass, typeSpec, structType)
	})

	// Retour de la fonction
	return nil, nil
}

// checkFieldOrder vérifie que les champs exportés sont avant les privés.
//
// Params:
//   - pass: contexte d'analyse
//   - typeSpec: spécification de type
//   - structType: type struct
//
// Returns: aucun
func checkFieldOrder(pass *analysis.Pass, typeSpec *ast.TypeSpec, structType *ast.StructType) {
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
	// Parcourir les champs
	for _, f := range fields {
		// Si on trouve un champ exporté après un privé
		if f.exported && foundPrivate {
			pass.Reportf(
				f.pos,
				"KTN-STRUCT-003: le champ exporté '%s' de la struct '%s' doit être placé avant les champs privés",
				f.name,
				typeSpec.Name.Name,
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
