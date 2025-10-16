package analyzer

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// AllocAnalyzer vérifie les règles d'allocation mémoire (make vs new).
	AllocAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktnalloc",
		Doc:  "Vérifie l'utilisation correcte de make() et new() pour l'allocation mémoire",
		Run:  runAllocAnalyzer,
	}
)

// runAllocAnalyzer exécute l'analyseur d'allocation mémoire.
//
// Params:
//   - pass: la passe d'analyse contenant les fichiers à vérifier
//
// Returns:
//   - interface{}: toujours nil car aucun résultat n'est nécessaire
//   - error: toujours nil, les erreurs sont rapportées via pass.Reportf
func runAllocAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		analyzeMakeAppendInFile(pass, file)
		analyzeAllocCallsInFile(pass, file)
	}
	// Retourne nil car l'analyseur rapporte via pass.Reportf
	return nil, nil
}

// analyzeMakeAppendInFile analyse les patterns make + append dans un fichier.
//
// Params:
//   - pass: la passe d'analyse
//   - file: le fichier à analyser
func analyzeMakeAppendInFile(pass *analysis.Pass, file *ast.File) {
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		checkMakeAppendPattern(pass, funcDecl)
	}
}

// analyzeAllocCallsInFile analyse les appels new() et make() dans un fichier.
//
// Params:
//   - pass: la passe d'analyse
//   - file: le fichier à analyser
func analyzeAllocCallsInFile(pass *analysis.Pass, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			// Retourne true pour continuer l'inspection
			return true
		}

		ident, ok := callExpr.Fun.(*ast.Ident)
		if !ok {
			// Retourne true pour continuer l'inspection
			return true
		}

		switch ident.Name {
		case "new":
			checkNewUsage(pass, callExpr)
		case "make":
			checkMakeUsage(pass, callExpr)
		}

		// Retourne true pour continuer l'inspection
		return true
	})
}

// checkNewUsage vérifie l'utilisation de new().
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - callExpr: l'expression d'appel à new()
func checkNewUsage(pass *analysis.Pass, callExpr *ast.CallExpr) {
	if len(callExpr.Args) != 1 {
		// Retourne car new() doit avoir exactement 1 argument
		return
	}

	arg := callExpr.Args[0]

	// KTN-ALLOC-001 : Interdire new() avec slice/map/chan
	if isReferenceType(arg) {
		reportReferenceTypeViolation(pass, callExpr, arg)
		// Retourne car la violation a été rapportée
		return
	}

	// KTN-ALLOC-004 : Préférer &struct{} à new(struct)
	if isStructType(arg) {
		reportStructTypeViolation(pass, callExpr, arg)
	}
}

// reportReferenceTypeViolation rapporte une violation KTN-ALLOC-001.
//
// Params:
//   - pass: la passe d'analyse
//   - callExpr: l'appel à new()
//   - arg: l'argument de type référence
func reportReferenceTypeViolation(pass *analysis.Pass, callExpr *ast.CallExpr, arg ast.Expr) {
	typeName := getTypeName(arg)
	pass.Reportf(callExpr.Pos(),
		"[KTN-ALLOC-001] Utilisation de new() avec un type référence (%s) interdite.\n"+
			"new() retourne un pointeur vers nil pour les types référence, ce qui cause des panics.\n"+
			"Utilisez make() à la place.\n"+
			"Exemple:\n"+
			"  // ❌ INTERDIT\n"+
			"  m := new(%s)  // *%s avec nil map/slice/chan\n"+
			"  (*m)[\"key\"] = value  // 💥 PANIC\n"+
			"\n"+
			"  // ✅ CORRECT\n"+
			"  m := make(%s)\n"+
			"  m[\"key\"] = value  // ✅ Fonctionne",
		typeName, typeName, typeName, typeName)
}

// reportStructTypeViolation rapporte une violation KTN-ALLOC-004.
//
// Params:
//   - pass: la passe d'analyse
//   - callExpr: l'appel à new()
//   - arg: l'argument de type struct
func reportStructTypeViolation(pass *analysis.Pass, callExpr *ast.CallExpr, arg ast.Expr) {
	typeName := getTypeName(arg)
	pass.Reportf(callExpr.Pos(),
		"[KTN-ALLOC-004] Utilisez le composite literal &%s{} au lieu de new(%s).\n"+
			"En Go idiomatique, on préfère les composite literals pour les structs.\n"+
			"Exemple:\n"+
			"  // ❌ NON-IDIOMATIQUE\n"+
			"  p := new(%s)\n"+
			"\n"+
			"  // ✅ IDIOMATIQUE GO\n"+
			"  p := &%s{}",
		typeName, typeName, typeName, typeName)
}

