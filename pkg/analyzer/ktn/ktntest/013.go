// Analyzer 013 for the ktntest package.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// INITIAL_MAP_CAPACITY est la capacité initiale des maps de signatures.
const INITIAL_MAP_CAPACITY int = 32

// testedFuncInfo contient les informations sur une fonction testée.
type testedFuncInfo struct {
	name         string
	returnsError bool
	hasReceiver  bool
	receiverName string
}

// Analyzer013 checks that tests cover error cases
var Analyzer013 = &analysis.Analyzer{
	Name:     "ktntest013",
	Doc:      "KTN-TEST-013: Les tests doivent couvrir les cas d'erreur et exceptions",
	Run:      runTest013,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runTest013 exécute l'analyse KTN-TEST-013.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest013(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecter toutes les fonctions du package avec leur signature
	funcSignatures := collectFuncSignatures(pass, insp)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Analyser les fonctions de test
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Vérification fichier de test
		if !shared.IsTestFile(filename) {
			return
		}

		// Skip exempt test files
		if shared.IsExemptTestFile(filename) {
			return
		}

		// Skip mock files
		if shared.IsMockFile(filename) {
			return
		}

		// Vérifier si c'est une fonction de test unitaire
		if !shared.IsUnitTestFunction(funcDecl) {
			return
		}

		// Skip exempt test names
		if shared.IsExemptTestName(funcDecl.Name.Name) {
			return
		}

		// Analyser le test par rapport à la fonction testée
		analyzeTestFunction(pass, funcDecl, funcSignatures)
	})

	// Retour de la fonction
	return nil, nil
}

// collectFuncSignatures collecte les signatures des fonctions.
//
// Params:
//   - pass: contexte d'analyse
//   - insp: inspecteur AST
//
// Returns:
//   - map[string]testedFuncInfo: map nom -> info fonction
func collectFuncSignatures(pass *analysis.Pass, insp *inspector.Inspector) map[string]testedFuncInfo {
	// Initialiser avec capacité estimée
	signatures := make(map[string]testedFuncInfo, INITIAL_MAP_CAPACITY)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Parcourir toutes les fonctions du pass
	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			return
		}

		// Ignorer les fichiers mock
		if shared.IsMockFile(filename) {
			return
		}

		// Skip mock functions
		if funcDecl.Name != nil && shared.IsMockName(funcDecl.Name.Name) {
			return
		}

		// Ajouter la signature
		addFuncSignature(signatures, funcDecl)
	})

	// Scanner les fichiers sources du répertoire (pour tests externes)
	collectExternalSourceSignatures(pass, signatures)

	// Retour des signatures
	return signatures
}

// addFuncSignature ajoute une signature de fonction.
//
// Params:
//   - signatures: map des signatures
//   - funcDecl: déclaration de fonction
func addFuncSignature(signatures map[string]testedFuncInfo, funcDecl *ast.FuncDecl) {
	info := extractFuncInfo(funcDecl)

	// Skip mock receiver types
	if info.hasReceiver && shared.IsMockName(info.receiverName) {
		return
	}

	// Stocker avec le nom simple
	signatures[info.name] = *info
	// Stocker aussi avec Receiver_Method si méthode
	if info.hasReceiver {
		signatures[info.receiverName+"_"+info.name] = *info
	}
}

// collectExternalSourceSignatures scanne les fichiers sources du répertoire.
//
// Params:
//   - pass: contexte d'analyse
//   - signatures: map des signatures à remplir
func collectExternalSourceSignatures(pass *analysis.Pass, signatures map[string]testedFuncInfo) {
	// Pas de fichiers dans le pass
	if len(pass.Files) == 0 {
		return
	}

	// Obtenir le répertoire du package
	firstFile := pass.Fset.Position(pass.Files[0].Pos()).Filename
	packageDir := filepath.Dir(firstFile)

	// Lire les fichiers du répertoire
	entries, err := os.ReadDir(packageDir)
	// Erreur de lecture
	if err != nil {
		return
	}

	// Parcourir les fichiers
	for _, entry := range entries {
		// Ignorer les répertoires
		if entry.IsDir() {
			continue
		}
		// Scanner les fichiers Go non-test
		scanSourceFile(packageDir, entry.Name(), signatures)
	}
}

