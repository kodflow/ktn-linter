package ktn_struct

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Rule004 vérifie le nombre de champs dans une struct.
//
// KTN-STRUCT-004: Maximum 15 champs par struct
var Rule004 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_STRUCT_004",
	Doc:  "Maximum 15 champs par struct",
	Run:  runRule004,
}

// runRule004 exécute la vérification KTN-STRUCT-004.
func runRule004(pass *analysis.Pass) (any, error) {
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

				checkStructFieldCount(pass, typeSpec, structType)
			}
		}
	}
	// Analysis completed successfully.
	return nil, nil
}

func checkStructFieldCount(pass *analysis.Pass, typeSpec *ast.TypeSpec, structType *ast.StructType) {
	structName := typeSpec.Name.Name
	if structType.Fields == nil {
		// Early return from function.
		return
	}

	fieldCount := 0
	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			// Champ embedded
			fieldCount++
		} else {
			fieldCount += len(field.Names)
		}
	}

	const maxFields = 15
	if fieldCount > maxFields {
		pass.Reportf(typeSpec.Name.Pos(),
			"[KTN-STRUCT-004] Struct '%s' a trop de champs (%d > %d).\n"+
				"Limitez à %d champs maximum. Si nécessaire, décomposez en plusieurs structs.\n"+
				"Exemple:\n"+
				"  type %sCore struct { ... }\n"+
				"  type %sMetadata struct { ... }\n"+
				"  type %s struct {\n"+
				"      Core %sCore\n"+
				"      Metadata %sMetadata\n"+
				"  }",
			structName, fieldCount, maxFields, maxFields,
			structName, structName, structName, structName, structName)
	}
}
