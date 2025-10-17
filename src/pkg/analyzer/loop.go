package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// LoopAnalyzer vérifie les patterns de boucles.
	LoopAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktnloop",
		Doc:  "Vérifie les patterns de boucles (range, for)",
		Run:  runLoopAnalyzer,
	}
)

// runLoopAnalyzer exécute l'analyseur de boucles.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runLoopAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch stmt := n.(type) {
			case *ast.RangeStmt:
				checkRangeKeyIgnored(pass, stmt)
				checkRangeVarCapture(pass, file, stmt)
			}
			// Retourne true pour continuer l'inspection
			return true
		})
	}
	// Retourne nil car l'analyse est terminée
	return nil, nil
}

// checkRangeKeyIgnored vérifie si la valeur est ignorée inutilement.
//
// En Go, `for i := range items` donne l'INDEX uniquement.
// Cette règle détecte `for i, _ := range items` (valeur ignorée avec underscore)
// et suggère `for i := range items` (plus propre, omettez la valeur).
//
// NE DÉTECTE PAS `for _, v := range items` qui est l'idiome correct pour itérer sur les valeurs.
//
// Params:
//   - pass: la passe d'analyse
//   - rangeStmt: le range statement
func checkRangeKeyIgnored(pass *analysis.Pass, rangeStmt *ast.RangeStmt) {
	// Vérifier si on a une valeur explicitement ignorée avec _
	if rangeStmt.Value == nil {
		// Pas de valeur déclarée, rien à vérifier
		return
	}

	valueIdent, ok := rangeStmt.Value.(*ast.Ident)
	if !ok {
		// Pas un identifiant simple pour la valeur
		return
	}

	// Détecter `for i, _ := range` ou `for _, _ := range`
	// (valeur ignorée avec underscore)
	if valueIdent.Name == "_" {
		// Vérifier si l'index est aussi ignoré (cas spécial)
		if rangeStmt.Key != nil {
			if keyIdent, ok := rangeStmt.Key.(*ast.Ident); ok && keyIdent.Name == "_" {
				// `for _, _ := range items` -> suggérer `for range items`
				reportBothIgnored(pass, rangeStmt)
				return
			}
		}

		// `for i, _ := range items` -> suggérer `for i := range items`
		reportValueIgnored(pass, rangeStmt)
	}
}

// checkRangeVarCapture vérifie la capture de variable range dans closure.
//
// Params:
//   - pass: la passe d'analyse
//   - file: le fichier
//   - rangeStmt: le range statement
func checkRangeVarCapture(pass *analysis.Pass, file *ast.File, rangeStmt *ast.RangeStmt) {
	// Rechercher les closures dans le body
	ast.Inspect(rangeStmt.Body, func(n ast.Node) bool {
		funcLit, ok := n.(*ast.FuncLit)
		if !ok {
			// Pas une fonction litérale
			// Retourne true pour continuer
			return true
		}

		// Vérifier si la closure capture des variables de range
		checkClosureCapturesRangeVar(pass, rangeStmt, funcLit)

		// Retourne true pour continuer l'inspection
		return true
	})
}

// checkClosureCapturesRangeVar vérifie si une closure capture une variable range.
//
// Params:
//   - pass: la passe d'analyse
//   - rangeStmt: le range statement
//   - funcLit: la fonction litérale (closure)
func checkClosureCapturesRangeVar(pass *analysis.Pass, rangeStmt *ast.RangeStmt, funcLit *ast.FuncLit) {
	// Récupérer les noms des variables de range
	rangeVars := getRangeVarNames(rangeStmt)
	if len(rangeVars) == 0 {
		// Pas de variables à capturer
		// Retourne
		return
	}

	// Vérifier si les variables ont été copiées localement (varName := varName)
	copiedVars := make(map[string]bool)
	ast.Inspect(rangeStmt.Body, func(n ast.Node) bool {
		// Arrêter avant d'entrer dans la closure
		if n == funcLit {
			return false
		}

		assignStmt, ok := n.(*ast.AssignStmt)
		if !ok || assignStmt.Tok.String() != ":=" {
			return true
		}

		// Vérifier le pattern varName := varName
		for i, lhs := range assignStmt.Lhs {
			if i >= len(assignStmt.Rhs) {
				continue
			}
			lhsIdent, ok1 := lhs.(*ast.Ident)
			rhsIdent, ok2 := assignStmt.Rhs[i].(*ast.Ident)
			if ok1 && ok2 && lhsIdent.Name == rhsIdent.Name {
				// C'est une copie locale
				copiedVars[lhsIdent.Name] = true
			}
		}
		return true
	})

	// Vérifier si la closure utilise des variables non copiées
	usesRangeVar := false
	capturedVar := ""

	ast.Inspect(funcLit.Body, func(n ast.Node) bool {
		ident, ok := n.(*ast.Ident)
		if !ok {
			// Pas un identifiant
			// Retourne true pour continuer
			return true
		}

		// Vérifier si c'est une variable de range non copiée
		for _, varName := range rangeVars {
			if ident.Name == varName && varName != "_" && !copiedVars[varName] {
				usesRangeVar = true
				capturedVar = varName
				// Retourne false pour arrêter
				return false
			}
		}

		// Retourne true pour continuer
		return true
	})

	if usesRangeVar {
		reportRangeVarCapture(pass, funcLit, capturedVar)
	}
}

