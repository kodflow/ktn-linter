package ktn_struct

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Rule002 vérifie la documentation godoc des structs.
//
// KTN-STRUCT-002: Documentation godoc des structs
var Rule002 = &analysis.Analyzer{
	Name: "KTN_STRUCT_002",
	Doc:  "Documentation godoc des structs",
	Run:  runRule002,
}

// runRule002 exécute la vérification KTN-STRUCT-002.
func runRule002(pass *analysis.Pass) (any, error) {
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

				checkStructDocumentation(pass, typeSpec, genDecl)
			}
		}
	}
	return nil, nil
}

func checkStructDocumentation(pass *analysis.Pass, typeSpec *ast.TypeSpec, genDecl *ast.GenDecl) {
	structName := typeSpec.Name.Name
	if genDecl.Doc == nil || len(genDecl.Doc.List) == 0 {
		pass.Reportf(typeSpec.Name.Pos(),
			"[KTN-STRUCT-002] Struct '%s' sans commentaire godoc.\n"+
				"Toute struct doit avoir un commentaire godoc.\n"+
				"Exemple:\n"+
				"  // %s représente...\n"+
				"  type %s struct { }",
			structName, structName, structName)
	}
}
