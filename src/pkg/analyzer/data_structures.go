package analyzer

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Analyzers
var (
	// DataStructuresAnalyzer vérifie les opérations sur les structures de données.
	DataStructuresAnalyzer *analysis.Analyzer = &analysis.Analyzer{
		Name: "ktndatastructures",
		Doc:  "Vérifie les opérations sur map, slice et array",
		Run:  runDataStructuresAnalyzer,
	}
)

// runDataStructuresAnalyzer exécute l'analyseur de structures de données.
//
// Params:
//   - pass: la passe d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: toujours nil
func runDataStructuresAnalyzer(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch node := n.(type) {
			case *ast.AssignStmt:
				checkMapAssignment(pass, node)
			case *ast.IndexExpr:
				checkSliceIndexing(pass, node)
			case *ast.CompositeLit:
				checkArrayLiteral(pass, node)
			}
			// Retourne true pour continuer l'inspection
			return true
		})
	}
	// Retourne nil car l'analyse est terminée
	return nil, nil
}

// checkMapAssignment vérifie l'écriture dans une map sans vérification.
//
// Params:
//   - pass: la passe d'analyse
//   - assign: l'assignation
func checkMapAssignment(pass *analysis.Pass, assign *ast.AssignStmt) {
	// Vérifier si Lhs contient un accès à map
	for _, lhs := range assign.Lhs {
		indexExpr, ok := lhs.(*ast.IndexExpr)
		if !ok {
			// Pas un accès indexé
			continue
		}

		// Vérifier si X est une map
		if isMapExpr(indexExpr.X) {
			// C'est une écriture dans une map
			// Vérifier si la map a été vérifiée pour nil
			if mapIdent, ok := indexExpr.X.(*ast.Ident); ok {
				if !isMapCheckedForNil(pass, mapIdent, assign) {
					reportUncheckedMapWrite(pass, indexExpr, mapIdent.Name)
				}
			}
		}
	}
}

// checkSliceIndexing vérifie l'indexation de slice sans vérification de bounds.
//
// Params:
//   - pass: la passe d'analyse
//   - index: l'expression d'indexation
func checkSliceIndexing(pass *analysis.Pass, index *ast.IndexExpr) {
	// Exclure les paramètres de type génériques (Container[T])
	if isGenericTypeInstantiation(pass, index) {
		// C'est un générique, pas une indexation de slice
		return
	}

	// Vérifier si X est un slice (et pas une map)
	if !isSliceExpr(index.X) {
		// Pas un slice
		// Retourne
		return
	}

	// Exclure les maps (l'indexation de map est sûre)
	if isLikelyMapUsage(index) {
		// Probablement une map, pas un slice
		return
	}

	// Vérifier si l'index est un littéral
	if _, ok := index.Index.(*ast.BasicLit); ok {
		// Index littéral, on ne peut pas vérifier statiquement
		// Retourne
		return
	}

	// Vérifier si l'index est une variable
	indexIdent, ok := index.Index.(*ast.Ident)
	if !ok {
		// Pas un identifiant simple
		// Retourne
		return
	}

	// Vérifier si l'index a été vérifié contre len()
	if !isIndexCheckedAgainstLen(pass, indexIdent, index) {
		if sliceIdent, ok := index.X.(*ast.Ident); ok {
			reportUncheckedSliceIndex(pass, index, sliceIdent.Name, indexIdent.Name)
		}
	}
}

// checkArrayLiteral vérifie les tableaux avec taille incohérente.
//
// Params:
//   - pass: la passe d'analyse
//   - lit: le composite literal
func checkArrayLiteral(pass *analysis.Pass, lit *ast.CompositeLit) {
	// Vérifier si c'est un array (pas un slice)
	arrayType, ok := lit.Type.(*ast.ArrayType)
	if !ok || arrayType.Len == nil {
		// Pas un array avec taille explicite
		// Retourne
		return
	}

	// Extraire la taille déclarée
	declaredLen := getArraySize(arrayType)
	if declaredLen == -1 {
		// Taille non déterminable
		// Retourne
		return
	}

	// Compter les éléments
	actualLen := len(lit.Elts)

	// Vérifier l'incohérence
	if actualLen > declaredLen {
		reportArraySizeMismatch(pass, lit, declaredLen, actualLen)
	}
}

