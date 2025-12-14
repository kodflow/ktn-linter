// Package ktntest implements KTN linter rules.
package ktntest

import (
	"go/ast"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
)

const (
	// ruleCode est le code de la règle.
	ruleCodeTest008 string = "KTN-TEST-008"
	// maxFuncsToDisplay nombre max de fonctions affichées.
	maxFuncsToDisplay int = 3
	// initialFuncsCap capacité initiale des listes.
	initialFuncsCap int = 8
)

// Analyzer008 checks that each source file has appropriate test files based on its content
var Analyzer008 *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktntest008",
	Doc:  "KTN-TEST-008: Chaque fichier .go doit avoir les fichiers de test appropriés (xxx_internal_test.go si fonctions privées, xxx_external_test.go si fonctions publiques)",
	Run:  runTest008,
}

// fileAnalysisResult contient le résultat de l'analyse d'un fichier.
type fileAnalysisResult struct {
	hasPublic    bool
	hasPrivate   bool
	publicFuncs  []string
	privateFuncs []string
}

// runTest008 exécute l'analyse KTN-TEST-008.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runTest008(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeTest008) {
		// Règle désactivée
		return nil, nil
	}

	// Itération sur les fichiers
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeTest008, filename) {
			// Fichier exclu
			continue
		}

		// Vérification si test
		if shared.IsTestFile(filename) {
			// Ignorer fichier de test
			continue
		}

		// Vérification si mock
		if shared.IsMockFile(filename) {
			// Ignorer fichier mock
			continue
		}

		// Vérification si package main
		if file.Name.Name == "main" {
			// Ignorer package main
			continue
		}

		// Analyser le fichier
		result := analyzeFileFunctions(file)

		// Vérification si fonctions présentes
		if len(result.publicFuncs) == 0 && len(result.privateFuncs) == 0 {
			// Pas de fonctions, continuer
			continue
		}

		// Vérifier les tests existants
		status := checkTestFilesExist(filename)

		// Reporter les problèmes
		reportTestFileIssues(pass, file, result, status)
	}

	// Retour de la fonction
	return nil, nil
}

// analyzeFileFunctions analyse un fichier pour détecter les fonctions/variables publiques/privées.
//
// Params:
//   - file: fichier AST à analyser
//
// Returns:
//   - *fileAnalysisResult: pointeur vers le résultat de l'analyse
func analyzeFileFunctions(file *ast.File) *fileAnalysisResult {
	// Initialisation du résultat avec capacité définie
	result := &fileAnalysisResult{
		publicFuncs:  make([]string, 0, initialFuncsCap),
		privateFuncs: make([]string, 0, initialFuncsCap),
	}

	// Parcourir l'AST
	ast.Inspect(file, func(n ast.Node) bool {
		// Vérification si fonction
		if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Name != nil {
			// Classifier la fonction
			classifyFunction(funcDecl, result)
			// Continuer la traversée
			return true
		}

		// Vérification si déclaration générique
		if genDecl, ok := n.(*ast.GenDecl); ok {
			// Vérifier les variables
			checkVariables(genDecl, result)
			// Vérifier les types
			checkTypes(genDecl, result)
			// Vérifier les constantes
			checkConsts(genDecl, result)
		}

		// Continuer la traversée
		return true
	})

	// Retour du pointeur vers le résultat
	return result
}

// classifyFunction classifie une fonction comme publique ou privée.
//
// Params:
//   - funcDecl: déclaration de fonction
//   - result: résultat à mettre à jour
func classifyFunction(funcDecl *ast.FuncDecl, result *fileAnalysisResult) {
	funcName := funcDecl.Name.Name
	// Vérification si exemptée
	if isExemptFunction(funcName) {
		// Retour si exemptée
		return
	}

	// Vérification si mock
	if shared.IsMockName(funcName) {
		// Retour si mock
		return
	}

	// Classifier la fonction
	meta := shared.ClassifyFunc(funcDecl)

	// Vérification receiver mock
	if meta.ReceiverName != "" && shared.IsMockName(meta.ReceiverName) {
		// Retour si mock
		return
	}

	// Construire le nom d'affichage
	displayName := buildFunctionDisplayName(funcDecl)

	// Vérification visibilité
	if meta.Visibility == shared.VisPublic {
		// Ajouter fonction publique
		result.hasPublic = true
		result.publicFuncs = append(result.publicFuncs, displayName)
	} else {
		// Cas alternatif: privée
		// Ajouter fonction privée
		result.hasPrivate = true
		result.privateFuncs = append(result.privateFuncs, displayName)
	}
}

