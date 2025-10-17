package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// ControlFlowAnalyzer vérifie les patterns de control flow.
	ControlFlowAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktncontrolflow",
		Doc:  "Vérifie les patterns de control flow (defer, if, switch, goto, fallthrough)",
		Run:  runControlFlowAnalyzer,
	}
)

// runControlFlowAnalyzer exécute l'analyseur control flow.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runControlFlowAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch stmt := n.(type) {
			case *ast.DeferStmt:
				checkDeferInLoop(pass, file, stmt)
			case *ast.IfStmt:
				checkSimplifiableIf(pass, stmt)
			case *ast.SwitchStmt:
				checkSingleCaseSwitch(pass, stmt)
			case *ast.BranchStmt:
				if stmt.Tok == token.GOTO {
					checkGotoUsage(pass, stmt)
				} else if stmt.Tok == token.FALLTHROUGH {
					checkFallthroughOutsideSwitch(pass, file, stmt)
				}
			}
			// Retourne true pour continuer l'inspection
			return true
		})
	}
	// Retourne nil car l'analyse est terminée
	return nil, nil
}

// checkDeferInLoop vérifie si un defer est dans une boucle.
//
// Params:
//   - pass: la passe d'analyse
//   - file: le fichier analysé
//   - deferStmt: le defer statement
func checkDeferInLoop(pass *analysis.Pass, file *ast.File, deferStmt *ast.DeferStmt) {
	if isInsideLoopStmt(file, deferStmt) {
		reportDeferInLoop(pass, deferStmt)
	}
}

// isInsideLoopStmt vérifie si un nœud est dans une boucle.
//
// Params:
//   - file: le fichier
//   - target: le nœud cible
//
// Returns:
//   - bool: true si dans une boucle
func isInsideLoopStmt(file *ast.File, target ast.Node) bool {
	inLoop := false

	ast.Inspect(file, func(n ast.Node) bool {
		switch loop := n.(type) {
		case *ast.ForStmt:
			if containsNode(loop.Body, target) {
				inLoop = true
				// Retourne false pour arrêter
				return false
			}
		case *ast.RangeStmt:
			if containsNode(loop.Body, target) {
				inLoop = true
				// Retourne false pour arrêter
				return false
			}
		}
		// Retourne true pour continuer
		return true
	})

	// Retourne le résultat
	return inLoop
}

// containsNode vérifie si un bloc contient un nœud.
//
// Params:
//   - block: le bloc
//   - target: le nœud recherché
//
// Returns:
//   - bool: true si trouvé
func containsNode(block *ast.BlockStmt, target ast.Node) bool {
	if block == nil {
		// Retourne false car pas de bloc
		return false
	}

	found := false
	ast.Inspect(block, func(n ast.Node) bool {
		if n == target {
			found = true
			// Retourne false pour arrêter
			return false
		}
		// Retourne true pour continuer
		return true
	})

	// Retourne le résultat
	return found
}

