// Package ktnvar provides analyzers for variable-related lint rules.
package ktnvar

import (
	"go/ast"
	"go/token"
	"go/types"
	"math"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeVar009 is the rule code for this analyzer
	ruleCodeVar009 string = "KTN-VAR-009"
	// defaultMaxStructBytes max bytes for struct without pointer
	defaultMaxStructBytes int = 64
)

// Analyzer009 checks for large struct passed by value in function parameters
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

	// Récupérer le seuil configuré en bytes
	maxBytes := cfg.GetThreshold(ruleCodeVar009, defaultMaxStructBytes)
	// Get verbose setting to pass down (avoid global dependency in helpers)
	verbose := cfg.Verbose

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip excluded files
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		// Vérification de l'exclusion
		if cfg.IsFileExcluded(ruleCodeVar009, filename) {
			// Fichier exclu
			return
		}

		// Vérifier les receivers (méthodes)
		if funcDecl.Recv != nil {
			// Analyse des receivers
			checkFuncParams009(pass, funcDecl.Recv, maxBytes, verbose)
		}

		// Vérifier les paramètres de la fonction
		if funcDecl.Type.Params != nil {
			// Analyse des paramètres
			checkFuncParams009(pass, funcDecl.Type.Params, maxBytes, verbose)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// checkFuncParams009 vérifie les paramètres d'une fonction.
//
// Params:
//   - pass: contexte d'analyse
//   - params: liste des paramètres
//   - maxBytes: taille max en bytes
//   - verbose: mode verbose pour les messages
func checkFuncParams009(pass *analysis.Pass, params *ast.FieldList, maxBytes int, verbose bool) {
	// Handle nil params gracefully
	if params == nil {
		return
	}
	// Clamp invalid threshold to default
	if maxBytes <= 0 {
		maxBytes = defaultMaxStructBytes
	}
	// Parcours des paramètres
	for _, param := range params.List {
		// Si le paramètre a des noms (ex: a, b T), vérifier chaque nom
		if len(param.Names) > 0 {
			// Pour (a, b T), ne reporter qu'une seule fois pour éviter les doublons
			pos := param.Names[0].NamePos
			checkParamType009(pass, param.Type, pos, maxBytes, verbose)
			continue
		}
		// Paramètre sans nom (ex: func f(T)), utiliser la position du type
		checkParamType009(pass, param.Type, param.Pos(), maxBytes, verbose)
	}
}

// checkParamType009 vérifie le type d'un paramètre.
//
// Params:
//   - pass: contexte d'analyse
//   - typ: type du paramètre
//   - pos: position du paramètre
//   - maxBytes: taille max en bytes
//   - verbose: mode verbose pour les messages
func checkParamType009(pass *analysis.Pass, typ ast.Expr, pos token.Pos, maxBytes int, verbose bool) {
	// Handle variadic params: `...T` should be checked as `T`
	if ell, ok := typ.(*ast.Ellipsis); ok && ell.Elt != nil {
		typ = ell.Elt
	}

	// Get struct size if applicable
	sizeBytes := getStructSize009(pass, typ)
	// Check if size exceeds threshold
	if sizeBytes > int64(maxBytes) {
		// Guard against int64 to int overflow on 32-bit systems
		displaySize := sizeBytes
		// Cap displaySize to math.MaxInt32 for safe int cast
		if displaySize > math.MaxInt32 {
			displaySize = math.MaxInt32
		}
		// Grande struct détectée
		msg, _ := messages.Get(ruleCodeVar009)
		pass.Reportf(
			pos,
			"%s: %s",
			ruleCodeVar009,
			msg.Format(verbose, int(displaySize), maxBytes, maxBytes),
		)
	}
}

// getStructSize009 returns the size of a struct type, or -1 if not applicable.
//
// Params:
//   - pass: contexte d'analyse
//   - typ: type expression
//
// Returns:
//   - int64: size in bytes, or -1 if not a local struct or can't determine
func getStructSize009(pass *analysis.Pass, typ ast.Expr) int64 {
	// Ignorer les pointeurs (déjà passés par référence)
	if _, isPointer := typ.(*ast.StarExpr); isPointer {
		// C'est un pointeur, OK
		return -1
	}

	// Récupération du type réel
	typeInfo := pass.TypesInfo.TypeOf(typ)
	// Vérification du type
	if typeInfo == nil {
		// Type inconnu
		return -1
	}

	// Ignorer les types externes (frameworks comme Terraform)
	if isExternalType009(typeInfo, pass) {
		// Retour de la fonction
		return -1
	}

	// Vérification que c'est une struct
	if _, ok := typeInfo.Underlying().(*types.Struct); !ok {
		// Pas une struct
		return -1
	}

	// Calcul de la taille en bytes
	// pass.TypesSizes is the only reliable source for the analysis target;
	// if unavailable, skip size-based reporting to avoid false positives.
	sizes := pass.TypesSizes
	// Vérifier si les informations de taille sont disponibles
	if sizes == nil {
		// Retour anticipé si taille indisponible
		return -1
	}

	// Calcul de la taille
	sz := sizes.Sizeof(typeInfo)
	// Skip reporting on unknown/invalid sizes
	if sz <= 0 {
		return -1
	}

	// Retour de la taille
	return sz
}

// isExternalType009 checks if type is from external package.
//
// Params:
//   - typeInfo: Type to check
//   - pass: Analysis pass
//
// Returns:
//   - bool: true if type is from external package
func isExternalType009(typeInfo types.Type, pass *analysis.Pass) bool {
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
