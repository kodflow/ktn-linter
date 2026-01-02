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
	// ruleCodeVar012 is the rule code for this analyzer
	ruleCodeVar012 string = "KTN-VAR-012"
)

// Analyzer012 checks for slice/map allocations inside loops
var Analyzer012 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar012",
	Doc:      "KTN-VAR-012: Évite les allocations de slices/maps dans les boucles chaudes",
	Run:      runVar012,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar012 exécute l'analyse KTN-VAR-012.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar012(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar012) {
		// Règle désactivée
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp := inspAny.(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.ForStmt)(nil),
		(*ast.RangeStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar012, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}
		// Récupération du corps de la boucle
		var body *ast.BlockStmt
		// Vérification du type de boucle
		switch loop := n.(type) {
		// Cas d'une boucle for classique
		case *ast.ForStmt:
			// Boucle for classique
			body = loop.Body
		// Cas d'une boucle range
		case *ast.RangeStmt:
			// Boucle range
			body = loop.Body
		// Cas par défaut
		default:
			// Type de boucle non supporté
			return
		}

		// Parcours des instructions du corps
		checkLoopBodyForAlloc(pass, body)
	})

	// Retour de la fonction
	return nil, nil
}

// checkLoopBodyForAlloc vérifie les allocations dans le corps d'une boucle.
//
// Params:
//   - pass: contexte d'analyse
//   - body: corps de la boucle à vérifier
func checkLoopBodyForAlloc(pass *analysis.Pass, body *ast.BlockStmt) {
	// Vérification du corps de la boucle
	if body == nil {
		// Corps de boucle vide
		return
	}

	// Parcours des instructions
	for _, stmt := range body.List {
		checkStmtForAlloc(pass, stmt)
	}
}

// checkStmtForAlloc vérifie une instruction pour détecter allocations.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: instruction à vérifier
func checkStmtForAlloc(pass *analysis.Pass, stmt ast.Stmt) {
	// Vérification du type d'instruction
	switch s := stmt.(type) {
	// Cas d'une affectation
	case *ast.AssignStmt:
		// Vérification des affectations
		checkAssignForAlloc(pass, s)
	// Cas d'une déclaration
	case *ast.DeclStmt:
		// Vérification des déclarations
		checkDeclForAlloc(pass, s)
		// Note: les boucles imbriquées sont déjà gérées par Preorder
	}
}

// checkAssignForAlloc vérifie une affectation pour détecter allocations.
//
// Params:
//   - pass: contexte d'analyse
//   - assign: affectation à vérifier
func checkAssignForAlloc(pass *analysis.Pass, assign *ast.AssignStmt) {
	// Parcours des valeurs affectées
	for _, rhs := range assign.Rhs {
		// Vérification si allocation de slice ou map
		if isSliceOrMapAlloc(rhs) {
			// Allocation de slice/map détectée
			msg, _ := messages.Get(ruleCodeVar012)
			pass.Reportf(
				rhs.Pos(),
				"%s: %s",
				ruleCodeVar012,
				msg.Format(config.Get().Verbose),
			)
		}
	}
}

// checkDeclForAlloc vérifie une déclaration pour détecter allocations.
//
// Params:
//   - pass: contexte d'analyse
//   - decl: déclaration à vérifier
func checkDeclForAlloc(pass *analysis.Pass, decl *ast.DeclStmt) {
	genDecl, ok := decl.Decl.(*ast.GenDecl)
	// Vérification du type de déclaration
	if !ok {
		// Pas une déclaration générale
		return
	}

	var valueSpec *ast.ValueSpec
	// Parcours des spécifications
	for _, spec := range genDecl.Specs {
		valueSpec, ok = spec.(*ast.ValueSpec)
		// Vérification de la spécification de valeur
		if !ok {
			// Pas une spécification de valeur
			continue
		}

		// Parcours des valeurs
		for _, value := range valueSpec.Values {
			// Vérification si allocation de slice ou map
			if isSliceOrMapAlloc(value) {
				// Allocation de slice/map détectée
				msg, _ := messages.Get(ruleCodeVar012)
				pass.Reportf(
					value.Pos(),
					"%s: %s",
					ruleCodeVar012,
					msg.Format(config.Get().Verbose),
				)
			}
		}
	}
}

// isSliceOrMapAlloc vérifie si une expression est une allocation de slice/map.
// Exclut les []byte qui sont gérés par VAR-010.
//
// Params:
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si allocation détectée
func isSliceOrMapAlloc(expr ast.Expr) bool {
	// Vérification du type d'expression
	switch e := expr.(type) {
	// Cas d'un littéral composite
	case *ast.CompositeLit:
		// Vérification du type composite (exclut []byte géré par VAR-010)
		if utils.IsSliceOrMapType(e.Type) && !utils.IsByteSlice(e.Type) {
			// Allocation de slice/map sous forme de littéral
			return true
		}
	// Cas d'un appel de fonction
	case *ast.CallExpr:
		// Vérification des appels make() (exclut []byte géré par VAR-010)
		if utils.IsMakeCall(e) && !isByteSliceMake(e) {
			// Appel à make() détecté
			return true
		}
	}
	// Pas d'allocation détectée
	return false
}

// isByteSliceMake vérifie si make crée un []byte.
//
// Params:
//   - call: expression d'appel make
//
// Returns:
//   - bool: true si make([]byte, ...)
func isByteSliceMake(call *ast.CallExpr) bool {
	// Vérification des arguments
	if len(call.Args) == 0 {
		// Pas d'arguments
		return false
	}
	// Vérification du type
	return utils.IsByteSlice(call.Args[0])
}
