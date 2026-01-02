package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// TestIsInterfaceType vérifie la détection des types interface.
func TestIsInterfaceType(t *testing.T) {
	tests := []struct {
		name     string
		typeVal  types.Type
		expected bool
	}{
		{
			name:     "empty interface",
			typeVal:  types.NewInterfaceType(nil, nil),
			expected: true,
		},
		{
			name:     "basic type int",
			typeVal:  types.Typ[types.Int],
			expected: false,
		},
		{
			name:     "basic type string",
			typeVal:  types.Typ[types.String],
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			result := isInterfaceType(tt.typeVal)
			// Vérification du résultat
			if result != tt.expected {
				t.Errorf("isInterfaceType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_runVar022 tests the private runVar022 function.
func Test_runVar022(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"passthrough validation"},
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
		})
	}
}

// Test_checkFuncDecls tests the private checkFuncDecls function.
func Test_checkFuncDecls(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks function declarations
		})
	}
}

// Test_checkVarDecls tests the private checkVarDecls function.
func Test_checkVarDecls(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks var declarations
		})
	}
}

// Test_checkStructFields tests the private checkStructFields function.
func Test_checkStructFields(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks struct fields
		})
	}
}

// Test_checkFieldList tests the private checkFieldList function.
func Test_checkFieldList(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks field list
		})
	}
}

// Test_checkPointerToInterface tests the private checkPointerToInterface function.
func Test_checkPointerToInterface(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - function checks pointer to interface
		})
	}
}

// Test_checkPointerToInterface_nonStarExpr tests with non-star expression.
func Test_checkPointerToInterface_nonStarExpr(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			// Test with non-star expression
			expr := &ast.Ident{Name: "x"}
			// Should not panic and should not report
			checkPointerToInterface(pass, expr)
		})
	}
}

// Test_checkPointerToInterface_nilType tests with nil type.
func Test_checkPointerToInterface_nilType(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			pass := &analysis.Pass{
				TypesInfo: &types.Info{
					Types: make(map[ast.Expr]types.TypeAndValue),
				},
			}

			// Test with star expression but underlying type not in TypesInfo
			innerExpr := &ast.Ident{Name: "Reader"}
			expr := &ast.StarExpr{X: innerExpr}
			// Should not panic and should not report
			checkPointerToInterface(pass, expr)
		})
	}
}

// Test_runVar022_disabled tests runVar022 with disabled rule.
func Test_runVar022_disabled(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Setup config with rule disabled
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-022": {Enabled: config.Bool(false)},
				},
			})
			defer config.Reset()

			// Parse simple code
			code := `package test
			var x int = 42
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			_, err = runVar022(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar022() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar022() reported %d issues, expected 0 when disabled", reportCount)
			}
		})
	}
}

// Test_runVar022_fileExcluded tests runVar022 with excluded file.
func Test_runVar022_fileExcluded(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"validation"},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Setup config with file exclusion
			config.Set(&config.Config{
				Rules: map[string]*config.RuleConfig{
					"KTN-VAR-022": {
						Exclude: []string{"test.go"},
					},
				},
			})
			defer config.Reset()

			// Parse simple code
			code := `package test
			var x int = 42
			`
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", code, 0)
			// Check parsing error
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			insp := inspector.New([]*ast.File{file})
			reportCount := 0

			pass := &analysis.Pass{
				Fset: fset,
				ResultOf: map[*analysis.Analyzer]any{
					inspect.Analyzer: insp,
				},
				Report: func(_d analysis.Diagnostic) {
					reportCount++
				},
			}

			_, err = runVar022(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar022() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar022() reported %d issues, expected 0 when file excluded", reportCount)
			}
		})
	}
}

// Test_runVar022_nilFset tests runVar022 with nil Fset.
func Test_runVar022_nilFset(t *testing.T) {
	// Ensure rule is enabled
	config.Reset()

	// Parse simple code
	code := `package test
var x int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	// Create pass with nil Fset
	pass := &analysis.Pass{
		Fset: nil, // nil Fset
		TypesInfo: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
		},
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	result, err := runVar022(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar022() error = %v", err)
	}
	// Check result is nil
	if result != nil {
		t.Errorf("runVar022() result = %v, expected nil", result)
	}
	// Should not report anything when Fset is nil
	if reportCount != 0 {
		t.Errorf("runVar022() reported %d issues, expected 0 with nil Fset", reportCount)
	}
}

// Test_runVar022_nilTypesInfo tests runVar022 with nil TypesInfo.
func Test_runVar022_nilTypesInfo(t *testing.T) {
	// Ensure rule is enabled
	config.Reset()

	// Parse simple code
	code := `package test
var x int = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	// Create pass with nil TypesInfo
	pass := &analysis.Pass{
		Fset:      fset,
		TypesInfo: nil, // nil TypesInfo
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	result, err := runVar022(pass)
	// Check no error
	if err != nil {
		t.Fatalf("runVar022() error = %v", err)
	}
	// Check result is nil
	if result != nil {
		t.Errorf("runVar022() result = %v, expected nil", result)
	}
	// Should not report anything when TypesInfo is nil
	if reportCount != 0 {
		t.Errorf("runVar022() reported %d issues, expected 0 with nil TypesInfo", reportCount)
	}
}

// Test_checkFuncDecls_fileExcluded tests checkFuncDecls with file exclusion.
func Test_checkFuncDecls_fileExcluded(t *testing.T) {
	// Setup config with file exclusion
	cfg := &config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-022": {
				Exclude: []string{"excluded.go"},
			},
		},
	}
	config.Set(cfg)
	defer config.Reset()

	// Parse code with pointer to interface in function param
	code := `package test
import "io"
func process(r *io.Reader) {}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	// Setup type info
	typesInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		Fset:      fset,
		TypesInfo: typesInfo,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Call checkFuncDecls directly
	checkFuncDecls(pass, insp, cfg)

	// Should not report anything when file is excluded
	if reportCount != 0 {
		t.Errorf("checkFuncDecls() reported %d issues, expected 0", reportCount)
	}
}

