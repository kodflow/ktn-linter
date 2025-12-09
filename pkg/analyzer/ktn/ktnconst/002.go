// Analyzer 002 for the ktnconst package.
package ktnconst

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Analyzer002 checks that constants are grouped together and placed at the top.
// Order must be: const → var → type → func
var Analyzer002 *analysis.Analyzer = &analysis.Analyzer{
	Name: "ktnconst002",
	Doc:  "KTN-CONST-002: Vérifie que les constantes sont groupées et placées en haut (ordre: const → var → type → func)",
	Run:  runConst002,
}

// fileDeclarations holds all declaration positions for a file.
type fileDeclarations struct {
	constDecls []token.Pos
	varDecls   []token.Pos
	typeDecls  []token.Pos
	funcDecls  []token.Pos
}

// runConst002 executes KTN-CONST-002 analysis.
//
// Params:
//   - pass: analysis context
//
// Returns:
//   - any: analysis result
//   - error: potential error
func runConst002(pass *analysis.Pass) (any, error) {
	// Analyze each file independently
	for _, file := range pass.Files {
		decls := collectDeclarations(file)
		checkConstOrder(pass, decls)
	}

	// Return result
	return nil, nil
}

// collectDeclarations gathers all declaration positions from a file.
//
// Params:
//   - file: AST file to analyze
//
// Returns:
//   - *fileDeclarations: collected declaration positions
func collectDeclarations(file *ast.File) *fileDeclarations {
	decls := &fileDeclarations{}

	// Iterate over all declarations
	for _, decl := range file.Decls {
		// Handle GenDecl (const, var, type)
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			// Switch on token type
			switch genDecl.Tok {
			// Const declaration
			case token.CONST:
				decls.constDecls = append(decls.constDecls, genDecl.Pos())
			// Var declaration
			case token.VAR:
				decls.varDecls = append(decls.varDecls, genDecl.Pos())
			// Type declaration
			case token.TYPE:
				decls.typeDecls = append(decls.typeDecls, genDecl.Pos())
			}
		}

		// Handle FuncDecl
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			decls.funcDecls = append(decls.funcDecls, funcDecl.Pos())
		}
	}

	// Return collected declarations
	return decls
}

// checkConstOrder verifies const declarations are properly ordered and grouped.
//
// Params:
//   - pass: analysis context
//   - decls: collected declarations
func checkConstOrder(pass *analysis.Pass, decls *fileDeclarations) {
	// No const declarations, nothing to check
	if len(decls.constDecls) == 0 {
		// Return early
		return
	}

	// Check scattered const blocks (all except first)
	checkScatteredConstBlocks(pass, decls.constDecls)

	// Check const vs var order
	checkConstBeforeVar(pass, decls)

	// Check const vs type order
	checkConstBeforeType(pass, decls)

	// Check const vs func order
	checkConstBeforeFunc(pass, decls)
}

// checkScatteredConstBlocks reports if const declarations are scattered.
//
// Params:
//   - pass: analysis context
//   - constDecls: positions of const declarations
func checkScatteredConstBlocks(pass *analysis.Pass, constDecls []token.Pos) {
	// If 0 or 1 const block, they're not scattered
	if len(constDecls) <= 1 {
		// Return early
		return
	}

	// Report all const groups except the first as scattered
	for i := 1; i < len(constDecls); i++ {
		pass.Reportf(
			constDecls[i],
			"KTN-CONST-002: les constantes doivent être groupées ensemble dans un seul bloc",
		)
	}
}

// checkConstBeforeVar ensures const declarations come before var declarations.
//
// Params:
//   - pass: analysis context
//   - decls: collected declarations
func checkConstBeforeVar(pass *analysis.Pass, decls *fileDeclarations) {
	// No var declarations, nothing to check
	if len(decls.varDecls) == 0 {
		// Return early
		return
	}

	// Get first var position
	firstVarPos := decls.varDecls[0]

	// Check each const declaration
	for _, constPos := range decls.constDecls {
		// Const after var is a violation
		if constPos > firstVarPos {
			pass.Reportf(
				constPos,
				"KTN-CONST-002: les constantes doivent être placées avant les déclarations var",
			)
		}
	}
}

// checkConstBeforeType ensures const declarations come before type declarations.
//
// Params:
//   - pass: analysis context
//   - decls: collected declarations
func checkConstBeforeType(pass *analysis.Pass, decls *fileDeclarations) {
	// No type declarations, nothing to check
	if len(decls.typeDecls) == 0 {
		// Return early
		return
	}

	// Get first type position
	firstTypePos := decls.typeDecls[0]

	// Check each const declaration
	for _, constPos := range decls.constDecls {
		// Const after type is a violation
		if constPos > firstTypePos {
			pass.Reportf(
				constPos,
				"KTN-CONST-002: les constantes doivent être placées avant les déclarations type",
			)
		}
	}
}

// checkConstBeforeFunc ensures const declarations come before func declarations.
//
// Params:
//   - pass: analysis context
//   - decls: collected declarations
func checkConstBeforeFunc(pass *analysis.Pass, decls *fileDeclarations) {
	// No func declarations, nothing to check
	if len(decls.funcDecls) == 0 {
		// Return early
		return
	}

	// Get first func position
	firstFuncPos := decls.funcDecls[0]

	// Check each const declaration
	for _, constPos := range decls.constDecls {
		// Const after func is a violation
		if constPos > firstFuncPos {
			pass.Reportf(
				constPos,
				"KTN-CONST-002: les constantes doivent être placées avant les déclarations func",
			)
		}
	}
}
