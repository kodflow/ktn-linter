package ktn_error

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Rule001 vérifie que les erreurs sont wrappées avec du contexte.
//
// KTN-ERROR-001: Les erreurs doivent être wrappées avec fmt.Errorf() et %w.
// Cela préserve la chaîne d'erreurs et améliore le debugging en production.
//
// Incorrect:
//   if err != nil {
//       return err  // Perd le contexte
//   }
//
// Correct:
//   if err != nil {
//       return fmt.Errorf("failed to process user %s: %w", userID, err)
//   }
var Rule001 = &analysis.Analyzer{
	Name: "KTN_ERROR_001",
	Doc:  "Vérifie que les erreurs sont wrappées avec du contexte",
	Run:  runRule001,
}

// runRule001 exécute la vérification KTN-ERROR-001.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - any: toujours nil
//   - error: toujours nil
func runRule001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		checkErrorWrapping(pass, file)
	}
	return nil, nil
}

// checkErrorWrapping vérifie que les erreurs sont wrappées.
//
// Params:
//   - pass: la passe d'analyse
//   - file: le fichier à analyser
func checkErrorWrapping(pass *analysis.Pass, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			return true
		}

		if !functionReturnsError(pass, funcDecl) {
			return true
		}

		checkReturnStatements(pass, funcDecl)
		return true
	})
}

// functionReturnsError vérifie si une fonction retourne un error.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
//
// Returns:
//   - bool: true si la fonction retourne error
func functionReturnsError(pass *analysis.Pass, funcDecl *ast.FuncDecl) bool {
	if funcDecl.Type == nil || funcDecl.Type.Results == nil {
		return false
	}

	for _, result := range funcDecl.Type.Results.List {
		if pass.TypesInfo != nil {
			resultType := pass.TypesInfo.TypeOf(result.Type)
			if resultType != nil && resultType.String() == "error" {
				return true
			}
		}

		if ident, ok := result.Type.(*ast.Ident); ok && ident.Name == "error" {
			return true
		}
	}

	return false
}

// checkReturnStatements vérifie les return statements.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
func checkReturnStatements(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		returnStmt, ok := n.(*ast.ReturnStmt)
		if !ok {
			return true
		}

		checkUnwrappedError(pass, returnStmt)
		return true
	})
}

// checkUnwrappedError vérifie si une erreur est retournée sans wrapping.
//
// Params:
//   - pass: la passe d'analyse
//   - returnStmt: le return statement
func checkUnwrappedError(pass *analysis.Pass, returnStmt *ast.ReturnStmt) {
	for _, result := range returnStmt.Results {
		if isNilIdent(result) {
			continue
		}

		ident, ok := result.(*ast.Ident)
		if !ok {
			continue
		}

		if !isErrorType(pass, result) {
			continue
		}

		if isErrorVariable(pass, ident) {
			reportUnwrappedError(pass, returnStmt, ident.Name)
		}
	}
}

// isNilIdent vérifie si une expression est nil.
//
// Params:
//   - expr: l'expression
//
// Returns:
//   - bool: true si c'est nil
func isNilIdent(expr ast.Expr) bool {
	ident, ok := expr.(*ast.Ident)
	return ok && ident.Name == "nil"
}

// isErrorType vérifie si une expression est de type error.
//
// Params:
//   - pass: la passe d'analyse
//   - expr: l'expression
//
// Returns:
//   - bool: true si c'est un error
func isErrorType(pass *analysis.Pass, expr ast.Expr) bool {
	if pass.TypesInfo == nil {
		return false
	}

	exprType := pass.TypesInfo.TypeOf(expr)
	if exprType == nil {
		return false
	}

	return exprType.String() == "error"
}

// isErrorVariable vérifie si un identifiant est une variable d'erreur.
//
// Params:
//   - pass: la passe d'analyse
//   - ident: l'identifiant
//
// Returns:
//   - bool: true si c'est une variable d'erreur
func isErrorVariable(pass *analysis.Pass, ident *ast.Ident) bool {
	if pass.TypesInfo == nil {
		return false
	}

	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		return false
	}

	_, ok := obj.(*types.Var)
	return ok
}

// reportUnwrappedError rapporte une violation KTN-ERROR-001.
//
// Params:
//   - pass: la passe d'analyse
//   - returnStmt: le return statement
//   - varName: le nom de la variable erreur
func reportUnwrappedError(pass *analysis.Pass, returnStmt *ast.ReturnStmt, varName string) {
	pass.Reportf(returnStmt.Pos(),
		"[KTN-ERROR-001] Erreur '%s' retournée sans contexte.\n"+
			"Les erreurs doivent être wrappées avec fmt.Errorf() pour préserver le contexte.\n"+
			"Cela améliore le debugging en production en traçant l'origine des erreurs.\n"+
			"Utilisez fmt.Errorf(\"contexte descriptif: %%w\", %s) pour wrapper l'erreur.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - perd le contexte\n"+
			"  if err != nil {\n"+
			"      return err\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - préserve le contexte\n"+
			"  if err != nil {\n"+
			"      return fmt.Errorf(\"failed to process user %%s: %%w\", userID, err)\n"+
			"  }",
		varName, varName)
}
