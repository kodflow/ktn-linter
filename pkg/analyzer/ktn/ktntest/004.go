// Analyzer 004 for the ktntest package.
package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

const (
	// MIN_PUBLIC_FUNCS est le nombre minimum de fonctions publiques
	MIN_PUBLIC_FUNCS int = 1
	// MAX_TEST_NAMES maximum number of test names to check per function
	MAX_TEST_NAMES int = 2
)

// funcInfo stores information about a function (public or private)
type funcInfo struct {
	name         string
	receiverName string // Nom du receiver pour les méthodes (vide pour les fonctions)
	isExported   bool   // true si fonction publique, false si privée
	pos          token.Pos
	filename     string
}

// Analyzer004 checks that all functions (public and private) have corresponding tests
var Analyzer004 = &analysis.Analyzer{
	Name: "ktntest004",
	Doc:  "KTN-TEST-004: Toutes les fonctions (publiques et privées) doivent avoir des tests",
	Run:  runTest004,
}

// collectFunctions collecte toutes les fonctions (publiques et privées) et les fonctions testées.
//
// Params:
//   - pass: contexte d'analyse
//   - funcs: pointeur vers slice de toutes les fonctions
//   - testedFuncs: map des fonctions testées
func collectFunctions(pass *analysis.Pass, funcs *[]funcInfo, testedFuncs map[string]bool) {
	// Parcourir tous les fichiers du pass
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename

		// Skip mock files
		if shared.IsMockFile(filename) {
			// Mock file, skip
			continue
		}

		// Parcourir l'AST du fichier pour trouver les fonctions
		ast.Inspect(file, func(n ast.Node) bool {
			// Vérifier si c'est une déclaration de fonction
			funcDecl, ok := n.(*ast.FuncDecl)
			// Si ce n'est pas une fonction, continuer
			if !ok {
				// Continue traversal
				return true
			}

			// Vérification de la condition
			if shared.IsTestFile(filename) {
				// Fichier de test - collecter les fonctions testées
				// On utilise collectExternalTestedFunctions ici aussi car
				// les fichiers _internal_test.go sont dans pass.Files mais
				// shared.IsUnitTestFunction peut échouer sans TypesInfo complet
				collectExternalTestedFunctions(funcDecl, testedFuncs)
			} else {
				// Fichier source - collecter TOUTES les fonctions (publiques et privées)
				meta := shared.ClassifyFunc(funcDecl)
				*funcs = append(*funcs, funcInfo{
					name:         meta.Name,
					receiverName: meta.ReceiverName,
					isExported:   meta.Visibility == shared.VIS_PUBLIC,
					pos:          funcDecl.Pos(),
					filename:     filename,
				})
			}
			// Continue traversal
			return true
		})
	}
}

// runTest004 exécute l'analyse KTN-TEST-004.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest004(pass *analysis.Pass) (any, error) {
	// Vérifier s'il y a des fichiers de test dans ce pass
	hasTestFiles, testFileCount := countTestFiles(pass)

	// Si pas de fichiers de test ou que des fichiers de test, skip
	if !hasTestFiles || testFileCount == len(pass.Files) {
		// Early return from function
		return nil, nil
	}

	// Collecter toutes les fonctions et les tests
	allFuncs, testedFuncs := collectAllFunctionsAndTests(pass)

	// Vérifier que chaque fonction a un test
	checkFunctionsHaveTests(pass, allFuncs, testedFuncs)

	// Retour de la fonction
	return nil, nil
}

// countTestFiles compte les fichiers de test dans le pass.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si des fichiers de test existent
//   - int: nombre de fichiers de test
func countTestFiles(pass *analysis.Pass) (bool, int) {
	hasTestFiles := false
	testFileCount := 0
	// Parcourir les fichiers pour compter les tests
	for _, file := range pass.Files {
		pos := pass.Fset.Position(file.Pos())
		// Vérification de la condition
		if shared.IsTestFile(pos.Filename) {
			hasTestFiles = true
			testFileCount++
		}
	}
	// Retour des compteurs
	return hasTestFiles, testFileCount
}