// buildFunctionDisplayName construit le nom d'affichage d'une fonction.
//
// Params:
//   - funcDecl: déclaration de fonction
//
// Returns:
//   - string: nom d'affichage (ex: "Func" ou "(*Type).Method")
func buildFunctionDisplayName(funcDecl *ast.FuncDecl) string {
	funcName := funcDecl.Name.Name

	// Vérification si méthode
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		// Extraire le type du receiver
		receiverType := extractReceiverTypeString(funcDecl.Recv.List[0].Type)
		// Vérification type valide
		if receiverType != "" {
			// Retour nom avec receiver
			return receiverType + "." + funcName
		}
	}

	// Retour nom simple
	return funcName
}

// extractReceiverTypeString extrait le type du receiver sous forme de string.
//
// Params:
//   - expr: expression du type
//
// Returns:
//   - string: type formaté (ex: "*App" ou "App")
func extractReceiverTypeString(expr ast.Expr) string {
	// Vérification si pointeur
	if starExpr, isPointer := expr.(*ast.StarExpr); isPointer {
		// Vérification si identifiant
		if innerIdent, isIdent := starExpr.X.(*ast.Ident); isIdent {
			// Retour nom avec pointeur
			return "(*" + innerIdent.Name + ")"
		}
		// Retour vide si non supporté
		return ""
	}

	// Vérification si identifiant
	if ident, isIdent := expr.(*ast.Ident); isIdent {
		// Retour nom simple
		return ident.Name
	}

	// Retour vide par défaut
	return ""
}

// checkVariables checks for public variables that require external testing.
// Public variables (exported, like Analyzer001) require external tests.
// Private variables are tested indirectly via the functions that use them.
//
// Params:
//   - genDecl: déclaration générique
//   - result: résultat de l'analyse
func checkVariables(genDecl *ast.GenDecl, result *fileAnalysisResult) {
	// Only process var declarations (not const or import)
	if genDecl.Tok != token.VAR {
		// Skip non-var declarations
		return
	}

	// Iterate over specs
	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		// Skip if not value spec
		if !ok {
			// Continue to next spec
			continue
		}

		// Check each variable name
		for _, name := range valueSpec.Names {
			// Skip blank identifier
			if name.Name == "_" {
				// Continue to next name
				continue
			}

			// Skip mock variables
			if shared.IsMockName(name.Name) {
				// Continue to next name
				continue
			}

			// Check if exported (public)
			if ast.IsExported(name.Name) {
				// Mark as having public element
				result.hasPublic = true
				// Public variables don't need to be listed (tested via their usage)
			} else {
				// Private variables may need internal tests to access them
				result.hasPrivate = true
			}
		}
	}
}

// checkTypes checks for public types that require external testing.
// Public types (exported, like MyStruct) require external tests.
// Private types don't require specific tests.
//
// Params:
//   - genDecl: déclaration générique
//   - result: résultat de l'analyse
func checkTypes(genDecl *ast.GenDecl, result *fileAnalysisResult) {
	// Only process type declarations
	if genDecl.Tok != token.TYPE {
		// Skip non-type declarations
		return
	}

	// Iterate over specs
	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		// Skip if not type spec
		if !ok {
			// Continue to next spec
			continue
		}

		// Skip mock types
		if shared.IsMockName(typeSpec.Name.Name) {
			// Continue to next spec
			continue
		}

		// Check if exported (public)
		if ast.IsExported(typeSpec.Name.Name) {
			// Mark as having public element
			result.hasPublic = true
		}
	}
}

