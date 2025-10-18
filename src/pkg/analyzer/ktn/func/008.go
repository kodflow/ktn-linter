package ktn_func

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Rule008 vérifie que tous les return statements ont des commentaires.
//
// KTN-FUNC-008: Tout return doit avoir un commentaire explicatif juste au-dessus.
//
// Incorrect:
//   if err != nil {
//       return err
//   }
//
// Correct:
//   if err != nil {
//       // Erreur de traitement
//       return err
//   }
var Rule008 = &analysis.Analyzer{
	Name: "KTN_FUNC_008",
	Doc:  "Vérifie que tous les return statements ont des commentaires",
	Run:  runRule008,
}

// runRule008 exécute la vérification KTN-FUNC-008.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule008(pass *analysis.Pass) (any, error) {
	if isTargetTestFile(pass) {
		return nil, nil
	}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			checkReturnComments(pass, file, funcDecl)
		}
	}

	return nil, nil
}

// checkReturnComments vérifie que tous les return ont des commentaires.
//
// Params:
//   - pass: la passe d'analyse
//   - file: le fichier AST
//   - funcDecl: la déclaration de fonction
func checkReturnComments(pass *analysis.Pass, file *ast.File, funcDecl *ast.FuncDecl) {
	if funcDecl.Body == nil {
		return
	}

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		returnStmt, ok := n.(*ast.ReturnStmt)
		if !ok {
			return true
		}

		if !hasCommentAbove(file, pass.Fset, returnStmt) {
			pass.Reportf(returnStmt.Pos(),
				"[KTN-FUNC-008] Return statement sans commentaire explicatif.\n"+
					"Tout return doit avoir un commentaire juste au-dessus expliquant ce qui est retourné.\n"+
					"Exemple:\n"+
					"  // Erreur de traitement\n"+
					"  return err\n"+
					"\n"+
					"  // Succès\n"+
					"  return nil")
		}
		return true
	})
}

// hasCommentAbove vérifie si un return a un commentaire juste au-dessus.
//
// Params:
//   - file: le fichier AST
//   - fset: le FileSet
//   - returnStmt: le statement return
//
// Returns:
//   - bool: true si un commentaire existe
func hasCommentAbove(file *ast.File, fset *token.FileSet, returnStmt *ast.ReturnStmt) bool {
	returnLine := fset.Position(returnStmt.Pos()).Line

	for _, commentGroup := range file.Comments {
		if commentGroup == nil || len(commentGroup.List) == 0 {
			continue
		}

		lastComment := commentGroup.List[len(commentGroup.List)-1]
		commentEndLine := fset.Position(lastComment.End()).Line

		if commentEndLine == returnLine-1 {
			return true
		}
	}

	return false
}
