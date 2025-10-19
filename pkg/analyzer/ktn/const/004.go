package ktnconst

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer004 checks that every constant has an associated comment
var Analyzer004 = &analysis.Analyzer{
	Name:     "ktnconst004",
	Doc:      "KTN-CONST-004: Vérifie que chaque constante a un commentaire associé",
	Run:      runConst004,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func runConst004(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.GenDecl)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		genDecl := n.(*ast.GenDecl)

		// Only check const declarations
		if genDecl.Tok != token.CONST {
			return
		}

		// Check if the GenDecl has a doc comment (applies to all constants in the group)
		// Filter out "want" directives used by analysistest
		hasGenDeclDoc := hasValidComment(genDecl.Doc)

		for _, spec := range genDecl.Specs {
			valueSpec := spec.(*ast.ValueSpec)

			// Check if this specific ValueSpec has a doc comment or line comment
			// Filter out "want" directives used by analysistest
			hasValueSpecDoc := hasValidComment(valueSpec.Doc)
			hasValueSpecComment := hasValidComment(valueSpec.Comment)

			// A constant is considered documented if:
			// 1. The GenDecl has a doc comment (group documentation), OR
			// 2. The ValueSpec has a doc comment (above the constant), OR
			// 3. The ValueSpec has a line comment (on the same line)
			hasComment := hasGenDeclDoc || hasValueSpecDoc || hasValueSpecComment

			if !hasComment {
				for _, name := range valueSpec.Names {
					pass.Reportf(
						name.Pos(),
						"KTN-CONST-004: la constante '%s' doit avoir un commentaire associé",
						name.Name,
					)
				}
			}
		}
	})

	return nil, nil
}

// hasValidComment checks if a comment group exists and contains valid comments
// (not just test directives like "want")
func hasValidComment(cg *ast.CommentGroup) bool {
	if cg == nil || len(cg.List) == 0 {
		return false
	}

	// Check if any comment is NOT a "want" directive
	for _, comment := range cg.List {
		text := comment.Text
		// Skip "want" directives used by analysistest
		// Line comment: // want "..."
		if len(text) >= 6 && text[2:6] == "want" {
			continue
		}
		// Block comment: /* want "..." */
		if len(text) >= 7 && text[2:7] == " want" {
			continue
		}
		// Found a valid comment
		return true
	}

	return false
}