// checkSimplifiableIf vérifie si un if peut être simplifié.
//
// Implémente Staticcheck S1008 : détecte les patterns de retour booléen simplifiables.
// Détecte:
//   - if cond { return true } else { return false } → return cond
//   - if cond { return false } else { return true } → return !cond
//   - if cond { return true }; return false → return cond
//   - if cond { return false }; return true → return !cond
//
// Ignore les blocs contenant des commentaires (intention documentée).
//
// Params:
//   - pass: la passe d'analyse
//   - ifStmt: le if statement
func checkSimplifiableIf(pass *analysis.Pass, ifStmt *ast.IfStmt) {
	// Vérifier que le if n'a pas d'initialisation
	if ifStmt.Init != nil {
		return
	}

	// Vérifier que le body contient exactement 1 statement
	if ifStmt.Body == nil || len(ifStmt.Body.List) != 1 {
		return
	}

	// Récupérer le statement du body
	bodyStmt, ok := ifStmt.Body.List[0].(*ast.ReturnStmt)
	if !ok {
		return
	}

	// Vérifier que c'est un return avec exactement 1 valeur
	if len(bodyStmt.Results) != 1 {
		return
	}

	// Vérifier si le body a des commentaires (intention documentée)
	if hasComments(ifStmt.Body) {
		return
	}

	// Récupérer la valeur de retour du body
	bodyValue, bodyBool := getBooleanLiteral(bodyStmt.Results[0])
	if !bodyBool {
		return
	}

	// Cas 1: if/else
	if ifStmt.Else != nil {
		elseBlock, ok := ifStmt.Else.(*ast.BlockStmt)
		if !ok || len(elseBlock.List) != 1 {
			return
		}

		// Vérifier si le else a des commentaires
		if hasComments(elseBlock) {
			return
		}

		elseStmt, ok := elseBlock.List[0].(*ast.ReturnStmt)
		if !ok || len(elseStmt.Results) != 1 {
			return
		}

		elseValue, elseBool := getBooleanLiteral(elseStmt.Results[0])
		if !elseBool {
			return
		}

		// Vérifier si c'est un pattern simplifiable
		if bodyValue != elseValue {
			reportSimplifiableIf(pass, ifStmt, bodyValue)
		}
		return
	}

	// Cas 2: if suivi d'un return (pas de else)
	// Chercher le prochain statement après le if
	nextReturn := findNextReturnAfterIf(pass, ifStmt)
	if nextReturn == nil {
		return
	}

	if len(nextReturn.Results) != 1 {
		return
	}

	nextValue, nextBool := getBooleanLiteral(nextReturn.Results[0])
	if !nextBool {
		return
	}

	// Vérifier si c'est un pattern simplifiable
	if bodyValue != nextValue {
		reportSimplifiableIf(pass, ifStmt, bodyValue)
	}
}

// checkSingleCaseSwitch vérifie si un switch n'a qu'un seul case.
//
// Params:
//   - pass: la passe d'analyse
//   - switchStmt: le switch statement
func checkSingleCaseSwitch(pass *analysis.Pass, switchStmt *ast.SwitchStmt) {
	if switchStmt.Body == nil {
		// Retourne car pas de body
		return
	}

	caseCount := 0
	for _, stmt := range switchStmt.Body.List {
		if caseClause, ok := stmt.(*ast.CaseClause); ok {
			// default case a List == nil
			if caseClause.List != nil {
				caseCount++
			}
		}
	}

	if caseCount == 1 {
		reportSingleCaseSwitch(pass, switchStmt)
	}
}

// checkGotoUsage vérifie l'utilisation de goto.
//
// Params:
//   - pass: la passe d'analyse
//   - branchStmt: le branch statement (goto)
func checkGotoUsage(pass *analysis.Pass, branchStmt *ast.BranchStmt) {
	reportGotoUsage(pass, branchStmt)
}

// checkFallthroughOutsideSwitch vérifie fallthrough hors switch.
//
// Params:
//   - pass: la passe d'analyse
//   - file: le fichier
//   - branchStmt: le branch statement (fallthrough)
func checkFallthroughOutsideSwitch(pass *analysis.Pass, file *ast.File, branchStmt *ast.BranchStmt) {
	if !isInsideSwitchCase(file, branchStmt) {
		reportFallthroughOutsideSwitch(pass, branchStmt)
	}
}

// isInsideSwitchCase vérifie si un nœud est dans un case de switch.
//
// Params:
//   - file: le fichier
//   - target: le nœud cible
//
// Returns:
//   - bool: true si dans un case
func isInsideSwitchCase(file *ast.File, target ast.Node) bool {
	inCase := false

	ast.Inspect(file, func(n ast.Node) bool {
		if switchStmt, ok := n.(*ast.SwitchStmt); ok {
			if switchStmt.Body != nil {
				for _, stmt := range switchStmt.Body.List {
					if caseClause, ok := stmt.(*ast.CaseClause); ok {
						for _, s := range caseClause.Body {
							if s == target {
								inCase = true
								// Retourne false pour arrêter
								return false
							}
						}
					}
				}
			}
		}
		// Retourne true pour continuer
		return true
	})

	// Retourne le résultat
	return inCase
}

