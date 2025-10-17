package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// ChannelOpsAnalyzer vérifie les opérations sur les channels.
	ChannelOpsAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktnchannelops",
		Doc:  "Vérifie les opérations sur les channels (close par receiver)",
		Run:  runChannelOpsAnalyzer,
	}
)

// runChannelOpsAnalyzer exécute l'analyseur d'opérations sur channels.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runChannelOpsAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				// Retourne true pour continuer
				return true
			}

			// Vérifier si c'est un appel à close()
			if isCloseCall(callExpr) {
				checkCloseInReceiver(pass, file, callExpr)
			}

			// Retourne true pour continuer l'inspection
			return true
		})
	}
	// Retourne nil car l'analyse est terminée
	return nil, nil
}

// isCloseCall vérifie si un appel est close().
//
// Params:
//   - call: l'appel de fonction
//
// Returns:
//   - bool: true si c'est close()
func isCloseCall(call *ast.CallExpr) bool {
	ident, ok := call.Fun.(*ast.Ident)
	return ok && ident.Name == "close" && len(call.Args) > 0
}

// checkCloseInReceiver vérifie si close() est appelé dans un receiver.
//
// Params:
//   - pass: la passe d'analyse
//   - file: le fichier
//   - closeCall: l'appel à close()
func checkCloseInReceiver(pass *analysis.Pass, file *ast.File, closeCall *ast.CallExpr) {
	// Trouver la fonction contenant le close()
	funcDecl := findEnclosingFunc(file, closeCall)
	if funcDecl == nil {
		// Pas dans une fonction
		// Retourne
		return
	}

	// Vérifier si la fonction reçoit d'un channel
	if functionReceivesFromChannel(funcDecl, closeCall.Args[0]) {
		reportCloseByReceiver(pass, closeCall)
	}
}

// findEnclosingFunc trouve la fonction englobante d'un nœud.
//
// Params:
//   - file: le fichier
//   - target: le nœud cible
//
// Returns:
//   - *ast.FuncDecl: la fonction, ou nil
func findEnclosingFunc(file *ast.File, target ast.Node) *ast.FuncDecl {
	var enclosing *ast.FuncDecl

	ast.Inspect(file, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			// Vérifier si target est dans cette fonction
			if funcDecl.Body != nil && containsNode(funcDecl.Body, target) {
				enclosing = funcDecl
				// Retourne false pour arrêter
				return false
			}
		}
		// Retourne true pour continuer
		return true
	})

	return enclosing
}

// functionReceivesFromChannel vérifie si une fonction reçoit d'un channel.
//
// Params:
//   - funcDecl: la déclaration de fonction
//   - chanExpr: l'expression du channel
//
// Returns:
//   - bool: true si la fonction reçoit du channel
func functionReceivesFromChannel(funcDecl *ast.FuncDecl, chanExpr ast.Expr) bool {
	if funcDecl.Body == nil {
		return false
	}

	// Extraire le nom du channel
	chanName := extractChannelName(chanExpr)
	if chanName == "" {
		return false
	}

	receives := false

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		// Chercher des opérations de réception <-chan
		unaryExpr, ok := n.(*ast.UnaryExpr)
		if !ok {
			// Retourne true pour continuer
			return true
		}

		if unaryExpr.Op.String() != "<-" {
			// Pas une réception
			// Retourne true pour continuer
			return true
		}

		// Vérifier si c'est une réception du même channel
		if receiveChanName := extractChannelName(unaryExpr.X); receiveChanName == chanName {
			receives = true
			// Retourne false pour arrêter
			return false
		}

		// Retourne true pour continuer
		return true
	})

	// Chercher aussi les range sur channel
	if !receives {
		ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
			rangeStmt, ok := n.(*ast.RangeStmt)
			if !ok {
				// Retourne true pour continuer
				return true
			}

			// Vérifier si c'est range sur le même channel
			if rangeChanName := extractChannelName(rangeStmt.X); rangeChanName == chanName {
				receives = true
				// Retourne false pour arrêter
				return false
			}

			// Retourne true pour continuer
			return true
		})
	}

	return receives
}

// extractChannelName extrait le nom d'un channel depuis une expression.
//
// Params:
//   - expr: l'expression
//
// Returns:
//   - string: nom du channel, ou ""
func extractChannelName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		// Pour chan.field
		if x, ok := e.X.(*ast.Ident); ok {
			return x.Name + "." + e.Sel.Name
		}
	}
	return ""
}

// reportCloseByReceiver rapporte une violation KTN-CHAN-002.
//
// Params:
//   - pass: la passe d'analyse
//   - closeCall: l'appel à close()
func reportCloseByReceiver(pass *analysis.Pass, closeCall *ast.CallExpr) {
	pass.Reportf(closeCall.Pos(),
		"[KTN-CHAN-002] close() appelé par le receiver du channel.\n"+
			"C'est une mauvaise pratique: seul le sender devrait fermer un channel.\n"+
			"Fermer côté receiver peut causer des panics si le sender écrit encore.\n"+
			"Pattern: sender ferme, receiver détecte la fermeture.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - receiver ferme\n"+
			"  func receive(ch chan int) {\n"+
			"      for v := range ch {\n"+
			"          process(v)\n"+
			"      }\n"+
			"      close(ch)  // 💥 Risque de panic si sender actif\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - sender ferme\n"+
			"  func send(ch chan int) {\n"+
			"      for i := 0; i < 10; i++ {\n"+
			"          ch <- i\n"+
			"      }\n"+
			"      close(ch)  // ✅ Sender ferme quand fini\n"+
			"  }\n"+
			"  func receive(ch chan int) {\n"+
			"      for v := range ch {  // ✅ Détecte automatiquement la fermeture\n"+
			"          process(v)\n"+
			"      }\n"+
			"  }")
}
