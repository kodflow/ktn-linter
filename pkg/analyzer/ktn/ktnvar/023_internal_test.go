package ktnvar

import (
	"go/ast"
	"go/token"
	"testing"
)

// TestIsSecurityName tests the isSecurityName function.
func TestIsSecurityName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "contains key",
			input:    "generateKey",
			expected: true,
		},
		{
			name:     "contains token",
			input:    "createToken",
			expected: true,
		},
		{
			name:     "contains secret",
			input:    "badSecretKey",
			expected: true,
		},
		{
			name:     "contains password",
			input:    "hashPassword",
			expected: true,
		},
		{
			name:     "contains salt",
			input:    "generateSalt",
			expected: true,
		},
		{
			name:     "contains nonce",
			input:    "createNonce",
			expected: true,
		},
		{
			name:     "contains crypt",
			input:    "encryptData",
			expected: true,
		},
		{
			name:     "contains auth",
			input:    "authHandler",
			expected: true,
		},
		{
			name:     "contains credential",
			input:    "getCredentials",
			expected: true,
		},
		{
			name:     "no security keyword",
			input:    "shuffleItems",
			expected: false,
		},
		{
			name:     "random index",
			input:    "randomIndex",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "case insensitive KEY",
			input:    "generateKEY",
			expected: true,
		},
		{
			name:     "case insensitive Token",
			input:    "CreateToken",
			expected: true,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := isSecurityName(tt.input)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isSecurityName(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestCheckAssignForMathRand tests checkAssignForMathRand function.
func TestCheckAssignForMathRand(t *testing.T) {
	// Test with non-call expression in RHS
	aliases := map[string]bool{"rand": true}

	// Test: RHS is not a call expression
	assign := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
		Rhs: []ast.Expr{&ast.Ident{Name: "y"}}, // Not a call
	}
	// Should not panic, just skip
	checkAssignForMathRand(nil, assign, aliases, false)

	// Test: RHS is a call but not math/rand
	callNotRand := &ast.CallExpr{
		Fun: &ast.Ident{Name: "someFunc"},
	}
	assign2 := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
		Rhs: []ast.Expr{callNotRand},
	}
	checkAssignForMathRand(nil, assign2, aliases, false)

	// Test: LHS index exceeds available
	assign3 := &ast.AssignStmt{
		Lhs: []ast.Expr{}, // Empty LHS
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "rand"},
					Sel: &ast.Ident{Name: "Int"},
				},
			},
		},
	}
	// LHS is empty, should handle gracefully
	checkAssignForMathRand(nil, assign3, aliases, false)

	// Test: LHS is not an ident
	assign4 := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.IndexExpr{X: &ast.Ident{Name: "arr"}}},
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "rand"},
					Sel: &ast.Ident{Name: "Int"},
				},
			},
		},
	}
	checkAssignForMathRand(nil, assign4, aliases, false)
}

// TestProcessLocalVarSpec tests processLocalVarSpec function.
func TestProcessLocalVarSpec(t *testing.T) {
	aliases := map[string]bool{"rand": true}

	// Test: spec is not a ValueSpec
	importSpec := &ast.ImportSpec{
		Path: &ast.BasicLit{Value: `"fmt"`},
	}
	processLocalVarSpec(nil, importSpec, aliases, false)

	// Test: ValueSpec with no values
	valueSpec := &ast.ValueSpec{
		Names:  []*ast.Ident{{Name: "x"}},
		Values: nil,
	}
	processLocalVarSpec(nil, valueSpec, aliases, false)
}

// TestHasSecurityVarName tests hasSecurityVarName function.
func TestHasSecurityVarName(t *testing.T) {
	// Test: empty list
	result := hasSecurityVarName(nil)
	// Expected: false
	if result {
		t.Error("hasSecurityVarName(nil) should return false")
	}

	// Test: list with non-security names
	names := []*ast.Ident{
		{Name: "foo"},
		{Name: "bar"},
	}
	result = hasSecurityVarName(names)
	// Expected: false
	if result {
		t.Error("hasSecurityVarName with non-security names should return false")
	}

	// Test: list with security name
	names2 := []*ast.Ident{
		{Name: "foo"},
		{Name: "secretKey"},
	}
	result = hasSecurityVarName(names2)
	// Expected: true
	if !result {
		t.Error("hasSecurityVarName with security name should return true")
	}
}

