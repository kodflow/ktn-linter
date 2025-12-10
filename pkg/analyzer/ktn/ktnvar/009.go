// Analyzer 009 for the ktnvar package.
package ktnvar

import (
	"go/ast"
	"go/token"
	"go/types"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar009 is the rule code for this analyzer
	ruleCodeVar009 string = "KTN-VAR-009"
	// defaultMaxStructFields max fields for struct without pointer
	defaultMaxStructFields int = 3
)

// Analyzer009 checks for large struct usage without pointers
var Analyzer009 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnvar009",
	Doc:      "KTN-VAR-009: Utilise des pointeurs pour les structs >64 bytes",
	Run:      runVar009,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runVar009 exécute l'analyse KTN-VAR-009.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runVar009(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeVar009) {
		// Règle désactivée
		return nil, nil
	}

	// Récupérer le seuil configuré
	maxFields := cfg.GetThreshold(ruleCodeVar009, defaultMaxStructFields)

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Vérifier si le fichier est exclu
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		if cfg.IsFileExcluded(ruleCodeVar009, filename) {
			// Fichier exclu
			return
		}

		// Vérification du corps de la fonction
		if funcDecl.Body == nil {
			// Pas de corps de fonction
			return
		}

		// Parcours des instructions du corps
		checkFuncBodyVar009(pass, funcDecl.Body, maxFields)
	})

	// Retour de la fonction
	return nil, nil
}

// checkFuncBodyVar009 vérifie le corps d'une fonction.
//
// Params:
//   - pass: contexte d'analyse
//   - body: corps de la fonction
//   - maxFields: nombre max de champs
func checkFuncBodyVar009(pass *analysis.Pass, body *ast.BlockStmt, maxFields int) {
	// Parcours des instructions
	for _, stmt := range body.List {
		checkStmtForLargeStruct(pass, stmt, maxFields)
	}
}

// checkStmtForLargeStruct vérifie une instruction.
//
// Params:
//   - pass: contexte d'analyse
//   - stmt: instruction à vérifier
//   - maxFields: nombre max de champs
func checkStmtForLargeStruct(pass *analysis.Pass, stmt ast.Stmt, maxFields int) {
	// Vérification du type d'instruction
	switch s := stmt.(type) {
	// Cas d'une affectation
	case *ast.AssignStmt:
		// Vérification des affectations
		checkAssignForLargeStruct(pass, s, maxFields)
	// Cas d'une déclaration
	case *ast.DeclStmt:
		// Vérification des déclarations
		checkDeclForLargeStruct(pass, s, maxFields)
	}
}

// checkAssignForLargeStruct vérifie une affectation.
//
// Params:
//   - pass: contexte d'analyse
//   - assign: affectation à vérifier
//   - maxFields: nombre max de champs
func checkAssignForLargeStruct(pass *analysis.Pass, assign *ast.AssignStmt, maxFields int) {
	// Parcours des valeurs affectées
	for _, rhs := range assign.Rhs {
		checkExprForLargeStruct(pass, rhs, maxFields)
	}
}

// checkDeclForLargeStruct vérifie une déclaration.
//
// Params:
//   - pass: contexte d'analyse
//   - decl: déclaration à vérifier
//   - maxFields: nombre max de champs
func checkDeclForLargeStruct(pass *analysis.Pass, decl *ast.DeclStmt, maxFields int) {
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
			checkTypeForLargeStruct(pass, valueSpec.Type, valueSpec.Pos(), maxFields)
		}

		// Parcours des valeurs
		for _, value := range valueSpec.Values {
			checkExprForLargeStruct(pass, value, maxFields)
		}
	}
}

// checkExprForLargeStruct vérifie une expression.
//
// Params:
//   - pass: contexte d'analyse
//   - expr: expression à vérifier
//   - maxFields: nombre max de champs
func checkExprForLargeStruct(pass *analysis.Pass, expr ast.Expr, maxFields int) {
	compositeLit, ok := expr.(*ast.CompositeLit)
	// Vérification du littéral composite
	if !ok {
		// Pas un littéral composite
		return
	}

	// Vérification du type
	checkTypeForLargeStruct(pass, compositeLit.Type, compositeLit.Pos(), maxFields)
}

// checkTypeForLargeStruct vérifie un type.
//
// Params:
//   - pass: contexte d'analyse
//   - typ: type à vérifier
//   - pos: position du type
//   - maxFields: nombre max de champs
func checkTypeForLargeStruct(pass *analysis.Pass, typ ast.Expr, pos token.Pos, maxFields int) {
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
		// Retour de la fonction
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
	if numFields > maxFields {
		// Grande struct détectée
		pass.Reportf(
			pos,
			"KTN-VAR-009: utilisez un pointeur pour les structs >64 bytes (%d champs détectés)",
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
		// Retour de la fonction
		return false
	}

	// Get package of the type
	obj := named.Obj()
	// Verification de la condition
	if obj == nil || obj.Pkg() == nil {
		// Retour de la fonction
		return false
	}

	// Check if type is from current package
	return obj.Pkg().Path() != pass.Pkg.Path()
}
