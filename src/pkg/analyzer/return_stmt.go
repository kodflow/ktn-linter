package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// ReturnStmtAnalyzer vérifie les statements de retour.
	ReturnStmtAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktnreturnstmt",
		Doc:  "Vérifie les statements de retour (naked returns)",
		Run:  runReturnStmtAnalyzer,
	}
)

// runReturnStmtAnalyzer exécute l'analyseur de return statements.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runReturnStmtAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			if !ok {
				// Retourne true pour continuer
				return true
			}

			checkNakedReturns(pass, funcDecl)

			// Retourne true pour continuer l'inspection
			return true
		})
	}
	// Retourne nil car l'analyse est terminée
	return nil, nil
}

// checkNakedReturns vérifie les naked returns dans une fonction.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
func checkNakedReturns(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	if funcDecl.Type == nil || funcDecl.Type.Results == nil {
		// Pas de valeurs de retour
		// Retourne
		return
	}

	// Vérifier si la fonction a des named returns
	hasNamedReturns := false
	for _, field := range funcDecl.Type.Results.List {
		if len(field.Names) > 0 {
			hasNamedReturns = true
			break
		}
	}

	if !hasNamedReturns {
		// Pas de named returns, pas de naked returns possibles
		// Retourne
		return
	}

	// Chercher les naked returns dans le body
	if funcDecl.Body != nil {
		findNakedReturns(pass, funcDecl)
	}
}

// findNakedReturns trouve les naked returns dans une fonction.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
func findNakedReturns(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		returnStmt, ok := n.(*ast.ReturnStmt)
		if !ok {
			// Retourne true pour continuer
			return true
		}

		// Vérifier si c'est un naked return (pas de valeurs)
		if len(returnStmt.Results) == 0 {
			// Vérifier si la fonction est "longue"
			if isFunctionLong(funcDecl) {
				reportNakedReturn(pass, returnStmt, funcDecl.Name.Name)
			}
		}

		// Retourne true pour continuer
		return true
	})
}

// isFunctionLong vérifie si une fonction est "longue".
//
// Params:
//   - funcDecl: la déclaration de fonction
//
// Returns:
//   - bool: true si la fonction est considérée longue (> 10 lignes)
func isFunctionLong(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Body == nil {
		return false
	}

	// Compter le nombre de statements dans le body
	statementCount := countStatements(funcDecl.Body)

	// Une fonction est considérée longue si elle a plus de 10 statements
	// ou si son body contient des structures de contrôle imbriquées
	return statementCount > 10 || hasNestedControl(funcDecl.Body)
}

// countStatements compte le nombre de statements dans un bloc.
//
// Params:
//   - block: le bloc
//
// Returns:
//   - int: nombre de statements
func countStatements(block *ast.BlockStmt) int {
	if block == nil {
		return 0
	}

	count := 0
	ast.Inspect(block, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.ExprStmt, *ast.AssignStmt, *ast.ReturnStmt,
			*ast.IfStmt, *ast.ForStmt, *ast.RangeStmt,
			*ast.SwitchStmt, *ast.SelectStmt, *ast.DeferStmt:
			count++
		}
		// Retourne true pour continuer
		return true
	})

	return count
}

// hasNestedControl vérifie si un bloc a des structures de contrôle imbriquées.
//
// Params:
//   - block: le bloc
//
// Returns:
//   - bool: true si structures imbriquées
func hasNestedControl(block *ast.BlockStmt) bool {
	if block == nil {
		return false
	}

	nestingLevel := 0
	maxNesting := 0

	ast.Inspect(block, func(n ast.Node) bool {
		switch n.(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt,
			*ast.SwitchStmt, *ast.SelectStmt:
			nestingLevel++
			if nestingLevel > maxNesting {
				maxNesting = nestingLevel
			}
		case *ast.BlockStmt:
			// En sortant d'un bloc
			if nestingLevel > 0 {
				nestingLevel--
			}
		}
		// Retourne true pour continuer
		return true
	})

	// Considéré comme imbriqué si nesting > 2
	return maxNesting > 2
}

// reportNakedReturn rapporte une violation KTN-RETURN-004.
//
// Params:
//   - pass: la passe d'analyse
//   - returnStmt: le return statement
//   - funcName: nom de la fonction
func reportNakedReturn(pass *analysis.Pass, returnStmt *ast.ReturnStmt, funcName string) {
	pass.Reportf(returnStmt.Pos(),
		"[KTN-RETURN-004] Naked return dans la fonction '%s'.\n"+
			"Les naked returns (return sans valeurs avec named returns) réduisent la lisibilité.\n"+
			"Dans les fonctions longues, il est difficile de savoir quelles valeurs sont retournées.\n"+
			"Retournez explicitement les valeurs.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - naked return dans fonction longue\n"+
			"  func process() (result int, err error) {\n"+
			"      // ... 50 lignes de code ...\n"+
			"      result = 42\n"+
			"      err = nil\n"+
			"      return  // Quelles valeurs? Pas clair!\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - return explicite\n"+
			"  func process() (result int, err error) {\n"+
			"      // ... code ...\n"+
			"      result = 42\n"+
			"      err = nil\n"+
			"      return result, err  // Clair!\n"+
			"  }\n"+
			"\n"+
			"  // ✅ ACCEPTABLE - naked return dans fonction courte\n"+
			"  func getDefault() (x int) {\n"+
			"      x = 42\n"+
			"      return  // OK: fonction courte, évident\n"+
			"  }",
		funcName)
}
