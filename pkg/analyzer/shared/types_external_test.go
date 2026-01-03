// External tests for types.go.
package shared_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/kodflow/ktn-linter/pkg/analyzer/shared"
)

// TestIsSerializableStruct tests the IsSerializableStruct function.
func TestIsSerializableStruct(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		structName string
		expected   bool
	}{
		{
			name: "DTO with yaml tag",
			code: `package test
type UserConfig struct {
	Name string ` + "`yaml:\"name\"`" + `
}`,
			structName: "UserConfig",
			expected:   true,
		},
		{
			name: "DTO with json tag",
			code: `package test
type User struct {
	Name string ` + "`json:\"name\"`" + `
}`,
			structName: "User",
			expected:   true,
		},
		{
			name: "Not a DTO",
			code: `package test
type User struct {
	Name string
}`,
			structName: "User",
			expected:   false,
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

			result := shared.IsSerializableStruct(structType, tt.structName)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("IsSerializableStruct() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIsPureDataStruct tests the IsPureDataStruct function.
func TestIsPureDataStruct(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name: "all public fields",
			code: `package test
type User struct {
	Name string
	Age int
}`,
			expected: true,
		},
		{
			name: "mixed fields",
			code: `package test
type User struct {
	Name string
	age int
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

			result := shared.IsPureDataStruct(structType)
			// Vérification résultat
			if result != tt.expected {
				t.Errorf("IsPureDataStruct() = %v, want %v", result, tt.expected)
			}
		})
	}
}