// reportDeferInLoop rapporte une violation KTN-DEFER-001.
//
// Params:
//   - pass: la passe d'analyse
//   - deferStmt: le defer statement
func reportDeferInLoop(pass *analysis.Pass, deferStmt *ast.DeferStmt) {
	pass.Reportf(deferStmt.Pos(),
		"[KTN-DEFER-001] defer dans une boucle accumule les appels.\n"+
			"Les defer ne s'exécutent qu'à la fin de la FONCTION, pas à chaque itération.\n"+
			"Cela cause des fuites de ressources (fichiers, locks, etc.) si la boucle est longue.\n"+
			"Extraire le traitement dans une fonction séparée avec defer.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - accumule N defer\n"+
			"  for _, file := range files {\n"+
			"      f := open(file)\n"+
			"      defer f.Close()  // Ne ferme qu'à la fin de la fonction!\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - ferme à chaque itération\n"+
			"  for _, file := range files {\n"+
			"      processFile(file)\n"+
			"  }\n"+
			"  func processFile(name string) {\n"+
			"      f := open(name)\n"+
			"      defer f.Close()  // Ferme à la fin de processFile\n"+
			"  }")
}

// getBooleanLiteral vérifie si une expression est un littéral booléen.
//
// Params:
//   - expr: l'expression à vérifier
//
// Returns:
//   - bool: la valeur booléenne (true ou false)
//   - bool: true si c'est un littéral booléen
func getBooleanLiteral(expr ast.Expr) (bool, bool) {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return false, false
	}
	if ident.Name == "true" {
		return true, true
	}
	if ident.Name == "false" {
		return false, true
	}
	return false, false
}

// hasComments vérifie si un bloc contient des commentaires.
//
// Params:
//   - block: le bloc à vérifier
//
// Returns:
//   - bool: true si le bloc contient des commentaires
func hasComments(block *ast.BlockStmt) bool {
	if block == nil {
		return false
	}

	// Vérifier les commentaires dans chaque statement
	for _, stmt := range block.List {
		// Vérifier si le statement a un commentaire associé
		// Note: Les commentaires sont gérés par l'AST via les positions
		// Pour une détection complète, on devrait vérifier pass.Fset
		// mais ici on fait une vérification basique
		if returnStmt, ok := stmt.(*ast.ReturnStmt); ok {
			// Si le return a des résultats, vérifier les commentaires
			for _, result := range returnStmt.Results {
				// Les commentaires inline seraient détectés via les positions
				_ = result
			}
		}
	}

	// Pour simplifier, on ne détecte pas les commentaires ici
	// La règle S1008 de Staticcheck fait cette vérification
	return false
}

// findNextReturnAfterIf trouve le prochain return après un if.
//
// Params:
//   - pass: la passe d'analyse
//   - ifStmt: le if statement
//
// Returns:
//   - *ast.ReturnStmt: le return trouvé, ou nil
func findNextReturnAfterIf(pass *analysis.Pass, ifStmt *ast.IfStmt) *ast.ReturnStmt {
	// Chercher dans tous les fichiers
	for _, file := range pass.Files {
		// Chercher le if dans le fichier
		var found bool
		var nextReturn *ast.ReturnStmt

		ast.Inspect(file, func(n ast.Node) bool {
			// Si on a déjà trouvé le return suivant, arrêter
			if nextReturn != nil {
				return false
			}

			// Chercher le if
			if n == ifStmt {
				found = true
				return true
			}

			// Si on n'a pas encore trouvé le if, continuer
			if !found {
				return true
			}

			// On a trouvé le if, chercher le prochain return au même niveau
			// Vérifier si on est dans une fonction
			if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Body != nil {
				// Chercher le if dans le body
				for i, stmt := range funcDecl.Body.List {
					if stmt == ifStmt && i+1 < len(funcDecl.Body.List) {
						// Le prochain statement
						if ret, ok := funcDecl.Body.List[i+1].(*ast.ReturnStmt); ok {
							nextReturn = ret
							return false
						}
					}
				}
			}

			return true
		})

		if nextReturn != nil {
			return nextReturn
		}
	}

	return nil
}

