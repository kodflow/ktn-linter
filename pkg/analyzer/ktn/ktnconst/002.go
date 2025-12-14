// Analyzer 002 for the ktnconst package.
package ktnconst

import (
	"go/ast"
	"go/token"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
)

const (
	// typeInfoCap initial capacity for type info maps
	typeInfoCap int = 8
	// ruleCodeConst002 est le code de la règle KTN-CONST-002.
	ruleCodeConst002 string = "KTN-CONST-002"
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
	// typeNames maps type name to its position
	typeNames map[string]token.Pos
	// constTypes maps const position to the type it uses (for iota detection)
	constTypes map[token.Pos]string
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
	// Récupération de la configuration
	cfg := config.Get()

	// Vérifier si la règle est activée
	if !cfg.IsRuleEnabled(ruleCodeConst002) {
		// Règle désactivée
		return nil, nil
	}

	// Analyze each file independently
	for _, file := range pass.Files {
		filename := pass.Fset.Position(file.Pos()).Filename
		// Skip excluded files
		if cfg.IsFileExcluded(ruleCodeConst002, filename) {
			// File is excluded
			continue
		}

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
	decls := &fileDeclarations{
		typeNames:  make(map[string]token.Pos, typeInfoCap),
		constTypes: make(map[token.Pos]string, typeInfoCap),
	}

	// Iterate over all declarations
	for _, decl := range file.Decls {
		// Handle GenDecl (const, var, type)
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			// Switch on token type
			switch genDecl.Tok {
			// Const declaration
			case token.CONST:
				decls.constDecls = append(decls.constDecls, genDecl.Pos())
				// Detect type used in const (for iota pattern)
				typeName := extractConstTypeName(genDecl)
				// Store type name if found
				if typeName != "" {
					decls.constTypes[genDecl.Pos()] = typeName
				}
			// Var declaration
			case token.VAR:
				decls.varDecls = append(decls.varDecls, genDecl.Pos())
			// Type declaration
			case token.TYPE:
				decls.typeDecls = append(decls.typeDecls, genDecl.Pos())
				// Store type name and position
				collectTypeNames(genDecl, decls.typeNames)
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

// extractConstTypeName extracts the type name from a const declaration.
// Used to detect iota pattern: const ( X TypeName = iota ).
//
// Params:
//   - genDecl: const declaration
//
// Returns:
//   - string: type name or empty if not found
func extractConstTypeName(genDecl *ast.GenDecl) string {
	// Iterate over specs to find typed const
	for _, spec := range genDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		// Skip if not value spec
		if !ok {
			continue
		}

		// Check if has explicit type
		if valueSpec.Type != nil {
			// Extract type name from ident
			if ident, ok := valueSpec.Type.(*ast.Ident); ok {
				// Return type name
				return ident.Name
			}
		}
	}

	// No type found
	return ""
}

// collectTypeNames stores type names and positions.
//
// Params:
//   - genDecl: type declaration
//   - typeNames: map to store names and positions
func collectTypeNames(genDecl *ast.GenDecl, typeNames map[string]token.Pos) {
	// Iterate over type specs
	for _, spec := range genDecl.Specs {
		typeSpec, ok := spec.(*ast.TypeSpec)
		// Skip if not type spec
		if !ok {
			continue
		}

		// Store type name and position
		typeNames[typeSpec.Name.Name] = genDecl.Pos()
	}
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

	// Check scattered const blocks (all except first, unless iota with custom type)
	checkScatteredConstBlocks(pass, decls)

	// Check const vs var order
	checkConstBeforeVar(pass, decls)

	// Check const vs type order
	checkConstBeforeType(pass, decls)

	// Check const vs func order
	checkConstBeforeFunc(pass, decls)
}

// checkScatteredConstBlocks reports if const declarations are scattered.
// Exception: const blocks using iota with custom types are allowed separate.
//
// Params:
//   - pass: analysis context
//   - decls: collected declarations (includes constTypes for iota detection)
func checkScatteredConstBlocks(pass *analysis.Pass, decls *fileDeclarations) {
	// If 0 or 1 const block, they're not scattered
	if len(decls.constDecls) <= 1 {
		// Return early
		return
	}

	// Report const groups except the first, unless they use iota with custom type
	for i := 1; i < len(decls.constDecls); i++ {
		constPos := decls.constDecls[i]

		// Skip if this const uses a custom type (iota pattern)
		if usedType, hasType := decls.constTypes[constPos]; hasType {
			// Check if the type is defined in this file
			if _, typeExists := decls.typeNames[usedType]; typeExists {
				// Valid iota pattern - skip
				continue
			}
		}

		// Report violation
		pass.Reportf(
			constPos,
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
// Exception: const blocks using iota with a custom type are allowed after that type.
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
		// Const after type is a violation, unless it uses iota with that type
		if constPos > firstTypePos {
			// Check if this const uses a type defined in this file (iota pattern)
			if usedType, hasType := decls.constTypes[constPos]; hasType {
				// Check if the type is defined in this file
				if typePos, typeExists := decls.typeNames[usedType]; typeExists {
					// If const comes after its type declaration, this is valid iota pattern
					if constPos > typePos {
						// Skip this const - it's a valid iota pattern
						continue
					}
				}
			}

			// Report violation
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