// collectAllFunctionsAndTests collecte les fonctions et les tests.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - []funcInfo: liste des fonctions
//   - map[string]bool: map des fonctions testées
func collectAllFunctionsAndTests(pass *analysis.Pass) ([]funcInfo, map[string]bool) {
	var allFuncs []funcInfo
	testedFuncs := make(map[string]bool, 0)

	// Collecter toutes les fonctions et les tests
	collectFunctions(pass, &allFuncs, testedFuncs)

	// Scanner les packages _test externes
	collectExternalTestFunctions(pass, testedFuncs)

	// Retour des collections
	return allFuncs, testedFuncs
}

// checkFunctionsHaveTests vérifie que chaque fonction a un test.
//
// Params:
//   - pass: contexte d'analyse
//   - allFuncs: liste des fonctions
//   - testedFuncs: map des fonctions testées
func checkFunctionsHaveTests(pass *analysis.Pass, allFuncs []funcInfo, testedFuncs map[string]bool) {
	// Build lookup map for functions
	funcLookup := make(map[string]bool)
	// Populate lookup map
	for _, fn := range allFuncs {
		key := buildFuncLookupKey(fn)
		funcLookup[key] = true
	}

	// Pré-allouer le slice pour les noms de test
	testNames := make([]string, 0, MAX_TEST_NAMES)

	// Vérifier chaque fonction
	for _, fn := range allFuncs {
		// Skip mocks
		if shared.IsMockName(fn.name) {
			// Mock function, skip
			continue
		}
		// Skip if receiver is a mock
		if fn.receiverName != "" && shared.IsMockName(fn.receiverName) {
			// Mock method, skip
			continue
		}

		// Réinitialiser et construire les noms de test possibles
		testNames = buildTestNames(testNames[:0], fn)

		// Vérifier si au moins un des noms possibles a un test
		if !hasMatchingTest(testNames, testedFuncs) && !isExemptFunction(fn.name) {
			// Report missing test
			reportMissingTest(pass, fn)
		}
	}
}

// buildFuncLookupKey builds the lookup key for a function.
//
// Params:
//   - fn: function info
//
// Returns:
//   - string: lookup key
func buildFuncLookupKey(fn funcInfo) string {
	// Handle methods
	if fn.receiverName != "" {
		// Method: ReceiverName_MethodName
		return fn.receiverName + "_" + fn.name
	}
	// Top-level function
	return fn.name
}

// buildTestNames construit les noms de test possibles pour une fonction.
//
// Params:
//   - testNames: slice à remplir (doit être vide)
//   - fn: information sur la fonction
//
// Returns:
//   - []string: noms de test possibles
func buildTestNames(testNames []string, fn funcInfo) []string {
	// Use shared helper to build lookup key
	meta := &shared.FuncMeta{
		Name:         fn.name,
		ReceiverName: fn.receiverName,
	}
	// Set kind based on receiver
	if fn.receiverName != "" {
		// Method
		meta.Kind = shared.FUNC_METHOD
	} else {
		// Top-level function
		meta.Kind = shared.FUNC_TOP_LEVEL
	}
	// Add lookup key
	testNames = append(testNames, shared.BuildTestLookupKey(meta))
	// Retour des noms
	return testNames
}

// hasMatchingTest vérifie si un test existe pour les noms donnés.
//
// Params:
//   - testNames: noms de test à chercher
//   - testedFuncs: map des fonctions testées
//
// Returns:
//   - bool: true si un test existe
func hasMatchingTest(testNames []string, testedFuncs map[string]bool) bool {
	// Parcours des noms de test possibles
	for _, testName := range testNames {
		// Vérification de la condition (case-sensitive)
		if testedFuncs[testName] {
			// Test trouvé
			return true
		}
		// Vérification case-insensitive pour les méthodes
		for testedName := range testedFuncs {
			// Comparaison insensible à la casse
			if strings.EqualFold(testName, testedName) {
				// Test trouvé
				return true
			}
		}
	}
	// Pas de test trouvé
	return false
}

