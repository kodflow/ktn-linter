// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar007 is the rule code for this analyzer
	ruleCodeVar007 string = "KTN-VAR-007"
)

// Analyzer007 checks that local variables use := instead of var
var Analyzer007 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar007",
	Doc:      "KTN-VAR-007: Vérifie que les variables locales utilisent ':=' au lieu de 'var'",
	Run:      runVar007,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar007 exécute l'analyse KTN-VAR-007.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar007(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar007) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// We need to track function bodies to check local variables only
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar007, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Check if function has a body
		if funcDecl.Body == nil {
			// Continue to next node
			return
		}

		// Inspect all statements in the function body
		checkFunctionBody(pass, funcDecl.Body)
	})

	// Return analysis result
	return nil, nil
}

// checkFunctionBody parcourt le corps d'une fonction pour détecter les var.
//
// Params:
//   - pass: contexte d'analyse
//   - body: corps de la fonction
func checkFunctionBody(pass *analysis.Pass, body *ast.BlockStmt) {
	// Iterate through all statements
	for _, stmt := range body.List {
		checkStatement(pass, stmt)
	}
}

// checkStatement vérifie si un statement contient un var avec initialisation.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: statement à vérifier
func checkStatement(pass *analysis.Pass, stmt ast.Stmt) {
	declStmt, ok := stmt.(*ast.DeclStmt)
	// If not a declaration, check for nested blocks
	if !ok {
		checkNestedBlocks(pass, stmt)
		// Early return for non-declarations
		return
	}

	genDecl, ok := declStmt.Decl.(*ast.GenDecl)
	// If not a GenDecl, return early
	if !ok {
		// Not a GenDecl, skip
		return
	}

	// Only check var declarations, skip others
	if genDecl.Tok != token.VAR {
		// Not a var declaration, skip
		return
	}

	// Check each variable specification
	checkVarSpecs(pass, genDecl)
}

// checkNestedBlocks vérifie les blocs imbriqués (if, for, switch, etc.).
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: statement à vérifier
func checkNestedBlocks(pass *analysis.Pass, stmt ast.Stmt) {
	// Check different types of statements with nested blocks
	switch s := stmt.(type) {
	// If statement: check body and else
	case *ast.IfStmt:
		checkIfStmt(pass, s)
	// For statement: check loop body
	case *ast.ForStmt:
		checkBlockIfNotNil(pass, s.Body)
	// Range statement: check loop body
	case *ast.RangeStmt:
		checkBlockIfNotNil(pass, s.Body)
	// Switch statement: check switch body
	case *ast.SwitchStmt:
		checkBlockIfNotNil(pass, s.Body)
	// Type switch: check switch body
	case *ast.TypeSwitchStmt:
		checkBlockIfNotNil(pass, s.Body)
	// Select statement: check select body
	case *ast.SelectStmt:
		checkBlockIfNotNil(pass, s.Body)
	// Nested block: check directly
	case *ast.BlockStmt:
		checkFunctionBody(pass, s)
	// Case clause: iterate through statements
	case *ast.CaseClause:
		checkCaseClause(pass, s)
	// Comm clause: iterate through statements
	case *ast.CommClause:
		checkCommClause(pass, s)
	}
}

// checkIfStmt vérifie un if statement.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: if statement
func checkIfStmt(pass *analysis.Pass, stmt *ast.IfStmt) {
	// Check if body exists
	if stmt.Body != nil {
		checkFunctionBody(pass, stmt.Body)
	}
	// Check else clause if exists
	if stmt.Else != nil {
		checkStatement(pass, stmt.Else)
	}
}

// checkBlockIfNotNil vérifie un bloc s'il n'est pas nil.
//
// Params:
//   - pass: contexte d'analyse
//   - block: bloc à vérifier
func checkBlockIfNotNil(pass *analysis.Pass, block *ast.BlockStmt) {
	// Check if block exists
	if block != nil {
		checkFunctionBody(pass, block)
	}
}

// checkCaseClause vérifie une case clause.
//
// Params:
//   - pass: contexte d'analyse
//   - clause: case clause
func checkCaseClause(pass *analysis.Pass, clause *ast.CaseClause) {
	// Iterate through case statements
	for _, caseStmt := range clause.Body {
		checkStatement(pass, caseStmt)
	}
}

// checkCommClause vérifie une comm clause.
//
// Params:
//   - pass: contexte d'analyse
//   - clause: comm clause
func checkCommClause(pass *analysis.Pass, clause *ast.CommClause) {
	// Iterate through comm statements
	for _, commStmt := range clause.Body {
		checkStatement(pass, commStmt)
	}
}

// checkVarSpecs vérifie les spécifications de variables.
//
// Params:
//   - pass: contexte d'analyse
//   - genDecl: déclaration générale
func checkVarSpecs(pass *analysis.Pass, genDecl *ast.GenDecl) {
	// Iterate through specifications
	for _, spec := range genDecl.Specs {
		valueSpec := spec.(*ast.ValueSpec)

		// Check if variable has initialization without explicit type
		if hasInitWithoutType(valueSpec) {
			// Report error for each variable
			reportVarErrors(pass, valueSpec)
		}
	}
}

// hasInitWithoutType vérifie si une variable a une initialisation explicite.
//
// Params:
//   - spec: spécification de variable
//
// Returns:
//   - bool: true si initialisation explicite présente
func hasInitWithoutType(spec *ast.ValueSpec) bool {
	// Return true if variable has explicit initialization
	// We want to report both "var x = value" and "var x Type = value"
	// But NOT "var x Type" (zero value declaration)
	return len(spec.Values) > 0
}

// reportVarErrors rapporte les erreurs pour chaque variable.
//
// Params:
//   - pass: contexte d'analyse
//   - spec: spécification de variable
func reportVarErrors(pass *analysis.Pass, spec *ast.ValueSpec) {
	// Iterate through variable names
	for _, name := range spec.Names {
		msg, ok := messages.Get(ruleCodeVar007)
		// Defensive: avoid panic if message is missing
		if !ok {
			pass.Reportf(name.Pos(), "%s: utiliser := au lieu de var pour les valeurs non-zéro", ruleCodeVar007)
			continue
		}
		pass.Reportf(
			name.Pos(),
			"%s: %s",
			ruleCodeVar007,
			msg.Format(config.Get().Verbose),
		)
	}
}
