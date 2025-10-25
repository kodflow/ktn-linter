package ktntest

import (
	"go/ast"
	"go/token"
	"strings"
	"unicode"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

const (
	// MIN_PUBLIC_FUNCS est le nombre minimum de fonctions publiques
	MIN_PUBLIC_FUNCS int = 1
)

// publicFuncInfo stores information about a public function
type publicFuncInfo struct {
	name     string
	pos      token.Pos
	filename string
}

// Analyzer003 checks that public functions have corresponding tests
var Analyzer003 *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktntest003",
	Doc:  "KTN-TEST-003: Toutes les fonctions publiques doivent avoir des tests",
	Run:  runTest003,
}

// collectFunctions collecte les fonctions publiques et testées.
//
// Params:
//   - pass: contexte d'analyse
//   - publicFuncs: pointeur vers slice de fonctions publiques
//   - testedFuncs: map des fonctions testées
func collectFunctions(pass *analysis.Pass, publicFuncs *[]publicFuncInfo, testedFuncs map[string]bool) {
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
				// Fichier source - collecter les fonctions publiques
				if isPublicFunction(funcDecl) {
					*publicFuncs = append(*publicFuncs, publicFuncInfo{
						name:     funcDecl.Name.Name,
						pos:      funcDecl.Pos(),
						filename: filename,
					})
				}
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

	// Analyze only packages with both source and test files
	// If we only have test files, this is probably a separate test package - skip it
	if testFileCount == len(pass.Files) {
		// Early return from function.
		return nil, nil
	}

	// Liste des fonctions publiques avec leurs positions
	var publicFuncs []publicFuncInfo
	// Map des fonctions testées
	testedFuncs := make(map[string]bool, 0)

	// Collecter les fonctions publiques et testées
	collectFunctions(pass, &publicFuncs, testedFuncs)

	// Vérifier que chaque fonction publique a un test
	for _, funcInfo := range publicFuncs {
		// Vérification de la condition
		if !testedFuncs[funcInfo.name] && !isExemptFunction(funcInfo.name) {
			// Fonction non testée - reporter à la position de la fonction
			pass.Reportf(
				funcInfo.pos,
				"KTN-TEST-003: fonction publique '%s' dans '%s' n'a pas de test correspondant",
				funcInfo.name,
				funcInfo.filename,
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
