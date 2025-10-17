package analyzer

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// BuiltinOpsAnalyzer vérifie les opérations built-in.
	BuiltinOpsAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktnbuiltinops",
		Doc:  "Vérifie les opérations built-in (append, new, division)",
		Run:  runBuiltinOpsAnalyzer,
	}
)

// runBuiltinOpsAnalyzer exécute l'analyseur builtin ops.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runBuiltinOpsAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.CallExpr:
				checkAppendWithoutAssign(pass, node)
				checkNewOnRefTypes(pass, node)
			case *ast.BinaryExpr:
				checkDivisionByZero(pass, node)
			}
			// Retourne true pour continuer
			return true
		})
	}
	// Retourne nil car terminé
	return nil, nil
}

// checkAppendWithoutAssign vérifie append sans assignment.
//
// Params:
//   - pass: la passe d'analyse
//   - call: l'appel de fonction
func checkAppendWithoutAssign(pass *analysis.Pass, call *ast.CallExpr) {
	// Vérifier si c'est un appel à append
	if !isBuiltinCall(call, "append") {
		// Retourne car pas append
		return
	}

	// Vérifier si le résultat est assigné
	if !isResultAssigned(pass, call) {
		reportAppendWithoutAssign(pass, call)
	}
}

// checkNewOnRefTypes vérifie new() sur types de référence.
//
// Params:
//   - pass: la passe d'analyse
//   - call: l'appel de fonction
func checkNewOnRefTypes(pass *analysis.Pass, call *ast.CallExpr) {
	if !isBuiltinCall(call, "new") {
		// Retourne car pas new
		return
	}

	if len(call.Args) == 0 {
		// Retourne car pas d'arguments
		return
	}

	// Vérifier le type
	arg := call.Args[0]
	if isRefType(arg) {
		reportNewOnRefType(pass, call)
	}
}

// checkDivisionByZero vérifie division par zéro.
//
// Params:
//   - pass: la passe d'analyse
//   - binary: l'expression binaire
func checkDivisionByZero(pass *analysis.Pass, binary *ast.BinaryExpr) {
	if binary.Op != token.QUO && binary.Op != token.REM {
		// Retourne car pas division/modulo
		return
	}

	// Vérifier si right est zéro
	if isZeroLiteral(binary.Y) {
		reportDivisionByZero(pass, binary)
	}
}

// isBuiltinCall vérifie si c'est un appel à une fonction built-in.
//
// Params:
//   - call: l'appel
//   - name: nom de la fonction
//
// Returns:
//   - bool: true si c'est la fonction
func isBuiltinCall(call *ast.CallExpr, name string) bool {
	ident, ok := call.Fun.(*ast.Ident)
	// Retourne true si c'est la bonne fonction
	return ok && ident.Name == name
}

// isResultAssigned vérifie si le résultat d'un appel est assigné.
//
// Params:
//   - pass: la passe d'analyse
//   - call: l'appel
//
// Returns:
//   - bool: true si assigné (mais pas à _)
func isResultAssigned(pass *analysis.Pass, call *ast.CallExpr) bool {
	for _, file := range pass.Files {
		found := false
		ast.Inspect(file, func(n ast.Node) bool {
			if assignStmt, ok := n.(*ast.AssignStmt); ok {
				for i, rhs := range assignStmt.Rhs {
					if rhs == call {
						// Vérifier si assigné à _
						if i < len(assignStmt.Lhs) {
							if ident, ok := assignStmt.Lhs[i].(*ast.Ident); ok {
								if ident.Name == "_" {
									// Assigné à _, c'est comme pas assigné
									return true
								}
							}
						}
						found = true
						// Retourne false pour arrêter
						return false
					}
				}
			}
			if valueSpec, ok := n.(*ast.ValueSpec); ok {
				for _, value := range valueSpec.Values {
					if value == call {
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
			// Retourne true car trouvé
			return true
		}
	}
	// Retourne false car pas assigné
	return false
}

// isRefType vérifie si c'est un type de référence.
//
// Params:
//   - expr: l'expression de type
//
// Returns:
//   - bool: true si slice/map/chan
func isRefType(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.ArrayType:
		// Slice si len == nil
		return t.Len == nil
	case *ast.MapType:
		// Retourne true pour map
		return true
	case *ast.ChanType:
		// Retourne true pour channel
		return true
	default:
		// Retourne false pour autres
		return false
	}
}

// reportAppendWithoutAssign rapporte une violation KTN-BUILTIN-APPEND-001.
//
// Params:
//   - pass: la passe d'analyse
//   - call: l'appel append
func reportAppendWithoutAssign(pass *analysis.Pass, call *ast.CallExpr) {
	pass.Reportf(call.Pos(),
		"[KTN-BUILTIN-APPEND-001] append sans assigner le résultat.\n"+
			"append ne modifie PAS le slice original, il retourne un nouveau slice.\n"+
			"Oublier d'assigner le résultat perd les données et cause des bugs.\n"+
			"Toujours assigner le résultat de append.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - résultat perdu\n"+
			"  append(slice, item)  // slice inchangé!\n"+
			"\n"+
			"  // ✅ CORRECT - assigner le résultat\n"+
			"  slice = append(slice, item)")
}

// reportNewOnRefType rapporte une violation KTN-BUILTIN-NEW-001.
//
// Params:
//   - pass: la passe d'analyse
//   - call: l'appel new
func reportNewOnRefType(pass *analysis.Pass, call *ast.CallExpr) {
	pass.Reportf(call.Pos(),
		"[KTN-BUILTIN-NEW-001] new() sur slice/map/channel.\n"+
			"new() alloue de la mémoire mais ne l'initialise pas correctement.\n"+
			"Pour slice/map/chan, utilisez make() qui initialise la structure interne.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - pointeur vers slice nil\n"+
			"  s := new([]int)  // *[]int avec slice nil\n"+
			"\n"+
			"  // ✅ CORRECT - slice initialisé\n"+
			"  s := make([]int, 0, 10)")
}

// reportDivisionByZero rapporte une violation KTN-OP-001.
//
// Params:
//   - pass: la passe d'analyse
//   - binary: l'expression binaire
func reportDivisionByZero(pass *analysis.Pass, binary *ast.BinaryExpr) {
	pass.Reportf(binary.Pos(),
		"[KTN-OP-001] Division ou modulo par zéro.\n"+
			"Division/modulo par zéro cause un panic immédiat en Go.\n"+
			"Vérifier que le diviseur n'est pas zéro avant l'opération.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - panic\n"+
			"  result := x / 0\n"+
			"\n"+
			"  // ✅ CORRECT - vérifier avant\n"+
			"  if divisor == 0 {\n"+
			"      return errors.New(\"division by zero\")\n"+
			"  }\n"+
			"  result := x / divisor")
}
