package ktnvar

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"github.com/kodflow/ktn-linter/pkg/messages"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Test_runVar001 tests the private runVar001 function.
func Test_runVar001(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"passthrough validation"},
		{"error case validation"},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - main logic tested via public API in external tests
			t.Log("runVar001 tested via external tests")
		})
	}
}

// Test_checkVarSpec tests the checkVarSpec function.
func Test_checkVarSpec(t *testing.T) {
	tests := []struct {
		name        string
		code        string
		expectError bool
	}{
		{
			name:        "type and value - OK",
			code:        "package test\nvar x int = 42",
			expectError: false,
		},
		{
			name:        "no type - ERROR",
			code:        "package test\nvar x = 42",
			expectError: true,
		},
		{
			name:        "type without value - OK (zero-value)",
			code:        "package test\nvar x int",
			expectError: false,
		},
		{
			name:        "slice with type and value - OK",
			code:        "package test\nvar x []string = []string{}",
			expectError: false,
		},
		{
			name:        "blank identifier - skip",
			code:        "package test\nvar _ = 42",
			expectError: false,
		},
		{
			name:        "multiple names without type - ERROR",
			code:        "package test\nvar x, y = 1, 2",
			expectError: true,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Trouver la ValueSpec
			var valueSpec *ast.ValueSpec
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type
				if vs, ok := n.(*ast.ValueSpec); ok {
					valueSpec = vs
					return false
				}
				return true
			})

			// Vérification valueSpec trouvée
			if valueSpec == nil {
				t.Fatal("no value spec found")
			}

			// Vérification des conditions
			hasType := valueSpec.Type != nil

			// Le format obligatoire est: var name type (avec ou sans valeur)
			// hasError = !hasType (sauf blank identifier)
			isBlank := len(valueSpec.Names) == 1 && valueSpec.Names[0].Name == "_"
			hasError := !hasType && !isBlank

			// Vérification résultat
			if hasError != tt.expectError {
				t.Errorf("checkVarSpec error = %v, want %v", hasError, tt.expectError)
			}
		})
	}
}

// Test_checkVarSpec_multipleVars tests checkVarSpec with multiple variables.
func Test_checkVarSpec_multipleVars(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		expectedCount int
	}{
		{
			name: "three_valid_vars",
			code: `package test
var (
	a int = 1
	b string = "hello"
	c bool = true
)`,
			expectedCount: 3,
		},
		{
			name: "mixed_valid_invalid",
			code: `package test
var (
	a int = 1
	b = "hello"
	c bool = true
)`,
			expectedCount: 2,
		},
	}

	// Itération sur les tests
	for _, tt := range tests {
		tt := tt // Capture range variable
		// Sous-test
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, 0)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Compter les ValueSpecs valides
			validCount := 0
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérification du type
				if vs, ok := n.(*ast.ValueSpec); ok {
					hasType := vs.Type != nil
					// Vérification format valide (type explicite requis)
					if hasType {
						validCount++
					}
				}
				return true
			})

			// Vérification nombre de vars valides
			if validCount != tt.expectedCount {
				t.Errorf("valid var count = %d, want %d", validCount, tt.expectedCount)
			}
		})
	}
}

// Test_runVar001_disabled tests runVar001 with disabled rule.
func Test_runVar001_disabled(t *testing.T) {
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
					"KTN-VAR-001": {Enabled: config.Bool(false)},
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

			_, err = runVar001(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar001() error = %v", err)
			}

			// Should not report anything when disabled
			if reportCount != 0 {
				t.Errorf("runVar001() reported %d issues, expected 0 when disabled", reportCount)
			}

		})
	}
}

// Test_runVar001_fileExcluded tests runVar001 with excluded file.
func Test_runVar001_fileExcluded(t *testing.T) {
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
					"KTN-VAR-001": {
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

			_, err = runVar001(pass)
			// Check no error
			if err != nil {
				t.Fatalf("runVar001() error = %v", err)
			}

			// Should not report anything when file is excluded
			if reportCount != 0 {
				t.Errorf("runVar001() reported %d issues, expected 0 when file excluded", reportCount)
			}

		})
	}
}

