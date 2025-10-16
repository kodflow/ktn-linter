package analyzer

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// PoolAnalyzer vérifie l'utilisation correcte de sync.Pool et la taille des structs.
	PoolAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktnpool",
		Doc:  "Vérifie l'utilisation de sync.Pool (defer Put) et détecte les grandes structs passées par valeur",
		Run:  runPoolAnalyzer,
	}
)

// Configuration thresholds
//
// Ces constantes définissent les seuils pour l'analyse des pools et structs.
const (
	// largeStructThreshold est le seuil en bytes au-delà duquel une struct devrait être passée par pointeur.
	largeStructThreshold int64 = 128
)

// runPoolAnalyzer exécute l'analyseur pool.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runPoolAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		checkSyncPoolUsage(pass, file)
		checkLargeStructByValue(pass, file)
	}
	// Retourne nil car l'analyse est terminée
	return nil, nil
}

// checkSyncPoolUsage vérifie l'utilisation de sync.Pool dans le fichier.
//
// Params:
//   - pass: la passe d'analyse
//   - file: le fichier à analyser
func checkSyncPoolUsage(pass *analysis.Pass, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			// Retourne true pour continuer l'inspection
			return true
		}

		checkPoolGetWithoutDeferPut(pass, funcDecl)
		// Retourne true pour continuer l'inspection
		return true
	})
}

// checkPoolGetWithoutDeferPut vérifie si Get() est utilisé sans defer Put().
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction à analyser
func checkPoolGetWithoutDeferPut(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	// Map pour tracker les variables issues de pool.Get()
	poolVars := make(map[string]ast.Expr) // varName -> pool expression
	deferredPuts := make(map[string]bool) // varName -> has defer Put

	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			trackPoolGetAssignment(stmt, poolVars, pass)
		case *ast.DeferStmt:
			trackDeferPut(stmt, deferredPuts)
		}
		// Retourne true pour continuer l'inspection
		return true
	})

	// Vérifier les variables sans defer Put()
	for varName, poolExpr := range poolVars {
		if !deferredPuts[varName] {
			reportMissingDeferPut(pass, poolExpr, varName)
		}
	}
}

// trackPoolGetAssignment détecte les assignations depuis pool.Get().
//
// Params:
//   - stmt: l'assignation à analyser
//   - poolVars: map des variables pool trackées
//   - pass: la passe d'analyse
func trackPoolGetAssignment(stmt *ast.AssignStmt, poolVars map[string]ast.Expr, pass *analysis.Pass) {
	if len(stmt.Lhs) != 1 || len(stmt.Rhs) != 1 {
		// Retourne car ce n'est pas une assignation simple
		return
	}

	// Extraire l'expression sous-jacente (peut être un type assertion)
	rhsExpr := unwrapTypeAssertion(stmt.Rhs[0])

	// Vérifier si RHS est un appel à pool.Get()
	if !isPoolGetCall(rhsExpr, pass) {
		// Retourne car ce n'est pas pool.Get()
		return
	}

	// Extraire le nom de la variable
	varName := extractVarName(stmt.Lhs[0])
	if varName != "" {
		poolVars[varName] = rhsExpr
	}
}

// unwrapTypeAssertion extrait l'expression d'un type assertion.
//
// Params:
//   - expr: l'expression (peut être TypeAssertExpr)
//
// Returns:
//   - ast.Expr: l'expression sous-jacente
func unwrapTypeAssertion(expr ast.Expr) ast.Expr {
	if typeAssert, ok := expr.(*ast.TypeAssertExpr); ok {
		// Retourne l'expression sous le type assertion
		return typeAssert.X
	}
	// Retourne l'expression telle quelle
	return expr
}

// trackDeferPut détecte les defer avec pool.Put().
//
// Params:
//   - stmt: le defer statement
//   - deferredPuts: map des Put() différés
func trackDeferPut(stmt *ast.DeferStmt, deferredPuts map[string]bool) {
	callExpr := stmt.Call
	if callExpr == nil {
		// Retourne car pas d'appel
		return
	}

	// Vérifier si c'est pool.Put(var)
	if !isPoolPutCall(callExpr) {
		// Retourne car ce n'est pas pool.Put()
		return
	}

	// Extraire le nom de la variable passée à Put()
	if len(callExpr.Args) > 0 {
		varName := extractVarName(callExpr.Args[0])
		if varName != "" {
			deferredPuts[varName] = true
		}
	}
}

