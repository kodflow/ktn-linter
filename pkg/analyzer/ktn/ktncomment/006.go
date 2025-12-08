// Analyzer 006 for the ktncomment package.
package ktncomment

import (
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var (
	// Analyzer006 checks that all functions have proper documentation.
	Analyzer006 *analysis.Analyzer = &analysis.Analyzer{
		Name:     "ktncomment006",
		Doc:      "KTN-COMMENT-006: Toutes les fonctions doivent avoir une documentation au format strict (description, Params, Returns)",
		Run:      runComment006,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}

	// paramsHeaderPattern matches Params section header.
	paramsHeaderPattern *regexp.Regexp = regexp.MustCompile(`^//\s*Params:\s*$`)

	// paramItemPattern matches individual Params items.
	paramItemPattern *regexp.Regexp = regexp.MustCompile(`^//\s*-\s*\w+:\s*.+`)

	// returnsHeaderPattern matches Returns section header.
	returnsHeaderPattern *regexp.Regexp = regexp.MustCompile(`^//\s*Returns:\s*$`)

	// returnItemPattern matches individual Returns items.
	returnItemPattern *regexp.Regexp = regexp.MustCompile(`^//\s*-\s*.+:\s*.+`)
)

// runComment006 description à compléter.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat
//   - error: erreur éventuelle
func runComment006(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		funcDecl := n.(*ast.FuncDecl)

		// Get function documentation
		doc := funcDecl.Doc

		// Check if documentation exists
		if doc == nil || len(doc.List) == 0 {
			pass.Reportf(
				funcDecl.Name.Pos(),
				"KTN-COMMENT-006: la fonction '%s' doit avoir une documentation",
				funcDecl.Name.Name,
			)
			// Retour de la fonction
			return
		}

		// Extract comment lines
		comments := extractCommentLines(doc)

		// Validate documentation format
		hasParams := funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) > 0
		hasReturns := funcDecl.Type.Results != nil && len(funcDecl.Type.Results.List) > 0

		err := validateDocFormat(comments, funcDecl.Name.Name, hasParams, hasReturns)
		// Vérification de la condition
		if err != "" {
			pass.Reportf(
				funcDecl.Name.Pos(),
				"KTN-COMMENT-006: %s",
				err,
			)
		}
	})

	// Retour de la fonction
	return nil, nil
}

// extractCommentLines extracts individual comment lines from a CommentGroup
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - []string: lignes de commentaires
func extractCommentLines(cg *ast.CommentGroup) []string {
	var lines []string
	// Itération sur les éléments
	for _, comment := range cg.List {
		text := comment.Text
		// Remove // prefix
		if strings.HasPrefix(text, "//") {
			lines = append(lines, text)
		}
	}
	// Retour de la fonction
	return lines
}

// validateDescriptionLine vérifie que la première ligne de la doc est correcte.
//
// Params:
//   - comments: lignes de commentaires
//   - funcName: nom de la fonction
//
// Returns:
//   - string: message d'erreur ou vide
func validateDescriptionLine(comments []string, funcName string) string {
	// Vérification de la présence de commentaires
	if len(comments) == 0 {
		// Retour si aucun commentaire
		return "documentation vide"
	}

	// Vérification du format de la première ligne
	if !strings.HasPrefix(comments[0], "// "+funcName) {
		// Retour si format incorrect
		return "la description doit commencer par '// " + funcName + " '"
	}

	// Validation réussie
	return ""
}