// getRangeVarNames extrait les noms des variables de range.
//
// Params:
//   - rangeStmt: le range statement
//
// Returns:
//   - []string: les noms des variables
func getRangeVarNames(rangeStmt *ast.RangeStmt) []string {
	var names []string

	if rangeStmt.Key != nil {
		if keyIdent, ok := rangeStmt.Key.(*ast.Ident); ok {
			names = append(names, keyIdent.Name)
		}
	}

	if rangeStmt.Value != nil {
		if valueIdent, ok := rangeStmt.Value.(*ast.Ident); ok {
			names = append(names, valueIdent.Name)
		}
	}

	// Retourne les noms des variables
	return names
}

// reportValueIgnored rapporte une violation KTN-FOR-001 (valeur ignorée).
//
// Params:
//   - pass: la passe d'analyse
//   - rangeStmt: le range statement
func reportValueIgnored(pass *analysis.Pass, rangeStmt *ast.RangeStmt) {
	pass.Reportf(rangeStmt.Pos(),
		"[KTN-FOR-001] Valeur de range ignorée inutilement avec _.\n"+
			"Si vous n'utilisez que l'index, omettez la valeur.\n"+
			"Utilisez `for i := range items` au lieu de `for i, _ := range items`.\n"+
			"Cela rend le code plus lisible et idiomatique.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - _ inutile\n"+
			"  for i, _ := range items {\n"+
			"      use(i)\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - omettez la valeur\n"+
			"  for i := range items {\n"+
			"      use(i)\n"+
			"  }\n"+
			"\n"+
			"Note: `for _, v := range items` est CORRECT pour itérer sur les valeurs.")
}

// reportBothIgnored rapporte une violation KTN-FOR-001 (index et valeur ignorés).
//
// Params:
//   - pass: la passe d'analyse
//   - rangeStmt: le range statement
func reportBothIgnored(pass *analysis.Pass, rangeStmt *ast.RangeStmt) {
	pass.Reportf(rangeStmt.Pos(),
		"[KTN-FOR-001] Index et valeur de range ignorés inutilement avec _.\n"+
			"Si vous n'utilisez ni l'index ni la valeur, omettez-les complètement.\n"+
			"Utilisez `for range items` au lieu de `for _, _ := range items`.\n"+
			"Cela rend le code plus lisible et idiomatique.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - _ _ inutiles\n"+
			"  for _, _ := range items {\n"+
			"      doSomething()\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - omettez tout\n"+
			"  for range items {\n"+
			"      doSomething()\n"+
			"  }")
}

// reportRangeVarCapture rapporte une violation KTN-RANGE-003.
//
// Params:
//   - pass: la passe d'analyse
//   - funcLit: la fonction litérale
//   - varName: nom de la variable capturée
func reportRangeVarCapture(pass *analysis.Pass, funcLit *ast.FuncLit, varName string) {
	pass.Reportf(funcLit.Pos(),
		"[KTN-RANGE-003] Variable de range '%s' capturée dans une closure.\n"+
			"Les variables de range sont réutilisées à chaque itération.\n"+
			"Capturer directement dans une goroutine/closure cause des bugs.\n"+
			"Créez une copie locale avant la closure.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - capture la même variable\n"+
			"  for _, v := range items {\n"+
			"      go func() {\n"+
			"          process(v)  // 💥 BUG: toutes les goroutines voient la dernière valeur\n"+
			"      }()\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - copie locale\n"+
			"  for _, v := range items {\n"+
			"      v := v  // Copie locale\n"+
			"      go func() {\n"+
			"          process(v)  // ✅ Chaque goroutine a sa propre copie\n"+
			"      }()\n"+
			"  }\n"+
			"\n"+
			"  // ✅ CORRECT - passage par paramètre\n"+
			"  for _, v := range items {\n"+
			"      go func(item Item) {\n"+
			"          process(item)  // ✅ Valeur passée explicitement\n"+
			"      }(v)\n"+
			"  }",
		varName)
}
