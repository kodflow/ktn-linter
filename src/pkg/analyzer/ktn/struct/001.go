package ktn_struct

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/src/pkg/analyzer/utils"
	"golang.org/x/tools/go/analysis"
)

// Rule001 vérifie le nommage MixedCaps des structs.
//
// KTN-STRUCT-001: Nommage MixedCaps des structs
var Rule001 = &analysis.Analyzer{
	Name: "KTN_STRUCT_001",
	Doc:  "Nommage MixedCaps des structs",
	Run:  runRule001,
}

// runRule001 exécute la vérification KTN-STRUCT-001.
func runRule001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.TYPE {
				continue
			}

			for _, spec := range genDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				_, ok = typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				checkStructNaming(pass, typeSpec)
			}
		}
	}
	return nil, nil
}

func checkStructNaming(pass *analysis.Pass, typeSpec *ast.TypeSpec) {
	structName := typeSpec.Name.Name
	if !utils.IsMixedCaps(structName) {
		pass.Reportf(typeSpec.Name.Pos(),
			"[KTN-STRUCT-001] Struct '%s' n'utilise pas la convention MixedCaps.\n"+
				"Utilisez MixedCaps pour les structs exportées ou mixedCaps pour les privées.\n"+
				"Exemple: UserConfig au lieu de user_config",
			structName)
	}
}
