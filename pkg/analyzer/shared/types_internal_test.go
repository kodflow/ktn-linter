// Internal tests for types.go.
package shared

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// Test_hasSerializableSuffix tests the hasSerializableSuffix function.
//
// Params:
//   - t: testing context
func Test_hasSerializableSuffix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Config suffix", "AppConfig", true},
		{"Request suffix", "UserRequest", true},
		{"Response suffix", "APIResponse", true},
		{"DTO suffix", "UserDTO", true},
		{"No suffix", "User", false},
		{"Empty", "", false},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasSerializableSuffix(tt.input)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("hasSerializableSuffix(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Test_hasSerializationTags tests the hasSerializationTags function.
//
// Params:
//   - t: testing context
func Test_hasSerializationTags(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name: "with json tag",
			code: `package test
type User struct {
	Name string ` + "`json:\"name\"`" + `
}`,
			expected: true,
		},
		{
			name: "without tags",
			code: `package test
type User struct {
	Name string
}`,
			expected: false,
		},
		{
			name: "empty struct",
			code: `package test
type User struct {}`,
			expected: false,
		},
	}

	// Parcourir les cas de test
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tt.code, parser.ParseComments)
			// Vérification erreur
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			// Trouver la struct
			var structType *ast.StructType
			ast.Inspect(file, func(n ast.Node) bool {
				// Vérifier TypeSpec
				if ts, ok := n.(*ast.TypeSpec); ok {
					// Vérifier StructType
					if st, ok := ts.Type.(*ast.StructType); ok {
						structType = st
						return false
					}
				}
				return true
			})

			// Vérification struct trouvée
			if structType == nil {
				t.Fatal("no struct found")
			}

			result := hasSerializationTags(structType)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("hasSerializationTags() = %v, want %v", result, tt.expected)
			}
		})
	}
}
