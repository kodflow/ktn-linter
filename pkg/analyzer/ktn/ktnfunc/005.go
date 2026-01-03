// Package ktnfunc implements KTN linter rules.
package ktnfunc

import (
	"go/ast"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeFunc005 is the rule code for this analyzer
	ruleCodeFunc005 string = "KTN-FUNC-005"
	// defaultMaxStatements is the default max statements in a function
	defaultMaxStatements int = 35
)

// Analyzer005 checks that functions don't exceed 35 statements
var Analyzer005 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc005",
	Doc:      "KTN-FUNC-005: Les fonctions ne doivent pas dépasser 35 statements (instructions logiques)",
	Run:      runFunc005,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc005 exécute l'analyse KTN-FUNC-005.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runFunc005(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeFunc005) {
		// Règle désactivée
		return nil, nil
	}

	// Récupérer le seuil configuré
	maxStmts := cfg.GetThreshold(ruleCodeFunc005, defaultMaxStatements)

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeFunc005, filename) {
			// Fichier exclu
			return
		}

		// Skip if no body (external functions)
		if funcDecl.Body == nil {
			// Retour de la fonction
			return
		}

		// Skip test functions (Test*, Benchmark*, Example*, Fuzz*)
		if shared.IsTestFunction(funcDecl) {
			// Retour de la fonction
			return
		}

		// Skip main function
		funcName := funcDecl.Name.Name
		// Vérification fonction main
		if funcName == "main" {
			// Retour de la fonction
			return
		}

		// Count statements (logical instructions)
		stmtCount := countStatements(funcDecl.Body)

		// Vérification de la condition
		if stmtCount > maxStmts {
			msg, _ := messages.Get(ruleCodeFunc005)
			pass.Reportf(
				funcDecl.Name.Pos(),
				"%s: %s",
				ruleCodeFunc005,
				msg.Format(config.Get().Verbose, funcName, stmtCount, maxStmts),
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// countStatements compte les statements (instructions logiques) dans un bloc.
// Chaque statement compte pour 1, peu importe le nombre de lignes qu'il occupe.
// Les blocs imbriqués (if, for, switch) ajoutent leurs statements au compte.
//
// Params:
//   - body: bloc de code à analyser
//
// Returns:
//   - int: nombre total de statements
func countStatements(body *ast.BlockStmt) int {
	// Vérification bloc nil
	if body == nil {
		// Retour 0 si pas de bloc
		return 0
	}

	count := 0
	// Parcours des statements du bloc
	for _, stmt := range body.List {
		// Ajout du compte de ce statement
		count += countStmtComplexity(stmt)
	}
	// Retour du compte total
	return count
}

// countStmtComplexity compte la complexité d'un statement.
// Un statement simple = 1, les blocs imbriqués ajoutent leurs statements.
//
// Params:
//   - stmt: statement à analyser
//
// Returns:
//   - int: complexité du statement
func countStmtComplexity(stmt ast.Stmt) int {
	// Analyse selon le type de statement
	switch s := stmt.(type) {
	// Bloc de code
	case *ast.BlockStmt:
		// Compte des statements du bloc
		return countStatements(s)
	// Structure if
	case *ast.IfStmt:
		// Compte du if avec ses branches
		return countIfStmt(s)
	// Boucle for classique
	case *ast.ForStmt:
		// 1 pour le for + statements du corps
		return 1 + countStatements(s.Body)
	// Boucle range
	case *ast.RangeStmt:
		// 1 pour le range + statements du corps
		return 1 + countStatements(s.Body)
	// Switch sur valeur
	case *ast.SwitchStmt:
		// Compte du switch avec ses cases
		return countSwitchStmt(s.Body)
	// Switch sur type
	case *ast.TypeSwitchStmt:
		// Compte du switch avec ses cases
		return countSwitchStmt(s.Body)
	// Select pour channels
	case *ast.SelectStmt:
		// Compte du select avec ses cases
		return countSwitchStmt(s.Body)
	// Statement simple (return, assign, expr, etc.)
	default:
		// 1 instruction logique
		return 1
	}
}

// countIfStmt compte la complexité d'un if statement.
//
// Params:
//   - s: if statement à analyser
//
// Returns:
//   - int: complexité totale
func countIfStmt(s *ast.IfStmt) int {
	// 1 pour le if lui-même
	count := 1
	// Ajout des statements du corps
	count += countStatements(s.Body)
	// Traitement du else si présent
	if s.Else != nil {
		// Vérification type du else
		switch elseStmt := s.Else.(type) {
		// Bloc else simple
		case *ast.BlockStmt:
			// Ajout des statements du else
			count += countStatements(elseStmt)
		// Else if
		case *ast.IfStmt:
			// Récursion sur le else if
			count += countIfStmt(elseStmt)
		}
	}
	// Retour du compte
	return count
}
// Fin de countIfStmt

// countSwitchStmt compte la complexité d'un switch/select statement.
//
// Params:
//   - body: corps du switch/select
//
// Returns:
//   - int: complexité totale
func countSwitchStmt(body *ast.BlockStmt) int {
	// Vérification bloc nil
	if body == nil {
		// Retour 1 pour le switch vide
		return 1
	}
	// 1 pour le switch lui-même
	count := 1
	// Parcours des cases
	for _, stmt := range body.List {
		// Vérification si c'est un case clause
		if cc, ok := stmt.(*ast.CaseClause); ok {
			// 1 pour le case + ses statements
			count += 1 + len(cc.Body)
		}
		// Vérification si c'est un comm clause (select)
		if cc, ok := stmt.(*ast.CommClause); ok {
			// 1 pour le case + ses statements
			count += 1 + len(cc.Body)
		}
	}
	// Retour du compte total
	return count
}
// Fin de countSwitchStmt