// checkConsts checks for public constants that require external testing.
// Public constants (exported) require external tests.
// Private constants don't require specific tests.
//
// Params:
//   - genDecl: déclaration générique
//   - result: résultat de l'analyse
func checkConsts(genDecl *ast.GenDecl, result *fileAnalysisResult) {
	// Only process const declarations
	if genDecl.Tok != token.CONST {
		// Skip non-const declarations
		return
	}

	// Iterate over specs
	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		// Skip if not value spec
		if !ok {
			// Continue to next spec
			continue
		}

		// Check each constant name
		for _, name := range valueSpec.Names {
			// Skip mock constants
			if shared.IsMockName(name.Name) {
				// Continue to next name
				continue
			}

			// Check if exported (public)
			if ast.IsExported(name.Name) {
				// Mark as having public element
				result.hasPublic = true
			}
		}
	}
}

// checkTestFilesExist vérifie l'existence des fichiers de test pour un fichier source.
//
// Params:
//   - filename: chemin du fichier source
//
// Returns:
//   - testFilesStatus: état des fichiers de test
func checkTestFilesExist(filename string) testFilesStatus {
	// Extraire le répertoire
	dir := filepath.Dir(filename)
	// Extraire le nom de base
	baseName := filepath.Base(filename)
	// Retirer extension .go
	fileBase := strings.TrimSuffix(baseName, ".go")

	// Retour de l'état des fichiers de test
	return testFilesStatus{
		baseName:    baseName,
		fileBase:    fileBase,
		hasInternal: fileExistsOnDisk(filepath.Join(dir, fileBase+"_internal_test.go")),
		hasExternal: fileExistsOnDisk(filepath.Join(dir, fileBase+"_external_test.go")),
	}
}

// reportTestFileIssues reporte les problèmes de fichiers de test manquants ou superflus.
//
// Params:
//   - pass: contexte d'analyse
//   - file: fichier AST
//   - result: pointeur vers le résultat de l'analyse des fonctions
//   - status: état des fichiers de test
func reportTestFileIssues(pass *analysis.Pass, file *ast.File, result *fileAnalysisResult, status testFilesStatus) {
	// Sélection selon le type
	switch {
	// Vérification cas mixte
	case result.hasPublic && result.hasPrivate:
		// Cas mixte
		// Reporter problèmes mixtes
		reportMixedFunctionsIssues(pass, file, result, status)
	// Vérification cas public
	case result.hasPublic:
		// Cas public uniquement
		// Reporter problèmes publics
		reportPublicOnlyIssues(pass, file, result, status)
	// Vérification cas privé
	case result.hasPrivate:
		// Cas privé uniquement
		// Reporter problèmes privés
		reportPrivateOnlyIssues(pass, file, result, status)
	}
}

// formatFuncList formate une liste de fonctions pour l'affichage.
//
// Params:
//   - funcs: liste des noms de fonctions
//
// Returns:
//   - string: liste formatée (ex: "Func1, Func2, Func3, ...")
func formatFuncList(funcs []string) string {
	// Vérification si vide
	if len(funcs) == 0 {
		// Retour chaîne vide
		return ""
	}

	// Vérification si dans limite
	if len(funcs) <= maxFuncsToDisplay {
		// Retour liste complète
		return strings.Join(funcs, ", ")
	}

	// Calculer restantes
	displayed := funcs[:maxFuncsToDisplay]
	remaining := len(funcs) - maxFuncsToDisplay
	// Retour liste tronquée
	return strings.Join(displayed, ", ") + ", ... (+" + formatCount(remaining) + ")"
}

// formatCount formate un nombre pour l'affichage.
//
// Params:
//   - count: nombre à formater
//
// Returns:
//   - string: nombre formaté
func formatCount(count int) string {
	// Retour du nombre sous forme de string
	return strconv.Itoa(count)
}

