package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// TypeOpsAnalyzer vérifie les opérations sur les types.
	TypeOpsAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktntypeops",
		Doc:  "Vérifie les opérations sur les types (assertions, conversions, pointeurs)",
		Run:  runTypeOpsAnalyzer,
	}
)

// runTypeOpsAnalyzer exécute l'analyseur d'opérations sur types.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runTypeOpsAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.TypeAssertExpr:
				checkTypeAssertion(pass, node)
			case *ast.CallExpr:
				checkTypeConversion(pass, node)
			case *ast.UnaryExpr:
				checkPointerDereference(pass, node)
			}
			// Retourne true pour continuer l'inspection
			return true
		})
	}
	// Retourne nil car l'analyse est terminée
	return nil, nil
}

// checkTypeAssertion vérifie les assertions de type sans vérification.
//
// Params:
//   - pass: la passe d'analyse
//   - typeAssert: l'assertion de type
func checkTypeAssertion(pass *analysis.Pass, typeAssert *ast.TypeAssertExpr) {
	// Exclure les type switches (ont Type == nil)
	if typeAssert.Type == nil {
		// C'est un type switch, c'est sûr
		return
	}

	// Vérifier si c'est dans un contexte d'assignation
	if !isInAssignment(pass, typeAssert) {
		// Pas dans une assignation, potentiellement dangereux
		reportUncheckedTypeAssertion(pass, typeAssert)
	}
}

// isInAssignment vérifie si une expression est dans une assignation.
//
// Params:
//   - pass: la passe d'analyse
//   - expr: l'expression
//
// Returns:
//   - bool: true si dans une assignation avec 2 valeurs de retour
func isInAssignment(pass *analysis.Pass, expr ast.Expr) bool {
	for _, file := range pass.Files {
		found := false
		ast.Inspect(file, func(n ast.Node) bool {
			assignStmt, ok := n.(*ast.AssignStmt)
			if !ok {
				// Pas une assignation
				// Retourne true pour continuer
				return true
			}

			// Vérifier si expr est dans Rhs
			for _, rhs := range assignStmt.Rhs {
				if rhs == expr {
					// Trouvé! Vérifier si c'est une assignation à 2 valeurs
					if len(assignStmt.Lhs) == 2 {
						found = true
						// Retourne false pour arrêter
						return false
					}
				}
			}

			// Retourne true pour continuer
			return true
		})

		if found {
			// Retourne true car trouvé dans assignation sûre
			return true
		}
	}

	// Retourne false car pas dans assignation sûre
	return false
}

// checkTypeConversion vérifie les conversions de type redondantes.
//
// Params:
//   - pass: la passe d'analyse
//   - call: l'appel (potentiellement une conversion)
func checkTypeConversion(pass *analysis.Pass, call *ast.CallExpr) {
	// Vérifier si c'est une conversion de type (pas un appel de fonction)
	if len(call.Args) != 1 {
		// Pas une conversion simple
		// Retourne
		return
	}

	// Vérifier si Fun est un type
	typeIdent, ok := call.Fun.(*ast.Ident)
	if !ok {
		// Pas un identifiant simple
		// Retourne
		return
	}

	// Vérifier si l'argument a potentiellement le même type
	arg := call.Args[0]
	if argIdent, ok := arg.(*ast.Ident); ok {
		// Si c'est T(x) où x est déjà de type T, c'est redondant
		// Note: on devrait utiliser pass.TypesInfo pour être sûr,
		// mais on fait une détection simple ici
		if typeIdent.Name == argIdent.Name {
			reportRedundantConversion(pass, call, typeIdent.Name)
		}
	}
}

// checkPointerDereference vérifie les déréférencements de pointeur potentiellement nil.
//
// Params:
//   - pass: la passe d'analyse
//   - unary: l'expression unaire
func checkPointerDereference(pass *analysis.Pass, unary *ast.UnaryExpr) {
	// On ne vérifie que les déréférencements *
	if unary.Op.String() != "*" {
		// Pas un déréférencement
		// Retourne
		return
	}

	// Vérifier si c'est un déréférencement direct d'un appel à new
	if call, ok := unary.X.(*ast.CallExpr); ok {
		if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "new" {
			// C'est *new(...) ce qui est sûr
			// Retourne
			return
		}
	}

	// Pour une détection plus robuste, on devrait tracker les nil assignments
	// Ici on fait une détection basique: si X est un identifiant récemment assigné à nil
	if ident, ok := unary.X.(*ast.Ident); ok {
		if isRecentlyNil(pass, ident, unary) {
			reportNilDereference(pass, unary, ident.Name)
		}
	}
}