// isMapExpr vérifie si une expression est une map.
//
// Params:
//   - expr: l'expression
//
// Returns:
//   - bool: true si c'est une map
func isMapExpr(expr ast.Expr) bool {
	// Pour une détection robuste, on devrait utiliser pass.TypesInfo
	// Ici on fait une détection basique sur la structure AST
	switch e := expr.(type) {
	case *ast.Ident:
		// Pourrait être une map, on suppose que oui
		return true
	case *ast.SelectorExpr:
		// Pourrait être une map
		return true
	case *ast.IndexExpr:
		// Résultat d'un accès, pourrait être une map
		return true
	case *ast.CallExpr:
		// Résultat d'un appel, pourrait retourner une map
		if ident, ok := e.Fun.(*ast.Ident); ok && ident.Name == "make" {
			// make() retourne une map si premier arg est MapType
			if len(e.Args) > 0 {
				_, isMap := e.Args[0].(*ast.MapType)
				return isMap
			}
		}
		return false
	default:
		return false
	}
}

// isGenericTypeInstantiation vérifie si une IndexExpr est une instanciation de type générique.
//
// Params:
//   - pass: la passe d'analyse
//   - index: l'expression d'indexation
//
// Returns:
//   - bool: true si c'est un générique (Container[T])
func isGenericTypeInstantiation(pass *analysis.Pass, index *ast.IndexExpr) bool {
	// Vérifier si l'index est un type (identifiant commençant par majuscule = type)
	indexIdent, ok := index.Index.(*ast.Ident)
	if !ok {
		// Pas un identifiant simple, probablement pas un paramètre de type
		return false
	}

	// Les paramètres de type commencent souvent par une majuscule (T, K, V, etc.)
	// ou sont des types connus
	indexName := indexIdent.Name
	if len(indexName) > 0 {
		firstChar := indexName[0]
		// Vérifier si c'est une majuscule (type) ou un nom de type connu
		if firstChar >= 'A' && firstChar <= 'Z' {
			// C'est probablement un paramètre de type générique
			return true
		}
		// Vérifier les types built-in Go
		if indexName == "string" || indexName == "int" || indexName == "any" ||
		   indexName == "bool" || indexName == "byte" || indexName == "rune" {
			// C'est un type built-in, donc un générique
			return true
		}
	}

	// Utiliser TypesInfo si disponible pour une détection plus précise
	if pass.TypesInfo != nil {
		// Vérifier si X est un nom de type (pas une valeur)
		if ident, ok := index.X.(*ast.Ident); ok {
			obj := pass.TypesInfo.ObjectOf(ident)
			if obj != nil {
				// Si c'est un TypeName, alors c'est une instanciation de type générique
				_, isTypeName := obj.(*types.TypeName)
				return isTypeName
			}
		}
	}

	return false
}

// isSliceExpr vérifie si une expression est un slice.
//
// Params:
//   - expr: l'expression
//
// Returns:
//   - bool: true si c'est un slice
func isSliceExpr(expr ast.Expr) bool {
	// Pour une détection robuste, on devrait utiliser pass.TypesInfo
	// Ici on fait une détection basique
	switch expr.(type) {
	case *ast.Ident:
		// Pourrait être un slice
		return true
	case *ast.SelectorExpr:
		// Pourrait être un slice
		return true
	case *ast.CallExpr:
		// Pourrait retourner un slice
		return true
	default:
		return false
	}
}

// isLikelyMapUsage détecte si une indexation est probablement sur une map.
//
// Params:
//   - index: l'expression d'indexation
//
// Returns:
//   - bool: true si c'est probablement une map
func isLikelyMapUsage(index *ast.IndexExpr) bool {
	// Si l'index est une string, c'est probablement une map
	if ident, ok := index.Index.(*ast.Ident); ok {
		// Si le nom de la variable d'index contient "key" ou "k", c'est probablement une map
		if ident.Name == "key" || ident.Name == "k" {
			return true
		}
	}

	// Si l'index est un basiclit de type string, c'est une map
	if lit, ok := index.Index.(*ast.BasicLit); ok {
		if lit.Kind.String() == "STRING" {
			return true
		}
	}

	return false
}

