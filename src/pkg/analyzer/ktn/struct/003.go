package ktn_struct

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Rule003 vérifie la documentation des champs exportés.
//
// KTN-STRUCT-003: Documentation des champs exportés
var Rule003 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_STRUCT_003",
	Doc:  "Documentation des champs exportés",
	Run:  runRule003,
}

// runRule003 exécute la vérification KTN-STRUCT-003.
func runRule003(pass *analysis.Pass) (any, error) {
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

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				checkStructFields(pass, typeSpec, structType)
			}
		}
	}
	// Analysis completed successfully.
	return nil, nil
}

func checkStructFields(pass *analysis.Pass, typeSpec *ast.TypeSpec, structType *ast.StructType) {
	structName := typeSpec.Name.Name
	if structType.Fields == nil || len(structType.Fields.List) == 0 {
		// Early return from function.
		return
	}

	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			if !ast.IsExported(name.Name) {
				continue
			}

			if field.Doc == nil || len(field.Doc.List) == 0 {
				pass.Reportf(name.Pos(),
					"[KTN-STRUCT-003] Champ exporté '%s.%s' sans commentaire.\n"+
						"Tous les champs exportés doivent être documentés.\n"+
						"Exemple:\n"+
						"  // %s description du champ\n"+
						"  %s string",
					structName, name.Name, name.Name, name.Name)
			}
		}
	}
}
