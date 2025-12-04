// Package ktntest implements KTN linter rules.
package ktntest

import (
	"go/ast"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"golang.org/x/tools/go/analysis"
)

const (
	// MAX_FUNCS_TO_DISPLAY nombre max de fonctions affichées.
	MAX_FUNCS_TO_DISPLAY int = 3
	// INITIAL_FUNCS_CAP capacité initiale des listes.
	INITIAL_FUNCS_CAP int = 8
)

// Analyzer008 checks that each source file has appropriate test files based on its content
var Analyzer008 = &analysis.Analyzer{
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
	// Parcourir tous les fichiers sources (non-test)
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename

		// Ignorer les fichiers de test
		if shared.IsTestFile(filename) {
			continue
		}

		// Ignorer le package main (ne peut pas avoir de tests externes)
		if file.Name.Name == "main" {
			continue
		}

		// Analyser le contenu du fichier
		result := analyzeFileFunctions(file)

		// Si le fichier n'a pas de FONCTIONS (publiques ou privées), pas de test requis
		// Note: les constantes/types seuls ne nécessitent pas de tests dédiés car
		// ils n'ont pas de comportement à tester (valeurs compile-time ou struct)
		if len(result.publicFuncs) == 0 && len(result.privateFuncs) == 0 {
			continue
		}

		// Vérifier les fichiers de test existants
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
		publicFuncs:  make([]string, 0, INITIAL_FUNCS_CAP),
		privateFuncs: make([]string, 0, INITIAL_FUNCS_CAP),
	}

	// Parcourir l'AST du fichier
	ast.Inspect(file, func(n ast.Node) bool {
		// Vérifier les déclarations de fonctions
		if funcDecl, ok := n.(*ast.FuncDecl); ok && funcDecl.Name != nil {
			classifyFunction(funcDecl, result)
			return true
		}

		// Vérifier les déclarations de variables (publiques et privées)
		if genDecl, ok := n.(*ast.GenDecl); ok {
			checkVariables(genDecl, result)
		}

		// Continuer l'itération
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
	// Ignorer les fonctions exemptées (init, main)
	if isExemptFunction(funcName) {
		return
	}

	// Construire le nom complet (avec receiver si méthode)
	displayName := buildFunctionDisplayName(funcDecl)

	// Classifier la fonction
	if len(funcName) > 0 && unicode.IsUpper(rune(funcName[0])) {
		result.hasPublic = true
		result.publicFuncs = append(result.publicFuncs, displayName)
	} else {
		// Fonction privée (commence par minuscule)
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

	// Vérifier si c'est une méthode avec receiver
	if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
		receiverType := extractReceiverTypeString(funcDecl.Recv.List[0].Type)
		// Vérifier si le receiver a un type valide
		if receiverType != "" {
			// Retour du nom avec receiver (ex: "(*Type).Method")
			return receiverType + "." + funcName
		}
	}

	// Retour du nom simple (fonction sans receiver)
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
	// Vérifier si c'est un pointeur (*Type)
	if starExpr, isPointer := expr.(*ast.StarExpr); isPointer {
		// Vérifier si l'expression pointée est un identifiant
		if innerIdent, isIdent := starExpr.X.(*ast.Ident); isIdent {
			// Retour du nom avec notation pointeur
			return "(*" + innerIdent.Name + ")"
		}
		// Type pointeur non supporté (ex: *pkg.Type)
		return ""
	}

	// Vérifier si c'est un identifiant simple (Type)
	if ident, isIdent := expr.(*ast.Ident); isIdent {
		// Retour du nom simple
		return ident.Name
	}

	// Type non géré (retour vide)
	return ""
}

// checkVariables vérifie si une déclaration contient des variables publiques ou privées.
//
// Params:
//   - genDecl: déclaration générique
//   - result: résultat de l'analyse à mettre à jour
func checkVariables(genDecl *ast.GenDecl, result *fileAnalysisResult) {
	// Parcourir les spécifications de la déclaration
	for _, spec := range genDecl.Specs {
		// Vérifier si c'est une spécification de valeur (var ou const)
		valueSpec, ok := spec.(*ast.ValueSpec)
		// Si ce n'est pas une valeur, continuer
		if !ok {
			continue
		}
		// Vérifier chaque nom dans la spécification
		for _, name := range valueSpec.Names {
			varName := name.Name
			// Ignorer les variables blank (_)
			if varName == "_" {
				continue
			}
			// Classifier la variable comme publique ou privée
			if len(varName) > 0 && unicode.IsUpper(rune(varName[0])) {
				// Variable publique (commence par majuscule)
				result.hasPublic = true
			} else {
				// Variable privée (commence par minuscule)
				result.hasPrivate = true
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
	dir := filepath.Dir(filename)
	baseName := filepath.Base(filename)
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
	// Vérification selon le type de fichier
	switch {
	// Cas fichier mixte avec fonctions publiques ET privées
	case result.hasPublic && result.hasPrivate:
		reportMixedFunctionsIssues(pass, file, result, status)
	// Cas fichier avec fonctions publiques uniquement
	case result.hasPublic:
		reportPublicOnlyIssues(pass, file, result, status)
	// Cas fichier avec fonctions privées uniquement
	case result.hasPrivate:
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
	// Vérifier si la liste est vide
	if len(funcs) == 0 {
		// Retour d'une chaîne vide pour liste vide
		return ""
	}

	// Vérifier si le nombre de fonctions est dans la limite
	if len(funcs) <= MAX_FUNCS_TO_DISPLAY {
		// Retour de la liste complète séparée par des virgules
		return strings.Join(funcs, ", ")
	}

	// Calcul du nombre de fonctions restantes
	displayed := funcs[:MAX_FUNCS_TO_DISPLAY]
	remaining := len(funcs) - MAX_FUNCS_TO_DISPLAY
	// Retour de la liste tronquée avec indicateur du reste
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
//   - result: pointeur vers le résultat de l'analyse
//   - status: état des fichiers de test
func reportMixedFunctionsIssues(pass *analysis.Pass, file *ast.File, result *fileAnalysisResult, status testFilesStatus) {
	pubList := formatFuncList(result.publicFuncs)
	privList := formatFuncList(result.privateFuncs)

	// Fichier avec fonctions publiques ET privées → besoin des deux fichiers
	if !status.hasInternal && !status.hasExternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: '%s' contient %d fonction(s) publique(s) [%s] et %d fonction(s) privée(s) [%s]. Créez '%s_external_test.go' (black-box, package xxx_test) ET '%s_internal_test.go' (white-box, package xxx)",
			status.baseName, len(result.publicFuncs), pubList, len(result.privateFuncs), privList, status.fileBase, status.fileBase)
		return
	}

	// Vérifier les fichiers manquants individuellement
	if !status.hasInternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: '%s' contient %d fonction(s) privée(s) [%s] → créez '%s_internal_test.go' (white-box, package xxx) pour les tester",
			status.baseName, len(result.privateFuncs), privList, status.fileBase)
	}
	// Vérification du fichier externe manquant
	if !status.hasExternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: '%s' contient %d fonction(s) publique(s) [%s] → créez '%s_external_test.go' (black-box, package xxx_test) pour les tester",
			status.baseName, len(result.publicFuncs), pubList, status.fileBase)
	}
}

