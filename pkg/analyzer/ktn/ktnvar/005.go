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
	// ruleCodeVar005 is the rule code for this analyzer
	ruleCodeVar005 string = "KTN-VAR-005"
	// maxVarNameLength005 is the maximum variable name length
	maxVarNameLength005 int = 30
)

// Analyzer005 checks that variable names are not too long.
var Analyzer005 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar005",
	Doc:      "KTN-VAR-005: Les noms de variables ne doivent pas dépasser 30 caractères",
	Run:      runVar005,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar005 exécute l'analyse KTN-VAR-005.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar005(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar005) {
		// Règle désactivée
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp, ok := inspAny.(*inspector.Inspector)
	// Defensive: ensure inspector is available
	if !ok || insp == nil {
		return nil, nil
	}

	// Check package-level variables
	checkVar005PackageLevel(pass, insp, cfg)

	// Check local variables in functions
	checkVar005LocalVars(pass, insp, cfg)

	// Return analysis result
	return nil, nil
}

// checkVar005PackageLevel vérifie les variables au niveau package.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - cfg: configuration
func checkVar005PackageLevel(
	pass *analysis.Pass,
	insp *inspector.Inspector,
	cfg *config.Config,
) {
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		file, ok := n.(*ast.File)
		// Defensive: ensure node type matches
		if !ok {
			return
		}

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar005, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Check package-level declarations
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			// Skip if not a GenDecl
			if !ok {
				// Not a general declaration
				continue
			}

			// Only check var declarations
			if genDecl.Tok != token.VAR {
				// Continue to next declaration
				continue
			}

			// Check each variable specification
			for _, spec := range genDecl.Specs {
				valueSpec := spec.(*ast.ValueSpec)
				checkVar005Spec(pass, valueSpec)
			}
		}
	})
}

// checkVar005LocalVars vérifie les variables locales dans les fonctions.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - cfg: configuration
func checkVar005LocalVars(
	pass *analysis.Pass,
	insp *inspector.Inspector,
	cfg *config.Config,
) {
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl, ok := n.(*ast.FuncDecl)
		// Defensive: ensure node type matches
		if !ok {
			return
		}

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar005, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Skip functions without body
		if funcDecl.Body == nil {
			// No body to analyze
			return
		}

		// Visit all statements in the function
		ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
			// Check the node type
			checkVar005Node(pass, node)
			// Continue traversal
			return true
		})
	})
}

// checkVar005Node vérifie un nœud AST pour les variables trop longues.
//
// Params:
//   - pass: contexte d'analyse
//   - node: nœud à vérifier
func checkVar005Node(pass *analysis.Pass, node ast.Node) {
	// Switch on node type
	switch stmt := node.(type) {
	// Handle assignment statements
	case *ast.AssignStmt:
		checkVar005AssignStmt(pass, stmt)
	// Handle range statements
	case *ast.RangeStmt:
		checkVar005RangeStmt(pass, stmt)
	// Handle var declarations in blocks
	case *ast.DeclStmt:
		checkVar005DeclStmt(pass, stmt)
	}
}

// checkVar005AssignStmt vérifie une assignation.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: statement d'assignation
func checkVar005AssignStmt(pass *analysis.Pass, stmt *ast.AssignStmt) {
	// Only check short declarations (:=)
	if stmt.Tok != token.DEFINE {
		// Not a short declaration
		return
	}

	// Check each left-hand side identifier
	for _, lhs := range stmt.Lhs {
		ident, ok := lhs.(*ast.Ident)
		// Skip if not an identifier
		if !ok {
			continue
		}

		// Check if name is too long
		checkVar005Name(pass, ident)
	}
}

// checkVar005RangeStmt vérifie une range statement.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: range statement
func checkVar005RangeStmt(pass *analysis.Pass, stmt *ast.RangeStmt) {
	// Only check if it's a definition
	if stmt.Tok != token.DEFINE {
		// Not defining new variables
		return
	}

	// Check key variable
	if key, ok := stmt.Key.(*ast.Ident); ok {
		checkVar005Name(pass, key)
	}

	// Check value variable
	if stmt.Value != nil {
		if value, ok := stmt.Value.(*ast.Ident); ok {
			checkVar005Name(pass, value)
		}
	}
}

// checkVar005DeclStmt vérifie une déclaration var dans un bloc.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: statement de déclaration
func checkVar005DeclStmt(pass *analysis.Pass, stmt *ast.DeclStmt) {
	genDecl, ok := stmt.Decl.(*ast.GenDecl)
	// Skip if not a GenDecl
	if !ok {
		// Not a general declaration
		return
	}

	// Only check var declarations
	if genDecl.Tok != token.VAR {
		// Continue
		return
	}

	// Check each variable specification
	for _, spec := range genDecl.Specs {
		valueSpec := spec.(*ast.ValueSpec)
		checkVar005Spec(pass, valueSpec)
	}
}

// checkVar005Spec vérifie une spécification de variable.
//
// Params:
//   - pass: contexte d'analyse
//   - valueSpec: spécification de variable
func checkVar005Spec(pass *analysis.Pass, valueSpec *ast.ValueSpec) {
	// Check each variable name
	for _, name := range valueSpec.Names {
		checkVar005Name(pass, name)
	}
}

// checkVar005Name vérifie si un nom de variable est trop long.
//
// Params:
//   - pass: contexte d'analyse
//   - ident: identifiant à vérifier
func checkVar005Name(pass *analysis.Pass, ident *ast.Ident) {
	varName := ident.Name

	// Skip blank identifier
	if varName == "_" {
		// Blank identifier is always allowed
		return
	}

	// Check if name is too long
	if len(varName) <= maxVarNameLength005 {
		// Name is within limit
		return
	}

	// Report error
	msg, ok := messages.Get(ruleCodeVar005)
	// Defensive: avoid panic if message is missing
	if !ok {
		pass.Reportf(ident.Pos(), "%s: nom de variable trop long: %q", ruleCodeVar005, varName)
		return
	}
	pass.Reportf(
		ident.Pos(),
		"%s: %s",
		ruleCodeVar005,
		msg.Format(config.Get().Verbose, varName),
	)
}
