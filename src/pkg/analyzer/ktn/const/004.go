package ktn_const

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

var Rule004 = &analysis.Analyzer{
	Name: "KTN_CONST_004",
	Doc:  "Vérifie que les constantes ont un type explicite (exception pour iota)",
	Run:  runRule004,
}

func runRule004(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok || genDecl.Tok != token.CONST {
				continue
			}

			// Vérifier uniquement les groupes (avec parenthèses)
			if genDecl.Lparen == token.NoPos {
				continue
			}

			// Détecter si le groupe utilise iota
			groupUsesIota := groupContainsIota(genDecl)

			for _, spec := range genDecl.Specs {
				valueSpec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				checkConstExplicitType(pass, valueSpec, groupUsesIota)
			}
		}
	}
	return nil, nil
}

// checkConstExplicitType vérifie qu'une constante a un type explicite.
func checkConstExplicitType(pass *analysis.Pass, spec *ast.ValueSpec, groupUsesIota bool) {
	for _, name := range spec.Names {
		if name.Name == "_" {
			continue
		}

		// Vérifier le type explicite (exception pour iota)
		if spec.Type == nil && !groupUsesIota {
			pass.Reportf(name.Pos(),
				"[KTN_CONST_004] Constante '%s' sans type explicite.\nSpécifiez toujours le type : bool, string, int, int8, uint, float64, etc.\nExemple:\n  %s int = ...",
				name.Name, name.Name)
		}
	}
}

// groupContainsIota vérifie si un groupe de constantes utilise iota.
func groupContainsIota(genDecl *ast.GenDecl) bool {
	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		if usesIota(valueSpec) {
			return true
		}
	}
	return false
}

// usesIota vérifie si une constante utilise iota directement.
func usesIota(spec *ast.ValueSpec) bool {
	for _, value := range spec.Values {
		if containsIota(value) {
			return true
		}
	}
	return false
}

// containsIota vérifie récursivement si une expression contient iota.
func containsIota(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name == "iota"
	case *ast.BinaryExpr:
		return containsIota(e.X) || containsIota(e.Y)
	case *ast.UnaryExpr:
		return containsIota(e.X)
	case *ast.ParenExpr:
		return containsIota(e.X)
	case *ast.CallExpr:
		for _, arg := range e.Args {
			if containsIota(arg) {
				return true
			}
		}
		return false
	}
	return false
}
