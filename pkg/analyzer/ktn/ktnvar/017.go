// Package ktnvar implements KTN linter rules.
package ktnvar

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// initialValueReceiversCap initial cap for value receivers
	initialValueReceiversCap int = 10
)

// Analyzer017 détecte les copies de mutex.
//
// Copier un mutex crée des bugs de concurrence car la copie ne partage pas
// le même état de verrouillage.
var Analyzer017 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar017",
	Doc:      "KTN-VAR-017: Vérifie les copies de mutex (sync.Mutex, sync.RWMutex, atomic.Value)",
	Run:      runVar017,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar017 exécute l'analyse de détection des copies de mutex.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - interface{}: toujours nil
//   - error: erreur éventuelle
func runVar017(pass *analysis.Pass) (any, error) {
	// Récupération de l'inspecteur AST
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecte des types avec receivers par valeur
	typesWithValueRecv := collectTypesWithValueReceivers(pass, insp)

	// Analyse des structs avec mutex (seulement celles avec value receivers)
	checkStructsWithMutex(pass, insp, typesWithValueRecv)

	// Analyse des receivers par valeur
	checkValueReceivers(pass, insp)

	// Analyse des paramètres par valeur
	checkValueParams(pass, insp)

	// Analyse des assignations
	checkAssignments(pass, insp)

	// Traitement
	return nil, nil
}

// collectTypesWithValueReceivers collecte les types ayant des receivers par valeur.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//
// Returns:
//   - map[string]bool: map des types avec receivers par valeur
func collectTypesWithValueReceivers(_pass *analysis.Pass, insp *inspector.Inspector) map[string]bool {
	// Map pour stocker les types
	typesWithValueRecv := make(map[string]bool, initialValueReceiversCap)

	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Parcours des fonctions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Cast en fonction
		funcDecl, ok := n.(*ast.FuncDecl)
		// Vérification de la condition
		if !ok || funcDecl.Recv == nil {
			// Traitement
			return
		}

		// Récupération du receiver
		if len(funcDecl.Recv.List) == 0 {
			// Traitement
			return
		}

		recv := funcDecl.Recv.List[0]

		// Vérification si c'est un receiver par valeur (pas un pointeur)
		if !isPointerType(recv.Type) {
			// Extraction du nom du type
			typeName := getTypeName(recv.Type)
			// Vérification de la condition
			if typeName != "" {
				typesWithValueRecv[typeName] = true
			}
		}
	})

	// Retour de la map
	return typesWithValueRecv
}

// checkStructsWithMutex vérifie les structs contenant des mutex.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//   - typesWithValueRecv: types ayant des receivers par valeur
func checkStructsWithMutex(pass *analysis.Pass, insp *inspector.Inspector, typesWithValueRecv map[string]bool) {
	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.TypeSpec)(nil),
	}

	// Parcours des types
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Cast en type spec
		typeSpec, ok := n.(*ast.TypeSpec)
		// Vérification de la condition
		if !ok {
			// Traitement
			return
		}

		// Vérification que c'est une struct
		structType, ok := typeSpec.Type.(*ast.StructType)
		// Vérification de la condition
		if !ok {
			// Traitement
			return
		}

		// Vérification si ce type a des receivers par valeur
		if !typesWithValueRecv[typeSpec.Name.Name] {
			// Si tous les receivers sont des pointeurs, pas de problème
			return
		}

		// Vérification des champs
		for _, field := range structType.Fields.List {
			// Vérification si le champ est un mutex
			if mutexType := getMutexType(pass, field.Type); mutexType != "" {
				// Rapport d'erreur seulement si le type a des receivers par valeur
				pass.Reportf(
					field.Pos(),
					"KTN-VAR-017: struct contient %s, utiliser *%s pour éviter les copies",
					mutexType,
					typeSpec.Name.Name,
				)
			}
		}
	})
}

