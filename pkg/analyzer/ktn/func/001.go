package ktnfunc

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// MAX_PURE_CODE_LINES définit le nombre maximum de lignes de code pur autorisées dans une fonction
	MAX_PURE_CODE_LINES int = 35
)

// Analyzer001 checks that functions don't exceed 35 lines of pure code
var Analyzer001 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktnfunc001",
	Doc:      "KTN-FUNC-001: Les fonctions ne doivent pas dépasser 35 lignes de code pur (hors commentaires et lignes vides)",
	Run:      runFunc001,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runFunc001 exécute l'analyse KTN-FUNC-001.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runFunc001(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Skip if no body (external functions)
		if funcDecl.Body == nil {
			// Retour de la fonction
			return
		}

		// Skip test functions (Test*, Benchmark*, Example*, Fuzz*)
		funcName := funcDecl.Name.Name
		// Vérification de la condition
		if isTestFunction(funcName) {
			// Retour de la fonction
			return
		}

		// Skip main function
		if funcName == "main" {
			// Retour de la fonction
			return
		}

		// Count pure code lines
		pureLines := countPureCodeLines(pass, funcDecl.Body)

		// Vérification de la condition
		if pureLines > MAX_PURE_CODE_LINES {
			pass.Reportf(
				funcDecl.Name.Pos(),
				"KTN-FUNC-001: la fonction '%s' contient %d lignes de code pur (max: %d)",
				funcName,
				pureLines,
				MAX_PURE_CODE_LINES,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// isTestFunction checks if a function name indicates a test function
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - bool: true si fonction de test
func isTestFunction(name string) bool {
	// Retour de la fonction
	return strings.HasPrefix(name, "Test") ||
		strings.HasPrefix(name, "Benchmark") ||
		strings.HasPrefix(name, "Example") ||
		strings.HasPrefix(name, "Fuzz")
}

// isLineToSkip détermine si une ligne doit être ignorée du compte.
//
// Params:
//   - trimmed: ligne trimée à vérifier
//   - inBlockComment: pointeur vers l'état du commentaire bloc
//
// Returns:
//   - bool: true si la ligne doit être ignorée
func isLineToSkip(trimmed string, inBlockComment *bool) bool {
	// Gestion des commentaires de bloc
	if strings.Contains(trimmed, "/*") {
		*inBlockComment = true
	}
	// Vérification fin de commentaire bloc
	if strings.Contains(trimmed, "*/") {
		*inBlockComment = false
		// Ligne de fin de bloc à ignorer
		return true
	}

	// Vérification si dans un commentaire bloc
	if *inBlockComment {
		// Ligne de commentaire bloc à ignorer
		return true
	}

	// Vérification ligne vide
	if trimmed == "" {
		// Ligne vide à ignorer
		return true
	}

	// Vérification commentaire ligne
	if strings.HasPrefix(trimmed, "//") {
		// Commentaire ligne à ignorer
		return true
	}

	// Vérification accolade seule
	if trimmed == "{" || trimmed == "}" {
		// Accolade seule à ignorer
		return true
	}

	// Ligne à compter
	return false
}

// countPureCodeLines compte les lignes de code pur dans le corps d'une fonction.
//
// Params:
//   - pass: contexte d'analyse
//   - body: corps de la fonction à analyser
//
// Returns:
//   - int: nombre de lignes de code pur
func countPureCodeLines(pass *analysis.Pass, body *ast.BlockStmt) int {
	startPos := pass.Fset.Position(body.Lbrace)
	endPos := pass.Fset.Position(body.Rbrace)

	// Lecture du fichier source
	filename := startPos.Filename
	// Vérification de ReadFile (peut être nil dans certains contextes)
	if pass.ReadFile == nil {
		// Retour si ReadFile indisponible
		return 0
	}
	content, err := pass.ReadFile(filename)
	// Vérification erreur lecture
	if err != nil {
		// Retour si erreur
		return 0
	}
	lines := strings.Split(string(content), "\n")
	pureCodeLines := 0
	inBlockComment := false

	// Itération sur les lignes de la fonction
	for i := startPos.Line + 1; i < endPos.Line; i++ {
		line := lines[i-1]
		trimmed := strings.TrimSpace(line)

		// Vérification si ligne doit être ignorée
		if !isLineToSkip(trimmed, &inBlockComment) {
			pureCodeLines++
		}
	}

	// Retour du compte
	return pureCodeLines
}