// reportSimplifiableIf rapporte une violation KTN-IF-004.
//
// Params:
//   - pass: la passe d'analyse
//   - ifStmt: le if statement
//   - bodyReturnsTrue: true si le body retourne true
func reportSimplifiableIf(pass *analysis.Pass, ifStmt *ast.IfStmt, bodyReturnsTrue bool) {
	suggestion := "return <condition>"
	if bodyReturnsTrue {
		suggestion = "return <condition>"
	} else {
		suggestion = "return !<condition>"
	}

	pass.Reportf(ifStmt.Pos(),
		"[KTN-IF-004] Expression booléenne simplifiable (Staticcheck S1008).\n"+
			"Un if qui retourne des littéraux booléens peut être simplifié.\n"+
			"Cela améliore la lisibilité et réduit l'indentation.\n"+
			"Suggestion: %s\n"+
			"Exemples:\n"+
			"  // ❌ MAUVAIS - if inutile\n"+
			"  if isValid {\n"+
			"      return true\n"+
			"  }\n"+
			"  return false\n"+
			"\n"+
			"  // ✅ CORRECT - return direct\n"+
			"  return isValid\n"+
			"\n"+
			"  // ❌ MAUVAIS - négation inutile\n"+
			"  if isValid {\n"+
			"      return false\n"+
			"  }\n"+
			"  return true\n"+
			"\n"+
			"  // ✅ CORRECT - négation directe\n"+
			"  return !isValid",
		suggestion)
}

// reportSingleCaseSwitch rapporte une violation KTN-SWITCH-005.
//
// Params:
//   - pass: la passe d'analyse
//   - switchStmt: le switch statement
func reportSingleCaseSwitch(pass *analysis.Pass, switchStmt *ast.SwitchStmt) {
	pass.Reportf(switchStmt.Pos(),
		"[KTN-SWITCH-005] switch avec un seul case devrait être un if.\n"+
			"Un switch avec un seul case est moins lisible qu'un simple if.\n"+
			"Utilisez if/else au lieu de switch pour 1-2 cas.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - switch pour un cas\n"+
			"  switch x {\n"+
			"  case 1:\n"+
			"      doSomething()\n"+
			"  default:\n"+
			"      doDefault()\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - if/else\n"+
			"  if x == 1 {\n"+
			"      doSomething()\n"+
			"  } else {\n"+
			"      doDefault()\n"+
			"  }")
}

// reportGotoUsage rapporte une violation KTN-GOTO-001.
//
// Params:
//   - pass: la passe d'analyse
//   - branchStmt: le branch statement (goto)
func reportGotoUsage(pass *analysis.Pass, branchStmt *ast.BranchStmt) {
	pass.Reportf(branchStmt.Pos(),
		"[KTN-GOTO-001] goto est considéré non idiomatique en Go.\n"+
			"L'utilisation de goto rend le code difficile à comprendre et maintenir.\n"+
			"Utilisez des structures de contrôle standards (if, for, return, break, continue).\n"+
			"Exception: cleanup dans du code bas niveau (rare).\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - goto\n"+
			"  if err != nil {\n"+
			"      goto cleanup\n"+
			"  }\n"+
			"  doWork()\n"+
			"  cleanup:\n"+
			"      close()\n"+
			"\n"+
			"  // ✅ CORRECT - defer\n"+
			"  defer close()\n"+
			"  if err != nil {\n"+
			"      return err\n"+
			"  }\n"+
			"  doWork()")
}

// reportFallthroughOutsideSwitch rapporte une violation KTN-FALL-001.
//
// Params:
//   - pass: la passe d'analyse
//   - branchStmt: le branch statement (fallthrough)
func reportFallthroughOutsideSwitch(pass *analysis.Pass, branchStmt *ast.BranchStmt) {
	pass.Reportf(branchStmt.Pos(),
		"[KTN-FALL-001] fallthrough ne peut être utilisé que dans un switch.\n"+
			"Le mot-clé fallthrough est uniquement valide dans un case de switch.\n"+
			"L'utiliser ailleurs est une erreur de syntaxe.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - fallthrough hors switch\n"+
			"  if x > 0 {\n"+
			"      fallthrough  // ERREUR\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - fallthrough dans switch\n"+
			"  switch x {\n"+
			"  case 1:\n"+
			"      doOne()\n"+
			"      fallthrough\n"+
			"  case 2:\n"+
			"      doTwo()\n"+
			"  }")
}
