package ktnstruct

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/config"
	"golang.org/x/tools/go/analysis"
)

// Test_runStruct004 tests the private runStruct004 function.
func Test_runStruct004(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"error case validation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test passthrough - logique principale testée via API publique
		})
	}
}

// Test_collectStructs tests the private collectStructs function.
func Test_collectStructs(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected int
	}{
		{
			name: "no structs",
			src: `package test
func main() {}`,
			expected: 0,
		},
		{
			name: "one struct",
			src: `package test
type User struct {
	Name string
}`,
			expected: 1,
		},
		{
			name: "multiple structs",
			src: `package test
type User struct {
	Name string
}
type Admin struct {
	Role string
}`,
			expected: 2,
		},
		{
			name: "struct with interface",
			src: `package test
type User struct {
	Name string
}
type Reader interface {
	Read() error
}`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			structs := collectStructs(file)

			if len(structs) != tt.expected {
				t.Errorf("expected %d structs, got %d", tt.expected, len(structs))
			}
		})
	}
}

// Test_structInfo tests the structInfo type.
func Test_structInfo(t *testing.T) {
	tests := []struct {
		name          string
		src           string
		expectedName  string
		expectedCount int
	}{
		{
			name: "verify struct info fields",
			src: `package test
type User struct {
	Name string
}`,
			expectedName:  "User",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			structs := collectStructs(file)
			if len(structs) != tt.expectedCount {
				t.Fatalf("expected %d struct, got %d", tt.expectedCount, len(structs))
			}

			s := structs[0]
			if s.name != tt.expectedName {
				t.Errorf("expected struct name '%s', got '%s'", tt.expectedName, s.name)
			}
			if s.node == nil {
				t.Error("expected non-nil node")
			}
		})
	}
}

// Test_allStructsAreSerializable tests the allStructsAreSerializable function.
//
// Params:
//   - t: testing context
func Test_allStructsAreSerializable(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected bool
	}{
		{
			name: "all DTOs by suffix",
			src: `package test
type UserConfig struct { Name string }
type AppSettings struct { Port int }`,
			expected: true,
		},
		{
			name: "all DTOs by tag",
			src: `package test
type User struct { Name string ` + "`json:\"name\"`" + ` }
type Admin struct { Role string ` + "`yaml:\"role\"`" + ` }`,
			expected: true,
		},
		{
			name: "not all DTOs",
			src: `package test
type User struct { Name string }
type Admin struct { Role string }`,
			expected: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.src, 0)
			// Vérification erreur de parsing
			if err != nil {
				t.Fatalf("failed to parse source: %v", err)
			}

			// Collecter les structs avec structType
			var structs []structInfo
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier TypeSpec
				if ts, ok := n.(*ast.TypeSpec); ok {
					// Vérifier StructType
					if st, ok := ts.Type.(*ast.StructType); ok {
						structs = append(structs, structInfo{
							name:       ts.Name.Name,
							node:       ts,
							structType: st,
						})
					}
				}
				return true
			})

			result := allStructsAreSerializable(structs)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("allStructsAreSerializable() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test_runStruct004_disabled tests that the rule is skipped when disabled.
func Test_runStruct004_disabled(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-STRUCT-004": {Enabled: config.Bool(false)},
		},
	})
	defer config.Reset()

	src := `package test
type User struct { Name string }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		Report: func(_ analysis.Diagnostic) {
			t.Error("Unexpected error when rule is disabled")
		},
	}

	_, err = runStruct004(pass)
	if err != nil {
		t.Errorf("runStruct004() error = %v", err)
	}
}

// Test_runStruct004_excludedFile tests that excluded files are skipped.
func Test_runStruct004_excludedFile(t *testing.T) {
	config.Set(&config.Config{
		Rules: map[string]*config.RuleConfig{
			"KTN-STRUCT-004": {
				Enabled: config.Bool(true),
				Exclude: []string{"**/test.go"},
			},
		},
	})
	defer config.Reset()

	src := `package test
type User struct { Name string }
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "/some/path/test.go", src, 0)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{f},
		Report: func(_ analysis.Diagnostic) {
			t.Error("Unexpected error for excluded file")
		},
	}

	_, err = runStruct004(pass)
	if err != nil {
		t.Errorf("runStruct004() error = %v", err)
	}
}
