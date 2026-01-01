// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/utils"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar014 is the rule code for this analyzer
	ruleCodeVar014 string = "KTN-VAR-014"
)

// Analyzer014 checks for repeated buffer allocations in loops
var Analyzer014 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar014",
	Doc:      "KTN-VAR-014: Vérifie que les buffers répétés utilisent sync.Pool",
	Run:      runVar014,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar014 exécute l'analyse KTN-VAR-014.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar014(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar014) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ForStmt)(nil),
		(*ast.RangeStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar014, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}
		// Check different loop types
		switch loop := n.(type) {
		// For loop: check body for buffer allocations
		case *ast.ForStmt:
			checkLoopBodyVar015(pass, loop.Body)
		// Range loop: check body for buffer allocations
		case *ast.RangeStmt:
			checkLoopBodyVar015(pass, loop.Body)
		}
	})

	// Return analysis result
	return nil, nil
}

// checkLoopBodyVar015 vérifie le corps d'une boucle.
//
// Params:
//   - pass: contexte d'analyse
//   - body: corps de la boucle
func checkLoopBodyVar015(pass *analysis.Pass, body *ast.BlockStmt) {
	// Check if body exists
	if body == nil {
		// No body to check
		return
	}

	// Inspect all statements in loop body
	ast.Inspect(body, func(n ast.Node) bool {
		// Check for assignment statements
		if assignStmt, ok := n.(*ast.AssignStmt); ok {
			checkAssignmentForBuffer(pass, assignStmt)
		}
		// Continue traversal
		return true
	})
}

// checkAssignmentForBuffer vérifie si une assignation crée un buffer.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: statement d'assignation
func checkAssignmentForBuffer(pass *analysis.Pass, stmt *ast.AssignStmt) {
	// Check each right-hand side value
	for _, rhs := range stmt.Rhs {
		// Check if it's a make call
		if callExpr, ok := rhs.(*ast.CallExpr); ok {
			checkMakeCallForByteSlice(pass, callExpr)
		}
	}
}

// checkMakeCallForByteSlice vérifie si un make crée un []byte.
//
// Params:
//   - pass: contexte d'analyse
//   - call: expression d'appel
func checkMakeCallForByteSlice(pass *analysis.Pass, call *ast.CallExpr) {
	// Check if it's a call to 'make'
	if !utils.IsMakeCall(call) {
		// Not a make call
		return
	}

	// Check if making a byte slice
	if len(call.Args) == 0 || !utils.IsByteSlice(call.Args[0]) {
		// Not a byte slice
		return
	}

	// Report the issue
	msg, _ := messages.Get(ruleCodeVar014)
	pass.Reportf(
		call.Pos(),
		"%s: %s",
		ruleCodeVar014,
		msg.Format(config.Get().Verbose),
	)
}