// validateParamsSection vérifie la section Params: de la documentation.
//
// Params:
//   - comments: lignes de commentaires
//   - startIdx: index de début dans comments
//
// Returns:
//   - error: message d'erreur ou vide
//   - newIdx: nouvel index après la section
func validateParamsSection(comments []string, startIdx int) (string, int) {
	idx := startIdx

	// Vérification de la présence du header Params:
	if idx >= len(comments) || !paramsHeaderPattern.MatchString(comments[idx]) {
		// Retour si header manquant avec message explicite
		return "section 'Params:' manquante. Ajouter après description: '// Params:' sur ligne seule, puis '//   - nomParam: description' (avec 2 espaces avant tiret)", idx
	}
	idx++

	// Validation des items de paramètres
	foundParam := false
	// Itération sur les items
	for idx < len(comments) {
		line := comments[idx]
		// Vérification si c'est un item de paramètre
		if paramItemPattern.MatchString(line) {
			foundParam = true
			idx++
		} else {
			// Sortie de boucle
			break
		}
	}

	// Vérification qu'au moins un paramètre est documenté
	if !foundParam {
		// Retour si aucun paramètre documenté
		return "au moins un paramètre doit être documenté dans 'Params:'", idx
	}

	// Skip blank line
	if idx < len(comments) && strings.TrimSpace(comments[idx]) == "//" {
		idx++
	}

	// Validation réussie
	return "", idx
}

// validateReturnsSection vérifie la section Returns: de la documentation.
//
// Params:
//   - comments: lignes de commentaires
//   - startIdx: index de début dans comments
//
// Returns:
//   - error: message d'erreur ou vide
//   - newIdx: nouvel index après la section
func validateReturnsSection(comments []string, startIdx int) (string, int) {
	idx := startIdx

	// Vérification de la présence du header Returns:
	if idx >= len(comments) || !returnsHeaderPattern.MatchString(comments[idx]) {
		// Retour si header manquant avec message explicite
		return "section 'Returns:' manquante. Ajouter après Params: '// Returns:' sur ligne seule, puis '//   - type: description' (ex: '//   - error: erreur éventuelle')", idx
	}
	idx++

	// Validation des items de retour
	foundReturn := false
	// Itération sur les items
	for idx < len(comments) {
		line := comments[idx]
		// Vérification si c'est un item de retour
		if returnItemPattern.MatchString(line) {
			foundReturn = true
			idx++
		} else {
			// Sortie de boucle
			break
		}
	}

	// Vérification qu'au moins un retour est documenté
	if !foundReturn {
		// Retour si aucun retour documenté
		return "au moins une valeur de retour doit être documentée dans 'Returns:'", idx
	}

	// Validation réussie
	return "", idx
}

// validateDocFormat vérifie le format de la documentation d'une fonction.
//
// Params:
//   - comments: lignes de commentaires de la fonction
//   - funcName: nom de la fonction
//   - hasParams: true si la fonction a des paramètres
//   - hasReturns: true si la fonction a des valeurs de retour
//
// Returns:
//   - string: message d'erreur ou chaîne vide si valide
func validateDocFormat(comments []string, funcName string, hasParams, hasReturns bool) string {
	// Validation de la ligne de description
	if err := validateDescriptionLine(comments, funcName); err != "" {
		// Retour si erreur de description
		return err
	}

	idx := 1

	// Skip description lines until we find blank line before Params/Returns
	for idx < len(comments) {
		line := strings.TrimSpace(comments[idx])
		// Stop at blank line separator
		if line == "//" {
			idx++
			break
		}
		// Stop if we find Params: or Returns: header
		if paramsHeaderPattern.MatchString(comments[idx]) || returnsHeaderPattern.MatchString(comments[idx]) {
			break
		}
		// Continue skipping description lines
		idx++
	}

	// Validation de la section Params si nécessaire
	if hasParams {
		var err string
		// Validation de la section
		err, idx = validateParamsSection(comments, idx)
		// Vérification d'erreur
		if err != "" {
			// Retour si erreur
			return err
		}
	}

	// Validation de la section Returns si nécessaire
	if hasReturns {
		var err string
		// Validation de la section (idx ignoré car dernière section)
		err, _ = validateReturnsSection(comments, idx)
		// Vérification d'erreur
		if err != "" {
			// Retour si erreur
			return err
		}
	}

	// Validation complète réussie
	return ""
}
