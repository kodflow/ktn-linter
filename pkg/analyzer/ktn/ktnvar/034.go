// Package ktnvar implements KTN linter rules.
package ktnvar

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar034 is the rule code for this analyzer
	ruleCodeVar034 string = "KTN-VAR-034"
)

// Analyzer034 detecte le pattern wg.Add(1) + go func() + defer wg.Done()
// qui devrait etre remplace par wg.Go() depuis Go 1.25.
//
// Returns:
//   - *analysis.Analyzer: analyseur pour KTN-VAR-034
var Analyzer034 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar034",
	Doc:      "KTN-VAR-034: Detecte wg.Add(1)+go func()+defer wg.Done() - utiliser wg.Go() (Go 1.25+)",
	Run:      runVar034,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar034 execute l'analyse de detection du pattern WaitGroup obsolete.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: erreur eventuelle
func runVar034(pass *analysis.Pass) (any, error) {
	// Recuperation de la configuration
	cfg := config.Get()

	// Verifier si la regle est activee
	if !cfg.IsRuleEnabled(ruleCodeVar034) {
		// Regle desactivee
		return nil, nil
	}

	// Get AST inspector
	inspAny := pass.ResultOf[inspect.Analyzer]
	insp := inspAny.(*inspector.Inspector)

	// Types de noeuds a analyser (blocs de code)
	nodeFilter := []ast.Node{
		(*ast.BlockStmt)(nil),
	}

	// Parcours des blocs de code
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeVar034, pass.Fset.Position(n.Pos()).Filename) {
			// Fichier exclu
			return
		}

		// Cast en bloc
		block := n.(*ast.BlockStmt)
		// Verification de la condition
		if block.List == nil {
			// Traitement
			return
		}

		// Analyse des statements consecutifs
		analyzeConsecutiveStatements(pass, block.List)
	})

	// Traitement
	return nil, nil
}

// analyzeConsecutiveStatements analyse les statements consecutifs pour detecter le pattern.
//
// Params:
//   - pass: contexte d'analyse
//   - stmts: liste de statements a analyser
func analyzeConsecutiveStatements(pass *analysis.Pass, stmts []ast.Stmt) {
	// Parcours des statements consecutifs
	for i := range len(stmts) - 1 {
		// Verification du pattern wg.Add(1) suivi de go func()
		addStmt := stmts[i]
		goStmt := stmts[i+1]

		// Verification si c'est le pattern recherche
		if isWaitGroupPattern(pass, addStmt, goStmt) {
			// Rapport d'erreur sur le go statement
			reportVar034(pass, goStmt)
		}
	}
}

// isWaitGroupPattern verifie si deux statements forment le pattern wg.Add(1) + go func().
//
// Params:
//   - pass: contexte d'analyse
//   - addStmt: statement candidat pour wg.Add(1)
//   - goStmt: statement candidat pour go func()
//
// Returns:
//   - bool: true si c'est le pattern
func isWaitGroupPattern(pass *analysis.Pass, addStmt, goStmt ast.Stmt) bool {
	// Extraction du WaitGroup depuis wg.Add(1)
	wgName := extractWaitGroupAdd1(pass, addStmt)
	// Verification de la condition
	if wgName == "" {
		// Pas un wg.Add(1)
		return false
	}

	// Verification du go statement avec defer wg.Done()
	return isGoWithDeferDone(pass, goStmt, wgName)
}

// extractWaitGroupAdd1 extrait le nom du WaitGroup depuis un statement wg.Add(1).
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: statement a analyser
//
// Returns:
//   - string: nom du WaitGroup ou "" si pas un wg.Add(1)
func extractWaitGroupAdd1(pass *analysis.Pass, stmt ast.Stmt) string {
	// Cast en expression statement
	exprStmt, ok := stmt.(*ast.ExprStmt)
	// Verification de la condition
	if !ok {
		// Pas un expression statement
		return ""
	}

	// Cast en call expression
	call, ok := exprStmt.X.(*ast.CallExpr)
	// Verification de la condition
	if !ok {
		// Pas un appel de fonction
		return ""
	}

	// Verification de l'argument (doit etre 1)
	if !isCallWithArg1(call) {
		// Pas un appel avec argument 1
		return ""
	}

	// Extraction du selecteur (wg.Add)
	return extractWaitGroupMethodCall(pass, call, "Add")
}

// isCallWithArg1 verifie si un appel a un seul argument egal a 1.
//
// Params:
//   - call: appel de fonction a verifier
//
// Returns:
//   - bool: true si l'argument est 1
func isCallWithArg1(call *ast.CallExpr) bool {
	// Verification du nombre d'arguments
	if len(call.Args) != 1 {
		// Pas exactement un argument
		return false
	}

	// Verification si l'argument est un literal entier
	lit, ok := call.Args[0].(*ast.BasicLit)
	// Verification de la condition
	if !ok || lit.Kind != token.INT {
		// Pas un literal entier
		return false
	}

	// Verification de la valeur
	return lit.Value == "1"
}

// extractWaitGroupMethodCall extrait le nom du WaitGroup depuis un appel de methode.
//
// Params:
//   - pass: contexte d'analyse
//   - call: appel de fonction
//   - methodName: nom de la methode attendue
//
// Returns:
//   - string: nom du WaitGroup ou ""
func extractWaitGroupMethodCall(pass *analysis.Pass, call *ast.CallExpr, methodName string) string {
	// Cast en selecteur (wg.Add ou wg.Done)
	selector, ok := call.Fun.(*ast.SelectorExpr)
	// Verification de la condition
	if !ok {
		// Pas un selecteur
		return ""
	}

	// Verification du nom de la methode
	if selector.Sel.Name != methodName {
		// Pas la methode attendue
		return ""
	}

	// Verification que c'est un WaitGroup
	if !isWaitGroupReceiver(pass, selector.X) {
		// Pas un WaitGroup
		return ""
	}

	// Extraction du nom de la variable
	return getReceiverName(selector.X)
}

