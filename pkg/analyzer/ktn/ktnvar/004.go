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
	// ruleCodeVar004 is the rule code for this analyzer
	ruleCodeVar004 string = "KTN-VAR-004"
	// minVarNameLength004 is the minimum variable name length
	minVarNameLength004 int = 2
)

var (
	// Analyzer004 checks that variable names are not too short.
	Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
		Name:     "ktnvar004",
		Doc:      "KTN-VAR-004: Les noms de variables doivent avoir au moins 2 caractères",
		Run:      runVar004,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	// idiomaticOneChar004 contains 1-char names allowed in function scope.
	idiomaticOneChar004 map[string]bool = map[string]bool{
		// Loop counters
		"i": true, "j": true, "k": true, "n": true,
		// Type hints
		"b": true, "c": true, "f": true, "m": true,
		"r": true, "s": true, "t": true, "w": true,
		// Results
		"_": true,
	}

	// idiomaticShort004 contains short idiomatic Go variable names.
	idiomaticShort004 map[string]bool = map[string]bool{
		"ok": true,
	}
)

// runVar004 exécute l'analyse KTN-VAR-004.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar004(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar004) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Check package-level variables
	checkVar004PackageLevel(pass, insp, cfg)

	// Check local variables in functions
	checkVar004LocalVars(pass, insp, cfg)

	// Return analysis result
	return nil, nil
}

// checkVar004PackageLevel vérifie les variables au niveau package.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - cfg: configuration
func checkVar004PackageLevel(
	pass *analysis.Pass,
	insp *inspector.Inspector,
	cfg *config.Config,
) {
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar004, pass.Fset.Position(n.Pos()).Filename) {
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
				valueSpec, ok := spec.(*ast.ValueSpec)
				// Defensive: only handle value specs
				if !ok {
					continue
				}
				checkVar004Spec(pass, valueSpec, true)
			}
		}
	})
}

// checkVar004LocalVars vérifie les variables locales dans les fonctions.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - cfg: configuration
func checkVar004LocalVars(
	pass *analysis.Pass,
	insp *inspector.Inspector,
	cfg *config.Config,
) {
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar004, pass.Fset.Position(n.Pos()).Filename) {
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
			checkVar004Node(pass, node)
			// Continue traversal
			return true
		})
	})
}

// checkVar004Node vérifie un nœud AST pour les variables courtes.
//
// Params:
//   - pass: contexte d'analyse
//   - node: nœud à vérifier
func checkVar004Node(pass *analysis.Pass, node ast.Node) {
	// Switch on node type
	switch stmt := node.(type) {
	// Handle assignment statements
	case *ast.AssignStmt:
		checkVar004AssignStmt(pass, stmt)
	// Handle var declarations in blocks
	case *ast.DeclStmt:
		checkVar004DeclStmt(pass, stmt)
	}
}

// checkVar004AssignStmt vérifie une assignation pour les noms courts.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: statement d'assignation
func checkVar004AssignStmt(pass *analysis.Pass, stmt *ast.AssignStmt) {
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

		// Check if name is too short (function-level)
		checkVar004Name(pass, ident, false)
	}
}

// checkVar004DeclStmt vérifie une déclaration var dans un bloc.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: statement de déclaration
func checkVar004DeclStmt(pass *analysis.Pass, stmt *ast.DeclStmt) {
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
		valueSpec, ok := spec.(*ast.ValueSpec)
		// Defensive: only handle value specs
		if !ok {
			continue
		}
		checkVar004Spec(pass, valueSpec, false)
	}
}

// checkVar004Spec vérifie une spécification de variable.
//
// Params:
//   - pass: contexte d'analyse
//   - valueSpec: spécification de variable
//   - isPackageLevel: indique si c'est une variable package-level
func checkVar004Spec(pass *analysis.Pass, valueSpec *ast.ValueSpec, isPackageLevel bool) {
	// Check each variable name
	for _, name := range valueSpec.Names {
		checkVar004Name(pass, name, isPackageLevel)
	}
}

// checkVar004Name vérifie si un nom de variable est trop court.
//
// Params:
//   - pass: contexte d'analyse
//   - ident: identifiant à vérifier
//   - isPackageLevel: indique si c'est une variable package-level
func checkVar004Name(pass *analysis.Pass, ident *ast.Ident, isPackageLevel bool) {
	varName := ident.Name

	// Skip blank identifier
	if varName == "_" {
		// Blank identifier is always allowed
		return
	}

	// Check if name is long enough
	if len(varName) >= minVarNameLength004 {
		// Name is long enough
		return
	}

	// Package-level: require min 2 chars always
	if isPackageLevel {
		// Report error for package-level short names
		msg, _ := messages.Get(ruleCodeVar004)
		pass.Reportf(
			ident.Pos(),
			"%s: %s",
			ruleCodeVar004,
			msg.Format(config.Get().Verbose, varName),
		)
		return
	}

	// Function-level: allow idiomatic 1-char names
	if idiomaticOneChar004[varName] {
		// Idiomatic 1-char name is allowed
		return
	}

	// Allow idiomatic short names like "ok"
	if idiomaticShort004[varName] {
		// Idiomatic name is allowed
		return
	}

	// Report error
	msg, _ := messages.Get(ruleCodeVar004)
	pass.Reportf(
		ident.Pos(),
		"%s: %s",
		ruleCodeVar004,
		msg.Format(config.Get().Verbose, varName),
	)
}
