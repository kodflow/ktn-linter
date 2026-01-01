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

	// loopVars004 contains single-letter loop variable names that are allowed.
	loopVars004 map[string]bool = map[string]bool{
		"i": true, "j": true, "k": true, "n": true,
		"x": true, "y": true, "z": true, "v": true,
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
				valueSpec := spec.(*ast.ValueSpec)
				checkVar004Spec(pass, valueSpec, false)
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

		// Track loop init positions to avoid double-processing
		loopInits := collectLoopInitPositions004(funcDecl.Body)

		// Visit all statements in the function
		ast.Inspect(funcDecl.Body, func(node ast.Node) bool {
			// Check the node type
			checkVar004Node(pass, node, loopInits)
			// Continue traversal
			return true
		})
	})
}

// collectLoopInitPositions004 collects positions of for/range loop init vars.
//
// Params:
//   - body: function body to analyze
//
// Returns:
//   - map[token.Pos]bool: set of positions that are loop inits
func collectLoopInitPositions004(body *ast.BlockStmt) map[token.Pos]bool {
	positions := make(map[token.Pos]bool)

	// Walk the AST to find loop init positions
	ast.Inspect(body, func(node ast.Node) bool {
		// Check for for statements
		switch stmt := node.(type) {
		// Handle for statements
		case *ast.ForStmt:
			// Add init positions
			if init, ok := stmt.Init.(*ast.AssignStmt); ok {
				if init.Tok == token.DEFINE {
					for _, lhs := range init.Lhs {
						if ident, ok := lhs.(*ast.Ident); ok {
							positions[ident.Pos()] = true
						}
					}
				}
			}
		// Handle range statements
		case *ast.RangeStmt:
			// Add key position
			if key, ok := stmt.Key.(*ast.Ident); ok {
				positions[key.Pos()] = true
			}
			// Add value position
			if stmt.Value != nil {
				if value, ok := stmt.Value.(*ast.Ident); ok {
					positions[value.Pos()] = true
				}
			}
		}
		// Continue traversal
		return true
	})

	// Return positions
	return positions
}

// checkVar004Node vérifie un nœud AST pour les variables courtes.
//
// Params:
//   - pass: contexte d'analyse
//   - node: nœud à vérifier
//   - loopInits: positions of loop init variables
func checkVar004Node(pass *analysis.Pass, node ast.Node, loopInits map[token.Pos]bool) {
	// Switch on node type
	switch stmt := node.(type) {
	// Handle assignment statements
	case *ast.AssignStmt:
		checkVar004AssignStmt(pass, stmt, loopInits)
	// Handle var declarations in blocks
	case *ast.DeclStmt:
		checkVar004DeclStmt(pass, stmt)
	}
	// Note: ForStmt and RangeStmt are handled via collectLoopInitPositions004
}

// checkVar004AssignStmt vérifie une assignation pour les noms courts.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: statement d'assignation
//   - loopInits: positions of loop init variables
func checkVar004AssignStmt(
	pass *analysis.Pass,
	stmt *ast.AssignStmt,
	loopInits map[token.Pos]bool,
) {
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

		// Check if this is a loop init variable
		isLoopVar := loopInits[ident.Pos()]

		// Check if name is too short
		checkVar004Name(pass, ident, isLoopVar)
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
		valueSpec := spec.(*ast.ValueSpec)
		checkVar004Spec(pass, valueSpec, false)
	}
}

// checkVar004Spec vérifie une spécification de variable.
//
// Params:
//   - pass: contexte d'analyse
//   - valueSpec: spécification de variable
//   - isLoop: indique si c'est dans un contexte de boucle
func checkVar004Spec(pass *analysis.Pass, valueSpec *ast.ValueSpec, isLoop bool) {
	// Check each variable name
	for _, name := range valueSpec.Names {
		checkVar004Name(pass, name, isLoop)
	}
}

// checkVar004Name vérifie si un nom de variable est trop court.
//
// Params:
//   - pass: contexte d'analyse
//   - ident: identifiant à vérifier
//   - isLoop: indique si c'est dans un contexte de boucle
func checkVar004Name(pass *analysis.Pass, ident *ast.Ident, isLoop bool) {
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

	// Allow loop variables in loop context
	if isLoop && loopVars004[varName] {
		// Loop variable is allowed
		return
	}

	// Allow idiomatic short names
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