// scanSourceFile scanne un fichier source pour les signatures.
//
// Params:
//   - dir: répertoire du fichier
//   - filename: nom du fichier
//   - signatures: map des signatures
func scanSourceFile(dir string, filename string, signatures map[string]testedFuncInfo) {
	// Ignorer les fichiers non-Go
	if !strings.HasSuffix(filename, ".go") {
		return
	}
	// Ignorer les fichiers de test
	if strings.HasSuffix(filename, "_test.go") {
		return
	}
	// Skip mock files
	if shared.IsMockFile(filename) {
		return
	}

	// Parser le fichier
	fullPath := filepath.Join(dir, filename)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fullPath, nil, 0)
	// Erreur de parsing
	if err != nil {
		return
	}

	// Extraire les signatures
	ast.Inspect(file, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		// Pas une fonction
		if !ok {
			return true
		}
		// Skip mock functions
		if funcDecl.Name != nil && shared.IsMockName(funcDecl.Name.Name) {
			return true
		}
		// Ajouter la signature
		addFuncSignature(signatures, funcDecl)
		return true
	})
}

// extractFuncInfo extrait les informations d'une fonction.
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - *testedFuncInfo: informations extraites
func extractFuncInfo(funcDecl *ast.FuncDecl) *testedFuncInfo {
	// Use shared helper to classify
	meta := shared.ClassifyFunc(funcDecl)

	// Créer les infos de base
	info := &testedFuncInfo{
		name:         meta.Name,
		returnsError: functionReturnsError(funcDecl),
		hasReceiver:  meta.Kind == shared.FUNC_METHOD,
		receiverName: meta.ReceiverName,
	}

	// Retour des infos
	return info
}

// functionReturnsError vérifie si une fonction retourne error.
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - bool: true si retourne error
func functionReturnsError(funcDecl *ast.FuncDecl) bool {
	// Pas de type de fonction
	if funcDecl.Type == nil {
		return false
	}
	// Pas de résultats
	if funcDecl.Type.Results == nil {
		return false
	}

	// Parcourir les types de retour
	for _, field := range funcDecl.Type.Results.List {
		// Vérifier si c'est "error"
		if isErrorType(field.Type) {
			// Type error trouvé
			return true
		}
	}

	// Aucun type error trouvé
	return false
}

// isErrorType vérifie si un type est error.
//
// Params:
//   - expr: expression du type
//
// Returns:
//   - bool: true si c'est error
func isErrorType(expr ast.Expr) bool {
	// Identifiant simple "error"
	if ident, ok := expr.(*ast.Ident); ok {
		// Comparer avec "error"
		return ident.Name == "error"
	}
	// Pas un type error
	return false
}

// analyzeTestFunction analyse une fonction de test.
//
// Params:
//   - pass: contexte d'analyse
//   - testFunc: fonction de test
//   - signatures: signatures des fonctions
func analyzeTestFunction(
	pass *analysis.Pass,
	testFunc *ast.FuncDecl,
	signatures map[string]testedFuncInfo,
) {
	// Use shared helper to parse test name
	target, ok := shared.ParseTestName(testFunc.Name.Name)
	// Parsing failed
	if !ok {
		return
	}

	// Build lookup key
	key := shared.BuildTestTargetKey(target)
	// Key empty
	if key == "" {
		return
	}

	// Chercher la fonction dans les signatures
	info, found := signatures[key]
	// Fonction non trouvée dans les signatures
	if !found {
		return
	}

	// Vérifier si la fonction retourne error
	if !info.returnsError {
		return
	}

	// Vérifier la couverture des cas d'erreur
	if !hasErrorCaseCoverage(testFunc) {
		pass.Reportf(
			testFunc.Pos(),
			"KTN-TEST-013: le test '%s' teste une fonction qui retourne error, "+
				"il devrait couvrir les cas d'erreur",
			testFunc.Name.Name,
		)
	}
}