// checkValueReceivers vérifie les receivers par valeur.
//
// Params:
//   - pass: contexte d'analyse
//   - inspect: inspecteur AST
func checkValueReceivers(pass *analysis.Pass, insp *inspector.Inspector) {
	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Parcours des fonctions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Cast en fonction
		funcDecl, ok := n.(*ast.FuncDecl)
		// Vérification de la condition
		if !ok || funcDecl.Recv == nil {
			// Traitement
			return
		}

		// Récupération du receiver
		if len(funcDecl.Recv.List) == 0 {
			// Traitement
			return
		}

		recv := funcDecl.Recv.List[0]

		// Vérification si c'est un receiver par valeur
		if !isPointerType(recv.Type) {
			// Vérification si le type contient un mutex
			if hasMutex(pass, recv.Type) {
				typeName := getTypeName(recv.Type)
				mutexType := getMutexTypeFromType(pass, recv.Type)

				pass.Reportf(
					recv.Pos(),
					"KTN-VAR-017: receiver par valeur copie %s, utiliser *%s",
					mutexType,
					typeName,
				)
			}
		}
	})
}

// checkValueParams vérifie les paramètres par valeur.
//
// Params:
//   - pass: contexte d'analyse
//   - inspect: inspecteur AST
func checkValueParams(pass *analysis.Pass, insp *inspector.Inspector) {
	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Parcours des fonctions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Cast en fonction
		funcDecl, ok := n.(*ast.FuncDecl)
		// Vérification de la condition
		if !ok || funcDecl.Type.Params == nil {
			// Traitement
			return
		}

		// Vérification des paramètres
		for _, param := range funcDecl.Type.Params.List {
			// Vérification si c'est un type mutex par valeur
			if mutexType := getMutexType(pass, param.Type); mutexType != "" {
				// Vérification de la condition
				if !isPointerType(param.Type) {
					pass.Reportf(
						param.Pos(),
						"KTN-VAR-017: passage de %s par valeur, utiliser *%s",
						mutexType,
						mutexType,
					)
				}
			}
		}
	})
}

// checkAssignments vérifie les assignations de mutex.
//
// Params:
//   - pass: contexte d'analyse
//   - inspect: inspecteur AST
func checkAssignments(pass *analysis.Pass, insp *inspector.Inspector) {
	// Types de nœuds à analyser
	nodeFilter := []ast.Node{
		(*ast.AssignStmt)(nil),
	}

	// Parcours des assignations
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Cast en assignation
		assign, ok := n.(*ast.AssignStmt)
		// Vérification de la condition
		if !ok {
			// Traitement
			return
		}

		// Vérification de chaque assignation
		for i, rhs := range assign.Rhs {
			// Vérification si on assigne un mutex
			if i < len(assign.Lhs) {
				// Vérification de la condition
				if isMutexCopy(pass, assign.Lhs[i], rhs) {
					pass.Reportf(
						assign.Pos(),
						"KTN-VAR-017: copie de sync.Mutex détectée, utiliser un pointeur",
					)
				}
			}
		}
	})
}

// getMutexType retourne le type de mutex ou "".
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression à vérifier
//
// Returns:
//   - string: type de mutex ou ""
func getMutexType(pass *analysis.Pass, expr ast.Expr) string {
	// Récupération du type
	tv, ok := pass.TypesInfo.Types[expr]
	// Vérification de la condition
	if !ok {
		// Traitement
		return ""
	}

	// Vérification du type
	return getMutexTypeName(tv.Type)
}

// getMutexTypeName retourne le nom du type mutex.
//
// Params:
//   - t: type à vérifier
//
// Returns:
//   - string: nom du type mutex ou ""
func getMutexTypeName(t types.Type) string {
	// Récupération du type nommé
	named, ok := t.(*types.Named)
	// Vérification de la condition
	if !ok {
		// Traitement
		return ""
	}

	// Vérification du package et du nom
	obj := named.Obj()
	// Vérification de la condition
	if obj.Pkg() == nil {
		// Traitement
		return ""
	}

	pkg := obj.Pkg().Path()
	name := obj.Name()

	// Vérification des types sync qui ne doivent pas être copiés
	if pkg == "sync" {
		// Liste des types sync non-copiables
		switch name {
		// sync.Mutex ne doit pas être copié
		case "Mutex", "RWMutex", "Cond", "Once", "WaitGroup", "Pool", "Map":
			// Retour du type qualifié
			return "sync." + name
		}
	}

	// Vérification des types atomic qui ne doivent pas être copiés
	if pkg == "sync/atomic" {
		// Tous les types atomic.* ne doivent pas être copiés
		// Inclut: Value, Bool, Int32, Int64, Uint32, Uint64, Uintptr, Pointer
		return "atomic." + name
	}

	// Type non-mutex
	return ""
}