// reportMixedFunctionsIssues reporte les problèmes pour fichiers avec fonctions publiques ET privées.
//
// Params:
//   - pass: contexte d'analyse
//   - file: fichier AST
//   - _result: pointeur vers le résultat de l'analyse (inutilisé avec messages package)
//   - status: état des fichiers de test
func reportMixedFunctionsIssues(pass *analysis.Pass, file *ast.File, _result *fileAnalysisResult, status testFilesStatus) {
	msg, _ := messages.Get(ruleCodeTest008)
	// Vérification absence totale fichiers
	if !status.hasInternal && !status.hasExternal {
		// Signaler manque des deux
		bothFiles := status.fileBase + "_external_test.go et " + status.fileBase + "_internal_test.go"
		pass.Reportf(file.Name.Pos(),
			"%s: %s",
			ruleCodeTest008,
			msg.Format(config.Get().Verbose, status.baseName, bothFiles))
		// Retour après signalement
		return
	}

	// Vérification fichier internal manquant
	if !status.hasInternal {
		// Signaler manque internal
		pass.Reportf(file.Name.Pos(),
			"%s: %s",
			ruleCodeTest008,
			msg.Format(config.Get().Verbose, status.baseName, status.fileBase+"_internal_test.go"))
	}
	// Vérification fichier external manquant
	if !status.hasExternal {
		// Signaler manque external
		pass.Reportf(file.Name.Pos(),
			"%s: %s",
			ruleCodeTest008,
			msg.Format(config.Get().Verbose, status.baseName, status.fileBase+"_external_test.go"))
	}
}

// reportPublicOnlyIssues reporte les problèmes pour fichiers avec UNIQUEMENT des fonctions publiques.
//
// Params:
//   - pass: contexte d'analyse
//   - file: fichier AST
//   - _result: pointeur vers le résultat de l'analyse (inutilisé avec messages package)
//   - status: état des fichiers de test
func reportPublicOnlyIssues(pass *analysis.Pass, file *ast.File, _result *fileAnalysisResult, status testFilesStatus) {
	// Vérification absence external
	if !status.hasExternal {
		// Signaler manque external
		msg, _ := messages.Get(ruleCodeTest008)
		pass.Reportf(file.Name.Pos(),
			"%s: %s",
			ruleCodeTest008,
			msg.Format(config.Get().Verbose, status.baseName, status.fileBase+"_external_test.go"))
		// Retour après signalement
		return
	}

	// Vérification internal superflu
	if status.hasInternal {
		// Signaler fichier superflu
		msg, _ := messages.Get(ruleCodeTest008)
		pass.Reportf(file.Name.Pos(),
			"%s: %s",
			ruleCodeTest008,
			msg.Format(config.Get().Verbose, status.baseName, status.fileBase+"_internal_test.go"))
	}
}

// reportPrivateOnlyIssues reporte les problèmes pour fichiers avec UNIQUEMENT des fonctions privées.
//
// Params:
//   - pass: contexte d'analyse
//   - file: fichier AST
//   - _result: pointeur vers le résultat de l'analyse (inutilisé avec messages package)
//   - status: état des fichiers de test
func reportPrivateOnlyIssues(pass *analysis.Pass, file *ast.File, _result *fileAnalysisResult, status testFilesStatus) {
	// Vérification absence internal
	if !status.hasInternal {
		// Signaler manque internal
		msg, _ := messages.Get(ruleCodeTest008)
		pass.Reportf(file.Name.Pos(),
			"%s: %s",
			ruleCodeTest008,
			msg.Format(config.Get().Verbose, status.baseName, status.fileBase+"_internal_test.go"))
		// Retour après signalement
		return
	}

	// Vérification external superflu
	if status.hasExternal {
		// Signaler fichier superflu
		msg, _ := messages.Get(ruleCodeTest008)
		pass.Reportf(file.Name.Pos(),
			"%s: %s",
			ruleCodeTest008,
			msg.Format(config.Get().Verbose, status.baseName, status.fileBase+"_external_test.go"))
	}
}

// fileExistsOnDisk vérifie si un fichier existe sur le disque.
// Cette fonction permet de détecter les fichiers _external_test.go même s'ils utilisent
// package xxx_test et sont dans XTestGoFiles (package séparé de Go).
//
// Params:
//   - path: chemin absolu du fichier à vérifier
//
// Returns:
//   - bool: true si le fichier existe, false sinon
func fileExistsOnDisk(path string) bool {
	// Obtenir les informations du fichier
	info, err := os.Stat(path)
	// Vérification de la condition
	if err != nil {
		// Erreur ou fichier n'existe pas
		return false
	}
	// Retour du résultat (fichier existe et n'est pas un répertoire)
	return !info.IsDir()
}