// reportMissingTest reporte une fonction sans test.
//
// Params:
//   - pass: contexte d'analyse
//   - fn: information sur la fonction
func reportMissingTest(pass *analysis.Pass, fn funcInfo) {
	// Build meta for suggested test name
	meta := &shared.FuncMeta{
		Name:         fn.name,
		ReceiverName: fn.receiverName,
	}
	// Set kind and visibility
	if fn.receiverName != "" {
		// Method
		meta.Kind = shared.FUNC_METHOD
	} else {
		// Top-level function
		meta.Kind = shared.FUNC_TOP_LEVEL
	}
	// Set visibility
	if fn.isExported {
		// Public
		meta.Visibility = shared.VIS_PUBLIC
	} else {
		// Private
		meta.Visibility = shared.VIS_PRIVATE
	}

	// Use shared helper for suggested name
	suggestedTestName := shared.BuildSuggestedTestName(meta)

	// Extraire le nom de base du fichier
	baseName := filepath.Base(fn.filename)
	fileBase := strings.TrimSuffix(baseName, ".go")

	// Déterminer le fichier et le type de test
	suggestedTestFile, testType, funcType := getTestFileInfo(fn.isExported, fileBase)

	// Reporter la fonction non testée
	pass.Reportf(
		fn.pos,
		"KTN-TEST-004: fonction %s '%s' n'a pas de test correspondant. Créer un test nommé '%s' dans le fichier '%s' (%s)",
		funcType, fn.name, suggestedTestName, suggestedTestFile, testType,
	)
}

// getTestFileInfo retourne les informations sur le fichier de test.
//
// Params:
//   - isExported: true si fonction publique
//   - fileBase: nom de base du fichier
//
// Returns:
//   - string: nom du fichier de test suggéré
//   - string: type de test
//   - string: type de fonction
func getTestFileInfo(isExported bool, fileBase string) (string, string, string) {
	// Vérification du type de fonction
	if isExported {
		// Fonction publique → test black-box
		return fileBase + "_external_test.go", "black-box testing avec package xxx_test", "publique"
	}
	// Fonction privée → test white-box
	return fileBase + "_internal_test.go", "white-box testing avec package xxx", "privée"
}

// findPackageDir trouve le répertoire du package à partir des fichiers du pass.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - string: chemin du répertoire ou chaîne vide si non trouvé
func findPackageDir(pass *analysis.Pass) string {
	// Parcourir les fichiers pour trouver un chemin valide
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		// Ignorer les fichiers du cache Go et les fichiers temporaires
		if isCacheOrTempFile(filename) {
			// Cache file, skip
			continue
		}
		// Utiliser le premier fichier source valide
		return filepath.Dir(filename)
	}
	// Pas de fichier source valide
	return ""
}

// isCacheOrTempFile vérifie si un fichier est dans le cache Go ou temporaire.
//
// Params:
//   - filename: chemin du fichier
//
// Returns:
//   - bool: true si c'est un fichier cache/temp
func isCacheOrTempFile(filename string) bool {
	// Vérifier les patterns de cache et temp
	return strings.Contains(filename, "/.cache/go-build/") ||
		strings.Contains(filename, "/tmp/") ||
		strings.Contains(filename, "\\cache\\go-build\\")
}

