package analyzer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// ErrorAnalyzer vérifie la gestion des erreurs.
	ErrorAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktnerror",
		Doc:  "Vérifie que les erreurs sont wrappées avec du contexte",
		Run:  runErrorAnalyzer,
	}
)

// runErrorAnalyzer exécute l'analyseur error.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runErrorAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		checkErrorWrapping(pass, file)
	}
	// Retourne nil car l'analyse est terminée
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
			// Retourne true pour continuer l'inspection
			return true
		}

		// Vérifier si la fonction retourne une erreur
		if !functionReturnsError(pass, funcDecl) {
			// Retourne true pour continuer l'inspection
			return true
		}

		checkReturnStatements(pass, funcDecl)
		// Retourne true pour continuer l'inspection
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
		// Retourne false car pas de résultats
		return false
	}

	for _, result := range funcDecl.Type.Results.List {
		if pass.TypesInfo != nil {
			resultType := pass.TypesInfo.TypeOf(result.Type)
			if resultType != nil && resultType.String() == "error" {
				// Retourne true car retourne error
				return true
			}
		}

		// Fallback: vérifier le nom du type
		if ident, ok := result.Type.(*ast.Ident); ok && ident.Name == "error" {
			// Retourne true car retourne error
			return true
		}
	}

	// Retourne false car ne retourne pas error
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
			// Retourne true pour continuer l'inspection
			return true
		}

		checkUnwrappedError(pass, returnStmt)
		// Retourne true pour continuer l'inspection
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
		// Ignorer les nil
		if isNilIdent(result) {
			// Continue car c'est nil
			continue
		}

		// Vérifier si c'est un identifiant simple (variable)
		ident, ok := result.(*ast.Ident)
		if !ok {
			// Continue car ce n'est pas un identifiant simple
			continue
		}

		// Vérifier que c'est bien un error type
		if !isErrorType(pass, result) {
			// Continue car ce n'est pas un error
			continue
		}

		// Vérifier si c'est une variable qui vient d'un appel de fonction
		// (donc potentiellement une erreur qu'on devrait wrapper)
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
	// Retourne true si c'est l'identifiant "nil"
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
		// Retourne false car pas d'info de type
		return false
	}

	exprType := pass.TypesInfo.TypeOf(expr)
	if exprType == nil {
		// Retourne false car type inconnu
		return false
	}

	// Retourne true si le type est "error"
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
		// Retourne false car pas d'info de type
		return false
	}

	obj := pass.TypesInfo.Uses[ident]
	if obj == nil {
		// Retourne false car objet inconnu
		return false
	}

	// Vérifier que c'est une variable (pas une fonction ou constante)
	_, ok := obj.(*types.Var)
	// Retourne true si c'est une variable
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