// isMapCheckedForNil vérifie si une map a été vérifiée pour nil.
//
// Params:
//   - pass: la passe d'analyse
//   - mapIdent: l'identifiant de la map
//   - usage: le nœud d'usage
//
// Returns:
//   - bool: true si vérifiée
func isMapCheckedForNil(pass *analysis.Pass, mapIdent *ast.Ident, usage ast.Node) bool {
	// Chercher une condition if mapName != nil avant usage
	for _, file := range pass.Files {
		checked := false

		ast.Inspect(file, func(n ast.Node) bool {
			// Si on atteint l'usage, on arrête
			if n == usage {
				// Retourne false pour arrêter
				return false
			}

			// Vérifier si la map est créée avec make()
			assignStmt, ok := n.(*ast.AssignStmt)
			if ok {
				for i, lhs := range assignStmt.Lhs {
					if lhsIdent, ok := lhs.(*ast.Ident); ok && lhsIdent.Name == mapIdent.Name {
						// Vérifier si Rhs est make()
						if i < len(assignStmt.Rhs) {
							if callExpr, ok := assignStmt.Rhs[i].(*ast.CallExpr); ok {
								if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == "make" {
									// Map créée avec make(), c'est sûr
									checked = true
									return false
								}
							}
						}
					}
				}
			}

			// Chercher un if qui vérifie != nil
			ifStmt, ok := n.(*ast.IfStmt)
			if !ok {
				// Retourne true pour continuer
				return true
			}

			// Vérifier si la condition est mapName != nil
			if isNilCheck(ifStmt.Cond, mapIdent.Name) {
				checked = true
			}

			// Retourne true pour continuer
			return true
		})

		if checked {
			// Retourne true car vérifiée
			return true
		}
	}

	// Retourne false car pas vérifiée
	return false
}

// isIndexCheckedAgainstLen vérifie si un index a été vérifié contre len().
//
// Params:
//   - pass: la passe d'analyse
//   - indexIdent: l'identifiant de l'index
//   - usage: le nœud d'usage
//
// Returns:
//   - bool: true si vérifié
func isIndexCheckedAgainstLen(pass *analysis.Pass, indexIdent *ast.Ident, usage ast.Node) bool {
	// Chercher une condition if index < len(...) avant usage
	for _, file := range pass.Files {
		checked := false

		ast.Inspect(file, func(n ast.Node) bool {
			// Si on atteint l'usage, on arrête
			if n == usage {
				// Retourne false pour arrêter
				return false
			}

			// Vérifier si l'index vient d'un range loop
			rangeStmt, ok := n.(*ast.RangeStmt)
			if ok && rangeStmt.Key != nil {
				if keyIdent, ok := rangeStmt.Key.(*ast.Ident); ok {
					if keyIdent.Name == indexIdent.Name {
						// L'index vient d'un range, c'est sûr
						checked = true
						return false
					}
				}
			}

			// Chercher un if qui vérifie < len()
			ifStmt, ok := n.(*ast.IfStmt)
			if !ok {
				// Retourne true pour continuer
				return true
			}

			// Vérifier si la condition est index < len(...)
			if isLenCheck(ifStmt.Cond, indexIdent.Name) {
				checked = true
			}

			// Retourne true pour continuer
			return true
		})

		if checked {
			// Retourne true car vérifié
			return true
		}
	}

	// Retourne false car pas vérifié
	return false
}

// isNilCheck vérifie si une expression est une vérification != nil.
//
// Params:
//   - expr: l'expression
//   - varName: nom de la variable
//
// Returns:
//   - bool: true si c'est varName != nil
func isNilCheck(expr ast.Expr, varName string) bool {
	binaryExpr, ok := expr.(*ast.BinaryExpr)
	if !ok {
		// Pas une expression binaire
		return false
	}

	if binaryExpr.Op.String() != "!=" {
		// Pas un !=
		return false
	}

	// Vérifier X == varName et Y == nil (ou inverse)
	xIdent, xOk := binaryExpr.X.(*ast.Ident)
	yIdent, yOk := binaryExpr.Y.(*ast.Ident)

	if xOk && yOk {
		return (xIdent.Name == varName && yIdent.Name == "nil") ||
			(yIdent.Name == varName && xIdent.Name == "nil")
	}

	return false
}

// isLenCheck vérifie si une expression est une vérification < len().
//
// Params:
//   - expr: l'expression
//   - indexName: nom de l'index
//
// Returns:
//   - bool: true si c'est index < len(...)
func isLenCheck(expr ast.Expr, indexName string) bool {
	binaryExpr, ok := expr.(*ast.BinaryExpr)
	if !ok {
		// Pas une expression binaire
		return false
	}

	if binaryExpr.Op.String() != "<" {
		// Pas un <
		return false
	}

	// Vérifier X == indexName
	xIdent, ok := binaryExpr.X.(*ast.Ident)
	if !ok || xIdent.Name != indexName {
		// X n'est pas l'index
		return false
	}

	// Vérifier Y == len(...)
	yCall, ok := binaryExpr.Y.(*ast.CallExpr)
	if !ok {
		// Y n'est pas un appel
		return false
	}

	yFunc, ok := yCall.Fun.(*ast.Ident)
	if !ok || yFunc.Name != "len" {
		// Y n'est pas len()
		return false
	}

	return true
}

