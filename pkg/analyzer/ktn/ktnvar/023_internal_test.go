package ktnvar

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
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
	aliases := map[string]bool{"rand": true}

	tests := []struct {
		name         string
		assign       *ast.AssignStmt
		secContext   bool
		expectReport bool
	}{
		{
			name: "RHS is not a call expression",
			assign: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
				Rhs: []ast.Expr{&ast.Ident{Name: "y"}},
			},
			secContext:   false,
			expectReport: false,
		},
		{
			name: "RHS is a call but not math/rand",
			assign: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.Ident{Name: "x"}},
				Rhs: []ast.Expr{&ast.CallExpr{Fun: &ast.Ident{Name: "someFunc"}}},
			},
			secContext:   false,
			expectReport: false,
		},
		{
			name: "LHS is empty",
			assign: &ast.AssignStmt{
				Lhs: []ast.Expr{},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "rand"},
							Sel: &ast.Ident{Name: "Int"},
						},
					},
				},
			},
			secContext:   false,
			expectReport: false,
		},
		{
			name: "LHS is not an ident",
			assign: &ast.AssignStmt{
				Lhs: []ast.Expr{&ast.IndexExpr{X: &ast.Ident{Name: "arr"}}},
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: "rand"},
							Sel: &ast.Ident{Name: "Int"},
						},
					},
				},
			},
			secContext:   false,
			expectReport: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			reported := false
			pass := &analysis.Pass{
				Fset: token.NewFileSet(),
				Report: func(_ analysis.Diagnostic) {
					reported = true
				},
			}
			// Call function
			checkAssignForMathRand(pass, tt.assign, aliases, tt.secContext)
			// Check result
			if reported != tt.expectReport {
				t.Errorf("checkAssignForMathRand() reported = %v, want %v", reported, tt.expectReport)
			}
		})
	}
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

// TestRunVar023_RuleDisabled tests runVar023 when the rule is disabled.
func TestRunVar023_RuleDisabled(t *testing.T) {
	// Save original config
	originalCfg := config.Get()
	defer config.Set(originalCfg)

	// Create config with rule disabled
	disabledCfg := config.DefaultConfig()
	disabled := false
	disabledCfg.Rules = map[string]*config.RuleConfig{
		"KTN-VAR-023": {Enabled: &disabled},
	}
	config.Set(disabledCfg)

	// Create a minimal pass
	pass := createTestPass023(t, `package test
import "math/rand"
func generateKey() int { return rand.Intn(100) }
`)
	// Run the analyzer
	result, err := runVar023(pass)
	// Check no error
	if err != nil {
		t.Errorf("runVar023() unexpected error: %v", err)
	}
	// Check result is nil
	if result != nil {
		t.Errorf("runVar023() expected nil result, got %v", result)
	}
}

// TestRunVar023_NoMathRandImport tests runVar023 when math/rand is not imported.
func TestRunVar023_NoMathRandImport(t *testing.T) {
	// Save original config
	originalCfg := config.Get()
	defer config.Set(originalCfg)

	// Ensure rule is enabled
	config.Reset()

	// Create a pass without math/rand import
	pass := createTestPass023(t, `package test
import "fmt"
func generateKey() int { fmt.Println("key"); return 100 }
`)
	// Run the analyzer
	result, err := runVar023(pass)
	// Check no error
	if err != nil {
		t.Errorf("runVar023() unexpected error: %v", err)
	}
	// Check result is nil
	if result != nil {
		t.Errorf("runVar023() expected nil result, got %v", result)
	}
}

// TestCheckMathRandUsage_FileExcluded tests file exclusion path.
func TestCheckMathRandUsage_FileExcluded(t *testing.T) {
	// Save original config
	originalCfg := config.Get()
	defer config.Set(originalCfg)

	// Create config with file exclusion
	excludedCfg := config.DefaultConfig()
	excludedCfg.Rules = map[string]*config.RuleConfig{
		"KTN-VAR-023": {Exclude: []string{"*_test.go", "excluded.go"}},
	}
	config.Set(excludedCfg)

	// Create a pass for excluded file
	pass := createTestPass023WithFilename(t, `package test
import "math/rand"
func generateKey() int { return rand.Intn(100) }
`, "excluded.go")

	var diagnostics int
	pass.Report = func(d analysis.Diagnostic) {
		diagnostics++
	}

	// Run the analyzer
	result, err := runVar023(pass)
	// Check no error
	if err != nil {
		t.Errorf("runVar023() unexpected error: %v", err)
	}
	// Check result is nil
	if result != nil {
		t.Errorf("runVar023() expected nil result, got %v", result)
	}
	// Check no diagnostics were reported
	if diagnostics != 0 {
		t.Errorf("runVar023() expected 0 diagnostics for excluded file, got %d", diagnostics)
	}
}

// TestCheckAssignForMathRand_FuncSecurityContext tests with funcIsSecurityContext = true.
func TestCheckAssignForMathRand_FuncSecurityContext(t *testing.T) {
	aliases := map[string]bool{"rand": true}

	var reported int
	mockPass := &analysis.Pass{
		Fset: token.NewFileSet(),
		Report: func(d analysis.Diagnostic) {
			reported++
		},
	}

	// Test: funcIsSecurityContext = true, variable name not security-related
	assign := &ast.AssignStmt{
		Lhs: []ast.Expr{&ast.Ident{Name: "value"}}, // Not a security name
		Rhs: []ast.Expr{
			&ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X:   &ast.Ident{Name: "rand"},
					Sel: &ast.Ident{Name: "Intn"},
				},
			},
		},
	}

	// Call with funcIsSecurityContext = true
	checkAssignForMathRand(mockPass, assign, aliases, true)

	// Check that diagnostic was reported
	if reported != 1 {
		t.Errorf("checkAssignForMathRand() expected 1 diagnostic with funcIsSecurityContext=true, got %d", reported)
	}
}

// createTestPass023 creates a test pass for testing runVar023.
func createTestPass023(t *testing.T, code string) *analysis.Pass {
	return createTestPass023WithFilename(t, code, "test.go")
}

// createTestPass023WithFilename creates a test pass with a specific filename.
func createTestPass023WithFilename(t *testing.T, code string, filename string) *analysis.Pass {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, code, parser.ParseComments)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse code: %v", err)
	}

	// Type check
	conf := &types.Config{
		Importer: importer.Default(),
		Error:    func(err error) {}, // Ignore errors
	}
	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	pkg, _ := conf.Check(file.Name.Name, fset, []*ast.File{file}, info)

	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		Report:    func(d analysis.Diagnostic) {},
		ResultOf:  make(map[*analysis.Analyzer]any),
	}

	// Run inspect.Analyzer to populate ResultOf
	inspResult, _ := inspect.Analyzer.Run(pass)
	pass.ResultOf[inspect.Analyzer] = inspResult

	return pass
}