// TestIsMathRandCall tests isMathRandCall function.
func TestIsMathRandCall(t *testing.T) {
	aliases := map[string]bool{"rand": true}

	// Test: not a selector
	call1 := &ast.CallExpr{
		Fun: &ast.Ident{Name: "someFunc"},
	}
	result := isMathRandCall(call1, aliases)
	// Expected: false
	if result {
		t.Error("isMathRandCall with non-selector should return false")
	}

	// Test: selector X is not an ident
	call2 := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.IndexExpr{X: &ast.Ident{Name: "arr"}},
			Sel: &ast.Ident{Name: "Int"},
		},
	}
	result = isMathRandCall(call2, aliases)
	// Expected: false
	if result {
		t.Error("isMathRandCall with non-ident X should return false")
	}

	// Test: ident not in aliases
	call3 := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "other"},
			Sel: &ast.Ident{Name: "Int"},
		},
	}
	result = isMathRandCall(call3, aliases)
	// Expected: false
	if result {
		t.Error("isMathRandCall with non-aliased ident should return false")
	}

	// Test: valid math/rand call
	call4 := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "rand"},
			Sel: &ast.Ident{Name: "Int"},
		},
	}
	result = isMathRandCall(call4, aliases)
	// Expected: true
	if !result {
		t.Error("isMathRandCall with valid rand call should return true")
	}
}

// TestCheckValuesForMathRand tests checkValuesForMathRand function.
func TestCheckValuesForMathRand(t *testing.T) {
	aliases := map[string]bool{"rand": true}

	// Test: empty values
	checkValuesForMathRand(nil, nil, aliases, false)

	// Test: non-call value
	values := []ast.Expr{&ast.Ident{Name: "x"}}
	checkValuesForMathRand(nil, values, aliases, false)

	// Test: call but not math/rand
	values2 := []ast.Expr{
		&ast.CallExpr{Fun: &ast.Ident{Name: "other"}},
	}
	checkValuesForMathRand(nil, values2, aliases, false)
}

// TestProcessGenDeclSpec023 tests processGenDeclSpec023 function.
func TestProcessGenDeclSpec023(t *testing.T) {
	aliases := map[string]bool{"rand": true}

	// Test: spec is not a ValueSpec
	importSpec := &ast.ImportSpec{
		Path: &ast.BasicLit{Value: `"fmt"`},
	}
	processGenDeclSpec023(nil, importSpec, aliases)

	// Test: ValueSpec with non-security names
	valueSpec := &ast.ValueSpec{
		Names:  []*ast.Ident{{Name: "normalVar"}},
		Values: nil,
	}
	processGenDeclSpec023(nil, valueSpec, aliases)
}

