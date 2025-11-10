package ktntest

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

const (
	// MIN_PUBLIC_FUNCS est le nombre minimum de fonctions publiques
	MIN_PUBLIC_FUNCS int = 1
)

// funcInfo stores information about a function (public or private)
type funcInfo struct {
	name         string
	receiverName string // Nom du receiver pour les méthodes (vide pour les fonctions)
	isExported   bool   // true si fonction publique, false si privée
	pos          token.Pos
	filename     string
}

// Analyzer003 checks that all functions (public and private) have corresponding tests
var Analyzer003 *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktntest003",
	Doc:  "KTN-TEST-003: Toutes les fonctions (publiques et privées) doivent avoir des tests",
	Run:  runTest003,
}

// collectFunctions collecte toutes les fonctions (publiques et privées) et les fonctions testées.
//
// Params:
//   - pass: contexte d'analyse
//   - funcs: pointeur vers slice de toutes les fonctions
//   - testedFuncs: map des fonctions testées
func collectFunctions(pass *analysis.Pass, funcs *[]funcInfo, testedFuncs map[string]bool) {
	// Parcourir tous les fichiers du pass
	// Pour chaque fichier, collecter les fonctions publiques et les tests
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename

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
				collectTestedFunctions(funcDecl, testedFuncs)
			} else {
				// Fichier source - collecter TOUTES les fonctions (publiques et privées)
				receiverName := ""
				// Vérification si c'est une méthode avec un receiver
				if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
					// Extraire le nom du type du receiver
					receiverName = extractReceiverTypeName(funcDecl.Recv.List[0].Type)
				}
				*funcs = append(*funcs, funcInfo{
					name:         funcDecl.Name.Name,
					receiverName: receiverName,
					isExported:   isPublicFunction(funcDecl),
					pos:          funcDecl.Pos(),
					filename:     filename,
				})
			}
			// Continue traversal
			return true
		})
	}
}

// runTest003 exécute l'analyse KTN-TEST-003.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest003(pass *analysis.Pass) (any, error) {
	// Vérifier s'il y a des fichiers de test dans ce pass
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

	// Si pas de fichiers de test, skip cette analyse
	// (cela arrive quand le package est analysé sans son variant de test)
	if !hasTestFiles {
		// Early return from function.
		return nil, nil
	}

	// If we only have test files, this is probably a separate test package - skip it
	if testFileCount == len(pass.Files) {
		// Early return from function.
		return nil, nil
	}

	// Liste de toutes les fonctions (publiques et privées) avec leurs positions
	var allFuncs []funcInfo
	// Map des fonctions testées
	testedFuncs := make(map[string]bool, 0)

	// Collecter toutes les fonctions et les tests
	collectFunctions(pass, &allFuncs, testedFuncs)

	// Toujours scanner les packages _test externes (XTestGoFiles)
	// pour détecter tous les tests, même si on a des tests internes
	collectExternalTestFunctions(pass, testedFuncs)

	// Vérifier que chaque fonction (publique ou privée) a un test
	for _, fn := range allFuncs {
		// Construire les noms possibles pour le test
		testNames := []string{fn.name}
		// Si c'est une méthode, ajouter aussi le pattern Receiver_Method
		if fn.receiverName != "" {
			testNames = append(testNames, fn.receiverName+"_"+fn.name)
		}

		// Vérifier si au moins un des noms possibles a un test
		hasTest := false
		// Parcours des noms de test possibles
		for _, testName := range testNames {
			// Vérification de la condition
			if testedFuncs[testName] {
				hasTest = true
				// Test trouvé
				break
			}
		}

		// Vérification de la condition
		if !hasTest && !isExemptFunction(fn.name) {
			// Construire le nom du test suggéré
			suggestedTestName := "Test" + fn.name
			// Si c'est une méthode, suggérer le pattern Type_Method
			if fn.receiverName != "" {
				suggestedTestName = "Test" + fn.receiverName + "_" + fn.name
			}

			// Extraire le nom de base du fichier pour construire le nom du fichier de test
			baseName := filepath.Base(fn.filename)
			fileBase := strings.TrimSuffix(baseName, ".go")

			// Déterminer le fichier et le type de test selon si la fonction est exportée
			var suggestedTestFile string
			var testType string
			var funcType string
			// Vérification du type de fonction
			if fn.isExported {
				// Fonction publique → test black-box dans _external_test.go
				suggestedTestFile = fileBase + "_external_test.go"
				testType = "black-box testing avec package xxx_test"
				funcType = "publique"
			} else {
				// Fonction privée → test white-box dans _internal_test.go
				suggestedTestFile = fileBase + "_internal_test.go"
				testType = "white-box testing avec package xxx"
				funcType = "privée"
			}

			// Fonction non testée - reporter à la position de la fonction
			pass.Reportf(
				fn.pos,
				"KTN-TEST-003: fonction %s '%s' n'a pas de test correspondant. Créer un test nommé '%s' dans le fichier '%s' (%s)",
				funcType,
				fn.name,
				suggestedTestName,
				suggestedTestFile,
				testType,
			)
		}
	}

	// Retour de la fonction
	return nil, nil
}