// parseTestFile parse un fichier de test et extrait les fonctions testées.
//
// Params:
//   - fullPath: chemin complet du fichier
//   - testedFuncs: map des fonctions testées à remplir
func parseTestFile(fullPath string, testedFuncs map[string]bool) {
	// Parser le fichier pour extraire les tests
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fullPath, nil, 0)
	// Si erreur de parsing, retourner
	if err != nil {
		// Parsing failed
		return
	}

	// Extraire les fonctions de test du fichier parsé
	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		// Si ce n'est pas une fonction, continuer
		if !ok {
			// Continue traversal
			return true
		}
		// Collecter les fonctions testées
		collectExternalTestedFunctions(funcDecl, testedFuncs)
		// Continue traversal
		return true
	})
}

// collectExternalTestFunctions scanne les packages _test externes.
// Cette fonction détecte les fichiers _external_test.go même s'ils utilisent
// package xxx_test et sont dans XTestGoFiles (package séparé de Go).
//
// Params:
//   - pass: contexte d'analyse
//   - testedFuncs: map des fonctions testées à remplir
func collectExternalTestFunctions(pass *analysis.Pass, testedFuncs map[string]bool) {
	// Extraire le répertoire du premier fichier du package
	if len(pass.Files) == 0 {
		// Pas de fichiers, rien à faire
		return
	}

	// Obtenir le répertoire du package
	packageDir := findPackageDir(pass)
	// Si pas de répertoire valide trouvé, retourner
	if packageDir == "" {
		// Pas de fichier source valide
		return
	}

	// Scanner tous les fichiers du répertoire
	entries, err := os.ReadDir(packageDir)
	// Si erreur de lecture, retourner sans erreur
	if err != nil {
		// Retour silencieux
		return
	}

	// Parcourir tous les fichiers de test
	for _, entry := range entries {
		// Ignorer les répertoires et non-test files
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			// Continuer avec le prochain
			continue
		}
		// Parser le fichier de test
		parseTestFile(filepath.Join(packageDir, entry.Name()), testedFuncs)
	}
}

// collectExternalTestedFunctions collecte les fonctions testées à partir
// d'un fichier _test.go parsé en dehors du Pass (sans TypesInfo).
// On ne peut pas utiliser shared.IsUnitTestFunction ici, on se base donc
// uniquement sur une heuristique syntaxique.
//
// Params:
//   - funcDecl: déclaration de fonction de test
//   - testedFuncs: map des fonctions testées
func collectExternalTestedFunctions(funcDecl *ast.FuncDecl, testedFuncs map[string]bool) {
	// Nom de test Go standard: doit commencer par "Test"
	if !strings.HasPrefix(funcDecl.Name.Name, "Test") {
		// Pas une fonction de test
		return
	}

	// Vérification minimale de la signature (func TestXxx(t *testing.T))
	if funcDecl.Type == nil || funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 {
		// Pas de paramètres
		return
	}

	// Vérifier que le premier paramètre est "t"
	firstParam := funcDecl.Type.Params.List[0]
	// Check if param has a name
	if len(firstParam.Names) == 0 || firstParam.Names[0].Name != "t" {
		// Premier paramètre n'est pas "t", probablement pas un test standard
		return
	}

	// Skip exempt test names (helper, main, etc.)
	if shared.IsExemptTestName(funcDecl.Name.Name) {
		// Exempt test name
		return
	}

	// On réutilise la même logique de parsing de nom et de clé cible
	target, ok := shared.ParseTestName(funcDecl.Name.Name)
	// Check if parse succeeded
	if !ok {
		// Could not parse
		return
	}

	// Build lookup key from target
	key := shared.BuildTestTargetKey(target)
	// Add to tested functions
	if key != "" {
		// Add key
		testedFuncs[key] = true
	}
}

// isExemptFunction vérifie si une fonction est exemptée.
//
// Params:
//   - funcName: nom de la fonction
//
// Returns:
//   - bool: true si la fonction est exemptée
func isExemptFunction(funcName string) bool {
	// Fonctions exemptées (uniquement init et main)
	exemptFuncs := []string{
		"init",
		"main",
	}

	// Vérifier si la fonction est exemptée
	return slices.Contains(exemptFuncs, funcName)
}