// isPoolGetCall vérifie si l'expression est un appel à pool.Get().
//
// Params:
//   - expr: l'expression à vérifier
//   - pass: la passe d'analyse
//
// Returns:
//   - bool: true si c'est pool.Get()
func isPoolGetCall(expr ast.Expr, pass *analysis.Pass) bool {
	callExpr, ok := expr.(*ast.CallExpr)
	if !ok {
		// Retourne false car ce n'est pas un appel
		return false
	}

	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok || selExpr.Sel.Name != "Get" {
		// Retourne false car ce n'est pas .Get()
		return false
	}

	// Vérifier le type de l'objet appelant
	if pass.TypesInfo != nil {
		if t := pass.TypesInfo.TypeOf(selExpr.X); t != nil {
			// Vérifier si c'est sync.Pool ou *sync.Pool
			typeStr := t.String()
			if strings.Contains(typeStr, "sync.Pool") {
				// Retourne true car c'est un sync.Pool
				return true
			}
		}
	}

	// Fallback: vérifier le nom de la variable
	if ident, ok := selExpr.X.(*ast.Ident); ok {
		varName := strings.ToLower(ident.Name)
		// Retourne true si le nom suggère un pool
		return strings.Contains(varName, "pool")
	}

	// Retourne false car ce n'est pas identifié comme pool
	return false
}

// isPoolPutCall vérifie si l'appel est pool.Put().
//
// Params:
//   - callExpr: l'appel à vérifier
//
// Returns:
//   - bool: true si c'est pool.Put()
func isPoolPutCall(callExpr *ast.CallExpr) bool {
	selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok || selExpr.Sel.Name != "Put" {
		// Retourne false car ce n'est pas .Put()
		return false
	}

	// Vérifier le nom de la variable
	if ident, ok := selExpr.X.(*ast.Ident); ok {
		varName := strings.ToLower(ident.Name)
		// Retourne true si le nom suggère un pool
		return strings.Contains(varName, "pool")
	}

	// Retourne false car ce n'est pas identifié comme pool
	return false
}

// extractVarName extrait le nom de variable d'une expression.
//
// Params:
//   - expr: l'expression
//
// Returns:
//   - string: le nom de la variable ou ""
func extractVarName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		// Retourne le nom de l'identifiant
		return e.Name
	case *ast.CallExpr:
		// Type assertion: buf := pool.Get().([]byte)
		// Retourne "" car on ne peut pas tracker facilement
		return ""
	}
	// Retourne "" car type non supporté
	return ""
}

// reportMissingDeferPut rapporte une violation KTN-POOL-001.
//
// Params:
//   - pass: la passe d'analyse
//   - expr: l'expression pool.Get()
//   - varName: le nom de la variable
func reportMissingDeferPut(pass *analysis.Pass, expr ast.Expr, varName string) {
	pass.Reportf(expr.Pos(),
		"[KTN-POOL-001] Variable '%s' obtenue via pool.Get() sans defer pool.Put().\n"+
			"Cela cause une fuite de ressources car l'objet ne retourne jamais au pool.\n"+
			"Utilisez 'defer pool.Put(%s)' immédiatement après Get().\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - fuite mémoire\n"+
			"  buf := bufferPool.Get().([]byte)\n"+
			"  process(buf)\n"+
			"  // buf n'est jamais retourné au pool\n"+
			"\n"+
			"  // ✅ CORRECT\n"+
			"  buf := bufferPool.Get().([]byte)\n"+
			"  defer bufferPool.Put(buf)\n"+
			"  process(buf)",
		varName, varName)
}

// checkLargeStructByValue vérifie les structs larges passées par valeur.
//
// Params:
//   - pass: la passe d'analyse
//   - file: le fichier à analyser
func checkLargeStructByValue(pass *analysis.Pass, file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		if !ok || funcDecl.Type == nil || funcDecl.Type.Params == nil {
			// Retourne true pour continuer l'inspection
			return true
		}

		checkFunctionParameters(pass, funcDecl)
		// Retourne true pour continuer l'inspection
		return true
	})
}

// checkFunctionParameters vérifie les paramètres d'une fonction.
//
// Params:
//   - pass: la passe d'analyse
//   - funcDecl: la déclaration de fonction
func checkFunctionParameters(pass *analysis.Pass, funcDecl *ast.FuncDecl) {
	for _, field := range funcDecl.Type.Params.List {
		checkParameterSize(pass, field)
	}
}

// checkParameterSize vérifie la taille d'un paramètre.
//
// Params:
//   - pass: la passe d'analyse
//   - field: le champ à vérifier
func checkParameterSize(pass *analysis.Pass, field *ast.Field) {
	if pass.TypesInfo == nil {
		// Retourne car pas d'info de type
		return
	}

	paramType := pass.TypesInfo.TypeOf(field.Type)
	if paramType == nil {
		// Retourne car type inconnu
		return
	}

	// Vérifier si c'est une struct (non-pointeur)
	if _, ok := paramType.Underlying().(*types.Struct); !ok {
		// Retourne car ce n'est pas une struct
		return
	}

	// Vérifier si c'est un pointeur
	if _, ok := paramType.(*types.Pointer); ok {
		// Retourne car c'est déjà un pointeur (correct)
		return
	}

	// Estimer la taille de la struct
	size := estimateStructSize(paramType)
	if size > largeStructThreshold {
		reportLargeStructByValue(pass, field, paramType, size)
	}
}

