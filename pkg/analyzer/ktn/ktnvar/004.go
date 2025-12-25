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
	// ruleCodeVar004 is the rule code for this analyzer
	ruleCodeVar004 string = "KTN-VAR-004"
	// minMakeArgs is the minimum number of arguments for make call
	minMakeArgs int = 2
	// initialAppendVarsCap initial capacity for append variables map
	initialAppendVarsCap int = 16
)

// Analyzer004 checks that slices are preallocated with capacity when known
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar004",
	Doc:      "KTN-VAR-004: Vérifie que les slices sont préalloués avec une capacité si elle est connue",
	Run:      runVar004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

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

	// Collecter les variables utilisées avec append
	appendVars := collectAppendVariables(insp)

	// Vérifier les make() sans capacité
	checkMakeCalls(pass, insp)

	// Vérifier les []T{} qui devraient être préalloués
	checkEmptySliceLiterals(pass, insp, appendVars)

	// Retour de la fonction
	return nil, nil
}

// collectAppendVariables collecte les variables utilisées avec append.
//
// Params:
//   - insp: inspecteur AST
//
// Returns:
//   - map[string]bool: map des noms de variables utilisées avec append
func collectAppendVariables(insp *inspector.Inspector) map[string]bool {
	appendVars := make(map[string]bool, initialAppendVarsCap)

	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
	}

	// Parcours des assignations pour trouver les appends
	insp.Preorder(nodeFilter, func(n ast.Node) {
		assign, ok := n.(*ast.AssignStmt)
		// Vérification de la condition
		if !ok {
			// Continue traversing AST nodes
			return
		}

		// Vérification de chaque expression à droite
		for _, rhs := range assign.Rhs {
			// Vérification si c'est un appel à append
			if isAppendCall(rhs) {
				// Récupération des variables à gauche
				for _, lhs := range assign.Lhs {
					// Extraction du nom de la variable
					if ident, isIdent := lhs.(*ast.Ident); isIdent {
						appendVars[ident.Name] = true
					}
				}
			}
		}
	})

	// Retour de la map
	return appendVars
}

// isAppendCall vérifie si une expression est un appel à append.
//
// Params:
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si c'est un appel à append
func isAppendCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	// Vérification de la condition
	if !ok {
		// Ce n'est pas un appel de fonction
		return false
	}

	// Vérification du nom de la fonction
	ident, ok := call.Fun.(*ast.Ident)
	// Vérification de la condition
	if !ok {
		// Ce n'est pas un identifiant simple
		return false
	}

	// Retour du résultat
	return ident.Name == "append"
}

// checkMakeCalls vérifie les appels à make sans capacité.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
func checkMakeCalls(pass *analysis.Pass, insp *inspector.Inspector) {
	// Récupération de la configuration
	cfg := config.Get()

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	// Parcours des appels de fonction
	insp.Preorder(nodeFilter, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		// Vérification de la condition
		if !ok {
			// Continue traversing AST nodes
			return
		}

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar004, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Vérification de l'appel make
		checkMakeCall(pass, call)
	})
}

// checkMakeCall vérifie un appel à make pour les slices sans capacité.
//
// Params:
//   - pass: contexte d'analyse
//   - call: appel de fonction à vérifier
func checkMakeCall(pass *analysis.Pass, call *ast.CallExpr) {
	// Vérification que c'est un appel à make
	if !utils.IsMakeCall(call) {
		// Continue traversing AST nodes
		return
	}

	// Vérification du nombre d'arguments (doit être 2: type et length)
	if len(call.Args) != minMakeArgs {
		// Continue traversing AST nodes
		return
	}

	// Vérification que le type est un slice
	if !utils.IsSliceTypeWithPass(pass, call.Args[0]) {
		// Continue traversing AST nodes
		return
	}

	// Skip si VAR-016 s'applique (constante <= 1024, suggère array)
	if utils.IsSmallConstantSize(pass, call.Args[1]) {
		// VAR-016 gère ce cas
		return
	}

	// Signalement de l'erreur
	msg, _ := messages.Get(ruleCodeVar004)
	pass.Reportf(
		call.Pos(),
		"%s: %s",
		ruleCodeVar004,
		msg.Format(config.Get().Verbose),
	)
}

// litCheckContext contains context for slice literal checking.
type litCheckContext struct {
	pass       *analysis.Pass
	appendVars map[string]bool
}

