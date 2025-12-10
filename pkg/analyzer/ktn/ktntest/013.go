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
	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	ruleCodeTest013 string = "KTN-TEST-013"
	// initialMapCapacity est la capacité initiale des maps de signatures.
	initialMapCapacity int = 32
)

// testedFuncInfo contient les informations sur une fonction testée.
type testedFuncInfo struct {
	name         string
	returnsError bool
	hasReceiver  bool
	receiverName string
}

// Analyzer013 checks that tests cover error cases
var Analyzer013 *analysis.Analyzer = &analysis.Analyzer{
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
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeTest013) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Collecter toutes les fonctions du package avec leur signature
	funcSignatures := collectFuncSignatures(pass, insp)

	// Définir le filtre de nœuds
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Parcourir les fonctions de test
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Conversion en FuncDecl
		funcDecl := n.(*ast.FuncDecl)
		// Obtenir le chemin du fichier
		filename := pass.Fset.Position(funcDecl.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeTest013, filename) {
			// Fichier exclu
			return
		}

		// Vérification si fichier de test
		if !shared.IsTestFile(filename) {
			// Retour si pas fichier de test
			return
		}

		// Vérification fichier exempté
		if shared.IsExemptTestFile(filename) {
			// Retour si exempté
			return
		}

		// Vérification fichier mock
		if shared.IsMockFile(filename) {
			// Retour si mock
			return
		}

		// Vérification test unitaire
		if !shared.IsUnitTestFunction(funcDecl) {
			// Retour si pas unitaire
			return
		}

		// Vérification nom exempté
		if shared.IsExemptTestName(funcDecl.Name.Name) {
			// Retour si exempté
			return
		}

		// Analyser le test
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
	signatures := make(map[string]testedFuncInfo, initialMapCapacity)

	// Définir le filtre de nœuds
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Parcourir les fonctions
	insp.Preorder(nodeFilter, func(n ast.Node) {
		// Conversion en FuncDecl
		funcDecl := n.(*ast.FuncDecl)
		// Obtenir le chemin du fichier
		filename := pass.Fset.Position(funcDecl.Pos()).Filename

		// Vérification fichier de test
		if shared.IsTestFile(filename) {
			// Retour si test
			return
		}

		// Vérification fichier mock
		if shared.IsMockFile(filename) {
			// Retour si mock
			return
		}

		// Vérification fonction mock
		if funcDecl.Name != nil && shared.IsMockName(funcDecl.Name.Name) {
			// Retour si mock
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
	// Extraire les informations de la fonction
	info := extractFuncInfo(funcDecl)

	// Vérification receiver mock
	if info.hasReceiver && shared.IsMockName(info.receiverName) {
		// Retour si mock
		return
	}

	// Stocker nom simple
	signatures[info.name] = *info
	// Vérification si méthode
	if info.hasReceiver {
		// Stocker avec receiver
		signatures[info.receiverName+"_"+info.name] = *info
	}
}

// collectExternalSourceSignatures scanne les fichiers sources du répertoire.
//
// Params:
//   - pass: contexte d'analyse
//   - signatures: map des signatures à remplir
func collectExternalSourceSignatures(pass *analysis.Pass, signatures map[string]testedFuncInfo) {
	// Vérification fichiers présents
	if len(pass.Files) == 0 {
		// Retour si pas de fichiers
		return
	}

	// Obtenir le répertoire
	firstFile := pass.Fset.Position(pass.Files[0].Pos()).Filename
	packageDir := filepath.Dir(firstFile)

	// Lire le répertoire
	entries, err := os.ReadDir(packageDir)
	// Vérification erreur
	if err != nil {
		// Retour si erreur
		return
	}

	// Itération sur les fichiers
	for _, entry := range entries {
		// Vérification si répertoire
		if entry.IsDir() {
			// Continuer si répertoire
			continue
		}
		// Scanner le fichier
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
	// Vérification extension Go
	if !strings.HasSuffix(filename, ".go") {
		// Retour si pas Go
		return
	}
	// Vérification fichier test
	if strings.HasSuffix(filename, "_test.go") {
		// Retour si test
		return
	}
	// Vérification fichier mock
	if shared.IsMockFile(filename) {
		// Retour si mock
		return
	}

	// Parser le fichier
	fullPath := filepath.Join(dir, filename)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, fullPath, nil, 0)
	// Vérification erreur parsing
	if err != nil {
		// Retour si erreur
		return
	}

	// Parcourir l'AST
	ast.Inspect(file, func(n ast.Node) bool {
		// Conversion en FuncDecl
		funcDecl, ok := n.(*ast.FuncDecl)
		// Vérification fonction
		if !ok {
			// Continuer si pas fonction
			return true
		}
		// Vérification mock
		if funcDecl.Name != nil && shared.IsMockName(funcDecl.Name.Name) {
			// Continuer si mock
			return true
		}
		// Ajouter la signature
		addFuncSignature(signatures, funcDecl)
		// Continuer la traversée
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
		hasReceiver:  meta.Kind == shared.FuncMethod,
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
	// Vérification type fonction
	if funcDecl.Type == nil {
		// Retour si pas de type
		return false
	}
	// Vérification résultats
	if funcDecl.Type.Results == nil {
		// Retour si pas de résultats
		return false
	}

	// Itération sur les résultats
	for _, field := range funcDecl.Type.Results.List {
		// Vérification type error
		if isErrorType(field.Type) {
			// Retour trouvé
			return true
		}
	}

	// Retour non trouvé
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
	// Vérification identifiant
	if ident, ok := expr.(*ast.Ident); ok {
		// Retour vérification error
		return ident.Name == "error"
	}
	// Retour faux par défaut
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
	// Parser le nom du test
	target, ok := shared.ParseTestName(testFunc.Name.Name)
	// Vérification parsing réussi
	if !ok {
		// Retour si échec
		return
	}

	// Construire la clé
	key := shared.BuildTestTargetKey(target)
	// Vérification clé vide
	if key == "" {
		// Retour si vide
		return
	}

	// Chercher la fonction
	info, found := signatures[key]
	// Vérification fonction trouvée
	if !found {
		// Retour si pas trouvée
		return
	}

	// Vérification retourne error
	if !info.returnsError {
		// Retour si pas error
		return
	}

	// Vérifier couverture erreur
	if !hasErrorCaseCoverage(testFunc) {
		// Signaler manque couverture
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
	// Initialiser drapeau
	found := false

	// Parcourir le corps
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		// Vérifier erreur dans nœud
		if checkErrorInNode(n) {
			// Marquer comme trouvé
			found = true
		}
		// Continuer l'inspection
		return true
	})

	// Retour du résultat
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
	// Sélection selon le type
	switch node := n.(type) {
	// Vérification composite literal
	case *ast.CompositeLit:
		// Cas composite literal
		// Vérifier cas d'erreur
		return hasErrorTestCases(node)
	// Vérification identifiant
	case *ast.Ident:
		// Cas identifiant
		// Vérifier nom indicateur
		return isErrorIndicatorName(node.Name)
	// Vérification literal basique
	case *ast.BasicLit:
		// Cas literal basique
		// Vérifier contenu
		return checkErrorInBasicLit(node)
	}
	// Retour faux par défaut
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
	// Vérification type string
	if node.Kind.String() != "STRING" {
		// Retour si pas string
		return false
	}
	// Vérifier le contenu
	value := strings.ToLower(node.Value)
	// Retour vérification mots-clés
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
	// Itération sur les éléments
	for _, elt := range lit.Elts {
		// Vérification KeyValueExpr
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			// Vérifier erreur dans kv
			if checkErrorInKeyValue(kv) {
				// Retour trouvé
				return true
			}
		}
	}
	// Retour non trouvé
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
	// Vérifier la clé
	ident, ok := kv.Key.(*ast.Ident)
	// Vérification clé name
	if !ok || ident.Name != "name" {
		// Retour si pas name
		return false
	}
	// Vérifier la valeur
	basic, basicOk := kv.Value.(*ast.BasicLit)
	// Vérification BasicLit
	if !basicOk {
		// Retour si pas BasicLit
		return false
	}
	// Vérifier mots-clés
	value := strings.ToLower(basic.Value)
	// Retour vérification
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
	// Convertir en minuscules
	lowerName := strings.ToLower(name)
	// Définir les indicateurs d'erreur
	indicators := []string{"err", "error", "invalid", "fail", "bad", "wrong"}

	// Itération sur les indicateurs
	for _, indicator := range indicators {
		// Vérification si contient
		if strings.Contains(lowerName, indicator) {
			// Retour trouvé
			return true
		}
	}
	// Retour non trouvé
	return false
}