// isRecentlyNil vérifie si une variable a été assignée à nil récemment.
//
// Params:
//   - pass: la passe d'analyse
//   - ident: l'identifiant
//   - deref: le déréférencement
//
// Returns:
//   - bool: true si potentiellement nil
func isRecentlyNil(pass *analysis.Pass, ident *ast.Ident, deref ast.Node) bool {
	// Rechercher dans le fichier les assignations à nil
	for _, file := range pass.Files {
		nilAssigned := false

		ast.Inspect(file, func(n ast.Node) bool {
			// Si on atteint le déréférencement, on arrête
			if n == deref {
				// Retourne false pour arrêter
				return false
			}

			// Chercher les assignations
			assignStmt, ok := n.(*ast.AssignStmt)
			if !ok {
				// Retourne true pour continuer
				return true
			}

			// Vérifier si c'est une assignation à notre variable
			for i, lhs := range assignStmt.Lhs {
				if lhsIdent, ok := lhs.(*ast.Ident); ok && lhsIdent.Name == ident.Name {
					// Vérifier si Rhs est nil
					if i < len(assignStmt.Rhs) {
						if rhsIdent, ok := assignStmt.Rhs[i].(*ast.Ident); ok && rhsIdent.Name == "nil" {
							nilAssigned = true
						}
					}
				}
			}

			// Retourne true pour continuer
			return true
		})

		if nilAssigned {
			// Retourne true car assigné à nil
			return true
		}
	}

	// Retourne false car pas de preuve que c'est nil
	return false
}

// reportUncheckedTypeAssertion rapporte une violation KTN-ASSERT-001.
//
// Params:
//   - pass: la passe d'analyse
//   - typeAssert: l'assertion de type
func reportUncheckedTypeAssertion(pass *analysis.Pass, typeAssert *ast.TypeAssertExpr) {
	pass.Reportf(typeAssert.Pos(),
		"[KTN-ASSERT-001] Assertion de type sans vérification du succès.\n"+
			"Une assertion de type non vérifiée cause un panic si le type est incorrect.\n"+
			"Utilisez toujours la forme à deux valeurs (value, ok) pour vérifier.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - panic si x n'est pas un int\n"+
			"  v := x.(int)\n"+
			"\n"+
			"  // ✅ CORRECT - vérification du type\n"+
			"  v, ok := x.(int)\n"+
			"  if !ok {\n"+
			"      // Gérer l'erreur\n"+
			"      return errors.New(\"wrong type\")\n"+
			"  }")
}

// reportRedundantConversion rapporte une violation KTN-CONV-002.
//
// Params:
//   - pass: la passe d'analyse
//   - call: l'appel de conversion
//   - typeName: nom du type
func reportRedundantConversion(pass *analysis.Pass, call *ast.CallExpr, typeName string) {
	pass.Reportf(call.Pos(),
		"[KTN-CONV-002] Conversion de type redondante: %s(%s).\n"+
			"Cette conversion est inutile car la variable a déjà le bon type.\n"+
			"Les conversions redondantes nuisent à la lisibilité.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - conversion inutile\n"+
			"  var x int = 5\n"+
			"  y := int(x)  // x est déjà int!\n"+
			"\n"+
			"  // ✅ CORRECT - pas de conversion\n"+
			"  var x int = 5\n"+
			"  y := x",
		typeName, typeName)
}

// reportNilDereference rapporte une violation KTN-POINTER-001.
//
// Params:
//   - pass: la passe d'analyse
//   - unary: l'expression unaire
//   - varName: nom de la variable
func reportNilDereference(pass *analysis.Pass, unary *ast.UnaryExpr, varName string) {
	pass.Reportf(unary.Pos(),
		"[KTN-POINTER-001] Déréférencement potentiel d'un pointeur nil '%s'.\n"+
			"Déréférencer un pointeur nil cause un panic immédiat.\n"+
			"Vérifiez toujours qu'un pointeur n'est pas nil avant de le déréférencer.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - déréférence sans vérification\n"+
			"  var p *int\n"+
			"  x := *p  // 💥 PANIC: nil pointer dereference\n"+
			"\n"+
			"  // ✅ CORRECT - vérifier avant\n"+
			"  var p *int\n"+
			"  if p != nil {\n"+
			"      x := *p\n"+
			"  }",
		varName)
}