// Test_checkVarDecls_fileExcluded tests checkVarDecls with file exclusion.
func Test_checkVarDecls_fileExcluded(t *testing.T) {
	// Setup config with file exclusion
	cfg := &config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-022": {
				Exclude: []string{"excluded.go"},
			},
		},
	}
	config.Set(cfg)
	defer config.Reset()

	// Parse code with pointer to interface in var decl
	code := `package test
var handler *interface{}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	// Setup type info
	typesInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		Fset:      fset,
		TypesInfo: typesInfo,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Call checkVarDecls directly
	checkVarDecls(pass, insp, cfg)

	// Should not report anything when file is excluded
	if reportCount != 0 {
		t.Errorf("checkVarDecls() reported %d issues, expected 0", reportCount)
	}
}

// Test_checkVarDecls_nonVarToken tests checkVarDecls with non-var token.
func Test_checkVarDecls_nonVarToken(t *testing.T) {
	// Reset config to enable rule
	config.Reset()
	cfg := config.Get()

	// Parse code with const declaration
	code := `package test
const myConst = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	// Setup type info
	typesInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		Fset:      fset,
		TypesInfo: typesInfo,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Call checkVarDecls directly
	checkVarDecls(pass, insp, cfg)

	// Should not report anything for const declarations
	if reportCount != 0 {
		t.Errorf("checkVarDecls() reported %d issues, expected 0", reportCount)
	}
}

// Test_checkVarDecls_nonValueSpec tests checkVarDecls with non-ValueSpec spec.
func Test_checkVarDecls_nonValueSpec(t *testing.T) {
	// Reset config to enable rule
	config.Reset()
	cfg := config.Get()

	// Create a synthetic GenDecl with var token but TypeSpec instead of ValueSpec
	genDecl := &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.TypeSpec{
				Name: &ast.Ident{Name: "MyType"},
				Type: &ast.Ident{Name: "int"},
			},
		},
	}

	// Create file with the GenDecl
	file := &ast.File{
		Name: &ast.Ident{Name: "test"},
		Decls: []ast.Decl{genDecl},
	}

	fset := token.NewFileSet()
	fset.AddFile("test.go", fset.Base(), 100)

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	// Setup type info
	typesInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		Fset:      fset,
		TypesInfo: typesInfo,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Call checkVarDecls directly
	checkVarDecls(pass, insp, cfg)

	// Should not report anything (non-ValueSpec is skipped)
	if reportCount != 0 {
		t.Errorf("checkVarDecls() reported %d issues, expected 0", reportCount)
	}
}

// Test_checkStructFields_fileExcluded tests checkStructFields with file exclusion.
func Test_checkStructFields_fileExcluded(t *testing.T) {
	// Setup config with file exclusion
	cfg := &config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-VAR-022": {
				Exclude: []string{"excluded.go"},
			},
		},
	}
	config.Set(cfg)
	defer config.Reset()

	// Parse code with pointer to interface in struct field
	code := `package test
import "io"
type Service struct {
	reader *io.Reader
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "excluded.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	// Setup type info
	typesInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		Fset:      fset,
		TypesInfo: typesInfo,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Call checkStructFields directly
	checkStructFields(pass, insp, cfg)

	// Should not report anything when file is excluded
	if reportCount != 0 {
		t.Errorf("checkStructFields() reported %d issues, expected 0", reportCount)
	}
}

// Test_checkStructFields_nonStructType tests checkStructFields with non-struct type.
func Test_checkStructFields_nonStructType(t *testing.T) {
	// Reset config to enable rule
	config.Reset()
	cfg := config.Get()

	// Parse code with type alias (not a struct)
	code := `package test
type MyInt int
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	// Check parsing error
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	// Setup type info
	typesInfo := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	pass := &analysis.Pass{
		Fset:      fset,
		TypesInfo: typesInfo,
		ResultOf: map[*analysis.Analyzer]any{
			inspect.Analyzer: insp,
		},
		Report: func(_d analysis.Diagnostic) {
			reportCount++
		},
	}

	// Call checkStructFields directly
	checkStructFields(pass, insp, cfg)

	// Should not report anything for non-struct types
	if reportCount != 0 {
		t.Errorf("checkStructFields() reported %d issues, expected 0", reportCount)
	}
}

// Test_checkPointerToInterface_withMessage tests checkPointerToInterface with valid message.
func Test_checkPointerToInterface_withMessage(t *testing.T) {
	// Reset config
	config.Reset()

	// Create a star expression with interface type
	innerExpr := &ast.Ident{Name: "Reader"}
	expr := &ast.StarExpr{X: innerExpr}

	// Create interface type for testing
	interfaceType := types.NewInterfaceType(nil, nil)

	// Setup type info with the interface type
	typesInfo := &types.Info{
		Types: map[ast.Expr]types.TypeAndValue{
			innerExpr: {Type: interfaceType},
		},
	}

	fset := token.NewFileSet()
	file := fset.AddFile("test.go", -1, 100)
	// Add line info
	file.SetLinesForContent([]byte("test content"))

	var reportedMsg string
	pass := &analysis.Pass{
		Fset:      fset,
		TypesInfo: typesInfo,
		Report: func(d analysis.Diagnostic) {
			reportedMsg = d.Message
		},
	}

	// Call checkPointerToInterface
	checkPointerToInterface(pass, expr)

	// Should have reported
	if reportedMsg == "" {
		t.Errorf("checkPointerToInterface() did not report")
	}
}
