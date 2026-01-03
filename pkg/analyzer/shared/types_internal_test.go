// Internal tests for types.go.
package shared

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// Test_hasSerializationTags tests the hasSerializationTags function.
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
		tt := tt // Capture range variable
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

// Test_hasSerializationTags_NilFields tests hasSerializationTags with nil fields.
func Test_hasSerializationTags_NilFields(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{
			name:     "nil fields should return false",
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt // Capture range variable
		t.Run(tt.name, func(t *testing.T) {
			// Create struct with nil fields
			structType := &ast.StructType{
				Fields: nil,
			}

			result := hasSerializationTags(structType)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("hasSerializationTags() with nil fields = %v, want %v", result, tt.expected)
			}
		})
	}
}