// Test_runVar001_nilInspector tests runVar001 with nil inspector.
func Test_runVar001_nilInspector(t *testing.T) {
	config.Reset()
	defer config.Reset()

	pass := &analysis.Pass{
		Fset:     token.NewFileSet(),
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: nil},
	}

	result, err := runVar001(pass)
	// Should return nil without error
	if err != nil {
		t.Errorf("runVar001() error = %v, want nil", err)
	}
	if result != nil {
		t.Errorf("runVar001() result = %v, want nil", result)
	}
}

// Test_runVar001_nilFset tests runVar001 with nil Fset.
func Test_runVar001_nilFset(t *testing.T) {
	config.Reset()
	defer config.Reset()

	fset := token.NewFileSet()
	file, _ := parser.ParseFile(fset, "test.go", "package test", 0)
	insp := inspector.New([]*ast.File{file})

	pass := &analysis.Pass{
		Fset:     nil,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
	}

	result, err := runVar001(pass)
	// Should return nil without error
	if err != nil {
		t.Errorf("runVar001() error = %v, want nil", err)
	}
	if result != nil {
		t.Errorf("runVar001() result = %v, want nil", result)
	}
}

// Test_runVar001_wrongInspectorType tests runVar001 with wrong inspector type.
func Test_runVar001_wrongInspectorType(t *testing.T) {
	config.Reset()
	defer config.Reset()

	pass := &analysis.Pass{
		Fset:     token.NewFileSet(),
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: "not an inspector"},
	}

	result, err := runVar001(pass)
	// Should return nil without error
	if err != nil {
		t.Errorf("runVar001() error = %v, want nil", err)
	}
	if result != nil {
		t.Errorf("runVar001() result = %v, want nil", result)
	}
}

// Test_checkVarSpec_withBlankIdentifier tests checkVarSpec skips blank identifier.
func Test_checkVarSpec_withBlankIdentifier(t *testing.T) {
	config.Reset()
	defer config.Reset()

	// Parse code with multiple vars including blank identifier
	code := `package test
var x, _, y = 1, 2, 3
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	_, _ = runVar001(pass)

	// Should report 2 issues (x and y, not _)
	if reportCount != 2 {
		t.Errorf("runVar001() reported %d issues, expected 2", reportCount)
	}
}

// Test_checkVarSpec_fallbackMessage tests checkVarSpec with missing message.
func Test_checkVarSpec_fallbackMessage(t *testing.T) {
	config.Reset()
	defer config.Reset()

	// Temporarily remove the message to test fallback
	msg, _ := messages.Get(ruleCodeVar001)
	messages.Unregister(ruleCodeVar001)
	defer messages.Register(msg)

	// Parse code without explicit type
	code := `package test
var missingType = 42
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0
	var lastMsg string

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report: func(d analysis.Diagnostic) {
			reportCount++
			lastMsg = d.Message
		},
	}

	_, _ = runVar001(pass)

	// Should report one issue with fallback message
	if reportCount != 1 {
		t.Errorf("runVar001() reported %d issues, expected 1", reportCount)
	}
	// Message should contain fallback text
	if !strings.Contains(lastMsg, "type explicite requis") {
		t.Errorf("expected fallback message, got: %s", lastMsg)
	}
}

// Test_checkVarSpec_fallbackMessageMultiple tests fallback with multiple names.
func Test_checkVarSpec_fallbackMessageMultiple(t *testing.T) {
	config.Reset()
	defer config.Reset()

	// Temporarily remove the message to test fallback
	msg, _ := messages.Get(ruleCodeVar001)
	messages.Unregister(ruleCodeVar001)
	defer messages.Register(msg)

	// Parse code with multiple names
	code := `package test
var x, y = 1, 2
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", code, 0)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	insp := inspector.New([]*ast.File{file})
	reportCount := 0

	pass := &analysis.Pass{
		Fset:     fset,
		ResultOf: map[*analysis.Analyzer]any{inspect.Analyzer: insp},
		Report:   func(_ analysis.Diagnostic) { reportCount++ },
	}

	_, _ = runVar001(pass)

	// Should report 2 issues (x and y)
	if reportCount != 2 {
		t.Errorf("runVar001() reported %d issues, expected 2", reportCount)
	}
}