// isPointerType vérifie si un type est un pointeur.
//
// Params:
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si pointeur
func isPointerType(expr ast.Expr) bool {
	// Vérification directe
	_, ok := expr.(*ast.StarExpr)
	// Traitement
	return ok
}

// hasMutex vérifie si un type contient un mutex.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression à vérifier
//
// Returns:
//   - bool: true si mutex trouvé
func hasMutex(pass *analysis.Pass, expr ast.Expr) bool {
	// Récupération du type
	tv, ok := pass.TypesInfo.Types[expr]
	// Vérification de la condition
	if !ok {
		// Traitement
		return false
	}

	// Vérification si c'est une struct
	return hasMutexInType(tv.Type)
}

// hasMutexInType vérifie si un type contient un mutex.
//
// Params:
//   - t: type à vérifier
//
// Returns:
//   - bool: true si mutex trouvé
func hasMutexInType(t types.Type) bool {
	var ok bool
	var ptr *types.Pointer
	// Déréférencement des pointeurs
	if ptr, ok = t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	var st *types.Struct
	// Vérification si c'est une struct
	st, ok = t.Underlying().(*types.Struct)
	// Vérification de la condition
	if !ok {
		// Traitement
		return false
	}

	// Parcours des champs avec itérateur standard
	for field := range st.Fields() {
		// Vérification de la condition
		if getMutexTypeName(field.Type()) != "" {
			// Traitement
			return true
		}
	}

	// Traitement
	return false
}

// getMutexTypeFromType retourne le type de mutex d'un type.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression à vérifier
//
// Returns:
//   - string: type de mutex ou ""
func getMutexTypeFromType(pass *analysis.Pass, expr ast.Expr) string {
	// Récupération du type
	tv, ok := pass.TypesInfo.Types[expr]
	// Vérification de la condition
	if !ok {
		// Traitement
		return ""
	}

	var ptr *types.Pointer
	// Déréférencement des pointeurs
	t := tv.Type
	// Vérification de la condition
	if ptr, ok = t.(*types.Pointer); ok {
		t = ptr.Elem()
	}

	var st *types.Struct
	// Vérification si c'est une struct
	st, ok = t.Underlying().(*types.Struct)
	// Vérification de la condition
	if !ok {
		// Traitement
		return ""
	}

	// Parcours des champs avec itérateur standard
	for field := range st.Fields() {
		// Vérification de la condition
		if mutexType := getMutexTypeName(field.Type()); mutexType != "" {
			// Traitement
			return mutexType
		}
	}

	// Traitement
	return ""
}

// getTypeName retourne le nom d'un type.
//
// Params:
//   - expr: expression à analyser
//
// Returns:
//   - string: nom du type
func getTypeName(expr ast.Expr) string {
	// Vérification si c'est un identifiant
	if ident, ok := expr.(*ast.Ident); ok {
		// Traitement
		return ident.Name
	}

	// Vérification si c'est une star expression
	if star, ok := expr.(*ast.StarExpr); ok {
		// Traitement
		return getTypeName(star.X)
	}

	// Traitement
	return ""
}

// isMutexCopy vérifie si une assignation copie un mutex.
//
// Params:
//   - pass: contexte d'analyse
//   - lhs: partie gauche de l'assignation
//   - rhs: partie droite de l'assignation
//
// Returns:
//   - bool: true si copie de mutex
func isMutexCopy(pass *analysis.Pass, _lhs, rhs ast.Expr) bool {
	// Récupération du type RHS
	tv, ok := pass.TypesInfo.Types[rhs]
	// Vérification de la condition
	if !ok {
		// Traitement
		return false
	}

	// Vérification si c'est un mutex
	if getMutexTypeName(tv.Type) != "" {
		// Vérification que ce n'est pas un pointeur
		_, isPointer := tv.Type.(*types.Pointer)
		// Traitement
		return !isPointer
	}

	// Traitement
	return false
}