// estimateStructSize estime la taille d'une struct en bytes.
//
// Params:
//   - t: le type à estimer
//
// Returns:
//   - int64: taille estimée en bytes
func estimateStructSize(t types.Type) int64 {
	structType, ok := t.Underlying().(*types.Struct)
	if !ok {
		// Retourne 0 car ce n'est pas une struct
		return 0
	}

	var totalSize int64
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		fieldSize := estimateTypeSize(field.Type())
		totalSize += fieldSize
	}

	// Retourne la taille totale estimée
	return totalSize
}

// estimateTypeSize estime la taille d'un type en bytes.
//
// Params:
//   - t: le type
//
// Returns:
//   - int64: taille estimée
func estimateTypeSize(t types.Type) int64 {
	switch typ := t.Underlying().(type) {
	case *types.Basic:
		// Retourne la taille des types de base
		return estimateBasicTypeSize(typ)
	case *types.Array:
		elemSize := estimateTypeSize(typ.Elem())
		// Retourne taille élément * longueur
		return elemSize * typ.Len()
	case *types.Slice:
		// Slice header: 24 bytes (ptr + len + cap)
		return 24
	case *types.Pointer:
		// Pointeur: 8 bytes (64-bit)
		return 8
	case *types.Map:
		// Map header: 8 bytes (pointeur)
		return 8
	case *types.Chan:
		// Channel header: 8 bytes (pointeur)
		return 8
	case *types.Struct:
		// Récursion pour structs imbriquées
		return estimateStructSize(t)
	case *types.Interface:
		// Interface: 16 bytes (type + data)
		return 16
	default:
		// Par défaut: 8 bytes
		return 8
	}
}

// estimateBasicTypeSize estime la taille d'un type de base.
//
// Params:
//   - t: le type de base
//
// Returns:
//   - int64: taille en bytes
func estimateBasicTypeSize(t *types.Basic) int64 {
	switch t.Kind() {
	case types.Bool, types.Int8, types.Uint8:
		// Retourne 1 byte pour bool, int8, uint8
		return 1
	case types.Int16, types.Uint16:
		// Retourne 2 bytes pour int16, uint16
		return 2
	case types.Int32, types.Uint32, types.Float32:
		// Retourne 4 bytes pour int32, uint32, float32
		return 4
	case types.Int64, types.Uint64, types.Float64, types.Complex64:
		// Retourne 8 bytes pour int64, uint64, float64, complex64
		return 8
	case types.Complex128:
		// Retourne 16 bytes pour complex128
		return 16
	case types.Int, types.Uint, types.Uintptr:
		// Retourne 8 bytes pour int, uint, uintptr (assume architecture 64-bit)
		return 8
	case types.String:
		// Retourne 16 bytes pour string header (ptr + len)
		return 16
	default:
		// Retourne 8 bytes par défaut pour types inconnus
		return 8
	}
}

// reportLargeStructByValue rapporte une violation KTN-STRUCT-004.
//
// Params:
//   - pass: la passe d'analyse
//   - field: le champ paramètre
//   - paramType: le type du paramètre
//   - size: la taille estimée
func reportLargeStructByValue(pass *analysis.Pass, field *ast.Field, paramType types.Type, size int64) {
	paramNames := ""
	if len(field.Names) > 0 {
		paramNames = field.Names[0].Name
	} else {
		paramNames = "param"
	}

	typeName := paramType.String()

	pass.Reportf(field.Pos(),
		"[KTN-STRUCT-004] Paramètre '%s' de type %s passé par valeur (taille ~%d bytes > %d bytes).\n"+
			"Les grandes structs devraient être passées par pointeur pour éviter des copies coûteuses.\n"+
			"Changez le type de '%s' en '*%s'.\n"+
			"Exemple:\n"+
			"  // ❌ NON-OPTIMAL - copie %d bytes à chaque appel\n"+
			"  func Process(%s %s) {}\n"+
			"\n"+
			"  // ✅ OPTIMAL - passe seulement un pointeur (8 bytes)\n"+
			"  func Process(%s *%s) {}",
		paramNames, typeName, size, largeStructThreshold,
		paramNames, typeName,
		size, paramNames, typeName, paramNames, typeName)
}