// getArraySize extrait la taille d'un array.
//
// Params:
//   - arrayType: le type array
//
// Returns:
//   - int: la taille, ou -1 si non déterminable
func getArraySize(arrayType *ast.ArrayType) int {
	if arrayType.Len == nil {
		// Retourne -1 car pas de taille (slice)
		return -1
	}

	basicLit, ok := arrayType.Len.(*ast.BasicLit)
	if !ok {
		// Retourne -1 car taille non littérale
		return -1
	}

	// Parser la valeur
	var size int
	_, err := fmt.Sscanf(basicLit.Value, "%d", &size)
	if err != nil {
		// Retourne -1 car parsing échoué
		return -1
	}

	// Retourne la taille
	return size
}

// reportUncheckedMapWrite rapporte une violation KTN-MAP-001.
//
// Params:
//   - pass: la passe d'analyse
//   - index: l'expression d'indexation
//   - mapName: nom de la map
func reportUncheckedMapWrite(pass *analysis.Pass, index *ast.IndexExpr, mapName string) {
	pass.Reportf(index.Pos(),
		"[KTN-MAP-001] Écriture dans la map '%s' sans vérification de nil.\n"+
			"Écrire dans une map nil cause un panic.\n"+
			"Vérifiez toujours qu'une map n'est pas nil avant d'y écrire.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - panic si m est nil\n"+
			"  var m map[string]int\n"+
			"  m[\"key\"] = 42  // 💥 PANIC\n"+
			"\n"+
			"  // ✅ CORRECT - initialiser avec make\n"+
			"  m := make(map[string]int)\n"+
			"  m[\"key\"] = 42\n"+
			"\n"+
			"  // ✅ CORRECT - vérifier si nil\n"+
			"  if m != nil {\n"+
			"      m[\"key\"] = 42\n"+
			"  }",
		mapName)
}

// reportUncheckedSliceIndex rapporte une violation KTN-SLICE-001.
//
// Params:
//   - pass: la passe d'analyse
//   - index: l'expression d'indexation
//   - sliceName: nom du slice
//   - indexName: nom de l'index
func reportUncheckedSliceIndex(pass *analysis.Pass, index *ast.IndexExpr, sliceName string, indexName string) {
	pass.Reportf(index.Pos(),
		"[KTN-SLICE-001] Indexation du slice '%s' avec '%s' sans vérification de bounds.\n"+
			"Accéder à un index hors limites cause un panic.\n"+
			"Vérifiez toujours que l'index est valide avant d'accéder.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - panic si i >= len(items)\n"+
			"  v := items[i]\n"+
			"\n"+
			"  // ✅ CORRECT - vérifier les bounds\n"+
			"  if i < len(items) {\n"+
			"      v := items[i]\n"+
			"  }",
		sliceName, indexName)
}

// reportArraySizeMismatch rapporte une violation KTN-ARRAY-002.
//
// Params:
//   - pass: la passe d'analyse
//   - lit: le composite literal
//   - declared: taille déclarée
//   - actual: taille actuelle
func reportArraySizeMismatch(pass *analysis.Pass, lit *ast.CompositeLit, declared int, actual int) {
	pass.Reportf(lit.Pos(),
		"[KTN-ARRAY-002] Taille d'array incohérente: déclaré %d, mais %d éléments fournis.\n"+
			"Un array ne peut pas contenir plus d'éléments que sa taille déclarée.\n"+
			"Soit augmentez la taille, soit utilisez un slice.\n"+
			"Exemple:\n"+
			"  // ❌ MAUVAIS - trop d'éléments\n"+
			"  arr := [2]int{1, 2, 3}  // ERREUR\n"+
			"\n"+
			"  // ✅ CORRECT - bonne taille\n"+
			"  arr := [3]int{1, 2, 3}\n"+
			"\n"+
			"  // ✅ CORRECT - utiliser un slice\n"+
			"  arr := []int{1, 2, 3}",
		declared, actual)
}
