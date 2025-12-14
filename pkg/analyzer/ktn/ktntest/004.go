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
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
)

const (
	// ruleCode est le code de la règle.
	ruleCodeTest004 string = "KTN-TEST-004"
	// minPublicFuncs est le nombre minimum de fonctions publiques
	minPublicFuncs int = 1
	// maxTestNames maximum number of test names to check per function
	maxTestNames int = 2
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
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
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
	// Récupération de la configuration
	cfg := config.Get()

	// Parcourir tous les fichiers du pass
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeTest004, filename) {
			// Fichier exclu
			continue
		}

		// Skip mock files
		if shared.IsMockFile(filename) {
			// Mock file, skip
			continue
		}

		// Parcourir l'AST du fichier
		ast.Inspect(file, func(n ast.Node) bool {
			funcDecl, ok := n.(*ast.FuncDecl)
			// Vérification si déclaration de fonction
			if !ok {
				// Continuer la traversée
				return true
			}

			// Vérification si fichier de test
			if shared.IsTestFile(filename) {
				// Collecter les fonctions testées
				collectExternalTestedFunctions(funcDecl, testedFuncs)
			} else {
				// Cas alternatif: fichier source
				// Collecter toutes les fonctions
				meta := shared.ClassifyFunc(funcDecl)
				*funcs = append(*funcs, funcInfo{
					name:         meta.Name,
					receiverName: meta.ReceiverName,
					isExported:   meta.Visibility == shared.VisPublic,
					pos:          funcDecl.Pos(),
					filename:     filename,
				})
			}
			// Continuer la traversée
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
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeTest004) {
		// Règle désactivée
		return nil, nil
	}

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
	// Parcourir les fichiers du pass
	for _, file := range pass.Files {
		pos := pass.Fset.Position(file.Pos())
		// Vérification si fichier de test
		if shared.IsTestFile(pos.Filename) {
			// Incrémenter les compteurs
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
	// Construire la map de recherche
	funcLookup := make(map[string]bool, len(allFuncs))
	// Parcourir les fonctions
	for _, fn := range allFuncs {
		key := buildFuncLookupKey(fn)
		funcLookup[key] = true
	}

	// Pré-allouer le slice pour les noms de test
	testNames := make([]string, 0, maxTestNames)

	// Parcourir toutes les fonctions
	for _, fn := range allFuncs {
		// Vérification si mock
		if shared.IsMockName(fn.name) {
			// Ignorer les mocks
			continue
		}
		// Vérification si receiver mock
		if fn.receiverName != "" && shared.IsMockName(fn.receiverName) {
			// Ignorer les méthodes mock
			continue
		}

		// Construire les noms de test possibles
		testNames = buildTestNames(testNames[:0], fn)

		// Vérification si test manquant
		if !hasMatchingTest(testNames, testedFuncs) && !isExemptFunction(fn.name) {
			// Signaler le test manquant
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
	// Vérification si méthode
	if fn.receiverName != "" {
		// Retour clé avec receiver
		return fn.receiverName + "_" + fn.name
	}
	// Retour clé simple
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
	// Vérification si méthode
	if fn.receiverName != "" {
		// Définir comme méthode
		meta.Kind = shared.FuncMethod
	} else {
		// Cas alternatif: fonction
		// Définir comme fonction
		meta.Kind = shared.FuncTopLevel
	}
	// Ajouter la clé de recherche
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
	// Itération sur les noms de test
	for _, testName := range testNames {
		// Vérification case-sensitive
		if testedFuncs[testName] {
			// Retour trouvé
			return true
		}
		// Itération sur les tests enregistrés
		for testedName := range testedFuncs {
			// Vérification case-insensitive
			if strings.EqualFold(testName, testedName) {
				// Retour trouvé
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
	// Vérification si méthode
	if fn.receiverName != "" {
		// Définir comme méthode
		meta.Kind = shared.FuncMethod
	} else {
		// Cas alternatif: fonction
		// Définir comme fonction
		meta.Kind = shared.FuncTopLevel
	}
	// Vérification si publique
	if fn.isExported {
		// Définir comme publique
		meta.Visibility = shared.VisPublic
	} else {
		// Cas alternatif: privée
		// Définir comme privée
		meta.Visibility = shared.VisPrivate
	}

	// Reporter la fonction non testée
	msg, _ := messages.Get(ruleCodeTest004)
	pass.Reportf(
		fn.pos,
		"%s: %s",
		ruleCodeTest004,
		msg.Format(config.Get().Verbose, fn.name),
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
	// Vérification si publique
	if isExported {
		// Retour info test external
		return fileBase + "_external_test.go", "black-box testing avec package xxx_test", "publique"
	}
	// Retour info test internal
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
	// Itération sur les fichiers
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		// Vérification si fichier cache
		if isCacheOrTempFile(filename) {
			// Ignorer fichier cache
			continue
		}
		// Retour du répertoire
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
	// Parser le fichier
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fullPath, nil, 0)
	// Vérification erreur parsing
	if err != nil {
		// Retour en cas d'erreur
		return
	}

	// Parcourir l'AST du fichier
	ast.Inspect(node, func(n ast.Node) bool {
		funcDecl, ok := n.(*ast.FuncDecl)
		// Vérification si fonction
		if !ok {
			// Continuer la traversée
			return true
		}
		// Collecter les fonctions testées
		collectExternalTestedFunctions(funcDecl, testedFuncs)
		// Continuer la traversée
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
	// Vérification si fichiers présents
	if len(pass.Files) == 0 {
		// Retour si pas de fichiers
		return
	}

	// Obtenir le répertoire
	packageDir := findPackageDir(pass)
	// Vérification répertoire valide
	if packageDir == "" {
		// Retour si pas de répertoire
		return
	}

	// Lire les entrées du répertoire
	entries, err := os.ReadDir(packageDir)
	// Vérification erreur lecture
	if err != nil {
		// Retour en cas d'erreur
		return
	}

	// Itération sur les entrées
	for _, entry := range entries {
		// Vérification si répertoire ou non-test
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), "_test.go") {
			// Continuer l'itération
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
	// Vérification si fonction Test
	if !strings.HasPrefix(funcDecl.Name.Name, "Test") {
		// Retour si pas Test
		return
	}

	// Vérification de la signature
	if funcDecl.Type == nil || funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) == 0 {
		// Retour si pas de paramètres
		return
	}

	// Vérifier le premier paramètre
	firstParam := funcDecl.Type.Params.List[0]
	// Vérification du nom du paramètre
	if len(firstParam.Names) == 0 || firstParam.Names[0].Name != "t" {
		// Retour si pas "t"
		return
	}

	// Vérification si test exempté
	if shared.IsExemptTestName(funcDecl.Name.Name) {
		// Retour si exempté
		return
	}

	// Parser le nom du test
	target, ok := shared.ParseTestName(funcDecl.Name.Name)
	// Vérification parsing réussi
	if !ok {
		// Retour si échec parsing
		return
	}

	// Construire la clé de recherche
	key := shared.BuildTestTargetKey(target)
	// Vérification clé valide
	if key != "" {
		// Ajouter la clé
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