// hasErrorCaseCoverage vérifie la couverture des cas d'erreur.
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - bool: true si cas d'erreur couverts
func hasErrorCaseCoverage(funcDecl *ast.FuncDecl) bool {
	found := false

	// Parcourir le corps de la fonction
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		// Vérifier les différents types de nœuds
		if checkErrorInNode(n) {
			found = true
		}
		// Continuer l'inspection
		return true
	})

	// Retour du résultat de l'analyse
	return found
}

// checkErrorInNode vérifie si un nœud contient un indicateur d'erreur.
//
// Params:
//   - n: nœud AST
//
// Returns:
//   - bool: true si indicateur trouvé
func checkErrorInNode(n ast.Node) bool {
	// Traiter selon le type de nœud
	switch node := n.(type) {
	// Composite literal avec cas d'erreur
	case *ast.CompositeLit:
		// Vérifier les cas de test d'erreur
		return hasErrorTestCases(node)
	// Identifiant indicateur d'erreur
	case *ast.Ident:
		// Vérifier le nom de l'identifiant
		return isErrorIndicatorName(node.Name)
	// String literal avec error/invalid/fail
	case *ast.BasicLit:
		// Vérifier le contenu du literal
		return checkErrorInBasicLit(node)
	}
	// Type de nœud non concerné
	return false
}

// checkErrorInBasicLit vérifie un literal string.
//
// Params:
//   - node: nœud BasicLit
//
// Returns:
//   - bool: true si indicateur trouvé
func checkErrorInBasicLit(node *ast.BasicLit) bool {
	// Vérifier si c'est une string
	if node.Kind.String() != "STRING" {
		return false
	}
	// Vérifier le contenu de la string
	value := strings.ToLower(node.Value)
	// Vérifier les mots-clés d'erreur
	return strings.Contains(value, "error") ||
		strings.Contains(value, "invalid") ||
		strings.Contains(value, "fail")
}

// hasErrorTestCases vérifie les cas d'erreur dans un composite.
//
// Params:
//   - lit: composite literal
//
// Returns:
//   - bool: true si cas d'erreur présents
func hasErrorTestCases(lit *ast.CompositeLit) bool {
	// Parcourir les éléments
	for _, elt := range lit.Elts {
		// Vérifier si c'est un KeyValueExpr
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			// Vérifier si c'est un cas d'erreur
			if checkErrorInKeyValue(kv) {
				// Cas d'erreur trouvé
				return true
			}
		}
	}
	// Pas de cas d'erreur
	return false
}

// checkErrorInKeyValue vérifie un KeyValueExpr.
//
// Params:
//   - kv: KeyValueExpr
//
// Returns:
//   - bool: true si indicateur trouvé
func checkErrorInKeyValue(kv *ast.KeyValueExpr) bool {
	// Vérifier si la clé est "name"
	ident, ok := kv.Key.(*ast.Ident)
	// Pas un identifiant ou pas "name"
	if !ok || ident.Name != "name" {
		return false
	}
	// Vérifier si la valeur est un BasicLit
	basic, basicOk := kv.Value.(*ast.BasicLit)
	// Pas un BasicLit
	if !basicOk {
		return false
	}
	// Vérifier les mots-clés d'erreur
	value := strings.ToLower(basic.Value)
	// Retour du résultat
	return strings.Contains(value, "error") ||
		strings.Contains(value, "invalid") ||
		strings.Contains(value, "fail")
}

// isErrorIndicatorName vérifie si un nom indique une erreur.
//
// Params:
//   - name: nom à vérifier
//
// Returns:
//   - bool: true si indicateur d'erreur
func isErrorIndicatorName(name string) bool {
	lowerName := strings.ToLower(name)
	indicators := []string{"err", "error", "invalid", "fail", "bad", "wrong"}

	// Parcourir les indicateurs
	for _, indicator := range indicators {
		// Vérifier si le nom contient l'indicateur
		if strings.Contains(lowerName, indicator) {
			// Indicateur trouvé
			return true
		}
	}
	// Aucun indicateur trouvé
	return false
}
