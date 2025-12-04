// Analyzer 010 for the ktnvar package.
package ktnvar

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// MAX_STRUCT_FIELDS définit le nombre maximal de champs pour une struct sans pointeur.
const MAX_STRUCT_FIELDS int = 3

// Analyzer010 checks for large struct usage without pointers
var Analyzer010 = &analysis.Analyzer{
	Name:     "ktnvar010",
	Doc:      "KTN-VAR-010: Utilise des pointeurs pour les structs >64 bytes",
	Run:      runVar010,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar010 exécute l'analyse KTN-VAR-010.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar010(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Vérification du corps de la fonction
		if funcDecl.Body == nil {
			// Pas de corps de fonction
			return
		}

		// Parcours des instructions du corps
		checkFuncBody(pass, funcDecl.Body)
	})

	// Retour de la fonction
	return nil, nil
}

// checkFuncBody vérifie le corps d'une fonction.
//
// Params:
//   - pass: contexte d'analyse
//   - body: corps de la fonction
func checkFuncBody(pass *analysis.Pass, body *ast.BlockStmt) {
	// Parcours des instructions
	for _, stmt := range body.List {
		checkStmtForLargeStruct(pass, stmt)
	}
}

// checkStmtForLargeStruct vérifie une instruction.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: instruction à vérifier
func checkStmtForLargeStruct(pass *analysis.Pass, stmt ast.Stmt) {
	// Vérification du type d'instruction
	switch s := stmt.(type) {
	// Cas d'une affectation
	case *ast.AssignStmt:
		// Vérification des affectations
		checkAssignForLargeStruct(pass, s)
	// Cas d'une déclaration
	case *ast.DeclStmt:
		// Vérification des déclarations
		checkDeclForLargeStruct(pass, s)
	}
}

// checkAssignForLargeStruct vérifie une affectation.
//
// Params:
//   - pass: contexte d'analyse
//   - assign: affectation à vérifier
func checkAssignForLargeStruct(pass *analysis.Pass, assign *ast.AssignStmt) {
	// Parcours des valeurs affectées
	for _, rhs := range assign.Rhs {
		checkExprForLargeStruct(pass, rhs)
	}
}

// checkDeclForLargeStruct vérifie une déclaration.
//
// Params:
//   - pass: contexte d'analyse
//   - decl: déclaration à vérifier
func checkDeclForLargeStruct(pass *analysis.Pass, decl *ast.DeclStmt) {
	genDecl, ok := decl.Decl.(*ast.GenDecl)
	// Vérification du type de déclaration
	if !ok {
		// Pas une déclaration générale
		return
	}

	var valueSpec *ast.ValueSpec
	// Parcours des spécifications
	for _, spec := range genDecl.Specs {
		valueSpec, ok = spec.(*ast.ValueSpec)
		// Vérification de la spécification de valeur
		if !ok {
			// Pas une spécification de valeur
			continue
		}

		// Vérification du type de variable
		if valueSpec.Type != nil {
			checkTypeForLargeStruct(pass, valueSpec.Type, valueSpec.Pos())
		}

		// Parcours des valeurs
		for _, value := range valueSpec.Values {
			checkExprForLargeStruct(pass, value)
		}
	}
}

// checkExprForLargeStruct vérifie une expression.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression à vérifier
func checkExprForLargeStruct(pass *analysis.Pass, expr ast.Expr) {
	compositeLit, ok := expr.(*ast.CompositeLit)
	// Vérification du littéral composite
	if !ok {
		// Pas un littéral composite
		return
	}

	// Vérification du type
	checkTypeForLargeStruct(pass, compositeLit.Type, compositeLit.Pos())
}

// checkTypeForLargeStruct vérifie un type.
//
// Params:
//   - pass: contexte d'analyse
//   - typ: type à vérifier
//   - pos: position du type
func checkTypeForLargeStruct(pass *analysis.Pass, typ ast.Expr, pos token.Pos) {
	// Vérification que ce n'est pas un pointeur
	if _, isPointer := typ.(*ast.StarExpr); isPointer {
		// C'est un pointeur, OK
		return
	}

	// Récupération du type réel
	typeInfo := pass.TypesInfo.TypeOf(typ)
	// Vérification du type
	if typeInfo == nil {
		// Type inconnu
		return
	}

	// Ignorer les types externes (frameworks comme Terraform)
	if isExternalType(typeInfo, pass) {
		return
	}

	// Vérification que c'est une struct
	structType, ok := typeInfo.Underlying().(*types.Struct)
	// Vérification du type struct
	if !ok {
		// Pas une struct
		return
	}

	// Comptage des champs
	numFields := structType.NumFields()
	// Vérification du nombre de champs
	if numFields > MAX_STRUCT_FIELDS {
		// Grande struct détectée
		pass.Reportf(
			pos,
			"KTN-VAR-010: utilisez un pointeur pour les structs >64 bytes (%d champs détectés)",
			numFields,
		)
	}
}

// isExternalType checks if type is from external package.
//
// Params:
//   - typeInfo: Type to check
//   - pass: Analysis pass
//
// Returns:
//   - bool: true if type is from external package
func isExternalType(typeInfo types.Type, pass *analysis.Pass) bool {
	// Check if it's a named type
	named, ok := typeInfo.(*types.Named)
	// Verification de la condition
	if !ok {
		return false
	}

	// Get package of the type
	obj := named.Obj()
	// Verification de la condition
	if obj == nil || obj.Pkg() == nil {
		return false
	}

	// Check if type is from current package
	return obj.Pkg().Path() != pass.Pkg.Path()
}