// checkMakeUsage vérifie l'utilisation de make() pour les slices.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - callExpr: l'expression d'appel à make()
func checkMakeUsage(pass *analysis.Pass, callExpr *ast.CallExpr) {
	if len(callExpr.Args) < 1 {
		// Retourne car make() doit avoir au moins 1 argument
		return
	}

	typeArg := callExpr.Args[0]

	// Vérifier si c'est un slice
	if !isSliceType(typeArg) {
		// Retourne car on ne vérifie que les slices
		return
	}

	// Vérifier si c'est make([]T, 0) ou make([]T, 0, 0)
	if len(callExpr.Args) == 2 {
		// make([]T, length)
		if isZeroLiteral(callExpr.Args[1]) {
			// C'est make([]T, 0) - potentiellement problématique
			// On marque cet appel pour vérification ultérieure dans checkMakeAppendPattern
		}
	} else if len(callExpr.Args) == 3 {
		// make([]T, length, capacity)
		if isZeroLiteral(callExpr.Args[1]) && isZeroLiteral(callExpr.Args[2]) {
			// C'est make([]T, 0, 0) - équivalent à make([]T, 0)
		}
	}
}

// checkMakeAppendPattern vérifie le pattern make(slice, 0) suivi d'append.
//
// Params:
//   - pass: la passe d'analyse pour rapporter l'erreur
//   - funcDecl: la déclaration de fonction à analyser
func checkMakeAppendPattern(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	makeSliceVars := make(map[string]token.Pos)

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		trackMakeSliceZeroAssignment(n, makeSliceVars)
		checkAppendOnTrackedSlice(pass, n, makeSliceVars)
		// Retourne true pour continuer l'inspection
		return true
	})
}

// trackMakeSliceZeroAssignment détecte et enregistre les assignations make([]T, 0).
//
// Params:
//   - n: le nœud AST à analyser
//   - makeSliceVars: map des variables trackées avec leur position
func trackMakeSliceZeroAssignment(n ast.Node, makeSliceVars map[string]token.Pos) {
	assignStmt, ok := n.(*ast.AssignStmt)
	if !ok || len(assignStmt.Lhs) != 1 || len(assignStmt.Rhs) != 1 {
		// Retourne car ce n'est pas une assignation simple
		return
	}

	if !isMakeSliceZero(assignStmt.Rhs[0]) {
		// Retourne car ce n'est pas make([]T, 0)
		return
	}

	if ident, ok := assignStmt.Lhs[0].(*ast.Ident); ok {
		makeSliceVars[ident.Name] = assignStmt.Pos()
	}
}

// checkAppendOnTrackedSlice vérifie si append() est appelé sur un slice tracké.
//
// Params:
//   - pass: la passe d'analyse
//   - n: le nœud AST à analyser
//   - makeSliceVars: map des variables trackées
func checkAppendOnTrackedSlice(pass *analysis.Pass, n ast.Node, makeSliceVars map[string]token.Pos) {
	callExpr, ok := n.(*ast.CallExpr)
	if !ok {
		// Retourne car ce n'est pas un appel de fonction
		return
	}

	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok || ident.Name != "append" {
		// Retourne car ce n'est pas append()
		return
	}

	if len(callExpr.Args) == 0 {
		// Retourne car append() sans arguments
		return
	}

	firstArg, ok := callExpr.Args[0].(*ast.Ident)
	if !ok {
		// Retourne car le premier argument n'est pas un identifiant
		return
	}

	makePos, found := makeSliceVars[firstArg.Name]
	if !found {
		// Retourne car la variable n'est pas trackée
		return
	}

	reportMakeAppendViolation(pass, makePos, firstArg.Name)
	delete(makeSliceVars, firstArg.Name)
}

// reportMakeAppendViolation rapporte une violation KTN-ALLOC-002.
//
// Params:
//   - pass: la passe d'analyse
//   - pos: position du make()
//   - sliceName: nom du slice
func reportMakeAppendViolation(pass *analysis.Pass, pos token.Pos, sliceName string) {
	pass.Reportf(pos,
		"[KTN-ALLOC-002] Slice '%s' créé avec make([]T, 0) puis utilisé avec append().\n"+
			"Cela force des réallocations coûteuses à chaque append.\n"+
			"Si la taille est connue, spécifiez la capacité.\n"+
			"Exemple:\n"+
			"  // ❌ INEFFICACE\n"+
			"  items := make([]Item, 0)\n"+
			"  for _, v := range source {\n"+
			"      items = append(items, v)  // Réallocation O(log n)\n"+
			"  }\n"+
			"\n"+
			"  // ✅ OPTIMISÉ\n"+
			"  items := make([]Item, 0, len(source))  // Préallocation\n"+
			"  for _, v := range source {\n"+
			"      items = append(items, v)  // Pas de réallocation\n"+
			"  }",
		sliceName)
}

