package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// ConstAnalyzer vérifie que les constantes respectent les règles KTN
var ConstAnalyzer = &analysis.Analyzer{
	Name: "ktnconst",
	Doc:  "Vérifie que les constantes sont regroupées, documentées et typées explicitement",
	Run:  runConstAnalyzer,
}

func runConstAnalyzer(pass *analysis.Pass) (interface{}, error) {
	// Pour détecter les constantes non groupées, on doit parcourir toutes les déclarations
	for _, file := range pass.Files {
		// Map pour suivre les commentaires avant chaque déclaration
		comments := make(map[token.Pos]*ast.CommentGroup)
		for _, cg := range file.Comments {
			if len(cg.List) > 0 {
				comments[cg.End()] = cg
			}
		}

		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.CONST {
				continue
			}

			// KTN-CONST-001: Vérifier si c'est une constante non groupée
			if genDecl.Lparen == token.NoPos {
				// Constante déclarée sans ()
				for _, spec := range genDecl.Specs {
					valueSpec := spec.(*ast.ValueSpec)
					for _, name := range valueSpec.Names {
						pass.Reportf(name.Pos(),
							"[KTN-CONST-001] Constante '%s' déclarée individuellement. Regroupez les constantes dans un bloc const ().\nExemple:\n  const (\n      %s %s = ...\n  )",
							name.Name, name.Name, getTypeString(valueSpec))
					}
				}
				continue
			}

			// C'est un groupe de constantes const ()
			// KTN-CONST-002: Vérifier le commentaire de groupe
			hasGroupComment := false
			if genDecl.Doc != nil && len(genDecl.Doc.List) > 0 {
				hasGroupComment = true
			}

			if !hasGroupComment {
				pass.Reportf(genDecl.Pos(),
					"[KTN-CONST-002] Groupe de constantes sans commentaire de groupe.\nAjoutez un commentaire avant le bloc const () pour décrire l'ensemble.\nExemple:\n  // Description du groupe de constantes\n  const (...)")
			}

			// Vérifier chaque constante dans le groupe
			for _, spec := range genDecl.Specs {
				valueSpec := spec.(*ast.ValueSpec)
				checkConstSpec(pass, valueSpec)
			}
		}
	}

	return nil, nil
}

func checkConstSpec(pass *analysis.Pass, spec *ast.ValueSpec) {
	for _, name := range spec.Names {
		if name.Name == "_" {
			continue
		}

		// KTN-CONST-003: Vérifier le commentaire individuel
		hasComment := false
		if spec.Doc != nil && len(spec.Doc.List) > 0 {
			hasComment = true
		} else if spec.Comment != nil && len(spec.Comment.List) > 0 {
			hasComment = true
		}

		if !hasComment {
			pass.Reportf(name.Pos(),
				"[KTN-CONST-003] Constante '%s' sans commentaire individuel.\nChaque constante doit avoir un commentaire explicatif.\nExemple:\n  // %s décrit son rôle\n  %s %s = ...",
				name.Name, name.Name, name.Name, getTypeString(spec))
		}

		// KTN-CONST-004: Vérifier le type explicite
		if spec.Type == nil {
			pass.Reportf(name.Pos(),
				"[KTN-CONST-004] Constante '%s' sans type explicite.\nSpécifiez toujours le type : bool, string, int, int8, uint, float64, etc.\nExemple:\n  %s int = ...",
				name.Name, name.Name)
		}
	}
}

func getTypeString(spec *ast.ValueSpec) string {
	if spec.Type != nil {
		return exprToString(spec.Type)
	}
	return "<type>"
}

func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	case *ast.ArrayType:
		return "[]" + exprToString(e.Elt)
	case *ast.MapType:
		return "map[" + exprToString(e.Key) + "]" + exprToString(e.Value)
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	default:
		return "unknown"
	}
}
