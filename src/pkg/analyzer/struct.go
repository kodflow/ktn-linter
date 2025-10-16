package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"

	"github.com/kodflow/ktn-linter/src/internal/naming"
)

// Analyzers
var (
	// StructAnalyzer vérifie que les structs respectent les règles KTN
	StructAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktnstruct",
		Doc:  "Vérifie que les structs sont bien nommés et documentés",
		Run:  runStructAnalyzer,
	}
)

// runStructAnalyzer vérifie que toutes les structs respectent les règles KTN.
//
// Params:
//   - pass: la passe d'analyse contenant les fichiers à vérifier
//
// Returns:
//   - any: toujours nil car aucun résultat n'est nécessaire
//   - error: toujours nil, les erreurs sont rapportées via pass.Reportf
func runStructAnalyzer(pass *analysis.Pass) (any, error) {
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

				checkStruct(pass, typeSpec, structType, genDecl)
			}
		}
	}

	// Retourne nil car l'analyseur rapporte via pass.Reportf
	return nil, nil
}

// checkStruct vérifie toutes les règles pour une struct.
//
// Params:
//   - pass: la passe d'analyse pour rapporter les erreurs
//   - typeSpec: la spécification de type de la struct
//   - structType: le type struct AST
//   - genDecl: la déclaration générale contenant la struct
func checkStruct(pass *analysis.Pass, typeSpec *ast.TypeSpec, structType *ast.StructType, genDecl *ast.GenDecl) {
	structName := typeSpec.Name.Name

	checkStructNaming(pass, typeSpec, structName)
	checkStructDocumentation(pass, typeSpec, structName, genDecl)
	checkStructFields(pass, structType, structName)
	checkStructFieldCount(pass, typeSpec, structType, structName)
}

// checkStructNaming vérifie le nommage de la struct.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - typeSpec: la spécification de type
//   - structName: le nom de la struct
func checkStructNaming(pass *analysis.Pass, typeSpec *ast.TypeSpec, structName string) {
	if !naming.IsMixedCaps(structName) {
		pass.Reportf(typeSpec.Name.Pos(),
			"[KTN-STRUCT-001] Struct '%s' n'utilise pas la convention MixedCaps.\n"+
				"Utilisez MixedCaps pour les structs exportées ou mixedCaps pour les privées.\n"+
				"Exemple: UserConfig au lieu de user_config",
			structName)
	}
}

// checkStructDocumentation vérifie la documentation de la struct.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - typeSpec: la spécification de type
//   - structName: le nom de la struct
//   - genDecl: la déclaration générale
func checkStructDocumentation(pass *analysis.Pass, typeSpec *ast.TypeSpec, structName string, genDecl *ast.GenDecl) {
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

// checkStructFields vérifie la documentation des champs exportés.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - structType: le type struct AST
//   - structName: le nom de la struct
func checkStructFields(pass *analysis.Pass, structType *ast.StructType, structName string) {
	if structType.Fields == nil || len(structType.Fields.List) == 0 {
		// Retourne car la struct n'a pas de champs
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

// checkStructFieldCount vérifie le nombre de champs dans la struct.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - typeSpec: la spécification de type
//   - structType: le type struct AST
//   - structName: le nom de la struct
func checkStructFieldCount(pass *analysis.Pass, typeSpec *ast.TypeSpec, structType *ast.StructType, structName string) {
	if structType.Fields == nil {
		// Retourne car la struct n'a pas de champs
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
