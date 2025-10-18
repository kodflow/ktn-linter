package ktn_control_flow

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// RuleDefer001 analyzer for defer statements.
var RuleDefer001 *analysis.Analyzer = &analysis.Analyzer{
	Name: "KTN_DEFER_001",
	Doc:  "Détecte les defer utilisés dans une boucle",
	Run:  runRuleDefer001,
}

func runRuleDefer001(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			deferStmt, ok := n.(*ast.DeferStmt)
			if !ok {
				// Continue traversing AST nodes.
				return true
			}

			if isInsideLoop(file, deferStmt) {
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
			// Continue traversing AST nodes.
			return true
		})
	}
	// Analysis completed successfully.
	return nil, nil
}

func isInsideLoop(file *ast.File, target ast.Node) bool {
	inLoop := false
	ast.Inspect(file, func(n ast.Node) bool {
		switch loop := n.(type) {
		case *ast.ForStmt:
			if loop.Body != nil && containsNode(loop.Body, target) {
				inLoop = true
				// Condition not met, return false.
				return false
			}
		case *ast.RangeStmt:
			if loop.Body != nil && containsNode(loop.Body, target) {
				inLoop = true
				// Condition not met, return false.
				return false
			}
		}
		// Continue traversing AST nodes.
		return true
	})
	// Early return from function.
	return inLoop
}

func containsNode(block *ast.BlockStmt, target ast.Node) bool {
	if block == nil {
		// Condition not met, return false.
		return false
	}
	found := false
	ast.Inspect(block, func(n ast.Node) bool {
		if n == target {
			found = true
			// Condition not met, return false.
			return false
		}
		// Continue traversing AST nodes.
		return true
	})
	// Early return from function.
	return found
}

// ContainsNodeExported est une version exportée pour les tests.
//
// Params:
//   - block: bloc d'instructions à analyser
//   - target: nœud AST à rechercher
//
// Returns:
//   - true si target est trouvé dans block
func ContainsNodeExported(block *ast.BlockStmt, target ast.Node) bool {
	// Early return from function.
	return containsNode(block, target)
}