// reportPublicOnlyIssues reporte les problèmes pour fichiers avec UNIQUEMENT des fonctions publiques.
//
// Params:
//   - pass: contexte d'analyse
//   - file: fichier AST
//   - result: pointeur vers le résultat de l'analyse
//   - status: état des fichiers de test
func reportPublicOnlyIssues(pass *analysis.Pass, file *ast.File, result *fileAnalysisResult, status testFilesStatus) {
	pubList := formatFuncList(result.publicFuncs)

	// Vérification du fichier externe manquant
	if !status.hasExternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: '%s' contient UNIQUEMENT %d fonction(s) publique(s) [%s] → créez '%s_external_test.go' (black-box, package xxx_test)",
			status.baseName, len(result.publicFuncs), pubList, status.fileBase)
		return
	}

	// Vérification du fichier interne superflu
	if status.hasInternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: '%s' contient UNIQUEMENT des fonctions publiques [%s]. Supprimez '%s_internal_test.go' (inutile) et utilisez '%s_external_test.go'",
			status.baseName, pubList, status.fileBase, status.fileBase)
	}
}

// reportPrivateOnlyIssues reporte les problèmes pour fichiers avec UNIQUEMENT des fonctions privées.
//
// Params:
//   - pass: contexte d'analyse
//   - file: fichier AST
//   - result: pointeur vers le résultat de l'analyse
//   - status: état des fichiers de test
func reportPrivateOnlyIssues(pass *analysis.Pass, file *ast.File, result *fileAnalysisResult, status testFilesStatus) {
	privList := formatFuncList(result.privateFuncs)

	// Vérification du fichier interne manquant
	if !status.hasInternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: '%s' contient UNIQUEMENT %d fonction(s) privée(s) [%s] → créez '%s_internal_test.go' (white-box, package xxx)",
			status.baseName, len(result.privateFuncs), privList, status.fileBase)
		return
	}

	// Vérification du fichier externe superflu
	if status.hasExternal {
		pass.Reportf(file.Name.Pos(),
			"KTN-TEST-008: '%s' contient UNIQUEMENT des fonctions privées [%s]. Supprimez '%s_external_test.go' (inutile) et utilisez '%s_internal_test.go'",
			status.baseName, privList, status.fileBase, status.fileBase)
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
	info, err := os.Stat(path)
	// Vérification de la condition
	if err != nil {
		// Erreur ou fichier n'existe pas
		return false
	}
	// Retour du résultat (fichier existe et n'est pas un répertoire)
	return !info.IsDir()
}