// checkEmptySliceLiterals vérifie les []T{} qui devraient être préalloués.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - appendVars: variables utilisées avec append
func checkEmptySliceLiterals(
	pass *analysis.Pass,
	insp *inspector.Inspector,
	appendVars map[string]bool,
) {
	// Récupération de la configuration
	cfg := config.Get()

	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
	}

	// Créer le contexte de vérification
	ctx := &litCheckContext{
		pass:       pass,
		appendVars: appendVars,
	}

	// Parcours des assignations
	insp.WithStack(nodeFilter, func(n ast.Node, push bool, stack []ast.Node) bool {
		// Ignorer le pop
		if !push {
			// Continuer le parcours
			return true
		}

		assign, ok := n.(*ast.AssignStmt)
		// Vérification de la condition
		if !ok {
			// Continuer le parcours
			return true
		}

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar004, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return true
		}

		// Vérification de chaque paire lhs/rhs
		for i, rhs := range assign.Rhs {
			// Vérification que c'est un composite literal
			lit, isLit := rhs.(*ast.CompositeLit)
			// Vérification de la condition
			if !isLit {
				// Continuer avec l'élément suivant
				continue
			}

			// Vérification de la slice vide
			checkCompositeLit(ctx, assign, i, lit, stack)
		}

		// Continuer le parcours
		return true
	})
}

// checkCompositeLit vérifie un composite literal pour les slices vides.
//
// Params:
//   - ctx: contexte de vérification
//   - assign: assignation contenant le literal
//   - index: index dans la liste des rhs
//   - lit: composite literal à vérifier
//   - stack: pile des nœuds parents
func checkCompositeLit(
	ctx *litCheckContext,
	assign *ast.AssignStmt,
	index int,
	lit *ast.CompositeLit,
	stack []ast.Node,
) {
	// Vérification que c'est un slice vide
	if len(lit.Elts) > 0 {
		// Le slice n'est pas vide
		return
	}

	// Vérification que le type est un slice
	if !utils.IsSliceTypeWithPass(ctx.pass, lit.Type) {
		// Ce n'est pas un slice
		return
	}

	// Vérification si on est dans un return statement
	if isInReturnStatement(stack) {
		// Pas de préallocation nécessaire pour un return
		return
	}

	// Vérification si on est dans un struct literal
	if isInStructLiteral(stack) {
		// Pas de préallocation nécessaire pour init de struct
		return
	}

	// Récupération du nom de la variable assignée
	if index >= len(assign.Lhs) {
		// Index invalide
		return
	}

	ident, ok := assign.Lhs[index].(*ast.Ident)
	// Vérification de la condition
	if !ok {
		// Ce n'est pas un identifiant simple
		return
	}

	// Vérification si la variable est utilisée avec append
	if !ctx.appendVars[ident.Name] {
		// La variable n'est jamais utilisée avec append
		return
	}

	// Signalement de l'erreur
	msg, _ := messages.Get(ruleCodeVar004)
	ctx.pass.Reportf(
		lit.Pos(),
		"%s: %s",
		ruleCodeVar004,
		msg.Format(config.Get().Verbose),
	)
}

// isInReturnStatement vérifie si le nœud est dans un return statement.
//
// Params:
//   - stack: pile des nœuds parents
//
// Returns:
//   - bool: true si dans un return
func isInReturnStatement(stack []ast.Node) bool {
	// Parcours de la pile des parents
	for _, node := range stack {
		// Vérification si c'est un return statement
		if _, ok := node.(*ast.ReturnStmt); ok {
			// Trouvé un return parent
			return true
		}
	}

	// Pas de return parent trouvé
	return false
}

// isInStructLiteral vérifie si le nœud est dans un struct literal.
//
// Params:
//   - stack: pile des nœuds parents
//
// Returns:
//   - bool: true si dans un struct literal
func isInStructLiteral(stack []ast.Node) bool {
	// Parcours de la pile des parents (en excluant le nœud courant)
	for i := len(stack) - 1; i >= 0; i-- {
		node := stack[i]

		// Vérification si c'est un composite literal (struct)
		if lit, ok := node.(*ast.CompositeLit); ok {
			// Vérification que ce n'est pas un slice/array/map
			if !isSliceArrayOrMap(lit.Type) {
				// C'est un struct literal
				return true
			}
		}

		// Vérification si c'est un key-value expression
		if _, ok := node.(*ast.KeyValueExpr); ok {
			// Dans une initialisation de champ
			return true
		}
	}

	// Pas de struct parent trouvé
	return false
}

// isSliceArrayOrMap vérifie si le type est un slice, array ou map.
//
// Params:
//   - typeExpr: expression de type
//
// Returns:
//   - bool: true si slice, array ou map
func isSliceArrayOrMap(typeExpr ast.Expr) bool {
	// Vérification du type nil
	if typeExpr == nil {
		// Type implicite (peut être struct)
		return false
	}

	// Vérification des différents types
	switch typeExpr.(type) {
	// Traitement des types slice/array/map
	case *ast.ArrayType, *ast.MapType:
		// C'est un slice, array ou map
		return true
	// Traitement des autres types
	default:
		// Ce n'est pas un slice, array ou map
		return false
	}
}