// isReferenceType vérifie si un type est un type référence (slice/map/chan).
//
// Params:
//   - expr: l'expression représentant le type
//
// Returns:
//   - bool: true si c'est un slice, map ou channel
func isReferenceType(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.ArrayType:
		// Retourne true si c'est un slice (ArrayType sans longueur)
		return t.Len == nil
	case *ast.MapType:
		// Retourne true car c'est une map
		return true
	case *ast.ChanType:
		// Retourne true car c'est un channel
		return true
	case *ast.Ident:
		// Vérifier les types natifs
		// Retourne true si c'est un alias pour slice/map/chan
		return strings.Contains(t.Name, "map") ||
			strings.Contains(t.Name, "chan") ||
			strings.Contains(t.Name, "slice")
	}
	// Retourne false car ce n'est pas un type référence
	return false
}

// isStructType vérifie si un type est une struct.
//
// Params:
//   - expr: l'expression représentant le type
//
// Returns:
//   - bool: true si c'est une struct
func isStructType(expr ast.Expr) bool {
	switch expr.(type) {
	case *ast.StructType:
		// Retourne true car c'est directement une struct
		return true
	case *ast.Ident:
		// Si c'est un identifiant, on suppose que c'est potentiellement une struct
		// (on ne peut pas vérifier le type exact sans type checker)
		// Retourne true pour les identifiants (types nommés)
		return true
	case *ast.SelectorExpr:
		// Type importé (ex: pkg.MyStruct)
		// Retourne true car c'est probablement une struct
		return true
	}
	// Retourne false car ce n'est clairement pas une struct
	return false
}

// isSliceType vérifie si un type est un slice.
//
// Params:
//   - expr: l'expression représentant le type
//
// Returns:
//   - bool: true si c'est un slice
func isSliceType(expr ast.Expr) bool {
	arrayType, ok := expr.(*ast.ArrayType)
	// Retourne true si c'est un ArrayType sans longueur (slice)
	return ok && arrayType.Len == nil
}

// isZeroLiteral vérifie si une expression est le littéral 0.
//
// Params:
//   - expr: l'expression à vérifier
//
// Returns:
//   - bool: true si c'est le littéral 0
func isZeroLiteral(expr ast.Expr) bool {
	basicLit, ok := expr.(*ast.BasicLit)
	// Retourne true si c'est le littéral "0"
	return ok && basicLit.Kind == token.INT && basicLit.Value == "0"
}

// isMakeSliceZero vérifie si une expression est make([]T, 0) ou make([]T, 0, 0).
//
// Params:
//   - expr: l'expression à vérifier
//
// Returns:
//   - bool: true si c'est make([]T, 0) ou make([]T, 0, 0)
func isMakeSliceZero(expr ast.Expr) bool {
	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		// Retourne false car ce n'est pas un appel
		return false
	}

	// Vérifier si c'est make()
	ident, ok := callExpr.Fun.(*ast.Ident)
	if !ok || ident.Name != "make" {
		// Retourne false car ce n'est pas make()
		return false
	}

	// Vérifier si c'est un slice
	if len(callExpr.Args) < 1 || !isSliceType(callExpr.Args[0]) {
		// Retourne false car ce n'est pas un slice
		return false
	}

	// Vérifier les arguments
	if len(callExpr.Args) == 2 {
		// make([]T, length)
		// Retourne true si length est 0
		return isZeroLiteral(callExpr.Args[1])
	} else if len(callExpr.Args) == 3 {
		// make([]T, length, capacity)
		// Retourne true si length et capacity sont 0
		return isZeroLiteral(callExpr.Args[1]) && isZeroLiteral(callExpr.Args[2])
	}

	// Retourne false car les arguments ne correspondent pas au pattern
	return false
}

// getTypeName extrait le nom d'un type depuis une expression.
//
// Params:
//   - expr: l'expression représentant le type
//
// Returns:
//   - string: le nom du type (ex: "map[string]int", "[]int", "chan int")
func getTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.ArrayType:
		// Slice
		elemType := getTypeName(t.Elt)
		// Retourne le nom du slice
		return "[]" + elemType
	case *ast.MapType:
		keyType := getTypeName(t.Key)
		valueType := getTypeName(t.Value)
		// Retourne le nom de la map
		return "map[" + keyType + "]" + valueType
	case *ast.ChanType:
		elemType := getTypeName(t.Value)
		// Retourne le nom du channel
		return "chan " + elemType
	case *ast.Ident:
		// Retourne le nom de l'identifiant
		return t.Name
	case *ast.SelectorExpr:
		pkg := getTypeName(t.X)
		// Retourne le nom qualifié (pkg.Type)
		return pkg + "." + t.Sel.Name
	case *ast.StarExpr:
		base := getTypeName(t.X)
		// Retourne le nom du pointeur
		return "*" + base
	}
	// Retourne un nom générique si le type est inconnu
	return "T"
}