// isPublicFunction vérifie si une fonction est publique.
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - bool: true si la fonction est publique
func isPublicFunction(funcDecl *ast.FuncDecl) bool {
	// Vérification du nom
	if funcDecl.Name == nil {
		// Pas de nom
		return false
	}

	name := funcDecl.Name.Name
	// Vérifier si le premier caractère est en majuscule
	if len(name) == 0 {
		// Nom vide
		return false
	}

	// Retour du résultat
	return unicode.IsUpper(rune(name[0]))
}

// collectTestedFunctions collecte les fonctions testées.
//
// Params:
//   - funcDecl: déclaration de fonction de test
//   - testedFuncs: map des fonctions testées
func collectTestedFunctions(funcDecl *ast.FuncDecl, testedFuncs map[string]bool) {
	// Vérifier si c'est une fonction de test unitaire (Test*)
	if !shared.IsUnitTestFunction(funcDecl) {
		// Pas une fonction de test unitaire
		return
	}

	// Extraire le nom de la fonction testée
	testedName := strings.TrimPrefix(funcDecl.Name.Name, "Test")
	// Vérification de la condition
	if testedName != "" {
		// Ajouter à la map
		testedFuncs[testedName] = true
	}
}

// extractReceiverTypeName extrait le nom du type du receiver.
//
// Params:
//   - expr: expression du type du receiver
//
// Returns:
//   - string: nom du type (sans * pour les pointeurs)
func extractReceiverTypeName(expr ast.Expr) string {
	// Gérer les pointeurs (*Type)
	if starExpr, ok := expr.(*ast.StarExpr); ok {
		// Expression pointeur
		return extractReceiverTypeName(starExpr.X)
	}

	// Gérer les identifiants simples (Type)
	if ident, ok := expr.(*ast.Ident); ok {
		// Retour du nom du type
		return ident.Name
	}

	// Type non géré
	return ""
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
	firstFile := pass.Fset.Position(pass.Files[0].Pos()).Filename
	packageDir := filepath.Dir(firstFile)

	// Scanner tous les fichiers du répertoire
	entries, err := os.ReadDir(packageDir)
	// Si erreur de lecture, retourner sans erreur
	if err != nil {
		// Retour silencieux
		return
	}

	// Parcourir tous les fichiers
	for _, entry := range entries {
		// Ignorer les répertoires
		if entry.IsDir() {
			// Continuer avec le prochain
			continue
		}

		filename := entry.Name()
		// Ne garder que les fichiers de test (*_test.go)
		if !strings.HasSuffix(filename, "_test.go") {
			// Continuer avec le prochain
			continue
		}

		// Construire le chemin complet
		fullPath := filepath.Join(packageDir, filename)

		// Parser le fichier pour extraire les tests
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, fullPath, nil, 0)
		// Si erreur de parsing, continuer avec le prochain fichier
		if err != nil {
			// Continuer avec le prochain
			continue
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
			collectTestedFunctions(funcDecl, testedFuncs)
			// Continue traversal
			return true
		})
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

	// Parcours des fonctions exemptées
	for _, exempt := range exemptFuncs {
		// Vérification de la condition
		if funcName == exempt {
			// Fonction exemptée
			return true
		}
	}

	// Fonction non exemptée
	return false
}
