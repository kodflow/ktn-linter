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
	// ruleCodeVar028 is the rule code for this analyzer
	ruleCodeVar028 string = "KTN-VAR-028"
)

// Analyzer028 checks for obsolete loop variable copying pattern (Go 1.22+)
var Analyzer028 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar028",
	Doc:      "KTN-VAR-028: Le pattern v := v dans les boucles est obsolete depuis Go 1.22",
	Run:      runVar028,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar028 execute l'analyse KTN-VAR-028.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: resultat de l'analyse
//   - error: erreur eventuelle
func runVar028(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeVar028) {
		// Regle desactivee
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp, ok := inspAny.(*inspector.Inspector)
	// Defensive: ensure inspector is available
	if !ok || insp == nil {
		return nil, nil
	}
	// Defensive: avoid nil dereference when resolving positions
	if pass.Fset == nil {
		return nil, nil
	}

	nodeFilter := []ast.Node{
		(*ast.RangeStmt)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		rangeStmt, ok := n.(*ast.RangeStmt)
		// Defensive: ensure node type matches
		if !ok {
			return
		}

		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar028, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Verification du pattern obsolete
		checkLoopVarCopyPattern(pass, rangeStmt)
	})

	// Retour de la fonction
	return nil, nil
}

// checkLoopVarCopyPattern verifie si une boucle range contient v := v.
//
// Params:
//   - pass: contexte d'analyse
//   - rangeStmt: boucle range a verifier
func checkLoopVarCopyPattern(pass *analysis.Pass, rangeStmt *ast.RangeStmt) {
	// Recuperer les noms des variables de range
	rangeVars := getRangeVariableNames(rangeStmt)
	// Aucune variable a verifier
	if len(rangeVars) == 0 {
		// Pas de variables de range
		return
	}

	// Parcourir le corps de la boucle
	if rangeStmt.Body == nil {
		// Corps vide
		return
	}

	// Verifier chaque instruction dans le corps
	for _, stmt := range rangeStmt.Body.List {
		// Verifier les declarations courtes
		checkShortVarDecl(pass, stmt, rangeVars)
	}
}

// getRangeVariableNames retourne les noms des variables de range.
//
// Params:
//   - rangeStmt: boucle range
//
// Returns:
//   - map[string]bool: ensemble des noms de variables
func getRangeVariableNames(rangeStmt *ast.RangeStmt) map[string]bool {
	vars := make(map[string]bool)

	// Ajouter la cle si presente
	if rangeStmt.Key != nil {
		// Verification de l'identifiant
		if keyIdent, ok := rangeStmt.Key.(*ast.Ident); ok && keyIdent.Name != "_" {
			vars[keyIdent.Name] = true
		}
	}

	// Ajouter la valeur si presente
	if rangeStmt.Value != nil {
		// Verification de l'identifiant
		if valIdent, ok := rangeStmt.Value.(*ast.Ident); ok && valIdent.Name != "_" {
			vars[valIdent.Name] = true
		}
	}

	// Retour des variables
	return vars
}

// checkShortVarDecl verifie si une instruction est une declaration v := v.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: instruction a verifier
//   - rangeVars: variables de range
func checkShortVarDecl(pass *analysis.Pass, stmt ast.Stmt, rangeVars map[string]bool) {
	// Verifier si c'est un AssignStmt
	assignStmt, ok := stmt.(*ast.AssignStmt)
	// Pas une affectation
	if !ok {
		// Pas un statement d'affectation
		return
	}

	// Doit etre une declaration courte (:=)
	if assignStmt.Tok != token.DEFINE {
		// Pas une declaration courte
		return
	}

	// Verifier chaque paire lhs := rhs
	for i := range assignStmt.Lhs {
		// Verification de la paire
		checkAssignmentPair(pass, assignStmt, i, rangeVars)
	}
}

// checkAssignmentPair verifie une paire d'affectation.
//
// Params:
//   - pass: contexte d'analyse
//   - assignStmt: statement d'affectation
//   - i: index de la paire
//   - rangeVars: variables de range
func checkAssignmentPair(
	pass *analysis.Pass,
	assignStmt *ast.AssignStmt,
	i int,
	rangeVars map[string]bool,
) {
	// Verifier les limites
	if i >= len(assignStmt.Rhs) {
		// Index hors limites
		return
	}

	// Lhs doit etre un identifiant
	lhsIdent, ok := assignStmt.Lhs[i].(*ast.Ident)
	// Verification du lhs
	if !ok {
		// Lhs n'est pas un identifiant
		return
	}

	// Rhs doit etre un identifiant avec le meme nom
	rhsIdent, ok := assignStmt.Rhs[i].(*ast.Ident)
	// Verification du rhs
	if !ok {
		// Rhs n'est pas un identifiant
		return
	}

	// Les noms doivent etre identiques
	if lhsIdent.Name != rhsIdent.Name {
		// Noms differents
		return
	}

	// La variable doit etre une variable de range
	if !rangeVars[lhsIdent.Name] {
		// Pas une variable de range
		return
	}

	// Pattern obsolete detecte
	reportLoopVarCopy(pass, assignStmt, lhsIdent.Name)
}

// reportLoopVarCopy signale un pattern v := v obsolete.
//
// Params:
//   - pass: contexte d'analyse
//   - node: noeud a signaler
//   - varName: nom de la variable
func reportLoopVarCopy(pass *analysis.Pass, node ast.Node, varName string) {
	msg, ok := messages.Get(ruleCodeVar028)
	// Defensive: avoid panic if message is missing
	if !ok {
		pass.Reportf(node.Pos(), "%s: %s := %s est obsolète depuis Go 1.22", ruleCodeVar028, varName, varName)
		return
	}
	pass.Reportf(
		node.Pos(),
		"%s: %s",
		ruleCodeVar028,
		msg.Format(config.Get().Verbose, varName),
	)
}
