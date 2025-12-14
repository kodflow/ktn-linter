// Analyzer 004 for the ktncomment package.
package ktncomment

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// ruleCodeComment004 is the rule code for this analyzer
	ruleCodeComment004 string = "KTN-COMMENT-004"
)

// Analyzer004 checks that every package-level variable has an associated comment
var Analyzer004 *analysis.Analyzer = &analysis.Analyzer{
	Name:     "ktncomment004",
	Doc:      "KTN-COMMENT-004: Vérifie que chaque variable de package a un commentaire associé",
	Run:      runComment004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// runComment004 exécute l'analyse KTN-COMMENT-004.
//
// Params:
//   - pass: contexte d'analyse
//
// Returns:
//   - any: résultat de l'analyse
//   - error: erreur éventuelle
func runComment004(pass *analysis.Pass) (any, error) {
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeComment004) {
		// Règle désactivée
		return nil, nil
	}

	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// Filter for File nodes to access package-level declarations only
	nodeFilter := []ast.Node{
		(*ast.File)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		file := n.(*ast.File)

		// Vérifier si le fichier est exclu
		filename := pass.Fset.Position(n.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeComment004, filename) {
			// File excluded by configuration
			return
		}

		// Check package-level declarations only
		checkFileDeclarations(pass, file)
	})

	// Retour de la fonction
	return nil, nil
}

// checkFileDeclarations checks all declarations in a file for var documentation.
//
// Params:
//   - pass: analysis pass
//   - file: file to check
func checkFileDeclarations(pass *analysis.Pass, file *ast.File) {
	// Check package-level declarations
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		// Skip if not a GenDecl
		if !ok {
			continue
		}

		// Only check var declarations
		if genDecl.Tok != token.VAR {
			continue
		}

		// Check var declaration
		checkVarDeclaration(pass, genDecl)
	}
}

// checkVarDeclaration checks a single var declaration for documentation.
//
// Params:
//   - pass: analysis pass
//   - genDecl: var declaration to check
func checkVarDeclaration(pass *analysis.Pass, genDecl *ast.GenDecl) {
	// Check if the GenDecl has a doc comment
	hasGenDeclDoc := shared.HasValidComment(genDecl.Doc)

	// Check each variable spec
	for _, spec := range genDecl.Specs {
		valueSpec := spec.(*ast.ValueSpec)

		// Check if this specific ValueSpec has documentation
		hasValueSpecDoc := shared.HasValidComment(valueSpec.Doc)
		hasValueSpecComment := shared.HasValidComment(valueSpec.Comment)

		// A variable is considered documented if it has any comment
		hasComment := hasGenDeclDoc || hasValueSpecDoc || hasValueSpecComment

		// Report missing documentation
		if !hasComment {
			reportMissingVarDoc(pass, valueSpec)
		}
	}
}

// reportMissingVarDoc reports missing documentation for a variable.
//
// Params:
//   - pass: analysis pass
//   - valueSpec: variable spec with missing documentation
func reportMissingVarDoc(pass *analysis.Pass, valueSpec *ast.ValueSpec) {
	msg, _ := messages.Get(ruleCodeComment004)
	// Report for each variable name
	for _, name := range valueSpec.Names {
		pass.Reportf(
			name.Pos(),
			"%s: %s",
			ruleCodeComment004,
			msg.Format(config.Get().Verbose, name.Name),
		)
	}
}