// isWaitGroupReceiver verifie si une expression est un sync.WaitGroup.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression a verifier
//
// Returns:
//   - bool: true si c'est un WaitGroup
func isWaitGroupReceiver(pass *analysis.Pass, expr ast.Expr) bool {
	// Recuperation du type
	tv, ok := pass.TypesInfo.Types[expr]
	// Verification de la condition
	if !ok {
		// Pas de type info
		return false
	}

	// Verification du type
	return isWaitGroupType(tv.Type)
}

// isWaitGroupType verifie si un type est sync.WaitGroup.
//
// Params:
//   - t: type a verifier
//
// Returns:
//   - bool: true si c'est un WaitGroup
func isWaitGroupType(t types.Type) bool {
	// Verification pour pointer vers WaitGroup
	if ptr, ok := t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	// Recuperation du type nomme
	named, ok := t.(*types.Named)
	// Verification de la condition
	if !ok {
		// Pas un type nomme
		return false
	}

	// Verification du package et du nom
	obj := named.Obj()
	// Verification de la condition
	if obj.Pkg() == nil {
		// Pas de package
		return false
	}

	// Verification sync.WaitGroup
	return obj.Pkg().Path() == "sync" && obj.Name() == "WaitGroup"
}

// getReceiverName extrait le nom d'une variable depuis une expression.
//
// Params:
//   - expr: expression a analyser
//
// Returns:
//   - string: nom de la variable ou ""
func getReceiverName(expr ast.Expr) string {
	// Cast en identifiant
	ident, ok := expr.(*ast.Ident)
	// Verification de la condition
	if !ok {
		// Pas un identifiant simple
		return ""
	}

	// Retour du nom
	return ident.Name
}

// isGoWithDeferDone verifie si un go statement a defer wg.Done() au debut.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: statement a verifier
//   - wgName: nom du WaitGroup attendu
//
// Returns:
//   - bool: true si c'est le pattern
func isGoWithDeferDone(pass *analysis.Pass, stmt ast.Stmt, wgName string) bool {
	// Cast en go statement
	goStmt, ok := stmt.(*ast.GoStmt)
	// Verification de la condition
	if !ok {
		// Pas un go statement
		return false
	}

	// Extraction de la fonction literal
	funcLit := extractFuncLiteral(goStmt)
	// Verification de la condition
	if funcLit == nil {
		// Pas une fonction literal
		return false
	}

	// Verification du defer wg.Done() au debut
	return hasDeferDoneFirst(pass, funcLit, wgName)
}

// extractFuncLiteral extrait la fonction literal d'un go statement.
//
// Params:
//   - goStmt: go statement
//
// Returns:
//   - *ast.FuncLit: fonction literal ou nil
func extractFuncLiteral(goStmt *ast.GoStmt) *ast.FuncLit {
	// Cast en call expression
	call, ok := goStmt.Call.Fun.(*ast.FuncLit)
	// Verification de la condition
	if !ok {
		// Pas une fonction literal
		return nil
	}

	// Retour de la fonction literal
	return call
}

// hasDeferDoneFirst verifie si la premiere instruction est defer wg.Done().
//
// Params:
//   - pass: contexte d'analyse
//   - funcLit: fonction literal
//   - wgName: nom du WaitGroup attendu
//
// Returns:
//   - bool: true si defer wg.Done() est la premiere instruction
func hasDeferDoneFirst(pass *analysis.Pass, funcLit *ast.FuncLit, wgName string) bool {
	// Verification du corps
	if funcLit.Body == nil || len(funcLit.Body.List) == 0 {
		// Corps vide
		return false
	}

	// Premiere instruction
	firstStmt := funcLit.Body.List[0]

	// Cast en defer statement
	deferStmt, ok := firstStmt.(*ast.DeferStmt)
	// Verification de la condition
	if !ok {
		// Pas un defer
		return false
	}

	// Extraction du WaitGroup depuis le defer call
	return isDeferDoneForWaitGroup(pass, deferStmt, wgName)
}

// isDeferDoneForWaitGroup verifie si un defer appelle wg.Done() pour le bon WaitGroup.
//
// Params:
//   - pass: contexte d'analyse
//   - deferStmt: defer statement
//   - wgName: nom du WaitGroup attendu
//
// Returns:
//   - bool: true si c'est defer wg.Done() pour le bon WaitGroup
func isDeferDoneForWaitGroup(pass *analysis.Pass, deferStmt *ast.DeferStmt, wgName string) bool {
	// Extraction du WaitGroup depuis l'appel
	doneWgName := extractWaitGroupMethodCall(pass, deferStmt.Call, "Done")

	// Verification que c'est le meme WaitGroup
	return doneWgName == wgName
}

// reportVar034 rapporte une erreur KTN-VAR-034.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: statement ou reporter l'erreur
func reportVar034(pass *analysis.Pass, stmt ast.Stmt) {
	// Recuperation du message
	msg, _ := messages.Get(ruleCodeVar034)

	// Rapport d'erreur
	pass.Reportf(
		stmt.Pos(),
		"%s: %s",
		ruleCodeVar034,
		msg.Format(config.Get().Verbose),
	)
}