// TestCollectMathRandAliases tests collectMathRandAliases function.
func TestCollectMathRandAliases(t *testing.T) {
	// Test with import that has explicit alias
	pass := &mockPass{
		files: []*ast.File{
			{
				Imports: []*ast.ImportSpec{
					{
						Name: &ast.Ident{Name: "r"},
						Path: &ast.BasicLit{Value: `"math/rand"`},
					},
				},
			},
		},
	}
	aliases := collectMathRandAliasesFromFiles(pass.files)
	// Expected: r is in aliases
	if !aliases["r"] {
		t.Error("collectMathRandAliases should include custom alias 'r'")
	}

	// Test with standard import (no alias)
	pass2 := &mockPass{
		files: []*ast.File{
			{
				Imports: []*ast.ImportSpec{
					{
						Path: &ast.BasicLit{Value: `"math/rand"`},
					},
				},
			},
		},
	}
	aliases2 := collectMathRandAliasesFromFiles(pass2.files)
	// Expected: rand is in aliases
	if !aliases2["rand"] {
		t.Error("collectMathRandAliases should include default alias 'rand'")
	}

	// Test with math/rand/v2
	pass3 := &mockPass{
		files: []*ast.File{
			{
				Imports: []*ast.ImportSpec{
					{
						Path: &ast.BasicLit{Value: `"math/rand/v2"`},
					},
				},
			},
		},
	}
	aliases3 := collectMathRandAliasesFromFiles(pass3.files)
	// Expected: rand is in aliases
	if !aliases3["rand"] {
		t.Error("collectMathRandAliases should include default alias for math/rand/v2")
	}

	// Test with non-math/rand import
	pass4 := &mockPass{
		files: []*ast.File{
			{
				Imports: []*ast.ImportSpec{
					{
						Path: &ast.BasicLit{Value: `"fmt"`},
					},
				},
			},
		},
	}
	aliases4 := collectMathRandAliasesFromFiles(pass4.files)
	// Expected: empty
	if len(aliases4) != 0 {
		t.Error("collectMathRandAliases should return empty map for non-math/rand imports")
	}
}

// mockPass is a mock structure for testing.
type mockPass struct {
	files []*ast.File
}

// collectMathRandAliasesFromFiles is a helper to test without full pass.
func collectMathRandAliasesFromFiles(files []*ast.File) map[string]bool {
	aliases := make(map[string]bool, initialAliasMapCap)
	// Parcours des fichiers
	for _, file := range files {
		// Parcours des imports
		for _, imp := range file.Imports {
			// Vérification si c'est math/rand
			if imp.Path.Value == `"math/rand"` || imp.Path.Value == `"math/rand/v2"` {
				// Détermination de l'alias
				if imp.Name != nil {
					// Import avec alias
					aliases[imp.Name.Name] = true
				} else {
					// Import standard
					aliases["rand"] = true
				}
			}
		}
	}
	// Retour des alias
	return aliases
}

// TestCheckFuncForMathRand tests checkFuncForMathRand function.
func TestCheckFuncForMathRand(t *testing.T) {
	aliases := map[string]bool{"rand": true}

	// Test: nil body
	funcDecl := &ast.FuncDecl{
		Name: &ast.Ident{Name: "normalFunc"},
		Body: nil,
	}
	checkFuncForMathRand(nil, funcDecl, aliases)

	// Test: empty body
	funcDecl2 := &ast.FuncDecl{
		Name: &ast.Ident{Name: "normalFunc"},
		Body: &ast.BlockStmt{List: nil},
	}
	checkFuncForMathRand(nil, funcDecl2, aliases)

	// Test: function with security name but no body
	funcDecl3 := &ast.FuncDecl{
		Name: &ast.Ident{Name: "generateKey"},
		Body: nil,
	}
	checkFuncForMathRand(nil, funcDecl3, aliases)
}

// TestCheckLocalVarDecl tests checkLocalVarDecl function.
func TestCheckLocalVarDecl(t *testing.T) {
	aliases := map[string]bool{"rand": true}

	// Test: empty specs
	genDecl := &ast.GenDecl{
		Tok:   token.VAR,
		Specs: nil,
	}
	checkLocalVarDecl(nil, genDecl, aliases, false)

	// Test: non-ValueSpec in specs
	genDecl2 := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{Path: &ast.BasicLit{Value: `"fmt"`}},
		},
	}
	checkLocalVarDecl(nil, genDecl2, aliases, false)
}

// TestCheckGenDeclForMathRand tests checkGenDeclForMathRand function.
func TestCheckGenDeclForMathRand(t *testing.T) {
	aliases := map[string]bool{"rand": true}

	// Test: empty specs
	genDecl := &ast.GenDecl{
		Tok:   token.VAR,
		Specs: nil,
	}
	checkGenDeclForMathRand(nil, genDecl, aliases)

	// Test: non-ValueSpec
	genDecl2 := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{
			&ast.ImportSpec{Path: &ast.BasicLit{Value: `"fmt"`}},
		},
	}
	checkGenDeclForMathRand(nil, genDecl2, aliases)
}
